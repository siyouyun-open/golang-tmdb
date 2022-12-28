package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	tmdb "github.com/siyouyun-open/golang-tmdb"
	"os"
)

func main() {

	tmdbClient, err := tmdb.Init(os.Getenv("APIKey"), resty.New())

	if err != nil {
		fmt.Println(err)
	}

	movie, err := tmdbClient.GetMovieDetails(299536, nil)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(movie.Title)

	fmt.Println("------")

	// With options
	options := make(map[string]string)
	options["append_to_response"] = "credits"
	options["language"] = "pt-BR"

	movie, err = tmdbClient.GetMovieDetails(299536, options)

	if err != nil {
		fmt.Println(err)
	}

	// Credits - Iterate cast from append to response.
	for _, v := range movie.MovieCreditsAppend.Credits.Cast {
		fmt.Println(v.Name)
	}
}
