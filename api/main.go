package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wesleywcr/dev-book/api/config"
	_ "github.com/wesleywcr/dev-book/api/docs" // Import generated Swagger docs
	"github.com/wesleywcr/dev-book/api/router"
)

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	config.Loading()

	r := router.InitRouter()

	fmt.Printf("Server ON %d\n", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
