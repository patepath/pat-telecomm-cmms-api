package main

import (
	"cmms-api/token"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Issue struct {
	Id               uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	IssueNo          string `gorm:"size:20" json:"issueno"`
	PhoneId          *uint64
	Phone            Phone  `gorm:"foreignKey:PhoneId" json:"phone"`
	PhoneBy          string `json:"phoneby"`
	TechId           *uint64
	Tech             User      `gorm:"foreignKey:TechId" json:"tech"`
	Created          time.Time `json:"created"`
	IssueType        int       `json:"issuetype"`
	IssueTypeOther   string    `json:"issuetypeother"`
	IssueBy          string    `json:"issueby"`
	IssueContactNo   string    `json:"issuecontactno"`
	IssueDescription string    `json:"issuedescription"`
	IssueLocation    string    `json:"issuelocation"`
	IssueCause       string    `json:"issuecause"`
	IssueSolution    string    `json:"issuesolution"`
	EngineerCode     string    `json:"engineercode"`
	Ext              string    `json:"ext"`
	FinishedDate     time.Time `json:"finisheddate"`
	Status           int       `json:"status"`
	Parts            []Part    `json:"parts" gorm:"foreignKey:IssueId;"`
}

type FileAttach struct {
	IssueNo  string `json:"issueno" gorm:"primaryKey"`
	Order    uint64 `json:"order" gorm:"primaryKey"`
	FilePath string `json:"filepath"`
	B64      string `json:"b64"`
}

//type IssuePart struct {
//	IssueId uint64 `gorm:"primaryKey"`
//	PartId  uint64 `gorm:"primaryKey"`
//	Qty     uint
//}

//type PartUsage struct {
//	IssueId uint64 `gorm:"primaryKey" json:"issueid"`
//	PartId  uint64 `gorm:"primaryKey" json:"partid"`
//	Rank    int    `json:"rank"`
//	Code    string `json:"code"`
//	Name    string `json:"name"`
//	Qty     int    `json:"qty"`
//	Unit    string `json:"unit"`
//	Remark  string `json:"remark"`
//}

type ReportBySummary struct {
	IssueType     string `json:"issue_type"`
	Proceeding    int    `json:"proceeding"`
	WaitForClosed int    `json:"waitforclosed"`
	Closed        int    `json:"closed"`
	Cancelled     int    `json:"cancelled"`
}

type IssueHandler struct {
	DB *gorm.DB
}

func (h *IssueHandler) Initialize(dsn string) {
	var err error
	var db *gorm.DB

	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(&Issue{}, &FileAttach{})
	h.DB = db
}

func (h *IssueHandler) CheckFileAttach(c *gin.Context) {
	var fileAttach []FileAttach
	issueno := c.Param("issueno")

	h.DB.Model(&FileAttach{}).Where("issue_no = ?", issueno).Find(&fileAttach)

	for i, file := range fileAttach {
		data, err := os.ReadFile(file.FilePath)
		if err != nil {
			log.Panic(err)
		}

		fileAttach[i].B64 = base64.StdEncoding.EncodeToString(data)
	}

	c.JSON(http.StatusOK, gin.H{"issue": issueno, "success": true, "data": fileAttach})
}

func (h *IssueHandler) Upload(c *gin.Context) {
	var t = c.Param("token")

	var claim, err = token.VerifyToken(t)
	if err != nil {
		log.Panic(err)
	}

	if claim.Role == 1 || claim.Role == 2 {
		var fileAtt FileAttach

		order := c.PostForm("order")
		issueno := c.PostForm("issueno")
		year := strings.Split(issueno, "-")[0]
		path := "attach/" + year + "/"
		r, _ := strconv.Atoi(order)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, os.ModePerm)
		}

		file, _ := c.FormFile("file")
		ext := filepath.Ext(file.Filename)
		order = "00" + order
		fullpath := path + issueno + "-" + order[len(order)-2:] + ext

		fileAtt.IssueNo = issueno
		fileAtt.Order = uint64(r)
		fileAtt.FilePath = fullpath

		h.DB.Create(&fileAtt)
		c.SaveUploadedFile(file, fullpath)
	}
}

func (h *IssueHandler) Download(c *gin.Context) {
	var t = c.Param("token")

	var claim, err = token.VerifyToken(t)
	if err != nil {
		log.Panic(err)
	}

	if claim.Role == 1 || claim.Role == 2 {
		issueno := c.Param("issueno")
		order := c.Param("order")
		order = "00" + order
		order = "-" + order[len(order)-2:]

		year := strings.Split(issueno, "-")[0]
		path := "attach/" + year + "/"

		files, _ := filepath.Glob(path + issueno + order + ".*")

		if len(files) > 0 {
			data, err := os.ReadFile(files[0])
			ext := filepath.Ext(files[0])

			if err != nil {
				fmt.Println(err)
			}

			base64Data := base64.StdEncoding.EncodeToString(data)
			c.JSON(http.StatusOK, gin.H{"issue": issueno, "success": true, "data": base64Data, "ext": ext})

		} else {
			c.JSON(http.StatusOK, gin.H{"issue": issueno, "success": false, "data": "", "ext": ""})
		}
	}
}

func (h *IssueHandler) Save(c *gin.Context) {
	var t = c.Param("token")
	var isattach = c.Param("isattach")

	var claim, err = token.VerifyToken(t)
	if err != nil {
		log.Panic(err)
	}

	if claim.Role == 1 || claim.Role == 2 {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Panic(err)
		}

		var issue Issue
		json.Unmarshal(body, &issue)

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

		if isattach == "true" {
			year := strings.Split(issue.IssueNo, "-")[0]
			path := "attach/" + year + "/"

			files, _ := filepath.Glob(path + issue.IssueNo + "*.*")

			for _, f := range files {
				if err = os.Remove(f); err != nil {
					panic(err)
				}
			}

			h.DB.Where("issue_no", issue.IssueNo).Delete(&FileAttach{})
		}

		c.JSON(http.StatusOK, gin.H{"issueno": issue.IssueNo, "success": true})
	}
}

func (h *IssueHandler) FindById(c *gin.Context) {
	var t = c.Param("token")
	var id = c.Param("id")

	var claim, err = token.VerifyToken(t)
	if err != nil {
		log.Panic(err)
	}

	if claim.Role == 1 || claim.Role == 2 || claim.Role == 4 {
		var issue Issue

		h.DB.Model(&Issue{}).Preload("Phone").Preload("Tech").Preload("Parts").Find(&Issue{}, id).First(&issue)
		c.JSON(http.StatusOK, issue)
	}
}

func (h *IssueHandler) FindByDate(c *gin.Context) {
	var t = c.Param("token")
	var frm = c.Param("frmdate")
	var to = c.Param("todate")

	var claim, err = token.VerifyToken(t)
	if err != nil {
		log.Panic(err)
	}

	if claim.Role == 1 || claim.Role == 2 || claim.Role == 4 {
		var issues []Issue

		h.DB.Model(&Issue{}).Preload("Phone").Where("status=1 and created between ? and adddate(?, interval 1 DAY)", frm, to).Find(&issues)
		c.JSON(http.StatusOK, issues)
	}
}

func (h *IssueHandler) FindToday(c *gin.Context) {
	var t = c.Param("token")
	var frmdate = c.Param("frmdate")

	var claim, err = token.VerifyToken(t)
	if err != nil {
		log.Panic(err)
	}

	if claim.Role == 1 || claim.Role == 2 {
		var issues []Issue

		h.DB.Model(&Issue{}).Preload("Phone").Where("status < 2 and date(created) = ?", frmdate).Find(&issues)
		c.JSON(http.StatusOK, issues)
	}
}

func (h *IssueHandler) FindOnProcess(c *gin.Context) {
	var t = c.Param("token")

	var claim, err = token.VerifyToken(t)
	if err != nil {
		log.Panic(err)
	}

	if claim.Role == 1 || claim.Role == 2 {
		var issues []Issue

		h.DB.Model(&Issue{}).Preload("Phone").Where("status=0").Find(&issues)
		c.JSON(http.StatusOK, issues)
	}
}

func (h *IssueHandler) FindWaitForClose(c *gin.Context) {
	var t = c.Param("token")

	var claim, err = token.VerifyToken(t)
	if err != nil {
		log.Panic(err)
	}

	if claim.Role == 1 || claim.Role == 2 || claim.Role == 4 {
		var issues []Issue

		h.DB.Model(&Issue{}).Preload("Phone").Where("status=2").Find(&issues)
		c.JSON(http.StatusOK, issues)
	}
}

func (h *IssueHandler) FindCompleted(c *gin.Context) {
	var t = c.Param("token")

	var claim, err = token.VerifyToken(t)
	if err != nil {
		log.Panic(err)
	}

	if claim.Role == 1 || claim.Role == 2 {
		var issues []Issue

		h.DB.Model(&Issue{}).Preload("Phone").Where("status=1").Find(&issues)
		c.JSON(http.StatusOK, issues)
	}
}

func (h *IssueHandler) FindAllByDate(c *gin.Context) {
	var t = c.Param("token")
	var frm = c.Param("frmdate")
	var to = c.Param("todate")

	var claim, err = token.VerifyToken(t)
	if err != nil {
		log.Panic(err)
	}

	if claim.Role == 1 || claim.Role == 2 || claim.Role == 4 {
		var issues []Issue

		h.DB.Model(&Issue{}).Preload("Phone").Where("created between ? and adddate(?, interval 1 DAY)", frm, to).Find(&issues)
		c.JSON(http.StatusOK, issues)
	}
}

func (h *IssueHandler) SummaryByDate(c *gin.Context) {
	var t = c.Param("token")
	var frm = c.Param("frmdate")
	var to = c.Param("todate")
	var result []ReportBySummary

	var claim, err = token.VerifyToken(t)
	if err != nil {
		log.Panic(err)
	}

	if claim.Role == 1 || claim.Role == 2 || claim.Role == 4 {
		var sql = `
			select 
				issue_type, 
				sum(case when status = 0 then  1 else 0 end) as proceeding,
				sum(case when status = 1 then  1 else 0 end) as closed,
				sum(case when status = 2 then  1 else 0 end) as waitforclosed,
				sum(case when status = 99 then  1 else 0 end) as cancelled
			from issues 
			where DATE(created) between ? and adddate(?, interval 1 DAY)
			group by issue_type; 
		`

		h.DB.Raw(sql, frm, to).Scan(&result)
		c.JSON(http.StatusOK, result)
	}
}
