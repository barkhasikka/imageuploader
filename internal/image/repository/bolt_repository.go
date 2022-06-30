package imagerepository

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

type boltRepository struct {
	db *bolt.DB
}

func NewBoltRepository(dbPath string) (*boltRepository, error) {
	imageDBPath := dbPath + "images.db"
	handle, err := bolt.Open(imageDBPath, os.ModePerm, nil)
	if err != nil || handle == nil {
		return nil, fmt.Errorf("%v [image db]: %v", ErrorOpenDB, err)
	}

	return &boltRepository{
		db: handle,
	}, nil
}

func (b *boltRepository) GetAll() ([]Image, error) {
	var images []Image
	if err := b.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Cursor()
		for imageID, _ := cursor.First(); imageID != nil; imageID, _ = cursor.Next() {
			if imageBucket := tx.Bucket(imageID); imageBucket != nil {
				imageMeta, err := get(imageBucket)

				if err != nil || imageMeta == nil {
					log.Printf("Error finding image %v \n", err)
				}
				imageMeta.ID = string(imageID)
				images = append(images, *imageMeta)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return images, nil
}

func (b *boltRepository) Get(id string) (*Image, error) {
	var imageMeta *Image
	if err := b.db.View(func(tx *bolt.Tx) error {
		if imageBucket := tx.Bucket([]byte(id)); imageBucket != nil {
			var err error
			imageMeta, err = get(imageBucket)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return imageMeta, nil

}

func (b *boltRepository) GetPath(id string) (string, error) {
	var path string
	if err := b.db.View(func(tx *bolt.Tx) error {
		if imageBucket := tx.Bucket([]byte(id)); imageBucket != nil {
			if details := getDetails(imageBucket); details != nil {
				path = details.Path
			}
		}
		return nil
	}); err != nil {
		return "", err
	}
	return path, nil
}

func (b *boltRepository) Add(image *Image) error {
	if err := b.db.Update(func(tx *bolt.Tx) error {
		if imageBucket, _ := tx.CreateBucketIfNotExists([]byte(image.ID)); imageBucket != nil {
			err := put(imageBucket, image)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (b *boltRepository) Update(id, path string, imageMeta *Meta) error {
	if err := b.db.Update(func(tx *bolt.Tx) error {
		if imageBucket := tx.Bucket([]byte(id)); imageBucket != nil {
			updateMeta(imageBucket, path, imageMeta)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
