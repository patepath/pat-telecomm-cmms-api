package main

import (
	"cmms-api/token"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Id        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:25;not null" json:"name"`
	Password  string    `gorm:"size:255;not null" json:"password"`
	Role      uint8     `gorm:"not null" json:"role"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Position  string    `gorm:"size:25;not null" json:"position"`
	Status    uint8     `gorm:"not null" json:"status"`
	Modified  time.Time `gorm:"autoUpdateTime" json:"modified"`
}

type LoginHandler struct {
	DB *gorm.DB
}

func (h *LoginHandler) Initialize(dsn string) {
	var err error
	var db *gorm.DB

	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(&User{})
	h.DB = db
}

func (h *LoginHandler) Login(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Panic(err)
	}

	var user User
	json.Unmarshal(body, &user)

	var username = user.Name
	var passwordMD5 = user.Password

	if h.DB.Where("name=? and password=?", username, passwordMD5).First(&user).RowsAffected == 1 {
		t, _ := token.GenerateToken(user.Name, user.Position, user.Role)
		c.JSON(http.StatusOK, gin.H{"fullname": user.FirstName + " " + user.LastName, "position": user.Position, "role": user.Role, "token": t})

	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
	}
}

func (h *LoginHandler) Add(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Panic(err)
	}

	var user, u User
	json.Unmarshal(body, &user)

	if h.DB.Where(&User{Name: user.Name}).First(&u).RowsAffected == 0 {
		h.DB.Create(&user)
		c.JSON(http.StatusOK, user)

	} else {
		c.JSON(http.StatusNotAcceptable, nil)
	}
}

func (h *LoginHandler) Save(c *gin.Context) {
	var t = c.Param("token")
	var claim, err = token.VerifyToken(t)
	if err != nil {
		log.Panic(err)
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Panic(err)
	}

	var user User
	var u User
	var count int64
	json.Unmarshal(body, &user)

	if user.Name == claim.Username || claim.Role == 1 {
		h.DB.Find(&User{Id: user.Id}).First(&u).Count(&count)

		if count == 0 {
			c.JSON(http.StatusOK, gin.H{"success": false})

		} else {
			user.Password = u.Password
			h.DB.Save(&user)
			c.JSON(http.StatusOK, gin.H{"success": true})
		}

	} else {
		c.JSON(http.StatusOK, gin.H{"success": false})
	}
}

func (h *LoginHandler) ChangePassword(c *gin.Context) {
	var t = c.Param("token")
	var claim, err = token.VerifyToken(t)
	if err != nil {
		log.Panic(err)
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Panic(err)
	}

	var user User
	json.Unmarshal(body, &user)

	if user.Name == claim.Username || claim.Role == 1 {
		h.DB.Model(&user).Updates(User{Password: user.Password})
		c.JSON(http.StatusOK, gin.H{"sueecess": true})

	} else {
		c.JSON(http.StatusOK, gin.H{"sueecess": false})
	}
}

func (h *LoginHandler) GetTech(c *gin.Context) {
	var t = c.Param("token")
	var claim, err = token.VerifyToken(t)
	if err != nil {
		log.Panic(err)
	}

	if claim.Role == 1 || claim.Role == 2 {
		var user []User

		h.DB.Where("role=2").Find(&user)
		c.JSON(http.StatusOK, user)
	}
}

func (h *LoginHandler) FindAll(c *gin.Context) {
	var t = c.Param("token")
	var claim, err = token.VerifyToken(t)
	if err != nil {
		log.Panic(err)
	}

	if claim.Role == 1 {
		var users []User
		h.DB.Where("status=1").Find(&users)
		c.JSON(http.StatusOK, users)

	} else {
		c.JSON(http.StatusUnauthorized, nil)
	}
}

func (h *LoginHandler) CheckToken(c *gin.Context) {
	var t = c.Param("token")
	var claim, err = token.VerifyToken(t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		log.Panic(err)
	}

	var user User
	err = h.DB.Where("name=?", claim.Username).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		log.Panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"fullname": user.FirstName + " " + user.LastName, "position": user.Position, "role": user.Role, "token": t})
}
