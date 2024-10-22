package main

import (
	"cmms-api/token"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Phone struct {
	Id       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Number   string `gorm:"size:25" json:"number"`
	Location string `gorm:"size:255" json:"location"`
	Hc       string `gorm:"size:30" json:"hc"`
	Kc       string `gorm:"size:30" json:"kc"`
	Tc1      string `gorm:"size:15" json:"tc1"`
	Tc2      string `gorm:"size:15" json:"tc2"`
	Tc3      string `gorm:"size:15" json:"tc3"`
	Tc4      string `gorm:"size:15" json:"tc4"`
	Tc5      string `gorm:"size:15" json:"tc5"`
}

type PhoneHandler struct {
	DB *gorm.DB
}

func (h *PhoneHandler) Initialize(dsn string) {
	var err error
	var db *gorm.DB

	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(&Phone{})
	h.DB = db
}

func (h *PhoneHandler) Save(c *gin.Context) {
	var t = c.Param("token")

	var claim, err = token.VerifyToken(t)
	if err != nil {
		log.Panic(err)
	}

	if claim.Role == 1 || claim.Role == 2 {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Panic(err)
		}

		var phone Phone
		json.Unmarshal(body, &phone)

		h.DB.Save(&phone)
		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func (h *PhoneHandler) FindAll(c *gin.Context) {
	var phones []Phone

	err := h.DB.Find(&phones).Error
	if err != nil {
		log.Panic(err)
	}

	c.JSON(http.StatusOK, phones)
}

func (h *PhoneHandler) FindByNumber(c *gin.Context) {
	var num = c.Param("num")
	var phones []Phone

	err := h.DB.Where("number like ?", num+"%").Find(&phones).Error
	if err != nil {
		log.Panic(err)
	}

	c.JSON(http.StatusOK, phones)
}
