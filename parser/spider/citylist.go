package spider

import (
	"regexp"

	"tesla/crawler/config"
	"tesla/crawler/parser"
)

const cityListRe = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`

func ParseCityList(contents []byte, _ string) parser.ParseResult {
	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(contents, -1)

	result := parser.ParseResult{}

	for _, m := range matches {
		result.Requests = append(result.Requests, parser.Request{
			Url: string(m[1]),
			Parser: parser.NewParser(ParseCity, config.ParseCity),
		})
	}

	return result
}
