package feature_utils

import (
	"encoding/base64"
	"encoding/binary"
	"math"
	"math/rand"
	"strconv"
)

func ConvFeat2Bytes(feature string) []byte {
	bytes, _ := base64.StdEncoding.DecodeString(feature)
	return bytes
}

func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func ByteToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func ByteToFloat32Array(featurebytes []byte, start int, length int) []float32 {
	len := length / 4
	feascore := make([]float32, len, len)

	for idx, _ := range feascore {
		tmp := featurebytes[start+idx*4:]
		feascore[idx] = ByteToFloat32(tmp)
	}

	return feascore
}

func ConvFeatureStr2FloatArr(feature string) []float32 {
	featureByte := ConvFeat2Bytes(feature)
	featureFloat := ByteToFloat32Array(featureByte, 0, len(featureByte))
	return featureFloat
}

func ConvFeatureFloat2Bytes(feat []float32) []byte{
	var res []byte
	for _, f := range feat {
		bytes := Float32ToByte(f)
		res = append(res, bytes...)
	}
	return res
}

func ConvId2Bytes(idStr string) []byte {
	id, _ := strconv.Atoi(idStr)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(id))
	return bytes
}

func MockFeature(dim int) string {
	var byteArr []byte
	for i := 0; i < dim; i++ {
		rfloat32 := rand.Float32()
		bytes := Float32ToByte(rfloat32)
		byteArr = append(byteArr, bytes...)
	}
	return base64.StdEncoding.EncodeToString(byteArr)
}
