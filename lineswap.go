package main

import (
	"cmms-api/token"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type LineSwap struct {
	Id               uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	IssueNo          string    `json:"issueno" gorm:"size:20"`
	Phone            Phone     `json:"phone" gorm:"foreignKey:PhoneId"`
	PhoneId          *uint64   `json:"phoneid"`
	PhoneBy          string    `json:"phoneby"`
	Tech             User      `json:"tech" gorm:"foreignKey:TechId"`
	TechId           *uint64   `json:"techid"`
	Created          time.Time `json:"created"`
	IssueType        int       `json:"issuetype"`
	IssueTypeOther   string    `json:"issuetypeother"`
	IssueBy          string    `json:"issueby"`
	IssueContactNo   string    `json:"issuecontactno"`
	IssueLocation    string    `json:"issuelocation"`
	IssueDescription string    `json:"issuedescription"`
	IssueCause       string    `json:"issuecause"`
	IssueSolution    string    `json:"issuesolution"`
	EngineerCode     string    `json:"engineercode"`
	Ext              string    `json:"ext"`
	FinishedDate     time.Time `json:"finisheddate"`
	Status           int       `json:"status"`
}

type LineSwapHandler struct {
	DB *gorm.DB
}

func (h *LineSwapHandler) Initialize(dsn string) {
	var err error
	var db *gorm.DB

	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(&LineSwap{}, &FileAttach{})
	h.DB = db
}

func (h *LineSwapHandler) Save(c *gin.Context) {
	var t = c.Param("token")

	var claim, err = token.VerifyToken(t)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Role: %d\n", claim.Role)

	if claim.Role == 1 || claim.Role == 4 {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Panic(err)
		}

		var issue LineSwap
		json.Unmarshal(body, &issue)

		fmt.Print("%V", issue)

		if issue.IssueNo == "" {
			var s Issue
			var count int64

			year, _, _ := time.Now().Date()
			h.DB.Where("year(created)=? and issue_no <> ''", year).Order("issue_no desc").Last(&s).Count(&count)

			if count == 0 {
				issue.IssueNo = fmt.Sprintf("%d-00001", year)

			} else {
				a, _ := strconv.Atoi(strings.Split(s.IssueNo, "-")[1])
				b := fmt.Sprintf("0000%d", a+1)
				issue.IssueNo = fmt.Sprintf("%d-%s", year, b[len(b)-5:])
			}
		}

		h.DB.Where("issue_id=?", issue.Id).Delete(&Part{})
		h.DB.Save(&issue)
		h.DB.Save(&issue.Phone)

		c.JSON(http.StatusOK, gin.H{"lineswap": issue.IssueNo, "success": true})
	}
}

func (h *LineSwapHandler) FindById(c *gin.Context) {
	idParam := c.Param("id")

	var lineSwap LineSwap
	if err := h.DB.Preload("Phone").Preload("Tech").First(&lineSwap, idParam).Error; err != nil {
		c.JSON(404, gin.H{"error": "LineSwap not found"})
		return
	}

	c.JSON(200, lineSwap)
}

func (h *LineSwapHandler) FindByDate(c *gin.Context) {
	frmDateParam := c.Param("frmdate")
	toDateParam := c.Param("todate")

	frmDate, err := time.Parse("2006-01-02", frmDateParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid from date"})
		return
	}

	toDate, err := time.Parse("2006-01-02", toDateParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid to date"})
		return
	}

	var lineSwaps []LineSwap
	if err := h.DB.Preload("Phone").Preload("Tech").Where("created BETWEEN ? AND ?", frmDate, toDate).Find(&lineSwaps).Error; err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	c.JSON(200, lineSwaps)
}

func (h *LineSwapHandler) FindToday(c *gin.Context) {
	frmDateParam := c.Param("frmdate")

	frmDate, err := time.Parse("2006-01-02", frmDateParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid date"})
		return
	}

	startOfDay := frmDate
	endOfDay := frmDate.Add(24 * time.Hour)

	var lineSwaps []LineSwap
	if err := h.DB.Preload("Phone").Preload("Tech").Where("created BETWEEN ? AND ?", startOfDay, endOfDay).Find(&lineSwaps).Error; err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	c.JSON(200, lineSwaps)
}
