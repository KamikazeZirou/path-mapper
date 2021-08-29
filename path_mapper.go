package path_mapper

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Mapping a URL or other path to a structure.
//goland:noinspection GoUnusedExportedFunction
func Mapping(pattern, path string, st interface{}) error {
	pathSegments := strings.Split(path, "/")
	patternSegments := strings.Split(pattern, "/")

	if len(pathSegments) != len(patternSegments) {
		return fmt.Errorf("pattern(%v) does not match path(%v)", pattern, path)
	}

	sv := reflect.ValueOf(st).Elem()
	for i := 0; i < len(pathSegments); i++ {
		pathSegment := pathSegments[i]
		patternSegment := patternSegments[i]

		if strings.HasPrefix(patternSegment, "{") && strings.HasSuffix(patternSegment, "}") {
			key := patternSegment[1 : len(patternSegment)-1]
			key = strings.Title(key)

			f := sv.FieldByName(key)
			switch f.Kind() {
			case reflect.Int:
				v, _ := strconv.Atoi(pathSegment)
				f.SetInt(int64(v))
			case reflect.String:
				f.SetString(pathSegment)
			}
		} else if pathSegment != patternSegment {
			return fmt.Errorf("pattern(%v) does not match path(%v)", pattern, path)
		}
	}
	return nil
}
