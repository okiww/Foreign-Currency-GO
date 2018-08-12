package currency

import (
	"errors"
	"fmt"
	"learn-viper/data"
	dataModel "learn-viper/data/model"
	"net/http"

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
		return nil, errors.New("failed to instantiate rate controller")
	}

	return &Controller{dbFactory: dbFactory}, nil
}

func (ctrl *Controller) ListCurrency(c *gin.Context) {
	db, err := ctrl.dbFactory.DBConnection()
	if err != nil {
		fmt.Println("err")
		glog.Errorf("Failed to open db connection: %s", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var currencies []dataModel.Currency

	// query get currency
	if err := db.Find(&currencies).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "cannot find currency",
		})
	}

	for i, _ := range currencies {
		db.Model(currencies[i]).Related(&currencies[i].Rate, "Rate")
		currencies[i].Rate = currencies[i].Rate
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   currencies,
	})
	return
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
	currency := dataModel.Currency{
		From: req.From,
		To:   req.To,
	}

	db.Save(&currency)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusCreated,
		"data":   currency,
	})
	return
}

func (ctrl *Controller) DeleteCurrency(c *gin.Context) {
	db, err := ctrl.dbFactory.DBConnection()
	if err != nil {
		fmt.Println("err")
		glog.Errorf("Failed to open db connection: %s", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var req deleteCurrencyRequest
	var currency dataModel.Currency
	var rate []dataModel.Rate

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
	db.Delete(&currency)

	if err := db.Where("currency_id = ?", req.CurrencyID).Find(&rate).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": "cannot find rates"})
		return
	}

	db.Delete(&rate)

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "success delete currency",
	})
}
