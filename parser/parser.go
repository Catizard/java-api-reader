package parser

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Param struct {
	pType string
}

type ApiInfo struct {
	apiName string
	params  []Param
	retType string
	usage   string
}

type ControllerApiInfo struct {
	ApiInfo
	Address string
}

type Parser struct {
	controllerApiInfos map[string][]ControllerApiInfo
}

func (p *Parser) Init() {
	p.controllerApiInfos = make(map[string][]ControllerApiInfo)
}

func (a *ApiInfo) Init() {
	a.params = make([]Param, 0)
}

func (c *ControllerApiInfo) Init() {
	c.params = make([]Param, 0)
}

func (p *Parser) Parse(file *os.File) error {
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Printf("read file %v failed with %v\n", file.Name(), err)
			return err
		}

		if len(line) != 0 {
			line = strings.Replace(line, "\r\n", "", -1)
			// log.Printf("[DEBUG] Parse reading line is {%v}\n", line)
			// log.Printf("[DEBUG] test len(line) > 0 && line[0] == '@' is %v", len(line) > 0 && line[0] == '@')
			// 如果当前行内容指明了当前类是我们感兴趣的，则将reader整个也就是剩余内容传给doParse处理函数
			if len(line) > 0 && line[0] == '@' {
				if isControllerClass(line) {
					p.doParseController(reader)
				}
			}
		}

		if err == io.EOF {
			break
		}
	}
	return nil
}

func extractMethod(signature string) (ApiInfo, error) {
	return ApiInfo{}, fmt.Errorf("TODO")
}

func isControllerClass(line string) bool {
	//TODO 增加处理方式
	if line == "@RestController" || line == "@Controller" {
		return true
	}
	return false
}

func (p *Parser) doParseController(reader *bufio.Reader) {
	// create a new ControllerApiInfo

	afterClass := false
	baseUrl := ""
	readyForMethod := false
	for {
		// TODO 如果一行内容被拆分？如何正确的组装
		line, err := reader.ReadString('\n')
		//log.Printf("[TODO] doParseController reading line is {%v}\n", line)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Printf("doParseController read failed", err)
			return
		}

		line = strings.Trim(line, " ")
		line = strings.Replace(line, "\r\n", "", -1)

		if readyForMethod {
			if !strings.HasPrefix(line, "public") {
				continue
			}
			// TODO extract method info
			log.Printf("caught method line is {%v}\n", line)
			apiInfo, err := extractMethod(line)
			if err != nil {
				fmt.Printf("extract method failed with %v\n", err)
				return
			}
			controllerApiInfo := &ControllerApiInfo{apiInfo, "[TODO]"}
			if _, ok := p.controllerApiInfos[baseUrl]; !ok {
				p.controllerApiInfos[baseUrl] = make([]ControllerApiInfo, 0)
			}
			p.controllerApiInfos[baseUrl] = append(p.controllerApiInfos[baseUrl], *controllerApiInfo)
			readyForMethod = false
			continue
		}

		if !afterClass {
			// catch class info
			if strings.HasPrefix(line, "@RequestMapping") {
				baseUrl = strings.Split(line, "\"")[1]
				log.Printf("[DEBUG] catch baseUrl = {%v}\n", baseUrl)
			}

			if strings.HasPrefix(line, "public class") {
				log.Printf("read class line is {%v}\n", line)
				afterClass = true
			}
		} else {
			// after class, try parse apis
			// every method can be an api
			// TODO 分割检查方法
			if strings.HasPrefix(line, "@RequestMapping") ||
				strings.HasPrefix(line, "@GetMapping") ||
				strings.HasPrefix(line, "@PostMapping") {
				log.Printf("[DEBUG] caught a mapping annotaion, line is {%v}\n", line)
				readyForMethod = true
			}
		}
	}
}
