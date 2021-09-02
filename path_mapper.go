package path_mapper

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/KamikazeZirou/path-mapper/internal/reflectx"
)

type MapperFunc func(v string) (interface{}, error)

var (
	Mapper = make(map[string]MapperFunc)
)

func LcFirst(s string) string {
	for i, v := range s {
		return string(unicode.ToLower(v)) + s[i+1:]
	}
	return ""
}

// Mapping a URL or other path to a structure.
//goland:noinspection GoUnusedExportedFunction
func Mapping(pattern, path string, dest interface{}) error {
	patternSegments := strings.Split(pattern, "/")
	pathSegments := strings.Split(path, "/")
	if len(pathSegments) != len(patternSegments) {
		return fmt.Errorf("pattern(%value) does not match path(%value)", pattern, path)
	}

	patterns := make([]string, 0, len(patternSegments))
	values := make([]string, 0, len(patternSegments))
	for i := 0; i < len(patternSegments); i++ {
		patternSegment := patternSegments[i]
		pathSegment := pathSegments[i]

		if strings.HasPrefix(patternSegment, "{") && strings.HasSuffix(patternSegment, "}") {
			patterns = append(patterns, patternSegment[1:len(patternSegment)-1])
			values = append(values, pathSegment)
		} else if pathSegment != patternSegment {
			return fmt.Errorf("pattern(%value) does not match path(%value)", pattern, path)
		}
	}

	v := reflect.ValueOf(dest)

	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to dest")
	}

	if v.IsNil() {
		return errors.New("must pass non-nil pointer to dest")
	}

	m := reflectx.NewMapperFunc("alias", LcFirst)
	traversals := m.TraversalsByName(v.Type(), patterns)
	fields := make([]interface{}, len(patterns))
	if err := fieldsByTraversal(v, traversals, fields, true); err != nil {
		return err
	}

	for i, value := range values {
		if len(traversals[i]) == 0 {
			// Allow missing fields
			continue
		}

		if err := convertAssign(value, fields[i]); err != nil {
			return fmt.Errorf("failed mapping %v into %v : %w", value, fields[i], err)
		}
	}

	return nil
}

func fieldsByTraversal(v reflect.Value, traversals [][]int, values []interface{}, ptrs bool) error {
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return errors.New("argument not a struct")
	}

	for i, traversal := range traversals {
		if len(traversal) == 0 {
			values[i] = new(interface{})
			continue
		}
		f := reflectx.FieldByIndexes(v, traversal)
		if ptrs {
			values[i] = f.Addr().Interface()
		} else {
			values[i] = f.Interface()
		}
	}
	return nil
}

func convertAssign(src string, dest interface{}) error {
	switch d := dest.(type) {
	case *string:
		*d = src
		return nil
	case *int:
		if v, err := strconv.ParseInt(src, 10, 0); err == nil {
			*d = int(v)
			return nil
		} else {
			return fmt.Errorf("%v is invalid as int", src)
		}
	case *int8:
		if v, err := strconv.ParseInt(src, 10, 0); err == nil {
			*d = int8(v)
			return nil
		} else {
			return fmt.Errorf("%v is invalid as int", src)
		}
	case *int16:
		if v, err := strconv.ParseInt(src, 10, 0); err == nil {
			*d = int16(v)
			return nil
		} else {
			return fmt.Errorf("%v is invalid as int", src)
		}
	case *int32:
		if v, err := strconv.ParseInt(src, 10, 0); err == nil {
			*d = int32(v)
			return nil
		} else {
			return fmt.Errorf("%v is invalid as int", src)
		}
	case *int64:
		if v, err := strconv.ParseInt(src, 10, 0); err == nil {
			*d = v
			return nil
		} else {
			return fmt.Errorf("%v is invalid as int", src)
		}
	case *uint:
		if v, err := strconv.ParseUint(src, 10, 0); err == nil {
			*d = uint(v)
			return nil
		} else {
			return fmt.Errorf("%v is invalid as uint", src)
		}
	case *uint8:
		if v, err := strconv.ParseUint(src, 10, 0); err == nil {
			*d = uint8(v)
			return nil
		} else {
			return fmt.Errorf("%v is invalid as uint", src)
		}
	case *uint16:
		if v, err := strconv.ParseUint(src, 10, 0); err == nil {
			*d = uint16(v)
			return nil
		} else {
			return fmt.Errorf("%v is invalid as uint", src)
		}
	case *uint32:
		if v, err := strconv.ParseUint(src, 10, 0); err == nil {
			*d = uint32(v)
			return nil
		} else {
			return fmt.Errorf("%v is invalid as uint", src)
		}
	case *uint64:
		if v, err := strconv.ParseUint(src, 10, 0); err == nil {
			*d = v
			return nil
		} else {
			return fmt.Errorf("%v is invalid as uint", src)
		}
	}

	dpv := reflect.ValueOf(dest)
	if dpv.Kind() != reflect.Ptr {
		return errors.New("destination not a pointer")
	}

	if dpv.IsNil() {
		return errors.New("destination pointer is nil")
	}

	dv := reflect.Indirect(dpv)
	switch dv.Kind() {
	case reflect.Ptr:
		dv.Set(reflect.New(dv.Type().Elem()))
		return convertAssign(src, dv.Interface())
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		if v, err := strconv.ParseInt(src, 10, 0); err == nil {
			dv.SetInt(v)
			return nil
		} else {
			return fmt.Errorf("%v is invalid as int", src)
		}
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		if v, err := strconv.ParseUint(src, 10, 0); err == nil {
			dv.SetUint(v)
			return nil
		} else {
			return fmt.Errorf("%v is invalid as uint", src)
		}
	case reflect.String:
		dv.SetString(src)
		return nil
	}

	return fmt.Errorf("unsupported conversion. Value type %T into type %T", src, dest)
}
