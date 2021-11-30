package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yescorihuela/hex-microservice/shortener"

	h "github.com/yescorihuela/hex-microservice/api"
	mr "github.com/yescorihuela/hex-microservice/repository/mongo"
	rr "github.com/yescorihuela/hex-microservice/repository/redis"
)

func main() {
	router := gin.Default()
	repo := chooseRepo()
	service := shortener.NewRedirectService(repo)
	handler := h.NewHandler(service)

	router.GET("/:url_code", handler.Get)
	router.POST("/", handler.Post)
	router.Run(":8081")
}

func httpPort() string {
	port := "8081"
	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = envPort
	}
	return fmt.Sprintf(":%s", port)
}

func chooseRepo() shortener.RedirectRepository {
	switch os.Getenv("URL_DB") {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := rr.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongodb := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mr.NewMongoRepository(mongoURL, mongodb, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil
}
