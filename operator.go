package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Operator struct {
	Id           uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	PhoneNumber  string `gorm:"size:50" json:"phonenumber"`
	NearByNumber string `gorm:"size:50" json:"nearbynumber"`
	StaffName    string `gorm:"size:100" json:"staffname"`
	Position     string `gorm:"size:100" json:"position"`
	Unit         string `gorm:"size:100" json:"unit"`
	Department   string `gorm:"size:100" json:"department"`
	Division     string `gorm:"size:100" json:"division"`
	Organization string `gorm:"size:100" json:"organization"`
}

type OperatorHandler struct {
	DB *gorm.DB
}

func (h *OperatorHandler) Initialize(dsn string) {
	var err error
	var db *gorm.DB

	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(&Operator{})
	h.DB = db
}

func (h *OperatorHandler) FindAll(c *gin.Context) {
	var operators []Operator

	err := h.DB.Find(&operators).Error
	if err != nil {
		log.Panic(err)
	}

	c.JSON(200, operators)
}

func (h *OperatorHandler) Save(c *gin.Context) {
	var operator Operator

	if err := c.ShouldBindJSON(&operator); err != nil {
		log.Panic(err)
	}

	err := h.DB.Save(&operator).Error
	if err != nil {
		log.Panic(err)
	}
}
