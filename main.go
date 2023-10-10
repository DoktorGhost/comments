// docker build -t two-db .
// docker-compose up -d
package main

import (
	"CommentsService/db"
	"CommentsService/pcg/comments"
	"CommentsService/pcg/types"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	db.InitDB()

	defer db.CloseDB()

	router := gin.Default()

	router.POST("/add-comment", func(c *gin.Context) {
		var request types.Request

		err := c.BindJSON(&request)
		fmt.Println(request)
		uniqueID := request.UniqueID
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d, Error: %v", time.Now().Format("2006-01-02 15:04:05"), uniqueID, c.ClientIP(), http.StatusBadRequest, err)
			return
		}

		commentID, err := comments.AddComment(request.NewsID, request.CommentText, request.ParentCommentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d, Error: %v", time.Now().Format("2006-01-02 15:04:05"), uniqueID, c.ClientIP(), http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"commentID": commentID, "message": "Comment added successfully"})
		log.Printf("Timestamp: %s, Request ID: %s, IP: %s, HTTP Code: %d", time.Now().Format("2006-01-02 15:04:05"), uniqueID, c.ClientIP(), http.StatusOK)
	})

	router.Run(":8082")
}
