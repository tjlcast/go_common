package stream_utils

import (
	"fmt"
	"strings"
)

//泛型类型定义
type T interface{}
type U interface{}

//流计算数据结构定义
type Stream struct {
	Head     T
	Tail     *Stream
	Length   int
	NotEmpty bool
}

var Nil = Stream{}

func Generate(r Stream, f func(Stream) T, m int) Stream {

	if m == 1 {
		return r
	} else {
		return Generate(New(f(r), &r), f, m-1)
	}
}

func New(head T, tail *Stream) Stream {
	return Stream{head, tail, tail.Length + 1, true}
}

func (s Stream) Add(i T) Stream {
	return New(i, &s)
}

func (s Stream) Addall(i ...T) Stream {
	for _, v := range i {
		s = s.Add(v)
	}
	return s
}

//左折叠 用于实现 reduce 的功能
func (s Stream) FoldLeft(i U, f func(U, T) U) U {
	if s.NotEmpty {
		return s.Tail.FoldLeft(f(i, s.Head), f)
	} else {
		return i
	}
}

//右折叠
func (s Stream) FoldRight(i U, f func(U, T) U) U {
	if s.NotEmpty {
		return f(s.Tail.FoldRight(i, f), s.Head)
	} else {
		return i
	}
}

//合并两个 Stream
func (s Stream) Merge(t Stream) Stream {
	if t.NotEmpty {
		return t.FoldRight(s, func(u U, t T) U {
			return u.(Stream).Add(t)
		}).(Stream)

	} else {
		return s
	}
}

//倒序
func (s Stream) Reverse() Stream {
	return s.FoldLeft(Nil, func(u U, t T) U {
		return u.(Stream).Add(t)
	}).(Stream)
}

//Map
func (s Stream) Map(f func(T) U) Stream {
	return s.FoldRight(Nil, func(u U, t T) U {
		return u.(Stream).Add(f(t))
	}).(Stream)
}

//Reduce
func (s Stream) Reduce(i T, f func(T, T) T) T {
	if s.NotEmpty {
		return s.Tail.Reduce(f(i, s.Head), f)
	} else {
		return i
	}
}

//过滤
func (s Stream) Filter(f func(T) bool) Stream {
	return s.FoldRight(Nil, func(u U, t T) U {
		if f(t) {
			return u.(Stream).Add(t)
		} else {
			return u
		}
	}).(Stream)
}

//归并排序
func (s Stream) Sort(c func(T,T) bool) Stream {

	n := s.Length / 2

	if n == 0 {
		return s
	}else{
		x,y := split(s, Nil, n)
		return merge(x.Sort(c),y.Sort(c),c)
	}
}


func split(x,y Stream , n int) (Stream,Stream) {

	if (n == 0 || !x.NotEmpty) {
		return x,y
	}
	return split(*x.Tail, y.Add(x.Head), n - 1);
}


func merge(x,y Stream , c func(T,T) bool) Stream {

	if (!x.NotEmpty){
		return y;
	}
	if (!y.NotEmpty){
		return x;
	}
	if c(x.Head,y.Head) {
		return merge(*x.Tail, y, c).Add(x.Head)
	}else{
		return merge(x, *y.Tail, c).Add(y.Head);
	}
}

//格式化显示 Stream 的所有项
func (s Stream) ToString() string {
	return "{" + strings.Join(s.FoldRight([]string{}, func(u U, t T) U {
		return append(u.([]string), fmt.Sprintf("%v", t))
	}).([]string), ",") + "}"
}
