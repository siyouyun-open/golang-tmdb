package main

import (
	"fmt"
	"os"

	tmdb "github.com/siyouyun-open/golang-tmdb"
)

func main() {
	tmdbClient, err := tmdb.Init(os.Getenv("APIKey"), nil)
	if err != nil {
		fmt.Println(err)
	}

	options := map[string]string{
		"append_to_response": "videos",
	}

	movie, err := tmdbClient.GetMovieDetails(299536, options)
	if err != nil {
		fmt.Println(err)
	}

	// Generating Image URLs
	fmt.Println(tmdb.GetImageURL(movie.BackdropPath, tmdb.W500))
	fmt.Println(tmdb.GetImageURL(movie.PosterPath, tmdb.Original))

	// Generating Video URLs
	for _, video := range movie.MovieVideosAppend.Videos.MovieVideos.Results {
		if video.Key != "" {
			fmt.Println(tmdb.GetVideoURL(video.Key))
		}
	}

}
