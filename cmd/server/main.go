package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/osang-school/backend/graph"
	"github.com/osang-school/backend/graph/generated"
	"github.com/osang-school/backend/internal/conf"
	"github.com/osang-school/backend/internal/db/mongodb"
)

func main() {
	rand.Seed(time.Now().Unix())
	conf.Init()
	port := conf.Get().Server.Port

	mongodb.Init()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%d/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
