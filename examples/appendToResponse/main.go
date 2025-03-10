package main

import (
	"fmt"
	"os"

	tmdb "github.com/siyouyun-open/golang-tmdb"
)

func main() {
	fmt.Println(os.Getenv("APIKey"))
	tmdbClient, err := tmdb.Init(os.Getenv("APIKey"), nil)

	if err != nil {
		fmt.Println(err)
	}

	// With options
	options := make(map[string]string)
	options["append_to_response"] = "watch/providers"

	movie, err := tmdbClient.GetMovieDetails(299536, options)

	if err != nil {
		fmt.Println(err)
	}

	for _, v := range movie.MovieWatchProvidersAppend.WatchProviders.Results {
		fmt.Println(v.Flatrate)
	}
}
