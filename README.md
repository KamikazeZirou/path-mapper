# path-mapper
This is a library that maps resource path into structures

## Usage

```go
import (
	mapper "github.com/KamikazeZirou/path-mapper"
)

type GitHubIssue struct {
	Owner      string
	Repository string
	Number     int
}

func main() {
  st := GitHubIssue{}
  _ = mapper.Mapping("/{owner}/{repository}/issues/{number}", "/KamikazeZirou/path-mapper/issues/1", &st)
  // Owner is set to "KamikazeZirou", Repository to "path-mapper", and Number to 1.
}
```
