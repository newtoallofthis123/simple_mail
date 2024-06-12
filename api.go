package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type ApiServer struct {
	Port       string
	Sender     *Sender
	Db         *DbInstance
	AuthSecret string
}

func NewApiServer(env *Env) *ApiServer {
	db := NewDbInstance(env.DatabaseUrl)
	db.Prep()

	sender := NewSender(env.Mail, env.MailPassword, "smtp.gmail.com", "smtp.gmail.com:587")

	return &ApiServer{Port: env.Port, Sender: sender, Db: db, AuthSecret: env.AuthSecret}
}

func (a *ApiServer) handleEmail(c *gin.Context) {
	// bearer token AuthSecret
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	if token != a.AuthSecret {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var e CreateEmailEntry
	if err := c.BindJSON(&e); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := a.Db.InsertEmail(&e); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	email := NewEmail(e.Email, e.Subject, e.Body)
	if err := a.Sender.Send(email); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Email sent!"})
}

func (a *ApiServer) handleGetEmails(c *gin.Context) {
	emails, err := a.Db.GetEmails()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, emails)
}

func (a *ApiServer) handleGetById(c *gin.Context) {
	id := c.Param("id")

	email, err := a.Db.GetEmail(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, email)
}

func (a *ApiServer) handleGetByMail(c *gin.Context) {
	mail := c.Param("mail")

	emails, err := a.Db.GetByEmail(mail)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, emails)
}

func (a *ApiServer) Start() {
	r := gin.Default()

	r.POST("/email", a.handleEmail)
	r.GET("/emails", a.handleGetEmails)
	r.GET("/email/:id", a.handleGetById)
	r.GET("/email/mail/:mail", a.handleGetByMail)

	r.Run(":" + a.Port)
}
