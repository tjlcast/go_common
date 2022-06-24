package feature_utils

import (
	"encoding/base64"
	"fmt"
	"runtime"
)

func DecodeBase64FeatureSlice(features []string) ([]float32, error) {
	var floats []float32 = make([]float32, 0)
	for _, feature := range features {
		fs, err := DecodeBase64Feature(feature)
		if err != nil {
			//log.Error(err.Error())
			fmt.Printf(err.Error())
			return nil, err
		}
		floats = append(floats, fs...)
	}
	return floats, nil
}

func DecodeBase64Feature(feature string) ([]float32, error) {
	bytes, err := base64.StdEncoding.DecodeString(feature)
	if err != nil {
		return nil, err
	}
	if len(bytes)%4 != 0 {
		return nil, fmt.Errorf("error feature length: %d", len(bytes))
	}
	floats := make([]float32, len(bytes)/4)
	parallel := runtime.GOMAXPROCS(0)
	dim := len(floats)
	if parallel > dim {
		parallel = dim
	}
	step := dim / parallel
	c := make(chan int, parallel)
	for i := 1; i <= parallel; i++ {
		start := (i - 1) * step
		end := i * step
		if i == parallel {
			end = dim
		}
		go func(i int, j int, c chan int) {
			for k := i; k < j; k++ {
				floats[k] = ByteToFloat32(bytes[k*4 : (k+1)*4])
			}
			c <- 1
		}(start, end, c)
	}
	for i := 0; i < parallel; i++ {
		<-c
	}

	return floats, nil
}
