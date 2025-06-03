package api

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

// Проксирование запроса к микросервису
func proxyTo(c *gin.Context, target string) (int, http.Header, []byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(c.Request.Method, target+c.Request.RequestURI, c.Request.Body)
	req.Header = c.Request.Header

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
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
