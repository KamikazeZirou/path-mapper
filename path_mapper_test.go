package path_mapper

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

type GitHubIssue struct {
	Owner      string
	Repository string
	Number     int
}

func TestMapping(t *testing.T) {
	type expected struct {
		st      GitHubIssue
		success bool
	}

	type args struct {
		pattern string
		path    string
		st      GitHubIssue
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "Basic",
			args: args{
				pattern: "/{owner}/{repository}/issues/{number}",
				path:    "/KamikazeZirou/path-mapper/issues/1",
				st:      GitHubIssue{},
			},
			expected: expected{
				st: GitHubIssue{
					Owner:      "KamikazeZirou",
					Repository: "path-mapper",
					Number:     1,
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
			expected: expected{
				st: GitHubIssue{
					Owner:      "guest",
					Repository: "sandbox",
					Number:     0,
				},
				success: true,
			},
		},
		{
			name: "Length of pattern is invalid",
			args: args{
				pattern: "/{owner}/{repository}/issues/{number}",
				path:    "/guest/sandbox/issues",
				st:      GitHubIssue{},
			},
			expected: expected{
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
			expected: expected{
				success: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Mapping(tt.args.pattern, tt.args.path, &(tt.args.st))
			if err != nil {
				if tt.expected.success {
					t.Errorf("Mapping() return %v, which is not what we expected.", err)
				}
				return
			} else {
				if !tt.expected.success {
					t.Errorf("Mapping() return %v, which is not what we expected.", err)
					return
				}

				if diff := cmp.Diff(tt.args.st, tt.expected.st); diff != "" {
					t.Errorf("Mapping() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
