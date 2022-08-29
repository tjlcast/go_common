package stream_utils

import (
	"fmt"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	x := Generate(Nil.Add(1), func(s Stream) T { return s.Head.(int) + 1 }, 50) // {1,2,3,...,48,49,50}
	x = x.Map(func(t T) U {
		p := t.(int) //平方映射
		return p * p
	}).Filter(func(t T) bool {
		return t.(int)%2 == 0 //偶数过滤
	})

	//计算所有项的和
	fmt.Printf("sum %s = %d\n", x.ToString(), x.FoldLeft(0, func(u U, t T) U {
		return u.(int) + t.(int)
	})) //22100

	//浮点数列表求和
	y := Nil.Addall(3.5, 4.3, 2.6, 1.1, 7.83, 4.42)
	fmt.Printf("%.2f\n", y.Reduce(0.0, func(t T, t2 T) T {
		return t.(float64) + t2.(float64)
	}))

	//排序
	z := Nil.Addall(4, 3, 7, 6, 2, 1, 9, 5, 8, 0).Sort(func(x, y T) bool {
		return x.(int) > y.(int)
	})
	fmt.Println(z.ToString()) //{0,1,2,3,4,5,6,7,8,9}

	//列出包含a字符的字符串
	g := Nil.Addall("aaa", "bbb", "aba", "ccc", "cbb", "cba")
	fmt.Println(g.Filter(func(t T) bool {
		return strings.Contains(t.(string), "a")
	}).ToString())

	//生成斐波拉契亚数列 的前 20 项
	fmt.Println(Generate(Nil.Addall(1, 1), func(s Stream) T {
		return s.Head.(int) + s.Tail.Head.(int)
	}, 19).ToString())

	//通过数列 π = 2 + 2/3 + 2/3*2/5 + 2/3*2/5*3/7  + ... + f(n-1) * n/(2*n+1) 计算圆周率的值
	fmt.Println(Generate(Nil.Add(2.0), func(s Stream) T {
		n := s.Length
		return s.Head.(float64) * float64(n) / (float64(n)*2 + 1)
	}, 51).Reduce(0.0, func(t T, t2 T) T {
		return t.(float64) + t2.(float64)
	}))
}
