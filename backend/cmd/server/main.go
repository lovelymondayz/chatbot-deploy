package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "/app/data/chatbot.db"
	}
	os.MkdirAll(filepath.Dir(dbPath), 0755)

	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	db.Exec(`CREATE TABLE IF NOT EXISTS clients (
		id TEXT PRIMARY KEY, name TEXT, status TEXT DEFAULT 'active', created_at TEXT
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS knowledge_docs (
		id TEXT PRIMARY KEY, client_id TEXT, content TEXT, doc_type TEXT, created_at TEXT
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS conversations (
		id TEXT PRIMARY KEY, client_id TEXT, user_message TEXT, bot_response TEXT, created_at TEXT
	)`)
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"service": "Chatbot Deploy API", "version": "1.0.0"})
}

type ClientReq struct {
	Name string `json:"name" binding:"required"`
}

func addClient(c *gin.Context) {
	var req ClientReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cid := uuid.New().String()
	_, err := db.Exec("INSERT INTO clients (id,name,status,created_at) VALUES (?,?,?,?)",
		cid, req.Name, "active", time.Now().Format(time.RFC3339))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"client_id":  cid,
		"message":    fmt.Sprintf("Client %s created", req.Name),
		"widget_code": fmt.Sprintf("<script src=\"https://chatbot.arjism.com/widget.js\" data-client-id=\"%s\"></script>", cid),
	})
}

func listClients(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM clients ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var clients []gin.H
	for rows.Next() {
		var id, name, status, createdAt string
		rows.Scan(&id, &name, &status, &createdAt)
		clients = append(clients, gin.H{"id": id, "name": name, "status": status, "created_at": createdAt})
	}
	c.JSON(http.StatusOK, gin.H{"clients": clients})
}

type TrainReq struct {
	ClientID string `json:"client_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
	DocType  string `json:"doc_type"`
}

func train(c *gin.Context) {
	var req TrainReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id, name, status, createdAt string
	err := db.QueryRow("SELECT * FROM clients WHERE id=?", req.ClientID).Scan(&id, &name, &status, &createdAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	docType := req.DocType
	if docType == "" {
		docType = "text"
	}

	did := uuid.New().String()
	_, err = db.Exec("INSERT INTO knowledge_docs (id,client_id,content,doc_type,created_at) VALUES (?,?,?,?,?)",
		did, req.ClientID, req.Content, docType, time.Now().Format(time.RFC3339))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"doc_id": did, "message": "Document trained"})
}

type ChatReq struct {
	ClientID string `json:"client_id" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

func chat(c *gin.Context) {
	var req ChatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id, name, status, createdAt string
	err := db.QueryRow("SELECT * FROM clients WHERE id=?", req.ClientID).Scan(&id, &name, &status, &createdAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	rows, err := db.Query("SELECT content FROM knowledge_docs WHERE client_id=? ORDER BY created_at", req.ClientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var docs []string
	for rows.Next() {
		var content string
		rows.Scan(&content)
		docs = append(docs, content)
	}

	context := "No knowledge base yet."
	if len(docs) > 0 {
		context = ""
		for i, d := range docs {
			if i > 0 {
				context += "\n"
			}
			context += d
		}
	}

	response := fmt.Sprintf("I understand you're asking about: %s\n\nBased on my knowledge: %s...", req.Message, context[:min(200, len(context))])

	convID := uuid.New().String()
	_, err = db.Exec("INSERT INTO conversations (id,client_id,user_message,bot_response,created_at) VALUES (?,?,?,?,?)",
		convID, req.ClientID, req.Message, response, time.Now().Format(time.RFC3339))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Next()
	})

	r.GET("/health", health)
	r.GET("/", root)
	r.POST("/clients", addClient)
	r.GET("/clients", listClients)
	r.POST("/clients/:cid/train", train)
	r.POST("/clients/:cid/chat", chat)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	r.Run(":" + port)
}
