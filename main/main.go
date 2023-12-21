package main

import "fmt"

var loggedIn bool = false

func main() {

	// if !loggedIn {
	// 	LoginProcess()
	// }

	profileInfo := []byte{}

	x, y := FetchWebAPI("GET", "https://api.spotify.com/v1/me", nil, profileInfo, "token")

	fmt.Print(x, y)

}
