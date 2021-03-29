package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/osang-school/backend/graph"
	"github.com/osang-school/backend/graph/generated"
	"github.com/osang-school/backend/graph/myerr"
	gqltools "github.com/osang-school/backend/graph/tools"
	"github.com/osang-school/backend/internal/conf"
	"github.com/osang-school/backend/internal/db/mongodb"
	"github.com/osang-school/backend/internal/db/redis"
	"github.com/osang-school/backend/internal/discord"
	"github.com/osang-school/backend/internal/session"
	"github.com/osang-school/backend/internal/user"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func main() {
	rand.Seed(time.Now().Unix())
	conf.Init()
	port := conf.Get().Server.Port

	mongodb.Init()
	redis.Init()

	c := generated.Config{Resolvers: &graph.Resolver{}}
	c.Directives.Auth = func(ctx context.Context, obj interface{}, next graphql.Resolver, getInfo *bool, reqPermission []string) (interface{}, error) {
		token := ctx.Value("authHeader")
		if token == nil {
			return nil, myerr.New(myerr.ErrAuth, "")
		}
		data, err := session.ParseToken(token.(string))
		if err != nil {
			return nil, err
		}

		hasItem := func(list []string, target string) bool {
			for _, v := range list {
				if v == target || v == "admin" {
					return true
				}
			}
			return false
		}

		for _, v := range reqPermission {
			if ok := hasItem(data.Permission, v); !ok {
				return nil, myerr.New(myerr.ErrPermission, "")
			}
		}

		if getInfo != nil && *getInfo {
			user, err := user.GetUserByID(data.ID)
			if err != nil {
				return nil, err
			}
			ctx = context.WithValue(ctx, "user", user)
		}

		ctx = context.WithValue(ctx, "data", data)
		return next(ctx)
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(c))
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)
		if parsed, ok := myerr.Parse(fmt.Errorf(err.Message)); ok {
			v := parsed.(*myerr.Error)
			err.Message = v.Message
			err.Extensions = map[string]interface{}{
				"code":        v.Code,
				"description": v.Description(),
			}
		}
		return err
	})

	if conf.Discord() != nil {
		discord.Init()
		defer discord.Close()
	}

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", gqltools.Middleware(srv))

	log.Printf("connect to http://localhost:%d/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
