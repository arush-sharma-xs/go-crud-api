package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xs-arush-0856/lms/models"
)

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {

        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Headers", "*")
        /*
            c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
            c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
            c.Writer.Header().Set("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
            c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH")
        */
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}

func main() {
	models.ConnectDB()
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.LoadHTMLGlob("pages/*")

		// endPoint
	r.POST("/library", handleCreateLibrary) // new library creation

	r.POST("/auth", auth)
	r.POST("/users", addUser)               // Creating new User

	// admin routes
	r.GET("/book", getBook)
	r.POST("/book", addBook)                // Add books
	r.DELETE("/book/:isbn", removeBook)     // Delete books
	r.PATCH("/book/:isbn", updateBook)      // Update Book
	r.POST("/raiseissue", raiseIssue)
	r.POST("/approveBookRequest/:reqId", handleApproveRequest)

	// reader routes
	r.GET("/issuerequest", getIssueRequest) // Get Issue List
	r.GET("/search/:key/:value", searchBook)

	r.Run(":5101")
}
