package api

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yescorihuela/hex-microservice/shortener"

	js "github.com/yescorihuela/hex-microservice/serializer/json"
	ms "github.com/yescorihuela/hex-microservice/serializer/msgpack"
)

type RedirectHandler interface {
	Get(c *gin.Context)
	Post(c *gin.Context)
}

type handler struct {
	redirectService shortener.RedirectService
}

func NewHandler(redirectService shortener.RedirectService) RedirectHandler {
	return &handler{redirectService: redirectService}
}

func setupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) serializer(contentType string) shortener.RedirectSerializer {
	if contentType == "application/x-msgpack" {
		return &ms.Redirect{}
	}
	return &js.Redirect{}
}

func (h *handler) Get(c *gin.Context) {
	code := c.Param("url_code")

	redirect, err := h.redirectService.Find(code)
	if err != nil {
		if errors.Cause(err) == shortener.ErrRedirectNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Redirect(http.StatusMovedPermanently, redirect.URL)
}

func (h *handler) Post(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")
	requestBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	redirect, err := h.serializer(contentType).Decode(requestBody)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = h.redirectService.Store(redirect)
	if err != nil {
		if errors.Cause(err) == shortener.ErrRedeirectInvalid {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	responseBody, err := h.serializer(contentType).Encode(redirect)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	setupResponse(c.Writer, contentType, responseBody, http.StatusCreated)
}
