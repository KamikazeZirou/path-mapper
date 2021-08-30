package path_mapper

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type MapperFunc func(v string) (interface{}, error)

var (
	Mapper = make(map[string]MapperFunc)
)

// Mapping a URL or other path to a structure.
//goland:noinspection GoUnusedExportedFunction
func Mapping(pattern, path string, st interface{}) error {
	patternSegments := strings.Split(pattern, "/")
	pathSegments := strings.Split(path, "/")

	if len(pathSegments) != len(patternSegments) {
		return fmt.Errorf("pattern(%v) does not match path(%v)", pattern, path)
	}

	v := reflect.ValueOf(st).Elem()
	sv := reflect.New(v.Elem().Type()).Elem()

	for i := 0; i < len(patternSegments); i++ {
		patternSegment := patternSegments[i]
		pathSegment := pathSegments[i]

		if strings.HasPrefix(patternSegment, "{") && strings.HasSuffix(patternSegment, "}") {
			err := setField(sv, patternSegment[1:len(patternSegment)-1], pathSegment)
			if err != nil {
				return fmt.Errorf(": %w", err)
			}
		} else if pathSegment != patternSegment {
			return fmt.Errorf("pattern(%v) does not match path(%v)", pattern, path)
		}
	}

	v.Set(sv)

	return nil
}

func setField(sv reflect.Value, n, v string) error {
	f := sv.FieldByName(strings.Title(n))
	if !f.IsValid() {
		return nil
	}

	if !f.CanSet() {
		return fmt.Errorf("%v cannot be set value", n)
	}

	if mapper, ok := Mapper[n]; ok {
		if mv, err := mapper(v); err == nil {
			f.Set(reflect.ValueOf(mv))
			return nil
		} else {
			return fmt.Errorf(": %w", err)
		}
	}

	switch f.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		if v, err := strconv.ParseInt(v, 10, 0); err == nil {
			f.SetInt(v)
		} else {
			return fmt.Errorf("failed mapping %v to %v because it is not int", v, f)
		}
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		if v, err := strconv.ParseUint(v, 10, 0); err == nil {
			f.SetUint(v)
		} else {
			return fmt.Errorf("failed mapping %v to %v because it is not int", v, f)
		}
	case reflect.String:
		f.SetString(v)
	case reflect.Ptr:
		t := deref(f.Type())
		switch t.Kind() {
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:
			if v, err := strconv.ParseInt(v, 10, 0); err == nil {
				setIntPtr(f, t.Kind(), v)
			} else {
				return fmt.Errorf("failed mapping %v to %v because it is not int", v, f)
			}
		case reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			if v, err := strconv.ParseUint(v, 10, 0); err == nil {
				setUintPtr(f, t.Kind(), v)
			} else {
				return fmt.Errorf("failed mapping %v to %v because it is not int", v, f)
			}
		case reflect.String:
			f.Set(reflect.ValueOf(&v))
		}
	}

	return nil
}

func setIntPtr(p reflect.Value, k reflect.Kind, x int64) {
	switch k {
	case reflect.Int:
		v := int(x)
		p.Set(reflect.ValueOf(&v))
	case reflect.Int8:
		v := int8(x)
		p.Set(reflect.ValueOf(&v))
	case reflect.Int16:
		v := int16(x)
		p.Set(reflect.ValueOf(&v))
	case reflect.Int32:
		v := int32(x)
		p.Set(reflect.ValueOf(&v))
	case reflect.Int64:
		p.Set(reflect.ValueOf(&x))
	}
}

func setUintPtr(p reflect.Value, k reflect.Kind, x uint64) {
	switch k {
	case reflect.Uint:
		v := uint(x)
		p.Set(reflect.ValueOf(&v))
	case reflect.Uint8:
		v := uint8(x)
		p.Set(reflect.ValueOf(&v))
	case reflect.Uint16:
		v := uint16(x)
		p.Set(reflect.ValueOf(&v))
	case reflect.Uint32:
		v := uint32(x)
		p.Set(reflect.ValueOf(&v))
	case reflect.Uint64:
		p.Set(reflect.ValueOf(&x))
	}
}

func deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
