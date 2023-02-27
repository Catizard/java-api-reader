package reader

import (
	"fmt"
	"testing"
)

func TestRegister(t *testing.T) {
	reader := &Reader{}
	reader.Init(".java")

	exts := []string{"..", "/..", "3", "."}
	for _, v := range exts {
		if err := reader.RegisterExt(v); err != nil {
			fmt.Printf(err.Error())
		}
	}
	registeredExt := reader.InterestingExt()
	fmt.Printf("registered exts = {%v}\n", registeredExt)
	if len(registeredExt) > 1 {
		t.Errorf("some malformed extension escaped")
	}
}

func TestRead(t *testing.T) {
	reader := &Reader{}
	reader.Init(".java", "txt")

	files, err := reader.Read("../test")
	if err != nil {
		t.Error(err)
	}
	if len(files) != 4 {
		t.Error("some files missed")
	}
	fmt.Println(files)
}
