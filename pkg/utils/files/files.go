package files

import (
	"io"
	"log"
	"mime/multipart"
	"os"
)

func CreateDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Println(" -- error creating " + dir)
			return err
		}
	}
	return nil
}

func SaveUploadedFile(file multipart.File, out *os.File) error {

	_, err := io.Copy(out, file)
	return err
}
