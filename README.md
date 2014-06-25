prostopleer
===========

Go api wrapper for pleer.com

Example:
```go
package main

import (
	"fmt"
	pp "github.com/t0pep0/prostopleer"
)

func main() {
	api := &pp.Api{
		User:     "user",
		Password: "password",
	}
	tracks, count, err := api.SearchTrack("Kruger", "all", 1, 10)
	fmt.Println("Error:", err)
	fmt.Println("Count:", count)
	for _, track := range tracks {
		track.GetLink("listen")
		track.GetLyrics()
		fmt.Println("Id:", track.Id)
		fmt.Println("Artist:", track.Artist)
		fmt.Println("Name:", track.Name)
		fmt.Println("Duration:", track.Duration)
		fmt.Println("Bitrate:", track.Bitrate)
		fmt.Println("Size:", track.Size)
		fmt.Println("Listen URL:", track.ListenUrl)
		fmt.Println("Lyrics:", track.Lyrics)
		err := track.Download()
		fmt.Println("Error on download:", err)
		fmt.Println("Size of track:", len(track.MP3))
		fmt.Println("===========================")
	}
	suggs, err := api.Autocomplete("Натиск")
	fmt.Println("Error:", err)
	for _, sugg := range suggs {
		fmt.Println("     ", sugg)
	}
}
```
