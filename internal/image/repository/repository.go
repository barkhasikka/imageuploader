package imagerepository

import (
	"fmt"
	"strings"
)

type Image struct {
	ID      string
	Meta    *Meta
	Details *Details
}

type Details struct {
	CreatedAt     string
	LastUpdatedAt string
	Title         string
	Path          string
}

type Meta struct {
	Filesize    uint64
	Height      uint16
	Width       uint16
	Extension   Extension
	ContentType string
}

type Extension int

const (
	Invalid Extension = iota
	JPG
	PNG
	GIF
	HEIC
)

var (
	imageTypeMap = map[Extension]string{
		Invalid: "invalid",
		JPG:     "jpg",
		PNG:     "png",
		GIF:     "gif",
		HEIC:    "heic",
	}
)

func (ex Extension) String() string {
	return imageTypeMap[ex]
}

func ExtensionFromString(extensionName string) (Extension, error) {
	for k, v := range imageTypeMap {
		if v == strings.ToLower(extensionName) {
			return k, nil
		}
	}
	return Invalid, fmt.Errorf("no extension found for '%s'", extensionName)
}

type Repository interface {
	GetAll() ([]Image, error)
	Get(id string) (*Image, error)
	GetPath(id string) (string, error)
	Add(image *Image) error
	Update(id, path string, imageMeta *Meta) error
}
