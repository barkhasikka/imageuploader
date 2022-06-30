package imagerepository

import (
	"log"
	"time"

	"imageuploader/pkg/bolt"
	"imageuploader/pkg/utils/conversion"

	"github.com/boltdb/bolt"
)

func put(imageBucket *bolt.Bucket, image *Image) error {
	if detailsBucket, _ := imageBucket.CreateBucketIfNotExists([]byte("details")); detailsBucket != nil {
		timeInString := conversion.GetTimeStringInUTC(time.Now())
		bolthelper.DBSerializeString(detailsBucket, "createdat", timeInString)
		bolthelper.DBSerializeString(detailsBucket, "lastupdatedat", timeInString)
		bolthelper.DBSerializeString(detailsBucket, "title", image.Details.Title)
		bolthelper.DBSerializeString(detailsBucket, "path", image.Details.Path)
	}

	if metaBucket, _ := imageBucket.CreateBucketIfNotExists([]byte("meta")); metaBucket != nil {
		bolthelper.DBSerializeUint64(metaBucket, "filesize", image.Meta.Filesize)
		bolthelper.DBSerializeUint16(metaBucket, "height", image.Meta.Height)
		bolthelper.DBSerializeUint16(metaBucket, "width", image.Meta.Width)
		bolthelper.DBSerializeString(metaBucket, "extension", image.Meta.Extension.String())
		bolthelper.DBSerializeString(metaBucket, "contenttype", image.Meta.ContentType)
	}

	return nil
}

func updateMeta(imageBucket *bolt.Bucket, path string, meta *Meta) {
	if detailsBucket := imageBucket.Bucket([]byte("details")); detailsBucket != nil {
		timeInString := conversion.GetTimeStringInUTC(time.Now())
		bolthelper.DBSerializeString(detailsBucket, "lastupdatedat", timeInString)
		bolthelper.DBSerializeString(detailsBucket, "path", path)
	}
	if metaBucket := imageBucket.Bucket([]byte("meta")); metaBucket != nil {
		bolthelper.DBSerializeUint64(metaBucket, "filesize", meta.Filesize)
		bolthelper.DBSerializeUint16(metaBucket, "height", meta.Height)
		bolthelper.DBSerializeUint16(metaBucket, "width", meta.Width)
		bolthelper.DBSerializeString(metaBucket, "extension", meta.Extension.String())
		bolthelper.DBSerializeString(metaBucket, "contenttype", meta.ContentType)
	}
}

func get(imageBucket *bolt.Bucket) (*Image, error) {
	return &Image{
		Meta:    getMeta(imageBucket),
		Details: getDetails(imageBucket),
	}, nil
}

func getDetails(imageBucket *bolt.Bucket) *Details {
	var details = &Details{}
	if detailsBucket := imageBucket.Bucket([]byte("details")); detailsBucket != nil {
		details.CreatedAt = bolthelper.DBDeSerializeString(detailsBucket, "createdat")
		details.LastUpdatedAt = bolthelper.DBDeSerializeString(detailsBucket, "lastupdatedat")
		details.Title = bolthelper.DBDeSerializeString(detailsBucket, "title")
		details.Path = bolthelper.DBDeSerializeString(detailsBucket, "path")
	}
	return details
}

func getMeta(imageBucket *bolt.Bucket) *Meta {
	var meta = &Meta{}
	if metaBucket := imageBucket.Bucket([]byte("meta")); metaBucket != nil {
		meta.Filesize = bolthelper.DBDeSerializeUint64(metaBucket, "filesize")
		meta.Height = bolthelper.DBDeSerializeUint16(metaBucket, "height")
		meta.Width = bolthelper.DBDeSerializeUint16(metaBucket, "width")
		var err error
		meta.Extension, err = ExtensionFromString(bolthelper.DBDeSerializeString(metaBucket, "extension"))
		meta.ContentType = bolthelper.DBDeSerializeString(metaBucket, "contenttype")
		if err != nil {
			log.Printf("Error in image extension %v \n", err)
		}
	}
	return meta
}
