package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/vmihailenco/msgpack"

	"github.com/yescorihuela/hex-microservice/shortener"
)

func httpPort() string {
	port := "8081"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func main() {
	address := fmt.Sprintf("http://localhost%s", httpPort())
	redirect := shortener.Redirect{}
	redirect.URL = "http://github.com/tensor-programming?tab=repositories"

	body, err := msgpack.Marshal(&redirect)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(address, "application/x-msgpack", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Request.Body.Close()
	body, err = ioutil.ReadAll(resp.Request.Body)
	if err != nil {
		log.Fatalln(err)
	}
	msgpack.Unmarshal(body, &redirect)
	log.Printf("%v\n", redirect)
}
