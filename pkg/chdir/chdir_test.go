package chdir

import (
	"testing"

	_ "github.com/jnovack/ipinfo/pkg/testing"
)

func TestEmptyWorkDir(t *testing.T) {
	testDir := ""
	workDir = testDir
	retval := WorkDir()
	expected := execDir + "/"
	if retval != expected {
		t.Errorf("working directory should be equal to executable directory plus trailing slash \nreturned:\n%v\nexpected:\n%v",
			retval, expected)
	}
}

func TestRelativeWorkDirWithoutSlash(t *testing.T) {
	testDir := "assets"
	workDir = testDir
	retval := WorkDir()
	expected := execDir + "/" + testDir + "/"
	if retval != expected {
		t.Errorf("working directory should be equal to executable directory plus trailing slash \nreturned:\n%v\nexpected:\n%v",
			retval, expected)
	}
}

func TestRelativeWorkDirWithSlash(t *testing.T) {
	testDir := "assets/"
	workDir = testDir
	retval := WorkDir()
	expected := execDir + "/" + testDir
	if retval != expected {
		t.Errorf("working directory should be equal to executable directory plus trailing slash \nreturned:\n%v\nexpected:\n%v",
			retval, expected)
	}
}

func TestAbsoluteWorkDirWithoutSlash(t *testing.T) {
	testDir := "/tmp"
	workDir = testDir
	retval := WorkDir()
	expected := testDir + "/"
	if retval != expected {
		t.Errorf("working directory should be equal to executable directory plus trailing slash \nreturned:\n%v\nexpected:\n%v",
			retval, expected)
	}
}

func TestAbsoluteWorkDirWithSlash(t *testing.T) {
	testDir := "/tmp/"
	workDir = testDir
	retval := WorkDir()
	expected := testDir
	if retval != expected {
		t.Errorf("working directory should be equal to executable directory plus trailing slash \nreturned:\n%v\nexpected:\n%v",
			retval, expected)
	}
}
