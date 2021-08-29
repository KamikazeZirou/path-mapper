package path_mapper

import (
	"reflect"
)

// Mapping a URL or other path to a structure.
//goland:noinspection GoUnusedExportedFunction
func Mapping(pattern, path string, st interface{}) {
	sv := reflect.ValueOf(st).Elem()
	sv.FieldByName("Owner").SetString("KamikazeZirou")
	sv.FieldByName("Repository").SetString("path-mapper")
	sv.FieldByName("Number").SetInt(1)
}
