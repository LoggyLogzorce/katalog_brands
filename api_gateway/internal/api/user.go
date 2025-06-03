package api

import "github.com/gin-gonic/gin"

func ProfileHandler(c *gin.Context) {
	userID := c.GetString("userID")
	role := c.GetString("role")

	c.Request.Header.Set("X-User-ID", userID)
	c.Request.Header.Set("X-Role", role)

	proxyTo(c, "http://localhost:8082")
}

func UpdateRoleHandler(c *gin.Context) {
	userID := c.GetString("userID")
	role := c.GetString("role")

	c.Request.Header.Set("X-User-ID", userID)
	c.Request.Header.Set("X-Role", role)

	proxyTo(c, "http://localhost:8082")
}
