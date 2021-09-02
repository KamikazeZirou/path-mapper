package path_mapper

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type GitHubIssue struct {
	Owner      string
	Repository string
	Number     int
	StrNumber  string
}

type Weather int32

const (
	WeatherUnknown Weather = 0
	WeatherFine    Weather = 1
)

func (w *Weather) Parse(s string) interface{} {
	switch s {
	case "fine":
		return WeatherFine
	default:
		return WeatherUnknown
	}
}

type Values struct {
	Int          int
	Int8         int8
	Int16        int16
	Int32        int32
	Int64        int64
	Uint         uint
	Uint8        uint8
	Uint16       uint16
	Uint32       uint32
	Uint64       uint64
	Str          string
	Weather      Weather
	MissingField int
}

type Pointers struct {
	Int          *int
	Int8         *int8
	Int16        *int16
	Int32        *int32
	Int64        *int64
	Uint         *uint
	Uint8        *uint8
	Uint16       *uint16
	Uint32       *uint32
	Uint64       *uint64
	Str          *string
	Weather      *Weather
	MissingField *int
}

type EmbedValues struct {
	Values
}

type EmbedPointers struct {
	*Pointers
}

func TestMapping(t *testing.T) {
	type want struct {
		st      interface{}
		success bool
	}

	type args struct {
		pattern string
		path    string
		st      interface{}
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Embed Struct",
			args: args{
				pattern: "/{int}/{int8}/{int16}/{int32}/{int64}/{uint}/{uint8}/{uint16}/{uint32}/{uint64}/{str}/{weather}",
				path:    "/1/2/3/4/5/6/7/8/9/10/abc/fine",
				st:      &EmbedValues{},
			},
			want: want{
				st: &EmbedValues{
					Values: Values{
						Int:     1,
						Int8:    2,
						Int16:   3,
						Int32:   4,
						Int64:   5,
						Uint:    6,
						Uint8:   7,
						Uint16:  8,
						Uint32:  9,
						Uint64:  10,
						Weather: WeatherFine,
						Str:     "abc",
					},
				},
				success: true,
			},
		},
		{
			name: "Embed Pointer Struct",
			args: args{
				pattern: "/{int}/{int8}/{int16}/{int32}/{int64}/{uint}/{uint8}/{uint16}/{uint32}/{uint64}/{str}/{weather}",
				path:    "/1/2/3/4/5/6/7/8/9/10/abc/fine",
				st:      &EmbedPointers{},
			},
			want: want{
				st: &EmbedPointers{
					Pointers: &Pointers{
						Int:     intAddr(1),
						Int8:    int8Addr(2),
						Int16:   int16Addr(3),
						Int32:   int32Addr(4),
						Int64:   int64Addr(5),
						Uint:    uintAddr(6),
						Uint8:   uint8Addr(7),
						Uint16:  uint16Addr(8),
						Uint32:  uint32Addr(9),
						Uint64:  uint64Addr(10),
						Str:     strAddr("abc"),
						Weather: weatherAddr(WeatherFine),
					},
				},
				success: true,
			},
		},
		{
			name: "Numbers",
			args: args{
				pattern: "/{int}/{int8}/{int16}/{int32}/{int64}/{uint}/{uint8}/{uint16}/{uint32}/{uint64}/{str}",
				path:    "/1/2/3/4/5/6/7/8/9/10/abc",
				st:      &Values{},
			},
			want: want{
				st: &Values{
					Int:    1,
					Int8:   2,
					Int16:  3,
					Int32:  4,
					Int64:  5,
					Uint:   6,
					Uint8:  7,
					Uint16: 8,
					Uint32: 9,
					Uint64: 10,
					Str:    "abc",
				},
				success: true,
			},
		},
		{
			name: "Number Pointers",
			args: args{
				pattern: "/{int}/{int8}/{int16}/{int32}/{int64}/{uint}/{uint8}/{uint16}/{uint32}/{uint64}/{str}",
				path:    "/1/2/3/4/5/6/7/8/9/10/abc",
				st:      &Pointers{},
			},
			want: want{
				st: &Pointers{
					Int:    intAddr(1),
					Int8:   int8Addr(2),
					Int16:  int16Addr(3),
					Int32:  int32Addr(4),
					Int64:  int64Addr(5),
					Uint:   uintAddr(6),
					Uint8:  uint8Addr(7),
					Uint16: uint16Addr(8),
					Uint32: uint32Addr(9),
					Uint64: uint64Addr(10),
					Str:    strAddr("abc"),
				},
				success: true,
			},
		},
		{
			name: "There is no field corresponding to the pattern in the structure to be mapped.",
			args: args{
				pattern: "/{owner}/{repository}/actions/runs/{buildNumber}",
				path:    "/guest/sandbox/actions/runs/1",
				st:      &GitHubIssue{},
			},
			want: want{
				st: &GitHubIssue{
					Owner:      "guest",
					Repository: "sandbox",
					Number:     0,
				},
				success: true,
			},
		},
		//{
		//	name: "Custom Mapper",
		//	args: args{
		//		pattern: "/{owner}/{repository}/issues/{strNumber}",
		//		path:    "/KamikazeZirou/path-mapper/issues/1",
		//		st:      &GitHubIssue{},
		//	},
		//	want: want{
		//		st: &GitHubIssue{
		//			Owner:      "KamikazeZirou",
		//			Repository: "path-mapper",
		//			Number:     0,
		//			StrNumber:  "#1",
		//		},
		//		success: true,
		//	},
		//},
		//{
		//	name: "Custom Mapper(Ptr)",
		//	args: args{
		//		pattern: "/{owner}/{repository}/issues/{strNumberPtr}",
		//		path:    "/KamikazeZirou/path-mapper/issues/1",
		//		st:      &GitHubIssuePtr{},
		//	},
		//	want: want{
		//		st: &GitHubIssuePtr{
		//			Owner:        strAddr("KamikazeZirou"),
		//			Repository:   strAddr("path-mapper"),
		//			Number:       nil,
		//			StrNumberPtr: strAddr("#1"),
		//		},
		//		success: true,
		//	},
		//},
		//{
		//	name: "Custom Mapper returns error",
		//	args: args{
		//		pattern: "/{owner}/{repository}/issues/{strNumber}",
		//		path:    "/KamikazeZirou/path-mapper/issues/abc",
		//		st:      &GitHubIssue{},
		//	},
		//	want: want{
		//		success: false,
		//	},
		//},
		{
			name: "Pattern and field types do not match",
			args: args{
				pattern: "/{owner}/{repository}/issues/{number}",
				path:    "/KamikazeZirou/path-mapper/issues/abc",
				st:      &GitHubIssue{},
			},
			want: want{
				success: false,
			},
		},
		{
			name: "Length of pattern is invalid",
			args: args{
				pattern: "/{owner}/{repository}/issues/{number}",
				path:    "/guest/sandbox/issues",
				st:      &GitHubIssue{},
			},
			want: want{
				success: false,
			},
		},
		{
			name: "The elements of pattern and path do not match.",
			args: args{
				pattern: "/{owner}/{repository}/issues/{number}",
				path:    "/KamikazeZirou/path-mapper/foobar/1",
				st:      &GitHubIssue{},
			},
			want: want{
				success: false,
			},
		},
		{
			name: "dest is nil",
			args: args{
				pattern: "/{owner}/{repository}/issues/{number}",
				path:    "/KamikazeZirou/path-mapper/foobar/1",
				st:      nil,
			},
			want: want{
				success: false,
			},
		},
		{
			name: "dest is value",
			args: args{
				pattern: "/{owner}/{repository}/issues/{number}",
				path:    "/KamikazeZirou/path-mapper/foobar/1",
				st: GitHubIssue{
					Owner:      "KamikazeZirou",
					Repository: "path-mapper",
					Number:     1,
				},
			},
			want: want{
				success: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Mapping(tt.args.pattern, tt.args.path, tt.args.st)
			if err != nil {
				if tt.want.success {
					t.Errorf("Mapping() return (%v), which is not what we want.", err)
				}
				return
			} else {
				if !tt.want.success {
					t.Errorf("Mapping() return %v, which is not what we want.", err)
					return
				}

				if diff := cmp.Diff(tt.want.st, tt.args.st); diff != "" {
					t.Errorf("Mapping() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func strAddr(s string) *string {
	return &s
}

func intAddr(i int) *int {
	return &i
}

func int8Addr(i int8) *int8 {
	return &i
}

func int16Addr(i int16) *int16 {
	return &i
}

func int32Addr(i int32) *int32 {
	return &i
}

func int64Addr(i int64) *int64 {
	return &i
}

func uintAddr(i uint) *uint {
	return &i
}

func uint8Addr(i uint8) *uint8 {
	return &i
}

func uint16Addr(i uint16) *uint16 {
	return &i
}

func uint32Addr(i uint32) *uint32 {
	return &i
}

func uint64Addr(i uint64) *uint64 {
	return &i
}

func weatherAddr(w Weather) *Weather {
	return &w
}
