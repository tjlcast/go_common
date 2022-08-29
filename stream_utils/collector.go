package stream_utils

func ToList(eles *Stream) []T {
	arr := make([]T, eles.Length)
	for i:=0; i<len(arr); i++ {
		arr[i] = eles.Head
		eles = eles.Tail
	}
	return arr
}

func ToBoolList(eles *Stream) []bool {
	arr := make([]bool, eles.Length)
	for i:=0; i<len(arr); i++ {
		arr[i] = eles.Head.(bool)
		eles = eles.Tail
	}
	return arr
}

func ToStringList(eles *Stream) []string {
	arr := make([]string, eles.Length)
	for i:=0; i<len(arr); i++ {
		arr[i] = eles.Head.(string)
		eles = eles.Tail
	}
	return arr
}

func ToIntList(eles *Stream) []int {
	arr := make([]int, eles.Length)
	for i:=0; i<len(arr); i++ {
		arr[i] = eles.Head.(int)
		eles = eles.Tail
	}
	return arr
}

func ToInt16List(eles *Stream) []int16 {
	arr := make([]int16, eles.Length)
	for i:=0; i<len(arr); i++ {
		arr[i] = eles.Head.(int16)
		eles = eles.Tail
	}
	return arr
}

func ToInt32List(eles *Stream) []int32 {
	arr := make([]int32, eles.Length)
	for i:=0; i<len(arr); i++ {
		arr[i] = eles.Head.(int32)
		eles = eles.Tail
	}
	return arr
}

func ToInt64List(eles *Stream) []int64 {
	arr := make([]int64, eles.Length)
	for i:=0; i<len(arr); i++ {
		arr[i] = eles.Head.(int64)
		eles = eles.Tail
	}
	return arr
}

func ToFloat32List(eles *Stream) []float32 {
	arr := make([]float32, eles.Length)
	for i:=0; i<len(arr); i++ {
		arr[i] = eles.Head.(float32)
		eles = eles.Tail
	}
	return arr
}

func ToFloat64List(eles *Stream) []float64 {
	arr := make([]float64, eles.Length)
	for i:=0; i<len(arr); i++ {
		arr[i] = eles.Head.(float64)
		eles = eles.Tail
	}
	return arr
}