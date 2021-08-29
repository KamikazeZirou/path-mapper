package path_mapper

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// Mapping a URL or other path to a structure.
//goland:noinspection GoUnusedExportedFunction
func Mapping(pattern, path string, st interface{}) error {
	s1 := strings.Split(path, "/")
	s2 := strings.Split(pattern, "/")

	if len(s1) != len(s2) {
		return errors.New("pattern does not match path")
	}

	sv := reflect.ValueOf(st).Elem()
	for i := 0; i < len(s1); i++ {
		value := s1[i]
		p2 := s2[i]

		if strings.HasPrefix(p2, "{") && strings.HasSuffix(p2, "}") {
			key := p2[1 : len(p2)-1]
			key = strings.Title(key)

			f := sv.FieldByName(key)
			switch f.Kind() {
			case reflect.Int:
				v, _ := strconv.Atoi(value)
				sv.FieldByName(key).SetInt(int64(v))
			case reflect.String:
				sv.FieldByName(key).SetString(value)
			}
		}
	}
	return nil
}
