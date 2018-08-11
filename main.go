package main

import (
	"flag"
	"fmt"
	"learn-viper/config"
	"learn-viper/data"
	dataModel "learn-viper/data/model"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

var DB = make(map[string]string)

var (
	runMigration  bool
	configuration config.Configuration
	dbFactory     *data.DBFactory
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
}
func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

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
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := DB[user]
		if ok {
			c.JSON(200, gin.H{"user": user, "value": value})
		} else {
			c.JSON(200, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			DB[user] = json.Value
			c.JSON(200, gin.H{"status": "ok"})
		}
	})

	return r
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
