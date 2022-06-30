package conversion

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
)

//ConvertByteArray ...
func ConvertByteArray(data []byte) uint16 {
	return binary.BigEndian.Uint16(data)
}

//UTILITIES

//Itob int to byte
func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(v))
	return b
}

// UI32tob ...
func UI32tob(v uint32) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint32(b, uint32(v))
	return b
}

//BtoI byte to int
func BtoI(v []byte) int {
	if len(v) < 8 {
		return 0
	}
	return int(binary.LittleEndian.Uint64(v))
}

// Btoui32 ....
func Btoui32(v []byte) uint32 {
	return (binary.LittleEndian.Uint32(v))
}

//IdFromString ...
func IdFromString(s *string) []byte {
	v, _ := strconv.Atoi(*s)
	return Itob(v)
}

//B16toi byte to uint16
func B16toi(v []byte) uint16 {
	if len(v) > 0 {
		return binary.LittleEndian.Uint16(v)
	}
	return 0
}

//Btoi64 byte to uint64
func Btoi64(v []byte) uint64 {
	return binary.LittleEndian.Uint64(v)
}

//Uitob16 uint16 to byte
//used for versioning
func Uitob16(v uint16) []byte {
	b := make([]byte, 8)

	binary.LittleEndian.PutUint16(b, v)

	return b
}

//Uitob64 uint to byte64
func Uitob64(v uint64) []byte {
	b := make([]byte, 8)

	binary.LittleEndian.PutUint64(b, v)
	return b
}

//BinBool boolean to byte
func BinBool(val bool) []byte {
	byteVal := []byte{}
	if val {
		return append(byteVal, 1&1)
	}
	return append(byteVal, 0&1)
}

//BoolInB byte to boolean
func BoolInB(v []byte) bool {
	if len(v) > 0 {
		return v[0]&0x1 == 0x1
	}
	return false
}

// Uint8ToByte to convert uint8 value to array of byte
func Uint8ToByte(f uint8) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, f)
	return buf.Bytes()
}

// UintToByte to convert uint value to array of byte
func UintToByte(f uint) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, f)
	return buf.Bytes()
}

// Float64ToByte to convert float64 to []byte
func Float64ToByte(f float64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

// ByteToFloat64 to convert []byte to float64
func ByteToFloat64(b []byte) float64 {
	bits := binary.LittleEndian.Uint64(b)
	float := math.Float64frombits(bits)
	return float
}

//Float32ToByte will convert float32 to []byte
func Float32ToByte(f float32) []byte {
	bits := math.Float32bits(f)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

//ByteToFloat32 will convert []byte to float32
func ByteToFloat32(b []byte) float32 {
	if len(b) == 0 {
		return 0
	}
	bits := binary.LittleEndian.Uint32(b)
	float := math.Float32frombits(bits)
	return float
}

//ConvertGBToMB function needs value in will accept gb value and will return mb
func ConvertGBToMB(maxGBSize int) uint64 {
	return uint64(maxGBSize * 1024)
}

//ConvertMinutesToSeconds will take count in minutes and will return seconds
func ConvertMinutesToSeconds(minute int) uint {
	return uint(minute * 60)
}

// Float64ToString will convert float 64 to string
func Float64ToString(input float64, prec int) string {
	return strconv.FormatFloat(input, 'f', prec, 64)
}

// Float32ToString will convert float 32 to string
func Float32ToString(input float32, prec int) string {
	format := "%." + strconv.Itoa(prec) + "f"
	return fmt.Sprintf(format, input)
}

// StringToFloat32 will convert float 32 to string
func StringToFloat32(input string) (float32, error) {
	value, err := strconv.ParseFloat(input, 32)
	if err != nil {
		return 0, err
	}
	return float32(value), nil
}
