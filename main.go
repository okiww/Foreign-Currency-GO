package main

import (
	"flag"
	"fmt"
	"learn-viper/config"
	"learn-viper/data"
	dataModel "learn-viper/data/model"
	"learn-viper/module/currency"
	"learn-viper/module/rate"
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
	flag.BoolVar(&runMigration, "migrate", true, "run db migration before starting the server")
	flag.BoolVar(&runSeeder, "seeder", false, "run db seeder before starting the server")

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
func setupRouter() *gin.Engine {
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
		c.String(200, "pong")
	})

	// api v1 endpoint
	apiv1 := router.Group("/api/v1")
	{
		currencyGroup := apiv1.Group("/currency")
		{
			currencyGroup.POST("/", currencyController.AddCurrency)
			currencyGroup.GET("/list", currencyController.ListCurrency)
			currencyGroup.DELETE("/delete", currencyController.DeleteCurrency)
		}

		rateGroup := apiv1.Group("/rate")
		{
			rateGroup.POST("/", rateController.AddRate)
			rateGroup.GET("/", rateController.GetListCurrencyByDate)
		}

	}

	return router
}

func main() {
	r := setupRouter()

	srv := &http.Server{
		Addr:    configuration.Server.Port,
		Handler: r,
	}
	// Listen and Server in 0.0.0.0:8080
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			glog.Fatalf("Failed to start server: %s", err)
			panic(fmt.Errorf("Fatal error failed to start server: %s", err))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	glog.Info("Shutting down server")
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
