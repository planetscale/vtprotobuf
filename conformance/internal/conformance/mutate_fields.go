// Copyright (c) 2022 PlanetScale Inc. All rights reserved.

package conformance

import "reflect"

// MutateFields modifies the value of each field of structs recursively.
func MutateFields(x interface{}) {
	mut(reflect.TypeOf(x), reflect.ValueOf(x))
}

func mut(t reflect.Type, v reflect.Value) {
	switch t.Kind() {
	case reflect.String:
		v.SetString("blip blop")

	case reflect.Uint, reflect.Uint64:
		v.SetUint(v.Uint() + 1)
	case reflect.Uint8:
		v.SetUint(uint64(v.Uint() + 1))
	case reflect.Uint16:
		v.SetUint(uint64(v.Uint() + 1))
	case reflect.Uint32:
		v.SetUint(uint64(v.Uint() + 1))

	case reflect.Int, reflect.Int64:
		v.SetInt(v.Int() + 1)
	case reflect.Int8:
		v.SetInt(int64(v.Int() + 1))
	case reflect.Int16:
		v.SetInt(int64(v.Int() + 1))
	case reflect.Int32:
		v.SetInt(int64(v.Int() + 1))

	case reflect.Float32, reflect.Float64:
		if x := v.Float(); x != 0.0 && x != -0.0 {
			v.SetFloat(-1 * x)
			return
		}
		v.SetFloat(1.0)

	case reflect.Bool:
		v.SetBool(!v.Bool())

	case reflect.Ptr:
		if v.IsNil() {
			return
		}
		mut(t.Elem(), v.Elem())

	case reflect.Array, reflect.Slice:
		n := v.Len()
		elemT := t.Elem()
		for i := 0; i < n; i++ {
			mut(elemT, v.Index(i))
		}
		if n == 0 {
			n = 3
			for i := 0; i < n; i++ {
				nv := reflect.New(elemT)
				mut(elemT, nv.Elem())
				v.Set(reflect.Append(reflect.Indirect(v), reflect.Indirect(nv)))
			}
		}

	case reflect.Struct:
		for i, n := 0, t.NumField(); i < n; i++ {
			elementT := t.Field(i)
			elementV := v.Field(i)
			if elementV.CanSet() || elementT.Anonymous {
				mut(elementT.Type, elementV)
			}
		}

	case reflect.Map:
		m := reflect.MakeMap(reflect.MapOf(t.Key(), t.Elem()))
		for _, mapKey := range v.MapKeys() {
			mapIndex := reflect.New(t.Key()).Elem()
			mapIndex.Set(mapKey)
			mut(t.Key(), mapIndex)
			mapValue := reflect.New(t.Elem()).Elem()
			mapValue.Set(v.MapIndex(mapKey))
			mut(t.Elem(), mapValue)
			m.SetMapIndex(mapIndex, mapValue)
		}
		v.Set(m)
	}
}

// VisitWithPredicate deep-visits the given struct
// and returns whether the predicate holds at least once.
func VisitWithPredicate(x interface{}, f func(w interface{}) bool) bool {
	return vwp(f, false, reflect.TypeOf(x), reflect.ValueOf(x))
}

func vwp(f func(w interface{}) bool, acc bool, t reflect.Type, v reflect.Value) bool {
	switch t.Kind() {
	case reflect.String:
		return acc || f(v.String())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return acc || f(v.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return acc || f(v.Int())
	case reflect.Float32, reflect.Float64:
		return acc || f(v.Float())
	case reflect.Bool:
		return acc || f(v.Bool())

	case reflect.Ptr:
		if v.IsNil() {
			return acc
		}
		return acc || vwp(f, acc, t.Elem(), v.Elem())

	case reflect.Array, reflect.Slice:
		for i, n, elemT := 0, v.Len(), t.Elem(); i < n; i++ {
			if acc || vwp(f, acc, elemT, v.Index(i)) {
				return true
			}
		}

	case reflect.Struct:
		for i, n := 0, t.NumField(); i < n; i++ {
			elementT := t.Field(i)
			elementV := v.Field(i)
			if elementV.CanSet() || elementT.Anonymous {
				if acc || vwp(f, acc, elementT.Type, elementV) {
					return true
				}
			}
		}

	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			if acc || VisitWithPredicate(iter.Key().Interface(), f) || VisitWithPredicate(iter.Value().Interface(), f) {
				return true
			}
		}
	}
	return acc
}
