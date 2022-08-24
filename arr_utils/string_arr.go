package arr_utils

type ArrString []string

func (arr *ArrString) Append(ele string) {
	*arr = append(*arr, ele)
}
