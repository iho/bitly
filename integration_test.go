package main_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iho/bitly/conf"
	"github.com/iho/bitly/handlers"
	"github.com/iho/bitly/shortener"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"

	"github.com/iho/bitly/storage"
	"go.uber.org/fx"
)

func registerHooks(lifecycle fx.Lifecycle, r *gin.Engine, c *conf.Config) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				go r.Run(c.HostName)
				return nil
			},
			OnStop: func(context.Context) error {
				return nil
			},
		},
	)
}

func TestShortenerHttp(t *testing.T) {
	app := fx.New(
		storage.Module,
		shortener.Module,
		handlers.Module,
		fx.Provide(func() *conf.Config {
			return &conf.Config{
				HostName: "localhost:3999",
				Host:     "localhost",
				Port:     dbPort,
				User:     "postgres",
				Password: "postgres",
				DbName:   "postgres",
			}
		}),
		fx.Invoke(registerHooks),
	)
	go app.Run()

	defer func() {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		app.Stop(ctx)
	}()

	t.Run("Save", func(t *testing.T) {
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		data := strings.NewReader(`{"url":"https://www.google.com"}`)
		req, err := http.NewRequest("POST", "http://localhost:3999/urls", data)
		if err != nil {
			log.Fatal(err)
		}
		resp := &http.Response{}
		req.Header.Set("Content-Type", "application/json")
		for i := 0; i < 100; i++ {
			time.Sleep(1 * time.Second)
			resp, err = client.Do(req)
			defer resp.Body.Close()
			if err == nil {
				break
			}
		}
		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		assert.NotEmpty(t, bodyText)

		response := handlers.UrlCreationResponse{}
		json.Unmarshal(bodyText, &response)
		assert.NotEmpty(t, response.Code)

		req, err = http.NewRequest("GET", fmt.Sprintf("http://localhost:3999/urls/%s/", response.Code), nil)
		if err != nil {
			log.Fatalln(err)
		}
		resp, err = client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		location := resp.Header.Get("Location")
		assert.Equal(t, "https://www.google.com", location)
		assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode)

	})
}

var db *sql.DB
var dbPort string

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.Run("postgres", "15.3", []string{"POSTGRES_USER=postgres", "POSTGRES_PASSWORD=postgres"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		dbPort = resource.GetPort("5432/tcp")
		db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			"localhost", dbPort, "postgres", "postgres", "postgres"))
		if err != nil {
			return err
		}
		db.Exec("create table urls ( id text primary key, url text not null);")
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
