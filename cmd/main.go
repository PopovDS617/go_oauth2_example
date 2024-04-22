package main

import (
	"fmt"
	"os"
)

func main() {
	v := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")

	fmt.Println(v)
}
