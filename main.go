package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

const (
	getKey = `SELECT value FROM keys WHERE key=$1 LIMIT 1;`
	putKey = `INSERT INTO keys(key,value) values($1,$2) RETURNING value;`
)

func main() {

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", getEnv("PG_HOST", "postgres"), getEnv("PG_USER", "postgres"), getEnv("PG_PASSWORD", "postgres"), getEnv("PG_DB", "postgres"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	p := ginprometheus.NewPrometheus("gin")
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.Path
		for _, p := range c.Params {
			if p.Key == "key" {
				url = strings.Replace(url, p.Value, ":key", 1)
				break
			}
		}
		return url
	}
	p.Use(r)

	r.GET("/:key", func(c *gin.Context) {
		var value string

		key := c.Param("key")

		row := db.QueryRow(getKey, key)

		switch err := row.Scan(&value); err {
		case sql.ErrNoRows:
			log.Printf("nil key requested: %s\n", key)
			c.AbortWithStatus(204)
		case nil:
			c.String(http.StatusOK, "%s\n", value)
		default:
			log.Println(err)
		}
	})

	r.POST("/:key", func(c *gin.Context) {
		var value string

		key := c.Param("key")

		bytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err)
		}

		err = db.QueryRow(putKey, key, string(bytes)).Scan(&value)
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				switch err.Code.Name() {
				case "unique_violation":
					c.AbortWithStatus(409)
				default:
					log.Println(err.Code.Name())
				}
			} else {
				log.Println(err)
				c.AbortWithStatus(500)
			}
		} else {
			c.String(http.StatusOK, "%s\n", value)
		}
	})

	r.Run()
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
