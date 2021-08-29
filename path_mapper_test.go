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
	type args struct {
		pattern  string
		path     string
		st       GitHubIssue
		expected GitHubIssue
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Basic",
			args: args{
				pattern: "/{owner}/{repository}/issues/{number}",
				path:    "/KamikazeZirou/path-mapper/issues/1",
				st:      GitHubIssue{},
				expected: GitHubIssue{
					Owner:      "KamikazeZirou",
					Repository: "path-mapper",
					Number:     1,
				},
			},
		},
		{
			name: "Basic2",
			args: args{
				pattern: "/{owner}/{repository}/issues/{number}",
				path:    "/guest/sandbox/issues/2",
				st:      GitHubIssue{},
				expected: GitHubIssue{
					Owner:      "guest",
					Repository: "sandbox",
					Number:     2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Mapping(tt.args.pattern, tt.args.path, &(tt.args.st))

			if diff := cmp.Diff(tt.args.st, tt.args.expected); diff != "" {
				t.Errorf("Mapping() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
