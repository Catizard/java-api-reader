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
	pAnno   string
	pType   string
	comment string
}

type ApiInfo struct {
	retType string
	apiName string
	params  []Param
	usage   string
}

type ControllerApiInfo struct {
	ApiInfo
	Address    string
	MethodType string
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
					log.Printf("[DEBUG] after doParseController, p.controllerApiInfos=%v\n", p.controllerApiInfos)
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
	apiInfo := ApiInfo{}
	apiInfo.Init()
	posL := strings.Index(signature, "(")
	posR := strings.LastIndex(signature, ")")
	if posL == -1 || posR == -1 {
		return apiInfo, fmt.Errorf("extractMethod failed with signature = %v, cannot find any of '(' or ')'\n", signature)
	}
	info := signature[:posL]
	params := signature[posL+1 : posR]
	doExtractMethodInfo(info, &apiInfo)
	doExtractMethodParams(params, &apiInfo)
	return apiInfo, nil
}

func doExtractMethodInfo(info string, apiInfo *ApiInfo) {
	ss := strings.Split(info, " ")
	// [access type] [return type] [method name]
	if len(ss) < 2 {
		log.Fatalf("I think something is wrong with doExtractMethodInfo?")
		return
	}

	// we dont need the public or private info
	if len(ss) > 2 {
		ss = ss[1:]
	}

	apiInfo.retType = ss[0]
	apiInfo.apiName = ss[1]
}

func doExtractMethodParams(params string, apiInfo *ApiInfo) {
	if params == "" {
		return
	}

	ss := strings.Split(params, ",")
	if len(ss) == 0 {
		return
	}

	for _, v := range ss {
		// every v of ss is a param
		param := Param{}
		sv := strings.Split(v, " ")
		if len(sv) < 2 || len(sv) > 3 {
			log.Fatalf("I think something is wrong with doExtractMethodParams")
			return
		}

		if len(sv) == 3 {
			param.pAnno = sv[0]
			sv = sv[1:]
		}

		param.pType = sv[0]
		apiInfo.params = append(apiInfo.params, param)
	}
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
	methodUrl := ""
	methodType := "any"
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

		// TODO 如果发现javadoc，或者其次只是普通的注释，应该怎样记录下来？
		if readyForMethod {
			if !strings.HasPrefix(line, "public") {
				continue
			}
			// TODO extract method info
			//log.Printf("caught method line is {%v}\n", line)
			apiInfo, err := extractMethod(line)
			if err != nil {
				fmt.Printf("extract method failed with %v\n", err)
				return
			}
			controllerApiInfo := &ControllerApiInfo{apiInfo, baseUrl + methodUrl, methodType}
			if _, ok := p.controllerApiInfos[baseUrl]; !ok {
				p.controllerApiInfos[baseUrl] = make([]ControllerApiInfo, 0)
			}
			p.controllerApiInfos[baseUrl] = append(p.controllerApiInfos[baseUrl], *controllerApiInfo)
			//log.Printf("caught a ControllerApiInfo with baseUrl=%v, info=%v\n", baseUrl, controllerApiInfo)
			readyForMethod = false
			continue
		}

		if !afterClass {
			// catch class info
			if strings.HasPrefix(line, "@RequestMapping") {
				baseUrl = strings.Split(line, "\"")[1]
			}

			if strings.HasPrefix(line, "public class") {
				//log.Printf("read class line is {%v}\n", line)
				afterClass = true
			}
		} else {
			// after class, try parse apis
			// every method can be an api
			// TODO 分割检查方法
			if strings.HasPrefix(line, "@RequestMapping") ||
				strings.HasPrefix(line, "@GetMapping") ||
				strings.HasPrefix(line, "@PostMapping") {
				//method type
				pos := strings.Index(line, "Mapping")
				methodType = line[1:pos]
				//method url
				ss := strings.Split(line, "\"")
				if len(ss) != 3 {
					log.Fatalf("split method url failed, len != 3")
				}
				methodUrl = ss[1]

				readyForMethod = true
			}
		}
	}
}

func (p Param) String() string {
	return "pAnno = " + p.pAnno + ", pType = " + p.pType
}
