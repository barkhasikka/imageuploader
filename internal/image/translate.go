package imageservice

import imagerepository "imageuploader/internal/image/repository"

func translateRepoToServiceModel(image *imagerepository.Image) Image {
	return Image{
		Filesize:    image.Meta.Filesize,
		Height:      image.Meta.Height,
		Width:       image.Meta.Width,
		Extension:   image.Meta.Extension.String(),
		CreatedAt:   image.Details.CreatedAt,
		Title:       image.Details.Title,
		Path:        image.Details.Path,
		ID:          image.ID,
		ContentType: image.Meta.ContentType,
	}
}

func translateServiceToRepoModel(image Image) *imagerepository.Image {
	extension, _ := imagerepository.ExtensionFromString(image.Extension)
	return &imagerepository.Image{
		ID: image.ID,
		Meta: &imagerepository.Meta{
			Filesize:    image.Filesize,
			Height:      image.Height,
			Width:       image.Width,
			Extension:   extension,
			ContentType: image.ContentType,
		},
		Details: &imagerepository.Details{
			Title: image.Title,
			Path:  image.Path,
		},
	}
}
