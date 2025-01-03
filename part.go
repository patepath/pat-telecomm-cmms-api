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

type Part struct {
	IssueId       uint64 `json:"issueid" gorm:"primaryKey"`
	PartProfileId uint64 `json:"partprofileid" gorm:"primaryKey"`
	Rank          int    `json:"rank"`
	Code          string `json:"code" gorm:"size:15"`
	Name          string `json:"name" gorm:"size:255"`
	Unit          string `json:"unit" gorm:"size:15"`
	Qty           int    `json:"qty"`
	Remark        string `json:"remark" gorm:"size:255"`
}

type PartHandler struct {
	DB *gorm.DB
}

func (h *PartHandler) Initialize(dsn string) {
	var err error
	var db *gorm.DB

	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(&Part{})
	h.DB = db
}

func (h *PartHandler) FindAll(c *gin.Context) {
	var parts []Part

	err := h.DB.Order("rank ASC").Find(&parts).Error
	if err != nil {
		log.Panic(err)
	}

	c.JSON(http.StatusOK, parts)
}

func (h *PartHandler) Save(c *gin.Context) {
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

	if claim.Position == "admin" {
		var part Part
		json.Unmarshal(body, &part)

		err := h.DB.Save(&part).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false})
			log.Panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}
