package rate

import (
	"errors"
	"fmt"
	"learn-viper/data"
	dataModel "learn-viper/data/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/glog"
	validator "gopkg.in/go-playground/validator.v8"
)

type Controller struct {
	dbFactory *data.DBFactory
}

func NewController(dbFactory *data.DBFactory) (*Controller, error) {
	if dbFactory == nil {
		return nil, errors.New("failed to instantiate survey controller")
	}

	return &Controller{dbFactory: dbFactory}, nil
}

func (ctrl *Controller) AddRate(c *gin.Context) {
	db, err := ctrl.dbFactory.DBConnection()
	if err != nil {
		fmt.Println("err")
		glog.Errorf("Failed to open db connection: %s", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var req rateRequest
	var currency dataModel.Currency

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		var errors []string
		ve, ok := err.(validator.ValidationErrors)
		if ok {
			for _, v := range ve {
				errors = append(errors, fmt.Sprintf("%s is %s", v.Field, v.Tag))
			}
		} else {
			errors = append(errors, err.Error())
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	if err := db.Where("id = ?", req.CurrencyID).Find(&currency).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": "cannot find currency"})
		return
	}

	//save rate
	date, _ := time.Parse("2006-01-02", req.Date)
	rate := dataModel.Rate{
		CurrencyID: req.CurrencyID,
		Rate:       req.Rate,
		Date:       date,
	}

	db.Save(&rate)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusCreated,
		"data":   rate,
	})
	return
}

func (ctrl *Controller) GetListCurrencyByDate(c *gin.Context) {
	db, err := ctrl.dbFactory.DBConnection()
	if err != nil {
		fmt.Println("err")
		glog.Errorf("Failed to open db connection: %s", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	param := c.Query("date")
	var rates []dataModel.Rate
	var currencies []dataModel.Currency
	var resp []rateResponse
	date, _ := time.Parse("2006-01-02", param)

	// query get currency
	if err := db.Find(&currencies).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "cannot find currency",
		})
	}

	//get rate and average
	for i, _ := range currencies {

		//query get exchange rate between param date
		if err := db.Where("date BETWEEN ? AND ? AND currency_id = ?", date.AddDate(0, 0, -7), date.AddDate(0, 0, 1), currencies[i].ID).Find(&rates).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": "cannot find currency"})
			return
		}

		rateData := "0"
		avg := 0.0
		avgData := ""

		//validate rate
		if len(rates) == 7 {
			for _, v := range rates {
				rateData = fmt.Sprintf("%f", v.Rate)
				avg = avg + v.Rate
			}
			avgData = fmt.Sprintf("%f", avg/7)
		} else {
			rateData = "insuficient data"
			avgData = ""
		}

		data := rateResponse{
			From: currencies[i].From,
			To:   currencies[i].To,
			Rate: rateData,
			Avg:  avgData,
		}
		resp = append(resp, data)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusCreated,
		"data":   resp,
	})
	return
}