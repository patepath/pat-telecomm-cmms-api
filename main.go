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

	h := LoginHandler{}
	h.Initialize(dsn)

	part := PartHandler{}
	part.Initialize(dsn)

	phone := PhoneHandler{}
	phone.Initialize(dsn)

	issue := IssueHandler{}
	issue.Initialize(dsn)

	r.POST("/login", h.Login)
	r.GET("/checktoken/:token", h.CheckToken)

	userGroup := r.Group("/user")
	{
		userGroup.GET("/findall/:token", h.FindAll)
		userGroup.GET("/gettech/:token", h.GetTech)
		userGroup.POST("/add", h.Add)
		userGroup.POST("/save/:token", h.Save)
		userGroup.POST("/changepassword/:token", h.ChangePassword)
	}

	partGroup := r.Group("/part")
	{
		partGroup.GET("/findall", part.FindAll)
		partGroup.POST("/save/:token", part.Save)
	}

	phoneGroup := r.Group("phone")
	{
		phoneGroup.GET("/findall", phone.FindAll)
		phoneGroup.GET("/findbynumber/:num", phone.FindByNumber)
		phoneGroup.POST("/save/:token", phone.Save)
	}

	issueGroup := r.Group("/issue")
	{
		issueGroup.GET(("/findbyid/:token/:id"), issue.FindById)
		issueGroup.GET(("/findtoday/:token"), issue.FindToday)
		issueGroup.GET(("/findbydate/:token/:frmdate/:todate"), issue.FindByDate)
		issueGroup.GET(("/findallbydate/:token/:frmdate/:todate"), issue.FindAllByDate)
		issueGroup.GET(("/findonprocess/:token"), issue.FindOnProcess)
		issueGroup.GET(("/findwaitforclose/:token"), issue.FindWaitForClose)
		issueGroup.GET(("/findcompleted/:token"), issue.FindCompleted)
		issueGroup.GET(("/summarybydate/:token/:frmdate/:todate"), issue.SummaryByDate)
		issueGroup.POST(("/save/:token/:isattach/:isparts"), issue.Save)
		issueGroup.GET(("/checkfileattach/:issueno"), issue.CheckFileAttach)
		issueGroup.POST(("/upload/:token"), issue.Upload)
		issueGroup.GET(("/download/:issueno/:order"), issue.Download)
	}

	profile := ProfileHandle{}
	profilegroup := r.Group("/profile")
	{
		profilegroup.GET(("/getone"), profile.GetOne)
		profilegroup.POST(("/save"), profile.Save)
	}

	r.Run(":8082")
}
