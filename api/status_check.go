package api

import (
	"fc-deal-making-service/core"
	"strings"

	"github.com/labstack/echo/v4"
)

type StatusCheckResponse struct {
	Content struct {
		ID      int64  `json:"id"`
		Name    string `json:"name"`
		Cid     string `json:"cid,omitempty"`
		Status  string `json:"status"`
		Message string `json:"message,omitempty"`
	} `json:"content"`
}

func ConfigureStatusCheckRouter(e *echo.Group, node *core.LightNode) {
	e.GET("/status/:id", func(c echo.Context) error {

		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		var content core.Content
		node.DB.Raw("select c.id, c.estuary_content_id, c.cid, c.status from contents as c where c.id = ? and c.requesting_api_key = ?", c.Param("id"), authParts[1]).Scan(&content)

		return c.JSON(200, StatusCheckResponse{
			Content: struct {
				ID      int64  `json:"id"`
				Name    string `json:"name"`
				Cid     string `json:"cid,omitempty"`
				Status  string `json:"status"`
				Message string `json:"message,omitempty"`
			}{ID: content.ID, Name: content.Name, Status: content.Status, Cid: content.Cid},
		})
	})

	e.GET("/list-all-cids", func(c echo.Context) error {

		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")

		var content []core.Content
		node.DB.Raw("select c.name, c.id, c.estuary_content_id, c.cid, c.status,c.created_at,c.updated_at from contents as c where requesting_api_key = ?", authParts[1]).Scan(&content)

		return c.JSON(200, content)

	})
}
