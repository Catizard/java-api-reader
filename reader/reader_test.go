package reader

import (
	"fmt"
	"github.com/catizard/java-api-reader/parser"
	"testing"
)

func TestRegister(t *testing.T) {
	reader := &Reader{}
	parser := &parser.Parser{}
	parser.Init()
	reader.Init(parser, ".java")

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
	parser := &parser.Parser{}
	parser.Init()
	reader.Init(parser, ".java", "txt")

	files, err := reader.Read("../test")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(files)
}
