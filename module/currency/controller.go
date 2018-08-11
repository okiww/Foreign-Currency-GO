package currency

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

func (ctrl *Controller) AddCurrency(c *gin.Context) {
	db, err := ctrl.dbFactory.DBConnection()
	if err != nil {
		fmt.Println("err")
		glog.Errorf("Failed to open db connection: %s", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var req currencyRequest
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
	//save currency
	date, err := time.Parse("2006-01-02", req.Date)
	currency := dataModel.Currency{
		From: req.From,
		To:   req.To,
		Date: date,
		Rate: req.Rate,
	}

	db.Save(&currency)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusCreated,
		"data":   currency,
	})
	return
}
