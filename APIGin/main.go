package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/supesharisuto/RestfulWebService-4/dbutils"
)

// DB Driver visible to whole program
var DB *sql.DB
var Sess *session.Session

// aws SDK
var Region = "us-east-1"

// VideoResource holds information about locations
type VideoResource struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	AssetID string `json:"assetid"`
}

// GetAllvideo returns the list of dropoff Videos
func GetAllAssetsDropoff(c *gin.Context) {
	var video VideoResource

	id := c.Param("asset_id")
	log.Println(id)
	video.AssetID = c.Param("asset_id")

	//err := listBucketsPrefix("video-sample")
	prefix := "dropoff" + "/" + c.Param("asset_id")
	log.Println(prefix)

	err := listObjects("video-sample", prefix)
	//err := listObjects("video-sample", "dropoff")
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"result": video,
		})
	}
}

// GetAllvideo returns the list of dropoff Videos
func GetAllBucket(c *gin.Context) {
	var video VideoResource

	err := listBuckets()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"result": video,
		})
	}
}

// Getvideo returns the video detail
func GetVideo(c *gin.Context) {
	var video VideoResource
	id := c.Param("video_id")
	err := DB.QueryRow("select ID, NAME, ASSETID from video where id=?", id).Scan(&video.ID, &video.Name, &video.AssetID)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"result": video,
		})
	}
}

// CreateVideo handles the POST
func CreateVideo(c *gin.Context) {
	var video VideoResource
	// Parse the body into our resrource
	if err := c.BindJSON(&video); err == nil {
		// Format Time to Go time format
		statement, _ := DB.Prepare("insert into video (NAME, ASSETID) values (?, ?)")
		result, _ := statement.Exec(video.Name, video.AssetID)
		if err == nil {
			newID, _ := result.LastInsertId()
			video.ID = int(newID)
			c.JSON(http.StatusOK, gin.H{
				"result": video,
			})
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}
	} else {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

// deleteVideo handles the removing of resource
func deleteVideo(c *gin.Context) {
	id := c.Param("video_id")
	statement, _ := DB.Prepare("delete from video where id=?")
	_, err := statement.Exec(id)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		c.String(http.StatusOK, "")
	}
}

func init() {

	Sess, _ = session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
}

func main() {
	var err error
	DB, err = sql.Open("sqlite3", "./videoapi.db")
	if err != nil {
		log.Println("Driver creation failed!")
	}
	dbutils.Initialize(DB)

	// AWS credentials from the shared credentials file ~/.aws/credentials.
	/*
		Sess, err = session.NewSession(&aws.Config{
			Region: aws.String("us-east-1")},
		)
	*/

	r := gin.Default()
	// Add routes to REST verbs
	r.GET("/v1/videos/:video_id", GetVideo)
	r.POST("/v1/videos", CreateVideo)
	r.DELETE("/v1/videos/:video_id", deleteVideo)
	r.GET("/v1/buckets/", GetAllBucket)
	r.GET("/v1/dropoff/:asset_id", GetAllAssetsDropoff)

	r.Run(":8000") // Default listen and serve on 0.0.0.0:8080
}
