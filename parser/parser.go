package parser

import "tesla/crawler/config"


type Request struct {
	Url    string
	Parser Parser
}

type Item struct {
	Url     string
	Type    string
	Id      string
	Payload interface{}
}


type ParseResult struct {
	Requests []Request
	Items    []Item
}


type ParserFunc func(contents []byte, url string) ParseResult

// Parser 接口
type Parser interface {
	Parse(contents []byte, url string) ParseResult
	Serialize() (name string, args interface{})
}

// Nil Parser实现Parser接口
type NilParser struct{}

func (NilParser) Parse(_ []byte, _ string) ParseResult {
	return ParseResult{}
}

func (NilParser) Serialize() (name string, args interface{}) {
	return config.NilParser, nil
}


type FuncParser struct {
	parser ParserFunc
	name   string
}

// Func Parser的指针实现Parser接口（指针接收者）
func (f *FuncParser) Parse(contents []byte, url string) ParseResult {
	return f.parser(contents, url) //Go语言自带隐式解引用
}

func (f *FuncParser) Serialize() (name string, args interface{}) {
	return f.name, nil //Go语言自带隐式解引用
}

// 17.6 工厂方法，创建/组装一个FuncParser指针，其实现了Parser接口
// 因为rpc，加入Serialize()
func NewParser(p ParserFunc, name string) *FuncParser {
	return &FuncParser{
		parser: p,
		name:   name,
	}
}

