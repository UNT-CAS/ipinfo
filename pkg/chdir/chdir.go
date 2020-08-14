package chdir

import (
	"log"
	"os"

	"github.com/kardianos/osext"
	"github.com/namsral/flag"
)

var (
	workDir = ""
	execDir = determineExecutableDirectory()
)

func init() {
	flag.StringVar(&workDir, "workdir", "", "set base path for assets")
}

func determineExecutableDirectory() string {
	// Get the absolute path this executable is located in.
	dir, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal("Error: Couldn't determine working directory: " + err.Error())
	}
	// Set the working directory to the path the executable is located in.
	os.Chdir(dir)
	return dir
}

// WorkDir returns the user-specified path. Empty string if no path was provided.
func WorkDir() string {

	// If no working directory, use executable directory
	if len(workDir) < 1 {
		return execDir + "/"
	}

	// Add trailing slash if not included
	// This modifies the variable permanently
	var last = workDir[len(workDir)-1:]
	if last != "/" {
		workDir = workDir + "/"
	}

	// Return full path of working directory, pad if relative
	var first = workDir[0:1]
	if first != "/" {
		// execDir never has trailing slash
		return execDir + "/" + workDir
	}
	return workDir
}

/*
https://www.kaihag.com/external-assets-working-directories-and-go/

We need the Go binary to be aware of the path it is in.

The way Go works, variables are always determined right after importing the
packages. After that, the init function is called and only after the init
function terminates does the main function execute. In chronological order:

    import -> var -> init() -> main()

That is why we are able to determine the directory our executable is in by
using the imported osext package. We then set the working directory to that
path by calling os.Chdir(executablePath).

*/
