package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wesleywcr/dev-book/api/router"
)

func main() {
	fmt.Println("Server ON 5000")
	r := router.InitRouter()

	log.Fatal(http.ListenAndServe(":5000", r))
}
