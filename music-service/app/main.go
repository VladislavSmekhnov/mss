package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"music-service/app/initializers"
	"music-service/app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("resources/templates/*")
	// Serve static files
	router.Static("/music_service/static", "./resources/static")

	defer initializers.DBClient.Disconnect(context.Background())
	// Get a handle to the "audio_db" database
	db := initializers.DBClient.Database("audio_db")

	// Retrieve data from the "audio_files" collection
	audioFilesCollection := db.Collection("audio_files")
	cursor, err := audioFilesCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var audioFiles []models.AudioFile
	if err := cursor.All(context.Background(), &audioFiles); err != nil {
		log.Fatal(err)
	}
	// Retrieve data from the "selections" collection
	selectionsCollection := db.Collection("selections")
	cursor, err = selectionsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var selections []models.Selection
	if err := cursor.All(context.Background(), &selections); err != nil {
		log.Fatal(err)
	}

	// Route to serve the HTML page
	router.GET("/music_service", func(c *gin.Context) {
		playerHTML, err := LoadHlsPlayer(c)
		if err != nil {
			log.Fatal(err)
			return
		}

		var selectionDtos = make(map[string]models.SelectionDto)

		for _, selection := range selections {
			var selectionDto models.SelectionDto
			selectionDto.ID = selection.ID
			selectionDto.SelectionName = selection.SelectionName

			for _, trackID := range selection.TrackIDs {
				for _, audioFile := range audioFiles {
					if audioFile.ID == trackID {
						selectionDto.Tracks = append(selectionDto.Tracks, audioFile)
						break
					}
				}
			}

			selectionDtos[selectionDto.SelectionName] = selectionDto
		}

		// Pass the content as a variable to the template
		c.HTML(http.StatusOK, "index.html", gin.H{
			"hls_player":    template.HTML(string(playerHTML)),
			"AudioFiles":    audioFiles,
			"Selections":    selections,
			"SelectionDtos": selectionDtos,
		})
	})

	router.GET("/selections", func(c *gin.Context) {
		playerHTML, err := LoadHlsPlayer(c)
		if err != nil {
			return
		}
		// Pass the content as a variable to the template
		c.HTML(http.StatusOK, "collection.html", gin.H{
			"hls_player": template.HTML(string(playerHTML)),
		})
	})

	router.GET("/podcasts_and_books", func(c *gin.Context) {
		playerHTML, err := LoadHlsPlayer(c)
		if err != nil {
			return
		}
		// Pass the content as a variable to the template
		c.HTML(http.StatusOK, "podcasts-and-books.html", gin.H{
			"hls_player": template.HTML(string(playerHTML)),
		})
	})

	//fmt.Println("Audio Files:")
	//for _, audioFile := range audioFiles {
	//	fmt.Printf("%+v\n", audioFile)
	//}

	//fmt.Println("Selections:")
	//for _, selection := range selections {
	//	fmt.Printf("%+v\n", selection)
	//}

	router.Run(":8083")
}

func LoadHlsPlayer(c *gin.Context) ([]byte, error) {
	// Read the content of the player.html file
	playerHTML, err := ioutil.ReadFile("resources/templates/player.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error: Cannot load hls-player (hls.js)")
		return nil, err
	}
	return playerHTML, nil
}
