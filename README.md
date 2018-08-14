# Foreign-Currency-GO
minimum build API for foreign currency BE Exercise Web API application using Go language powered by [Gin web framework](https://github.com/gin-gonic/gin) and [GORM](https://github.com/jinzhu/gorm) as backend DAL layer.

## Run the code

The first is, clone or download this project:
`$ git clone git@github.com:okiww/Foreign-Currency-GO.git && cd Foreign-Currency-GO`

There are two ways to run the code
### The easy way using docker

If you have docker and docker-compose installed, simply run:

`$ docker-compose up`

Wait until you see line similar to `Starting Foreign-Currency-GO server version 1.0.0 at :8080`. The app is ready for you to use and listening on port 8080 by default.

### The hard way

You will need Go installed in your local machine

* Install dep dependency manager

  `$ go get -u github.com/golang/dep/cmd/dep`

* Run dep ensure to downloads all dependencies

  `$ dep ensure`

* Copy file default.yml into .env.yml and modify the config to suit your environment

  `$ cp default.yml .env.yml`

* Ensure your database server is running and application table of your choice (by default it is Foreign-Currency-GO, you can change it in .env.yml file) is exist

* Run the app. For first run you may want to add `-migrate` switch to run auto db migration.

  `$ go run main.go -migrate`
