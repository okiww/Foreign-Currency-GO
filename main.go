package main

import (
	"flag"
	"fmt"
	"learn-viper/config"
	"learn-viper/data"
	dataModel "learn-viper/data/model"
	"learn-viper/module/currency"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

var (
	runMigration       bool
	configuration      config.Configuration
	dbFactory          *data.DBFactory
	currencyController *currency.Controller
)

func init() {
	//flag for migration if set true then running migration
	flag.BoolVar(&runMigration, "migrate", true, "run db migration before starting the server")
	cfg, err := config.New()
	if err != nil {
		glog.Fatalf("Failed to load configuration: %s", err)
		panic(fmt.Errorf("Fatal error loading configuration: %s", err))
	}

	configuration = *cfg
	dbFactory = data.NewDbFactory(configuration.Database)

	currencyController, err = currency.NewController(dbFactory)
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
		apiv1.POST("/currency", currencyController.AddCurrency)
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
	)
	glog.Info("Done running db migration")
}
