package path_mapper

import (
	"github.com/google/go-cmp/cmp"
	"strconv"
	"testing"
)

type GitHubIssue struct {
	Owner      string
	Repository string
	Number     int
	StrNumber  string
}

type GitHubIssuePtr struct {
	Owner      *string
	Repository *string
	Number     *int
	StrNumber  *string
}

func strAddr(s string) *string {
	return &s
}

func intAddr(i int) *int {
	return &i
}

func TestMapping(t *testing.T) {
	Mapper["strNumber"] = func(v string) (interface{}, error) {
		if _, err := strconv.Atoi(v); err == nil {
			return "#" + v, nil
		} else {
			return nil, err
		}
	}

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
			name: "Basic",
			args: args{
				pattern: "/{owner}/{repository}/issues/{number}",
				path:    "/KamikazeZirou/path-mapper/issues/1",
				st:      GitHubIssue{},
			},
			want: want{
				st: GitHubIssue{
					Owner:      "KamikazeZirou",
					Repository: "path-mapper",
					Number:     1,
				},
				success: true,
			},
		},
		{
			name: "Pointer Field",
			args: args{
				pattern: "/{owner}/{repository}/issues/{number}",
				path:    "/KamikazeZirou/path-mapper/issues/1",
				st:      GitHubIssuePtr{},
			},
			want: want{
				st: GitHubIssuePtr{
					Owner:      strAddr("KamikazeZirou"),
					Repository: strAddr("path-mapper"),
					Number:     intAddr(1),
					StrNumber:  nil,
				},
				success: true,
			},
		},
		{
			name: "There is no field corresponding to the pattern in the structure to be mapped.",
			args: args{
				pattern: "/{owner}/{repository}/actions/runs/{buildNumber}",
				path:    "/guest/sandbox/actions/runs/1",
				st:      GitHubIssue{},
			},
			want: want{
				st: GitHubIssue{
					Owner:      "guest",
					Repository: "sandbox",
					Number:     0,
				},
				success: true,
			},
		},
		{
			name: "Custom Mapper",
			args: args{
				pattern: "/{owner}/{repository}/issues/{strNumber}",
				path:    "/KamikazeZirou/path-mapper/issues/1",
				st:      GitHubIssue{},
			},
			want: want{
				st: GitHubIssue{
					Owner:      "KamikazeZirou",
					Repository: "path-mapper",
					Number:     0,
					StrNumber:  "#1",
				},
				success: true,
			},
		},
		{
			name: "Custom Mapper returns error",
			args: args{
				pattern: "/{owner}/{repository}/issues/{strNumber}",
				path:    "/KamikazeZirou/path-mapper/issues/abc",
				st:      GitHubIssue{},
			},
			want: want{
				success: false,
			},
		},
		{
			name: "Pattern and field types do not match",
			args: args{
				pattern: "/{owner}/{repository}/issues/{number}",
				path:    "/KamikazeZirou/path-mapper/issues/abc",
				st:      GitHubIssue{},
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
				st:      GitHubIssue{},
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
				st:      GitHubIssue{},
			},
			want: want{
				success: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Mapping(tt.args.pattern, tt.args.path, &(tt.args.st))
			if err != nil {
				if tt.want.success {
					t.Errorf("Mapping() return %v, which is not what we want.", err)
				}
				return
			} else {
				if !tt.want.success {
					t.Errorf("Mapping() return %v, which is not what we want.", err)
					return
				}

				if diff := cmp.Diff(tt.args.st, tt.want.st); diff != "" {
					t.Errorf("Mapping() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
