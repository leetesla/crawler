package spider

import (
	"regexp"

	"tesla/crawler/config"
	"tesla/crawler/parser"
)

var (
	profileRe = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`)
	cityUrlRe = regexp.MustCompile(`href="(http://www.zhenai.com/zhenghun/[^"]+)"`)
)

func ParseCity(contents []byte, _ string) parser.ParseResult {
	matches := profileRe.FindAllSubmatch(contents, -1)

	result := parser.ParseResult{}
	for _, m := range matches {
		result.Requests = append(result.Requests, parser.Request{
			Url: string(m[1]),
			Parser: NewProfileParser(string(m[2])),
		})
	}

	matches = cityUrlRe.FindAllSubmatch(contents, -1)

	for _, m := range matches {
		result.Requests = append(result.Requests, parser.Request{
			Url: string(m[1]),
			Parser: parser.NewParser(ParseCity, config.ParseCity),
		})
	}

	return result
}
