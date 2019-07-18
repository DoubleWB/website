package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/DoubleWB/website/signatures"
)

//TODO: replace with an actual database
//File in which signatures are stored
//Only used to ensure signature persistence between shutdowns
const STORAGE = "config/signatures"

//Signatures currently loaded from the file
var currentSignatures []signatures.Signature

//Used to keep the current signatures synced with the file
var signaturesFetched = false

//Struct to unmarshal requests into
type req struct {
	//Expected field name required in all signature requests
	Name string `json:"name"`
}

//Return all signatures that are currently loaded, or make sure to load them if it is the first time this function is called
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

//Create a new signature with the name from the request and the current time,
//and save it to our file
func createSignature(c *gin.Context) {
	//Ensure the request is valid
	var requestDecoded req
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&requestDecoded); err != nil {
		fmt.Printf(err.Error() + "\n")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	//Ensure the signature is valid
	if !validSignature(requestDecoded.Name) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	s := signatures.NewSignature(requestDecoded.Name)
	currentSignatures = append(currentSignatures, s)

	//Ensure the signature is saved
	if err := signatures.SaveOverFile(currentSignatures, STORAGE); err == nil {
		c.JSON(http.StatusOK, &s)
	} else {
		fmt.Printf(err.Error() + "\n")
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

//Return whether or not a signature is valid according to the underlying criteria
func validSignature(sign string) bool {
	//No Repeats allowed
	for _, signature := range currentSignatures {
		if sign == signature.Name {
			return false
		}
	}

	//No Blanks allowed
	return sign != ""
}

//Middleware to ensure that all requests have CORS enabled
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Currently allowing all domains
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

//Run our server
func main() {
	r := gin.Default()

	//Enable our middleware
	r.Use(CORSMiddleware())

	//connect routing to handlers, along with a default debugging message at /api/
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

	r.Run()
}
