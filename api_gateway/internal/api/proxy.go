package api

import (
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// Проксирование запроса к микросервису
func proxyTo(c *gin.Context, target string, requestURI string, body io.Reader) (int, http.Header, []byte, error) {
	if body == nil {
		body = c.Request.Body
	}

	if requestURI == "" {
		requestURI = c.Request.RequestURI
	}
	client := &http.Client{}
	req, _ := http.NewRequest(c.Request.Method, target+requestURI, body)
	req.Header = c.Request.Header

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		if strings.EqualFold(key, "Content-Length") {
			continue
		}
		for _, v := range values {
			c.Writer.Header().Add(key, v)
		}
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, resp.Header, nil, err
	}

	return resp.StatusCode, resp.Header, respBytes, nil
}
