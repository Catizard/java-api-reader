package reader

import (
	"fmt"
	"github.com/catizard/java-api-reader/parser"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type Reader struct {
	interestingExt map[string]bool
	parser         *parser.Parser
}

func (r *Reader) Init(parser *parser.Parser, exts ...string) {
	r.interestingExt = make(map[string]bool)
	for _, v := range exts {
		if err := r.RegisterExt(v); err != nil {

		}
	}

	r.parser = parser
}

func (r *Reader) RegisterExt(ext string) error {
	if ext == "" {
		return fmt.Errorf("extension cannot be empty\n")
	}
	if ext == "." {
		return fmt.Errorf("extension cannot be one forwarding comma\n")
	}
	if ext[0] == '.' {
		ext = ext[1:]
	}
	for _, v := range ext {
		if (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') {
			// do nothing
		} else {
			return fmt.Errorf("extension cannot contain non-alphabet expect forwarding comma\n")
		}
	}
	if r.interestingExt[ext] {
		return fmt.Errorf("registering a same extension: %v\n", ext)
	}

	r.interestingExt[ext] = true
	fmt.Printf("registered a new extension: %v\n", ext)
	return nil
}

func (r *Reader) InterestingExt() map[string]bool {
	return r.interestingExt
}

func (r *Reader) ContainExt(ext string) bool {
	if ext == "" {
		return false
	}
	if ext[0] == '.' {
		ext = ext[1:]
	}

	return r.interestingExt[ext]
}

func (r *Reader) Read(path string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure from accessing a path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext == "" {
			return nil
		}

		fmt.Printf("ext={%v}\n", ext)
		if r.ContainExt(ext) {
			//TODO read file and pass contents to parser
			file, err := os.Open(path)
			if err != nil {
				fmt.Println("open file failed with ", err)
				return err
			}
			defer file.Close()
			log.Printf("ready to read file %v\n", file.Name())
			r.parser.Parse(file)
			log.Printf("read file end\n")
		}

		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
	}
	return files, err
}
