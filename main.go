package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/mail.v2"
	"net/http"
	"strings"
)

type task struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type team struct {
	ID    string          `json:"id"`
	Pass  string          `json:"password"`
	Email string          `json:"email"`
	Tasks map[string]task `json:"tasks"`
}

type req struct {
	Pass   string `json:"password"`
	TaskID string `json:"task_id"`
}

type ansReq struct {
	Pass        string `json:"password"`
	TaskID      string `json:"task_id"`
	InputAnswer string `json:"input_answer"`
}

type loginReq struct {
	Pass string `json:"password"`
}

var teams = make(map[string]team)

func getTeams(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, teams)
}

func checkAnsByID(c *gin.Context) {
	id := c.Param("id")
	var newReq ansReq
	if err := c.BindJSON(&newReq); err != nil {
		return
	}
	if !validatePassword(id, newReq.Pass, teams[id].Pass) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}
	if strings.ToLower(newReq.InputAnswer) != strings.ToLower(teams[id].Tasks[newReq.TaskID].Answer) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"answerStatus": 0})
		return
	}
	go func() {
		sub := "Team no.: " + id + " Task No.: " + newReq.TaskID
		if err := sendMail(teams[id].Email, sub, newReq.InputAnswer); err != nil {
			fmt.Println("error")
		} else {
			fmt.Println(sub+" Email sent successfully. " + newReq.InputAnswer)
		}
	}()
	c.IndentedJSON(http.StatusOK, gin.H{"answerStatus": 1})
}

func getQuesByID(c *gin.Context) {
	id := c.Param("id")
	var newReq req
	if err := c.BindJSON(&newReq); err != nil {
		return
	}
	if !validatePassword(id, newReq.Pass, teams[id].Pass) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"question": teams[id].Tasks[newReq.TaskID].Question, "ansLen": len(teams[id].Tasks[newReq.TaskID].Answer)})
}

func teamById(c *gin.Context) {
	id := c.Param("id")
	var newReq req
	if err := c.BindJSON(&newReq); err != nil {
		return
	}
	if !validatePassword(id, newReq.Pass, teams[id].Pass) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": teams[id].Tasks[newReq.TaskID]})

}

func setTeams(c *gin.Context) {
	if err := c.BindJSON(&teams); err != nil {
		return
	}
	c.IndentedJSON(http.StatusOK, teams)
}

func login(c *gin.Context) {
	id := c.Param("id")
	var newReq loginReq
	if err := c.BindJSON(&newReq); err != nil {
		return
	}
	if !validatePassword(id, newReq.Pass, teams[id].Pass) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"login": 1})
}

func validatePassword(id, inputPassword, hashedPassword string) bool {
	return id+inputPassword+id == hashedPassword
}

func sendMail(toAddress, subject, body string) error {
	m := mail.NewMessage()
	m.SetHeader("From", "ankithemapark@gmail.com")
	m.SetHeader("To", toAddress)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := mail.NewDialer("smtp.gmail.com", 587, "ankithemapark@gmail.com", "vcwglezyadtrzbqu") // Replace with your SMTP server and credentials

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func firstPage(c *gin.Context){
	c.IndentedJSON(http.StatusOK, gin.H{"success": 1})
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/",firstPage)
	router.GET("/teams", getTeams)
	router.POST("/teams", setTeams)
	router.POST("/teams/:id", teamById)
	router.POST("/teams/login/:id", login)
	router.POST("/teams/ques/:id", getQuesByID)
	router.POST("/teams/ans/:id", checkAnsByID)
	router.Run("localhost:8080")
}
