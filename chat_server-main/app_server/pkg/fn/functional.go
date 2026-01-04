// Package glib provides generic functional programming utilities.
package fn

import (
	"golang.org/x/exp/constraints"
)

// Map returns a new slice with the results of applying the given function to each element of the input slice.
// 接受一个slice 和一个函数，返回一个新slice，新slice的元素是原slice的元素经过函数处理后的结果
func Map[T any, R any](arr []T, fn func(T) R) []R {
	result := make([]R, 0, len(arr))
	for _, item := range arr {
		result = append(result, fn(item))
	}
	return result
}

// Partial returns new function that, when called, has its first argument set to the provided value.
// 接受一个二维的函数，返回一个已经传入第一个参数的偏函数，偏函数为一维函数。新函数在调用时，第一个参数被设置为提供的值
// 适用于与其他functional的函数组合使用 如 Map(numbers, Partial(math.Max,0))
func Partial[T1, T2, R any](f func(a T1, b T2) R, arg1 T1) func(T2) R {
	return func(t2 T2) R {
		return f(arg1, t2)
	}
}

// Partial returns new function that, when called, has its first argument set to the provided value.
// 接受一个二维的函数，返回一个已经传入第二个参数的偏函数，偏函数为一维函数。新函数在调用时，第二个参数被设置为提供的值
// 适用于与其他functional的函数组合使用 如 Map(numbers, PartialR(strconv.FormatInt,10))
func PartialR[T1, T2, R any](f func(a T1, b T2) R, arg1 T2) func(T1) R {
	return func(t1 T1) R {
		return f(t1, arg1)
	}
}

// DropLast
// 接受一个一出二的函数 返回一个丢弃第二个结果的一维函数
// 适用于与其他functional的函数组合使用 如 Map(req.Msg.TaskId, DropLast(strconv.Atoi))
func DropLast[P, R1, R2 any](f func(p P) (R1, R2)) func(P) R1 {
	return func(p P) R1 {
		r1, _ := f(p)
		return r1
	}
}

// Cast 将一个类型转换为另一个类型
// 适用于与其他functional的函数组合使用 Map(req.Status, Cast[MaterialPushStatus, int])
func Cast[T, R constraints.Signed](item T) R {
	return R(item)
}

// Ptr 将一个值转换为指针
// 适用于与其他functional的函数组合使用 pointers := Map(values, Ptr)
func Ptr[T any](item T) *T {
	return &item
}

// DerefOr0 取指针，如果指针为nil，则返回零值
// 适用于与其他functional的函数组合使用 Map(pointers, DerefOr0)
func DerefOr0[T any](item *T) T {
	if item == nil {
		return *new(T)
	}
	return *item
}

// Filter 接受一个slice和一个函数，返回一个新slice，新slice的元素是原slice中满足函数条件的元素
// 适用于与其他functional的函数组合使用 如 Filter(numbers, lo.IsNotEmpty)
func Filter[T any](arr []T, fn func(T) bool) []T {
	result := make([]T, 0, len(arr))
	for _, item := range arr {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}

func Drop1Of2[T1, T2 any](v1 T1, v2 T2) T1 {
	return v1
}
