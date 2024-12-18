package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panic(err)
	}

	dsn := os.Getenv("dsn")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(CORSMiddleware())

	login := LoginHandler{}
	login.Initialize(dsn)
	r.POST("/login", login.Login)
	r.GET("/checktoken/:token", login.CheckToken)

	userGroup := r.Group("/user")
	{
		userGroup.GET("/findall/:token", login.FindAll)
		userGroup.GET("/gettech/:token", login.GetTech)
		userGroup.POST("/add", login.Add)
		userGroup.POST("/save/:token", login.Save)
		userGroup.POST("/changepassword/:token", login.ChangePassword)
	}

	issue := IssueHandler{}
	issue.Initialize(dsn)
	issueGroup := r.Group("/issue")
	{
		issueGroup.GET(("/findbyid/:token/:id"), issue.FindById)
		issueGroup.GET(("/findtoday/:token/:frmdate"), issue.FindToday)
		issueGroup.GET(("/findbydate/:token/:frmdate/:todate"), issue.FindByDate)
		issueGroup.GET(("/findallbydate/:token/:frmdate/:todate"), issue.FindAllByDate)
		issueGroup.GET(("/findonprocess/:token"), issue.FindOnProcess)
		issueGroup.GET(("/findwaitforclose/:token"), issue.FindWaitForClose)
		issueGroup.GET(("/findcompleted/:token"), issue.FindCompleted)
		issueGroup.GET(("/summarybydate/:token/:frmdate/:todate"), issue.SummaryByDate)
		issueGroup.POST(("/save/:token/:isattach"), issue.Save)
		issueGroup.GET(("/checkfileattach/:issueno"), issue.CheckFileAttach)
		issueGroup.POST(("/upload/:token"), issue.Upload)
		issueGroup.GET(("/download/:issueno/:order"), issue.Download)
	}

	phone := PhoneHandler{}
	phone.Initialize(dsn)
	phoneGroup := r.Group("phone")
	{
		phoneGroup.GET("/findall", phone.FindAll)
		phoneGroup.GET("/findbynumber/:num", phone.FindByNumber)
		phoneGroup.POST("/save/:token", phone.Save)
	}

	part := PartHandler{}
	part.Initialize(dsn)
	partGroup := r.Group("/part")
	{
		partGroup.GET("/findall", part.FindAll)
		partGroup.POST("/save/:token", part.Save)
	}

	partprofile := PartProfileHandler{}
	partprofile.Initialize(dsn)
	partprofileGroup := r.Group("/partprofile")
	{
		partprofileGroup.GET("/findall", partprofile.FindAll)
		partprofileGroup.POST("/save/:token", partprofile.Save)
	}

	profile := ProfileHandle{}
	profilegroup := r.Group("/profile")
	{
		profilegroup.GET(("/getone"), profile.GetOne)
		profilegroup.POST(("/save"), profile.Save)
	}

	r.Run("192.168.0.10:8082")
}
