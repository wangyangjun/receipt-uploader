package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/wangyangjun/receipt-uploader/graph"
	"github.com/wangyangjun/receipt-uploader/graph/generated"
	"github.com/wangyangjun/receipt-uploader/graph/service/auth"
)

const defaultPort = "8080"
const defaultJwtKey = "this-is-not-safe"
const defaultPwdKey = "this-is-not-safe"

func preenv() {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Println("No environment file found will use default values instead!")
	} else {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}
	// set default env variables if not present
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", defaultPort)
	}
	if os.Getenv("JWT_KEY") == "" {
		os.Setenv("JWT_KEY", defaultJwtKey)
	}
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", defaultPwdKey)
	}
}

func main() {
	preenv()
	port := os.Getenv("PORT")

	router := chi.NewRouter()

	router.Use(auth.Middleware())

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	fs := http.FileServer(http.Dir("./images"))
	router.Handle("/image/", http.StripPrefix("/image/", fs))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
