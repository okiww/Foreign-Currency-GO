# Foreign-Currency-GO
Minimum build API for foreign currency BE Exercise using Go language powered by [Gin web framework](https://github.com/gin-gonic/gin) and [GORM](https://github.com/jinzhu/gorm) as backend DAL layer.

## Run the code

The first clone or download this project:
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

## API Documentation

By default the app will listen on all interface at port 8080. Here is the list of endpoint curently available

* Ping endpoint `GET /ping`
* Add data currency `POST /api/v1/currency/add`
* List all currency `GET /api/v1/currency/list`
* Delete currency  `DELETE /api/v1/currency/delete`
* Add rate currency `POST /api/v1/rate/add`
* List currency by date `GET /api/v1/rate/list` make sure you add query date example `?date=2018-01-01"`
* List most 7 data point of rate by date `POST /api/v1/rate/most`

## DB SCHEME Documentation

We have simply 2 table for this API. here :

![alt text](https://github.com/okiww/Foreign-Currency-GO/blob/master/db_scheme.png)