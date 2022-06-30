package apihandlers

import (
	imageservice "imageuploader/internal/image"

	"github.com/gin-gonic/gin"
)

type Api struct {
	ImageHandlers *imageservice.Handler
}

func (a *Api) CreateRoutes(r *gin.Engine) {
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	v1 := r.Group("/v1")

	{
		// `GET /v1/images`
		// List metadata for stored images.
		v1.GET("/images", a.ImageHandlers.GetAll)

		// 	`GET /v1/images/<id>`
		// Get metadata for image with id `<id>`.
		v1.GET("/images/:id", a.ImageHandlers.Get)

		// `GET /v1/images/<id>/data`
		// Get image data for image with id `<id>`.
		// Optional GET parameter: `?bbox=<x>,<y>,<w>,<h>` to get a cutout of the image.
		v1.GET("/images/:id/data", a.ImageHandlers.ServeImageData)

		// `POST /v1/images`
		// Upload new image. Request body should be image data.
		v1.POST("/images", a.ImageHandlers.Add)

		// `PUT /v1/images/<id>`
		// Update image. Request body should be image data.
		v1.PUT("/images/:id", a.ImageHandlers.Update)

	}

}
