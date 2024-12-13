package main

import (
	"cmms-api/token"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type PartProfile struct {
	Id   int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Rank int    `json:"rank"`
	Code string `json:"code" gorm:"size:15"`
	Name string `json:"name" gorm:"size:255"`
	Unit string `json:"unit" gorm:"size:15"`
}

type PartProfileHandler struct {
	DB *gorm.DB
}

func (h *PartProfileHandler) Initialize(dsn string) {
	var err error
	var db *gorm.DB

	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(&PartProfile{})
	h.DB = db
}

func (h *PartProfileHandler) FindAll(c *gin.Context) {
	var partprofiles []PartProfile

	err := h.DB.Order("rank ASC").Find(&partprofiles).Error
	if err != nil {
		log.Panic(err)
	}

	c.JSON(http.StatusOK, partprofiles)
}

func (h *PartProfileHandler) Save(c *gin.Context) {

	var t = c.Param("token")
	var claim, err = token.VerifyToken(t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		log.Panic(err)
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		log.Panic(err)
	}

	if claim.Role == 1 {
		var partprofile PartProfile
		json.Unmarshal(body, &partprofile)

		fmt.Printf("%v", partprofile)

		err := h.DB.Save(&partprofile).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false})
			log.Panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}
