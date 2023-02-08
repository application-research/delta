package api

import (
	"fc-deal-making-service/core"
	"fc-deal-making-service/jobs"
	"strings"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/labstack/echo/v4"
)

type CidRequest struct {
	Cids []string `json:"cids"`
}

type UploadResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	ID      int64  `json:"id,omitempty"`
}

func ConfigurePinningRouter(e *echo.Group, node *core.LightNode) {

	content := e.Group("/content")
	content.POST("/add-car", func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")
		file, err := c.FormFile("data")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}

		addNode, err := node.Node.AddPinFile(c.Request().Context(), src, nil)

		// save the file to the database.
		content := core.Content{
			Name:             file.Filename,
			Size:             file.Size,
			Cid:              addNode.Cid().String(),
			RequestingApiKey: authParts[1],
			Status:           "pinned",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		node.DB.Create(&content)

		if err != nil {
			return err
		}

		c.JSON(200, UploadResponse{
			Status:  "success",
			Message: "Car uploaded and pinned successfully to the network.",
			ID:      content.ID,
		})

		d := jobs.NewItemContentProcessor(node, content)
		node.Dispatcher.AddJob(d) // add the job so we can process it later

		return nil
	})

	content.POST("/add", func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")
		file, err := c.FormFile("data")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}

		addNode, err := node.Node.AddPinFile(c.Request().Context(), src, nil)

		// get available staging buckets.
		// save the file to the database.
		content := core.Content{
			Name:             file.Filename,
			Size:             file.Size,
			Cid:              addNode.Cid().String(),
			RequestingApiKey: authParts[1],
			Status:           "pinned",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		node.DB.Create(&content)

		if err != nil {
			c.JSON(500, UploadResponse{
				Status:  "error",
				Message: "Error pinning the file" + err.Error(),
			})
		}

		c.JSON(200, UploadResponse{
			Status:  "success",
			Message: "File uploaded and pinned successfully",
			ID:      content.ID,
		})

		// prepare to dispatch
		d := jobs.NewItemContentProcessor(node, content)
		node.Dispatcher.AddJob(d) // add the job so we can process it later

		return nil
	})

	content.POST("/cid/:cid", func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")
		cidFromForm := c.Param("cid")
		cidNode, err := cid.Decode(cidFromForm)
		if err != nil {
			return err
		}

		//	 get the node
		addNode, err := node.Node.Get(c.Request().Context(), cidNode)

		// get available staging buckets.
		// save the file to the database.
		size, err := addNode.Size()

		content := core.Content{
			Name:             addNode.Cid().String(),
			Size:             int64(size),
			Cid:              addNode.Cid().String(),
			RequestingApiKey: authParts[1],
			Status:           "pinned",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		node.DB.Create(&content)

		if err != nil {
			c.JSON(500, UploadResponse{
				Status:  "error",
				Message: "Error pinning the cid" + err.Error(),
			})
		}

		c.JSON(200, UploadResponse{
			Status:  "success",
			Message: "CID uploaded and pinned successfully",
			ID:      content.ID,
		})

		d := jobs.NewItemContentProcessor(node, content)
		node.Dispatcher.AddJob(d) // add the job so we can process it later

		return nil
	})

	content.POST("/cids", func(c echo.Context) error {
		authorizationString := c.Request().Header.Get("Authorization")
		authParts := strings.Split(authorizationString, " ")
		var cidRequest CidRequest
		c.Bind(&cidRequest)
		for _, cidFromForm := range cidRequest.Cids {
			cidNode, err := cid.Decode(cidFromForm)
			if err != nil {
				return err
			}

			//	 get the node and save on the database
			addNode, err := node.Node.Get(c.Request().Context(), cidNode)

			// get available staging buckets.
			// save the file to the database.
			size, err := addNode.Size()

			content := core.Content{
				Name:             addNode.Cid().String(),
				Size:             int64(size),
				Cid:              addNode.Cid().String(),
				RequestingApiKey: authParts[1],
				Status:           "pinned",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}

			node.DB.Create(&content)
			d := jobs.NewItemContentProcessor(node, content)
			node.Dispatcher.AddJob(d) // add the job so we can process it later
		}
		return nil
	})
}
