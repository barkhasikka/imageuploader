package bolthelper

import (
	"time"

	"imageuploader/pkg/utils/conversion"

	"github.com/boltdb/bolt"
)

func DBSerializeString(bucket *bolt.Bucket, key, value string) {
	bucket.Put([]byte(key), []byte(value))
}

func DBSerializeBool(bucket *bolt.Bucket, key string, value bool) {
	valueByte := conversion.BinBool(value)
	bucket.Put([]byte(key), valueByte)
}

func DBSerializeUint16(bucket *bolt.Bucket, key string, value uint16) {
	valueByte := conversion.Uitob16(value)
	bucket.Put([]byte(key), valueByte)
}

func DBSerializeUint64(bucket *bolt.Bucket, key string, value uint64) {
	valueByte := conversion.Uitob64(value)
	bucket.Put([]byte(key), valueByte)
}

func DBSerializeInt(bucket *bolt.Bucket, key string, value int) {
	valueByte := conversion.Itob(int(value))
	bucket.Put([]byte(key), valueByte)
}

func DBSerializeUTCTime(bucket *bolt.Bucket, key string, time time.Time) {
	timeInString := conversion.GetTimeStringInUTC(time)
	DBSerializeString(bucket, key, timeInString)
}

func DBSerializeFloat64(bucket *bolt.Bucket, key string, value float64) {
	bucket.Put([]byte(key), conversion.Float64ToByte(value))
}

func DBSerializeFloat32(bucket *bolt.Bucket, key string, value float32) {
	bucket.Put([]byte(key), conversion.Float32ToByte(value))
}
