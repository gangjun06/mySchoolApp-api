package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/osang-school/backend/graph"
	"github.com/osang-school/backend/graph/errors"
	"github.com/osang-school/backend/graph/generated"
	"github.com/osang-school/backend/internal/conf"
	"github.com/osang-school/backend/internal/db/mongodb"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func main() {
	rand.Seed(time.Now().Unix())
	conf.Init()
	port := conf.Get().Server.Port

	mongodb.Init()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {

		err := graphql.DefaultErrorPresenter(ctx, e)
		if v, ok := e.(*errors.Error); ok {
			err.Message = v.Message
			err.Extensions["code"] = v.Code
			err.Extensions["description"] = v.Description()
		}

		return err
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%d/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
