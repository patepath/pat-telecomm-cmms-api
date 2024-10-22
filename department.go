package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Department struct {
	Id   int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Code string `gorm:"size:10" json:"code"`
	Name string `gorm:"size:255" json:"name"`
}

type DeptHandler struct {
	DB *gorm.DB
}

func (h *DeptHandler) Initialize(dsn string) {
	var err error
	var db *gorm.DB

	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(&Department{})
	h.DB = db
}

func (h *DeptHandler) FindAll(c *gin.Context) {
	var depts []Department
	err := h.DB.Find(&depts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		log.Panic(err)
	}

	c.JSON(http.StatusOK, depts)
}
