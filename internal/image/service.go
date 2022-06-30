package imageservice

import (
	"errors"
	"image"
	"imageuploader/config"
	imagerepository "imageuploader/internal/image/repository"
	"imageuploader/pkg/utils/files"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
)

type Image struct {
	Filesize    uint64
	Height      uint16
	Width       uint16
	Extension   string
	ContentType string
	CreatedAt   string
	Title       string
	Path        string
	ID          string
}

type Service interface {
	GetAll() ([]Image, error)
	Get(id string) (Image, error)
	GetPath(id string) (string, error)
	Add(image multipart.File, file *os.File, contentType *mimetype.MIME) error
	Update(id string, image multipart.File, file *os.File, contentType *mimetype.MIME) error
}

type service struct {
	db imagerepository.Repository
}

func NewService(repo imagerepository.Repository) *service {
	return &service{
		db: repo,
	}
}

func (s *service) GetAll() ([]Image, error) {
	var response []Image
	images, err := s.db.GetAll()
	if err != nil {
		return nil, err
	}
	for _, image := range images {
		response = append(response, translateRepoToServiceModel(&image))
	}
	return response, nil
}

func (s *service) Get(id string) (Image, error) {
	image, err := s.db.Get(id)
	if err != nil {
		return Image{}, err
	}
	return translateRepoToServiceModel(image), nil
}

func (s *service) GetPath(id string) (string, error) {
	return s.db.GetPath(id)
}

func (s *service) Add(multipartFile multipart.File, savedImageFile *os.File, contentType *mimetype.MIME) error {

	if err := files.SaveUploadedFile(multipartFile, savedImageFile); err == nil {
		fileStats, err := savedImageFile.Stat()
		log.Printf("File name: %s\n File size: %d\n File path: %s", fileStats.Name(), uint64(fileStats.Size()), config.LocalImagesDirectory+fileStats.Name())

		savedImageFile.Seek(0, 0)

		var height, width uint16
		imageConfig, _, err := image.DecodeConfig(savedImageFile)
		if err != nil {
			// this is to support other image formats as well.. In case if it is not image/png pr image/jpg
			if !errors.Is(err, image.ErrFormat) {
				log.Printf("error decoding image file %v\n", err)
				return err
			}
		} else {
			height = uint16(imageConfig.Height)
			width = uint16(imageConfig.Width)
		}

		log.Printf("Image height: %d\n Image width: %d\n", height, width)

		image := Image{
			Filesize:    uint64(fileStats.Size()),
			Height:      height,
			Width:       width,
			Extension:   strings.Trim(contentType.Extension(), "."),
			Title:       fileStats.Name(),
			Path:        config.LocalImagesDirectory + fileStats.Name(),
			ID:          uuid.NewString(),
			ContentType: contentType.String(),
		}
		repoImage := translateServiceToRepoModel(image)
		return s.db.Add(repoImage)
	} else {
		log.Printf("error saving image to filesystem %v\n", err)
	}
	return nil
}

func (s *service) Update(id string, multipartFile multipart.File, savedImageFile *os.File, contentType *mimetype.MIME) error {
	if err := files.SaveUploadedFile(multipartFile, savedImageFile); err == nil {
		fileStats, err := savedImageFile.Stat()
		log.Printf("File name: %s\n File size: %d\n File path: %s", fileStats.Name(), uint64(fileStats.Size()), config.LocalImagesDirectory+fileStats.Name())

		savedImageFile.Seek(0, 0)

		var height, width uint16
		imageConfig, _, err := image.DecodeConfig(savedImageFile)
		if err != nil {
			// this is to support other image formats as well.. In case if it is not image/png pr image/jpg
			if !errors.Is(err, image.ErrFormat) {
				log.Printf("error decoding image file %v\n", err)
				return err
			}
		} else {
			height = uint16(imageConfig.Height)
			width = uint16(imageConfig.Width)
		}

		log.Printf("Image height: %d\n Image width: %d\n", height, width)

		image := Image{
			Filesize:    uint64(fileStats.Size()),
			Height:      height,
			Width:       width,
			Extension:   strings.Trim(contentType.Extension(), "."),
			Path:        config.LocalImagesDirectory + fileStats.Name(),
			ID:          id,
			ContentType: contentType.String(),
		}
		repoImage := translateServiceToRepoModel(image)
		return s.db.Update(id, config.LocalImagesDirectory+fileStats.Name(), repoImage.Meta)
	} else {
		log.Printf("error saving image to filesystem %v\n", err)
	}
	return nil
}
