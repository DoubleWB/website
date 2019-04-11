package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/DoubleWB/website/hci/room"
	"github.com/DoubleWB/website/signatures"
)

const STORAGE = "config/signatures"

var currentSignatures []signatures.Signature
var signaturesFetched = false

type req struct {
	Name string `json:"name"`
}

func fetchAllSignatures(c *gin.Context) {
	if !signaturesFetched {
		var err error
		currentSignatures, err = signatures.ReadFromFile(STORAGE)
		if err != nil {
			fmt.Printf(err.Error() + "\n")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		signaturesFetched = true
	}
	c.JSON(http.StatusOK, &currentSignatures)
}

func createSignature(c *gin.Context) {
	var requestDecoded req
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&requestDecoded); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !validSignature(requestDecoded.Name) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	s := signatures.NewSignature(requestDecoded.Name)
	currentSignatures = append(currentSignatures, s)

	if err := signatures.SaveOverFile(currentSignatures, STORAGE); err == nil {
		c.JSON(http.StatusOK, &s)
	} else {
		fmt.Printf(err.Error() + "\n")
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func validSignature(sign string) bool {
	//temporary measure to help prevent vulgarity
	bannedWords := []string{"fuck", "shit", "ass", "cunt", "fag", "pussy", "bitch"}
	for _, bannedWord := range bannedWords {
		if strings.Contains(strings.ToLower(sign), bannedWord) {
			return false
		}
	}

	for _, signature := range currentSignatures {
		if sign == signature.Name {
			return false
		}
	}

	return sign != ""
}

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
	r := gin.Default()

	r.Use(static.Serve("/", static.LocalFile("./views", true)), CORSMiddleware())

	api := r.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		api.GET("/signs", fetchAllSignatures)
		api.POST("/sign", createSignature)
	}

	hci := r.Group("/hci")
	{
		hci.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		hci.POST("/rooms", room.CreateRoom)
		hci.GET("/rooms/:code", room.GetRoom)
		hci.POST("/finished", room.ChangeFinished)
		hci.DELETE("/rooms", room.DeleteRooms)
		hci.POST("/join_room", room.JoinRoom)
		hci.POST("/items", room.CreateItem)
		hci.PUT("/items", room.EditItem)
		hci.POST("/amounts", room.AddParticipation)
		hci.PUT("/amounts", room.EditParticipation)
		hci.DELETE("/amounts", room.RemoveParticipation)
		hci.GET("/bill/user/:user/room/:code", room.GetBill)
	}

	r.GET("/js/script.js", func(c *gin.Context) {
		c.File("views/js/script.js")
	})

	log.Fatal(autotls.Run(r, "doublewb.xyz", "www.doublewb.xyz")) // listen and serve on 0.0.0.0:443
}
