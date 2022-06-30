package imageservice

import (
	"bufio"
	"fmt"
	"image"
	"imageuploader/config"
	imagerepository "imageuploader/internal/image/repository"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandlers(service Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	images, err := h.service.GetAll()
	if err != nil {
		c.String(http.StatusInternalServerError, "get all images error: %s", err.Error())
	}
	c.JSON(http.StatusOK, images)
}

func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	image, err := h.service.Get(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "get image error: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, image)
}

func (h *Handler) ServeImageData(c *gin.Context) {
	id := c.Param("id")
	image, err := h.service.Get(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "get image error: %v", err)
		return
	}
	if query, queryExists := c.GetQuery("bbox"); queryExists {
		serveCroppedImage(c, image, query)
	} else {
		serveImage(c, image)
	}
}

func serveCroppedImage(c *gin.Context, imageObj Image, query string) {
	file, err := os.Open(imageObj.Path)
	if err != nil {
		c.String(http.StatusNotFound, "requested resource does not exist: %v", err)
		return
	}
	defer file.Close()

	dimensions := strings.Split(query, ",")
	if len(dimensions) != 4 {
		c.String(http.StatusBadRequest, "bbox dimensions are not in correct format (bbox=<x>,<y>,<w>,<h>)")
		return
	}
	x, err := strconv.Atoi(dimensions[0])
	if err != nil {
		c.String(http.StatusBadRequest, "bbox dimensions are not in correct format (bbox=<x>,<y>,<w>,<h>)")
		return
	}
	y, err := strconv.Atoi(dimensions[1])
	if err != nil {
		c.String(http.StatusBadRequest, "bbox dimensions are not in correct format (bbox=<x>,<y>,<w>,<h>)")
		return
	}
	w, err := strconv.Atoi(dimensions[2])
	if err != nil {
		c.String(http.StatusBadRequest, "bbox dimensions are not in correct format (bbox=<x>,<y>,<w>,<h>)")
		return
	}
	h, err := strconv.Atoi(dimensions[3])
	if err != nil {
		c.String(http.StatusBadRequest, "bbox dimensions are not in correct format (bbox=<x>,<y>,<w>,<h>)")
		return
	}
	img, _, err := image.Decode(file)
	if err != nil {
		c.String(http.StatusBadRequest, "error decoding file to an image: %v", err)
		return
	}
	croppedImage, err := cropImage(img, image.Rectangle{
		Min: image.Point{
			X: x,
			Y: y,
		},
		Max: image.Point{
			X: x + w,
			Y: y + h,
		},
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "error cropping the image: %v", err)
		return
	}
	ext, err := imagerepository.ExtensionFromString(imageObj.Extension)
	if err != nil {
		c.String(http.StatusInternalServerError, "error encoding the image: %v", err)
		return

	}
	switch ext {
	case imagerepository.PNG:
		err = png.Encode(c.Writer, croppedImage)
		if err != nil {
			c.String(http.StatusInternalServerError, "error encoding the image: %v", err)
			return
		}
	case imagerepository.JPG:
		err = jpeg.Encode(c.Writer, croppedImage, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, "error encoding the image: %v", err)
			return
		}
	default:
		c.String(http.StatusBadRequest, "unsupported format: %v", err)
		return
	}

}

func serveImage(c *gin.Context, imageObj Image) {
	file, err := os.Open(imageObj.Path)
	if err != nil {
		c.String(http.StatusNotFound, "requested resource does not exist: %v", err)
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileModTime := fileInfo.ModTime().String()
	fileSize := fileInfo.Size()
	eTag := fileModTime + strconv.Itoa(int(fileSize))

	c.Writer.Header().Add("Content-Type", imageObj.ContentType)
	c.Writer.Header().Add("ETag", eTag)

	if match := c.Request.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, eTag) {
			c.Writer.WriteHeader(http.StatusNotModified)
			return
		}
	}
	reader := bufio.NewReader(file)
	reader.WriteTo(c.Writer)
	return
}

// Even though the assignment statement does not mention anyng about single/ multiple upload I am just using multiple upload
// For how single upload should work update can be seen since that is uploading single image
func (h *Handler) Add(c *gin.Context) {
	name := c.PostForm("name")

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}
	uploadedImages := form.File["files"]

	for _, uploadedImage := range uploadedImages {
		multipartFile, err := uploadedImage.Open()
		if err != nil {
			log.Printf("Error opening multipart file: %v\n", err)
			continue
		}
		defer multipartFile.Close()

		mType, err := mimetype.DetectReader(multipartFile)

		if err != nil || !isValidImage(mType) {
			log.Printf("mType: %s, err: %v", mType.Extension(), err)
			continue
		}

		// In ideal scenario its better to store images in cloud
		savedImageFile, err := os.CreateTemp(config.LocalImagesDirectory, "image.*"+mType.Extension())
		if err != nil {
			log.Printf("error creating image file %v\n", err)
			continue
		}
		defer savedImageFile.Close()

		// Required to do so if you read file earlier
		multipartFile.Seek(0, 0)

		err = h.service.Add(multipartFile, savedImageFile, mType)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("upload image err: %v", err))
			return
		}

	}

	c.String(http.StatusOK, "Uploaded successfully %d files with fields name=%s.", len(uploadedImages), name)

	return
}

func (h *Handler) Update(c *gin.Context) {
	// Ideal way could be storing the old image somewhere so user can easily revert to the image he/she is interested in
	// For now I am just deleting the old image saved to work on saving space
	id := c.Param("id")
	image, err := h.service.Get(id)
	if err != nil {
		c.String(http.StatusNotFound, "the resource you are trying to update does not exist: %v", err)
		return
	}
	err = os.Remove(image.Path)
	if err != nil {
		log.Printf("Error deleting the old image file: %v\n", err)
	}
	uploadedImage, err := c.FormFile("file")

	multipartFile, err := uploadedImage.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, "error opening multipart file: %v\n", err)
		return
	}
	defer multipartFile.Close()

	mType, err := mimetype.DetectReader(multipartFile)
	if err != nil || !isValidImage(mType) {
		c.String(http.StatusInternalServerError, "unsupported/invalid file format: %v\n", err)
		return
	}

	// In ideal scenario its better to store images in cloud
	savedImageFile, err := os.CreateTemp(config.LocalImagesDirectory, "image.*"+mType.Extension())
	if err != nil {
		c.String(http.StatusInternalServerError, "error creating image file: %v\n", err)
		return
	}
	defer savedImageFile.Close()

	// Required to do so if you read file earlier
	multipartFile.Seek(0, 0)

	err = h.service.Update(id, multipartFile, savedImageFile, mType)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("update image err: %v", err))
		return
	}

	c.String(http.StatusOK, "updated successfully")
}

func isValidImage(mimeType *mimetype.MIME) bool {
	for _, allowedType := range allowedMimeTypes {
		// Using this mimetype library since it is supporting so many formats.
		if mimeType.Is(allowedType) {
			return true
		}
	}
	return false
}

type subImager interface {
	SubImage(r image.Rectangle) image.Image
}

func cropImage(img image.Image, boundary image.Rectangle) (image.Image, error) {
	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}

	return simg.SubImage(boundary), nil
}

var allowedMimeTypes []string = []string{"image/png", "image/jpeg", "image/tiff", "image/heic", "image/svg+xml"}
