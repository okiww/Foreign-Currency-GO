package main

import (
	"Foreign-Currency-GO/config"
	"Foreign-Currency-GO/data"
	dataModel "Foreign-Currency-GO/data/model"
	"Foreign-Currency-GO/module/currency"
	"Foreign-Currency-GO/module/rate"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

var (
	runMigration       bool
	runSeeder          bool
	configuration      config.Configuration
	dbFactory          *data.DBFactory
	currencyController *currency.Controller
	rateController     *rate.Controller
)

func init() {
	//flag for migration and seeder if set true then running migration and seeder
	flag.BoolVar(&runMigration, "migrate", false, "run db migration before starting the server")
	flag.BoolVar(&runSeeder, "seeder", false, "run db seeder before starting the server")
	flag.Parse()

	cfg, err := config.New()
	if err != nil {
		glog.Fatalf("Failed to load configuration: %s", err)
		panic(fmt.Errorf("Fatal error loading configuration: %s", err))
	}

	configuration = *cfg
	dbFactory = data.NewDbFactory(configuration.Database)

	//inject dbFactory to currency controller
	currencyController, err = currency.NewController(dbFactory)
	if err != nil {
		glog.Fatal(err.Error())
		panic(fmt.Errorf("Fatal error: %s", err))
	}

	//inject dbFactory to rate controller
	rateController, err = rate.NewController(dbFactory)
	if err != nil {
		glog.Fatal(err.Error())
		panic(fmt.Errorf("Fatal error: %s", err))
	}

}
func SetupRouter() *gin.Engine {
	// Disable Console Color
	router := gin.New()

	db, err := dbFactory.DBConnection()
	if err != nil {
		glog.Fatalf("Failed to open database connection: %s", err)
		panic(fmt.Errorf("Fatal error connecting to database: %s", err))
	}
	defer db.Close()
	if runMigration {
		runDBMigration()
	}
	// Ping test
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// api v1 endpoint
	apiv1 := router.Group("/api/v1")
	{
		//grouping by currency
		currencyGroup := apiv1.Group("/currency")
		{
			currencyGroup.POST("/add", currencyController.AddCurrency)
			currencyGroup.GET("/list", currencyController.ListCurrency)
			currencyGroup.DELETE("/delete", currencyController.DeleteCurrency)
		}
		//grouping by rate
		rateGroup := apiv1.Group("/rate")
		{
			rateGroup.POST("/add", rateController.AddRate)
			rateGroup.GET("/list", rateController.GetListCurrencyByDate)
			rateGroup.POST("/most", rateController.GetMost7DataPointByCurrency)
		}

	}

	return router
}

func main() {
	r := SetupRouter()

	srv := &http.Server{
		Addr:    configuration.Server.Port,
		Handler: r,
	}
	// Listen and Serve in 0.0.0.0:8080
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			glog.Fatalf("Failed to start server: %s", err)
			panic(fmt.Errorf("Fatal error failed to start server: %s", err))
		}
	}()

	// wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	glog.Info("Server shutted down")
}

func runDBMigration() {
	glog.Info("Running db migration")
	db, err := dbFactory.DBConnection()
	if err != nil {
		glog.Fatalf("Failed to open database connection: %s", err)
		panic(fmt.Errorf("Fatal error connecting to database: %s", err))
	}
	defer db.Close()

	db.AutoMigrate(
		&dataModel.Currency{},
		&dataModel.Rate{},
	)
	glog.Info("Done running db migration")

	if runSeeder {
		glog.Info("Running db seeder")
		var count int
		db.Model(&dataModel.Currency{}).Count(&count)
		if count == 0 {
			glog.V(1).Info("Running db seeder for table Currency")
			currency := dataModel.Currency{
				From: "USD",
				To:   "IDR",
			}
			db.Create(&currency)

			db.Model(&dataModel.Rate{}).Count(&count)
			rates := []float64{0.008, 0.007, 0.009, 0.009, 0.089, 0.079, 0.010}
			if count == 0 {
				glog.V(1).Info("Running db seeder for table Rate")
				date := time.Now().UTC()
				i := 0
				for _, c := range rates {
					i++
					rate := strconv.FormatFloat(c, 'f', 3, 32)
					f, _ := strconv.ParseFloat(rate, 64)
					rates := dataModel.Rate{
						CurrencyID: currency.ID,
						Rate:       f,
						Date:       date.AddDate(0, 0, i),
					}
					db.Create(&rates)
				}
			}
		}
	}
}
