package file_utils

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func IsDir(path string) bool {
	fi, e := os.Stat(path)
	if e != nil {
		return false
	}
	return fi.IsDir()
}

func IsExist(path string) bool {
	_, e := os.Stat(path)
	if e != nil {
		return false
	}
	return true
}

func CreateDir(path string) bool {
	err := os.Mkdir(path, 0755)
	if err != nil {
		return false
	}
	return true
}

func ListDir(path string) []string {
	var s []string
	paths, _ := ioutil.ReadDir(path)

	for _, fi := range paths {
		fullName := path + "/" + fi.Name()
		s = append(s, fullName)
	}
	return s
}

func GetFileNameInPath(path string) string {
	words := strings.Split(path, string(os.PathSeparator))
	return words[len(words)-1]
}

func OpenClearFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
}

func OpenAppendFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
}

func SaveFile(filename string, jsonBytes []byte) (bool, error) {
	f, err := OpenClearFile(filename)
	if err != nil {
		return false, errors.New("Fail to open: " + filename)
	}
	_, err = f.Write(jsonBytes)
	if err != nil {
		return false, err
	}
	defer f.Close()
	return true, nil
}

func ReadFile(filename string) ([]byte, error) {
	byteValue, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.New("Fail to open: " + filename)
	}
	return byteValue, nil
}

type SimpleFileReader struct {
	data []byte
}

func NewSimpleFileReader(data []byte) *SimpleFileReader {
	return &SimpleFileReader{
		data: data,
	}
}

func (reader *SimpleFileReader) Close() {

}

func (reader *SimpleFileReader) Empty() bool {
	return len(reader.data) <= 0
}

func (reader *SimpleFileReader) Size() int {
	return len(reader.data)
}

func (reader *SimpleFileReader) ReadId() string {
	bytes := reader.readUint64()
	return strconv.Itoa(int(bytes))
}

func (reader *SimpleFileReader) ReadFeature(dim uint32) string {
	bytes := reader.readBytes(dim * 4)
	featStr := base64.StdEncoding.EncodeToString(bytes)
	return featStr
}

func (reader *SimpleFileReader) readUint32() uint32 {
	val := binary.LittleEndian.Uint32(reader.data)
	reader.data = reader.data[4:]
	return val
}

func (reader *SimpleFileReader) readBytes(length uint32) []byte {
	bytes := reader.data[:length]
	reader.data = reader.data[length:]
	return bytes
}

func (reader *SimpleFileReader) readUint64() uint64 {
	val := binary.LittleEndian.Uint64(reader.data)
	reader.data = reader.data[8:]
	return val
}

type SimpleFileReaderV2 struct {
	fileName string
	file     *os.File
	isEmpty  bool
	dim      int

	buff64  []byte
	buffDim []byte
}

func NewSimpleFileReaderV2(fileName string, dim int) *SimpleFileReaderV2 {
	f, err := os.Open(fileName)
	if err != nil {
		panic("Fail to open " + err.Error())
	}
	return &SimpleFileReaderV2{
		fileName: fileName,
		file:     f,
		dim:      dim,
		isEmpty:  false,
		buff64:   make([]byte, 8),
		buffDim:  make([]byte, dim*4),
	}
}

func (reader *SimpleFileReaderV2) Empty() bool {
	return reader.isEmpty
}

func (reader *SimpleFileReaderV2) Size() int {
	panic("Not support")
}

func (reader *SimpleFileReaderV2) ReadIdBytes() []byte {
	nr, _ := reader.file.Read(reader.buff64)
	if nr <= 0 {
		reader.isEmpty = true
	}
	return reader.buff64
}

func (reader *SimpleFileReaderV2) ReadId() string {
	bytes := reader.ReadIdBytes()
	return strconv.Itoa(int(binary.LittleEndian.Uint64(bytes)))
}

func (reader *SimpleFileReaderV2) ReadFeatureBytes(dim uint32) []byte {
	bytes := reader.readBytes(dim * 4)
	return bytes
}

func (reader *SimpleFileReaderV2) ReadFeature(dim uint32) string {
	bytes := reader.ReadFeatureBytes(dim)
	featStr := base64.StdEncoding.EncodeToString(bytes)
	return featStr
}

func (reader *SimpleFileReaderV2) readBytes(length uint32) []byte {
	nr, _ := reader.file.Read(reader.buffDim)
	if nr <= 0 {
		reader.isEmpty = true
	}

	return reader.buffDim
}

func (reader *SimpleFileReaderV2) Close() {
	reader.file.Close()
}
