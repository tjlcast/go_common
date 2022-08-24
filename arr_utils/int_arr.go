package arr_utils

type ArrInt []int

func (arr *ArrInt) Append(ele int) {
	*arr = append(*arr, ele)
}

type ArrInt32 []int32

func (arr *ArrInt32) Append(ele int32) {
	*arr = append(*arr, ele)
}

type ArrInt64 []int64

func (arr *ArrInt64) Append(ele int64) {
	*arr = append(*arr, ele)
}
