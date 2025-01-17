package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	docs "main/docs"

	"github.com/gin-gonic/gin"
	"github.com/kkdai/youtube/v2"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			YouTube Video Downloader API
//	@version		1.0
//	@description	API for downloading YouTube videos using kkdai/youtube.
//	@host			localhost:8080
//	@BasePath		/

var ytClient = youtube.Client{}

// CustomDuration wraps time.Duration to enable custom JSON marshalling
type CustomDuration struct {
	time.Duration
}

// MarshalJSON converts the duration to a human-readable string (e.g., "2h45m")
func (d CustomDuration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

type VideoDetails struct {
	Title      string             `json:"title"`
	Author     string             `json:"author"`
	Duration   time.Duration      `json:"duration"`
	Thumbnails youtube.Thumbnails `json:"thumbnails"`
	Formats    youtube.FormatList `json:"formats"`
}

type DownloadRequest struct {
	VideoID   string `json:"video_id" binding:"required"`
	AudioOnly bool   `json:"audio_only"`
	Quality   string `json:"quality"`
}

type DownloadProgress struct {
	VideoID   string
	Progress  int
	Completed bool
	Error     string
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	progressMap = make(map[string]*DownloadProgress)
	progressMux sync.RWMutex
)

// ValidateByVideoId godoc
//
//	@Summary		Get video details
//	@Description	Get details about a YouTube video by video ID
//	@Tags			Video
//	@Param			id	path		string	true	"YouTube Video ID"
//	@Success		200	{object}	main.VideoDetails
//	@Failure		400	{object}	main.ErrorResponse
//	@Router			/video/{id} [get]
func ValidateByVideoId(c *gin.Context) {
	log.Println("GET /video/:id endpoint hit")
	videoID := c.Param("id")
	video, err := ytClient.GetVideo(videoID)
	if err != nil {
		log.Printf("Failed to fetch video details: %v", err)
		errorResponse := ErrorResponse{
			Error: "Invalid video ID",
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	videoDetails := VideoDetails{
		Title:      video.Title,
		Author:     video.Author,
		Duration:   video.Duration,
		Thumbnails: video.Thumbnails,
		Formats:    video.Formats,
	}
	c.JSON(http.StatusOK, videoDetails)
}

// Download godoc
//
//	@Summary		Download a YouTube video
//	@Description	Download a YouTube video by specifying the video ID, quality, and audio-only option
//	@Tags			Download
//	@Accept			json
//	@Produce		json
//	@Param			downloadRequest	body		DownloadRequest	true	"Download request payload"
//	@Success		200				{object}	map[string]string
//	@Failure		400				{object}	map[string]string
//	@Router			/download [post]
func Download(c *gin.Context) {
	log.Println("POST /download endpoint hit")
	var req DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	video, err := ytClient.GetVideo(req.VideoID)
	if err != nil {
		log.Printf("Failed to fetch video: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	progressMux.Lock()
	log.Printf("Locked progress map for update")
	progressMap[req.VideoID] = &DownloadProgress{
		VideoID:  req.VideoID,
		Progress: 0,
	}
	progressMux.Unlock()
	log.Printf("Unlocked progress map after update")

	// Start download asynchronously
	go func(video *youtube.Video, req DownloadRequest) {
		defer func() {
			progressMux.Lock()
			log.Printf("Locked progress map for update")
			progressMap[req.VideoID].Completed = true
			progressMux.Unlock()
			log.Printf("Unlocked progress map after update")
		}()

		// Simulate download and set progress
		for i := 0; i <= 100; i += 20 {
			progressMux.Lock()
			log.Printf("Locked progress map for update")
			progressMap[req.VideoID].Progress = i
			progressMux.Unlock()
			log.Printf("Unlocked progress map after update")
		}
	}(video, req)

	c.JSON(http.StatusOK, gin.H{"message": "Download started"})
}

func main() {
	router := gin.Default()

	// Swagger setup
	docs.SwaggerInfo.BasePath = "/"
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler, url))

	router.GET("/video/:id", ValidateByVideoId)

	router.POST("/download", Download)

	//	@Summary		Get download progress
	//	@Description	Get the download progress of a YouTube video by video ID
	//	@Tags			Progress
	//	@Param			id	path		string	true	"YouTube Video ID"
	//	@Success		200	{object}	DownloadProgress
	//	@Failure		404	{object}	map[string]string
	//	@Router			/progress/{id} [get]
	router.GET("/progress/:id", func(c *gin.Context) {
		log.Println("GET /progress/:id endpoint hit")
		videoID := c.Param("id")

		progressMux.RLock()
		log.Printf("Read-locked progress map")
		progress, exists := progressMap[videoID]
		progressMux.RUnlock()
		log.Printf("Released read-lock on progress map")

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "No progress found for the given video ID"})
			return
		}

		c.JSON(http.StatusOK, progress)
	})

	router.Run(":8080")
}
