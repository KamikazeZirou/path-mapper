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
	case reflect.Int:
		if v, err := strconv.ParseInt(v, 10, 0); err == nil {
			f.SetInt(v)
		} else {
			return fmt.Errorf("failed mapping %v to %v because it is not int", v, f)
		}
	case reflect.String:
		f.SetString(v)
	}

	return nil
}
