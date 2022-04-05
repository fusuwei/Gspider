package xpath

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"log"
	"regexp"
	"strings"
)

type SelectInterface interface {
	Xpath(strXpath string) (*Select, error)
	ExtractFirst() string
	Extract() string
}

type Select struct {
	doc     *html.Node
	attr    string
	preNode string
	xpath   string
}

type NextSelect struct {
	NextDocs []*Select
}

func Parse(body string) *Select {
	s := &Select{}
	a := strings.NewReader(body)
	doc, err := htmlquery.Parse(a)
	if err != nil {
		log.Fatal(err)
	}
	s.doc = doc
	return s
}

func getAttrAndNode(xPath string) (string, string) {
	spaceRe, _ := regexp.Compile(`/{1,2}`)
	nodes := spaceRe.Split(xPath, -1)
	return nodes[len(nodes)-1], nodes[len(nodes)-2]
}

func isInFuncString(xpath string) string {
	if ok, _ := regexp.MatchString("string\\(.*?\\)", xpath); ok {
		return xpath
	}
	xpath = strings.Replace(strings.Replace(xpath, "string(", "", -1), ")", "", -1)
	return xpath
}

func (s *Select) Xpath(strXpath string) *NextSelect {
	nextSelect := &NextSelect{}
	attr, node := getAttrAndNode(strXpath)
	strXpath = isInFuncString(strXpath)
	docs, err := htmlquery.QueryAll(s.doc, strXpath)
	if err != nil {
		return nextSelect
	}
	for _, doc := range docs {
		s := &Select{doc: doc, attr: attr, preNode: node}
		nextSelect.NextDocs = append(nextSelect.NextDocs, s)
	}
	return nextSelect
}

func (nextS *NextSelect) ExtractFirst() string {
	if len(nextS.NextDocs) < 0 || nextS.NextDocs == nil {
		return ""
	}
	return htmlquery.InnerText(nextS.NextDocs[0].doc)

}

func (nextS *NextSelect) Extract() []string {
	ls := make([]string, len(nextS.NextDocs))
	if nextS.NextDocs != nil && len(nextS.NextDocs) < 0 {
		return ls
	}
	for index, row := range nextS.NextDocs {
		ls[index] = htmlquery.InnerText(row.doc)
	}
	return ls
}
