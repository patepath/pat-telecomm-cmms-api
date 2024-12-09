package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Profile struct {
	Name1     string `json:"name1"`
	Position1 string `json:"position1"`
	Name2     string `json:"name2"`
	Position2 string `json:"position2"`
	Name3     string `json:"name3"`
	Position3 string `json:"position3"`
	Name4     string `json:"name4"`
	Position4 string `json:"position4"`
	Name5     string `json:"name5"`
	Position5 string `json:"position5"`
	Name6     string `json:"name6"`
	Position6 string `json:"position6"`
}

type ProfileHandle struct {
}

func (p *ProfileHandle) GetOne(c *gin.Context) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print(err)
		return
	}

	var pf Profile
	pfstr := os.Getenv("profile")

	err = json.Unmarshal([]byte(pfstr), &pf)
	if err != nil {
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, pf)
}

func (p *ProfileHandle) Save(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Print(err)
		return
	}

	var pf Profile
	err = json.Unmarshal([]byte(body), &pf)
	if err != nil {
		log.Print(err)
		return
	}

	v, err := json.Marshal(pf)
	if err != nil {
		log.Panic(err)
		return
	}

	err = godotenv.Load(".env")
	if err != nil {
		log.Print(err)
		return
	}

	os.Setenv("profile", string(v))
	updateEnvVariable(".env", "profile", string(v))

	c.JSON(http.StatusOK, pf)
}

func updateEnvVariable(filePath, key, newValue string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			lines[i] = fmt.Sprintf("%s=%s", key, newValue)
			break
		}
	}

	output := strings.Join(lines, "\n")
	err = os.WriteFile(filePath, []byte(output), 0644)
	if err != nil {
		log.Fatalf("Error writing to .env file: %v", err)
	}
}
