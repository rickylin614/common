package utils

import (
	"reflect"
	"sort"
)

type List struct {
	Data     []reflect.Value
	sorted   []string
	reversed []bool
}

func NewList(list interface{}) *List {
	f := reflect.ValueOf(list)
	kind := f.Kind()
	data := make([]reflect.Value, 0)
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < f.Len(); i++ {
			data = append(data, f.Index(i))
		}
		return &List{Data: data}
	}
	return &List{}
}

func (l *List) SortBy(filename string, reverse ...bool) *List {
	// 陣列小於2不作排序
	if len(l.Data) < 2 {
		return l
	}
	// 紀錄排序過程
	l.sorted = []string{filename}
	if len(reverse) > 0 {
		l.reversed = append(l.reversed, reverse[0])
	} else {
		l.reversed = append(l.reversed, false)
	}
	// 排序
	sort.Slice(l.Data, func(i, j int) bool {
		return compare(l.Data[i], l.Data[j], l, 0)
	})
	return l
}

func (l *List) ThenSort(filename string, reverse ...bool) *List {
	// 陣列小於2不作排序
	if len(l.Data) < 2 {
		return l
	}
	// 紀錄排序過程
	l.sorted = append(l.sorted, filename)
	if len(reverse) > 0 {
		l.reversed = append(l.reversed, reverse[0])
	} else {
		l.reversed = append(l.reversed, false)
	}
	// 排序
	sort.Slice(l.Data, func(i, j int) bool {
		return compare(l.Data[i], l.Data[j], l, 0)
	})
	return l
}

func (l *List) Reverse(filename string) *List {
	for i, j := 0, len(l.Data)-1; i < j; i, j = i+1, j-1 {
		l.Data[i], l.Data[j] = l.Data[j], l.Data[i]
	}
	return l
}

func (l *List) ToList() []interface{} {
	res := make([]interface{}, 0)
	for _, v := range l.Data {
		if v.CanInterface() {
			res = append(res, v.Interface())
		}
	}
	return res
}

func (l *List) Contains(value interface{}) bool {
	if len(l.Data) == 0 {
		return false
	}
	rv := reflect.ValueOf(value)
	if l.Data[0].Kind() != rv.Kind() {
		return false
	}

	for _, v := range l.Data {
		switch rv.Kind() {
		case reflect.String:
			if v.String() == rv.String() {
				return true
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if uint64(v.Int()) == uint64(rv.Int()) {
				return true
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if v.Uint() == rv.Uint() {
				return true
			}
		case reflect.Float32, reflect.Float64:
			if v.Float() == rv.Float() {
				return true
			}
		case reflect.Ptr, reflect.Interface:
			if v.Addr() == rv.Addr() {
				return true
			}
		}
	}
	return false
}

func (l *List) Len() int {
	return len(l.Data)
}

/*
	@Author: ansel
	compare with reflect.Value
*/
func compare(d1, d2 reflect.Value, l *List, level int) bool {
	// 是否最後一層
	sorted := l.sorted
	if len(sorted) <= level {
		return false
	}
	reversed := l.reversed[level]
	// 抓取映射欄位資料
	v1 := d1.FieldByName(sorted[level])
	v2 := d2.FieldByName(sorted[level])
	if v1.Kind().String() != v2.Kind().String() {
		return false
	}
	// 判斷欄位種類
	switch v1.Kind() {
	case reflect.String:
		// 相同情形 下一層繼續判斷
		if v1.String() == v2.String() {
			return compare(d1, d2, l, level+1)
		}
		// 是否逆序
		if reversed {
			return v1.String() > v2.String()
		}
		return v1.String() < v2.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// 相同情形 下一層繼續判斷
		if uint64(v1.Int()) == uint64(v2.Int()) {
			return compare(d1, d2, l, level+1)
		}
		// 是否逆序
		if reversed {
			return uint64(v1.Int()) > uint64(v2.Int())
		}
		return uint64(v1.Int()) < uint64(v2.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		// 相同情形 下一層繼續判斷
		if v1.Uint() == v2.Uint() {
			return compare(d1, d2, l, level+1)
		}
		// 是否逆序
		if reversed {
			return v1.Uint() > v2.Uint()
		}
		return v1.Uint() < v2.Uint()
	case reflect.Float32, reflect.Float64:
		// 相同情形 下一層繼續判斷
		if v1.Float() == v2.Float() {
			return compare(d1, d2, l, level+1)
		}
		// 是否逆序
		if reversed {
			return v1.Float() > v2.Float()
		}
		return v1.Float() < v2.Float()
	}
	return false
}
