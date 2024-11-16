package main

import (
	"crypto/sha256"
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

	"github.com/rssnyder/keys/db"
)

type Keys struct {
	Database *db.Database
}

func main() {

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", getEnv("PG_HOST", "postgres"), getEnv("PG_USER", "postgres"), getEnv("PG_PASSWORD", "postgres"), getEnv("PG_DB", "postgres"))
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	// defer dbConn.Close()

	keys := &Keys{
		Database: &db.Database{dbConn},
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

	r = keys.getValue(r)
	r = keys.SetValue(r)
	r = keys.SetKey(r)

	r.Run()
}

func (k *Keys) getValue(r *gin.Engine) *gin.Engine {
	r.GET("/:key", func(c *gin.Context) {
		key := c.Param("key")

		value, err := k.Database.GetValue(key)

		switch err {
		case sql.ErrNoRows:
			c.AbortWithStatus(http.StatusNoContent)
		case nil:
			c.String(http.StatusOK, "%s", value)
		default:
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	})

	return r
}

func (k *Keys) SetValue(r *gin.Engine) *gin.Engine {
	r.POST("/", func(c *gin.Context) {
		bytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		value := string(bytes)

		key := generateKey(value)

		_, err = k.Database.SetKey(key, value)
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				switch err.Code.Name() {
				case "unique_violation":
					c.AbortWithStatus(http.StatusConflict)
				default:
					log.Println(err.Code.Name())
				}
			} else {
				log.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		} else {
			c.String(http.StatusOK, "%s", key)
		}
	})

	return r
}

func (k *Keys) SetKey(r *gin.Engine) *gin.Engine {
	r.POST("/:key", func(c *gin.Context) {
		key := c.Param("key")

		bytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err)
		}
		value := string(bytes)

		_, err = k.Database.SetKey(key, value)
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				switch err.Code.Name() {
				case "unique_violation":
					c.AbortWithStatus(http.StatusConflict)
				default:
					log.Println(err.Code.Name())
				}
			} else {
				log.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		} else {
			c.String(http.StatusOK, "%s", value)
		}
	})

	return r
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func generateKey(value string) (hashString string) {
	h := sha256.New()
	h.Write([]byte(value))
	hashBytes := h.Sum(nil)
	hashString = fmt.Sprintf("%x", hashBytes)

	return
}
