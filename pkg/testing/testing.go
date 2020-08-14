package testing

import (
	"os"
	"path"
	"runtime"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

/*
https://brandur.org/fragments/testing-go-project-root

The filename returned by runtime.Caller(0) will be the one currently
executing; testing.go in this case.  Which should reside two directories
deep under the project root (e.g. [project_root]/pkg/testing/testing.go).

We hard code jumping two directories up ("../..") and that gets the
project’s root, and a call to os.Chdir() takes the running program there.

From any other package in the project, import the testing package with a
blank identifier, which will run its init function, but including any of
its symbols.

This works because init() is always executed in the deepest package possible
first, then works it's way back up the tree.

E.g.

```
package module_name

import (
	"testing"

	_ "github.com/jnovack/ipinfo/internal/testing"
)
```

*/
