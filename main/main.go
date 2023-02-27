package main

import (
	"flag"
	"fmt"
	"github.com/catizard/java-api-reader/reader"
)

func main() {
	// TODO 增加默认值
	var pFlag = flag.String("p", "../test", "the path of directory to read")
	flag.Parse()
	fmt.Println(*pFlag)

	r := &reader.Reader{}
	r.Init(".java")
	files, err := r.Read(*pFlag)
	if err != nil {
		fmt.Printf("read failed with %v\n", err)
	}
	fmt.Printf("read result = {%v}\n", files)
}
