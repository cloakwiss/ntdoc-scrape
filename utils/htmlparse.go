// File contains Generic function which will be shared in parsing of all symbols types
package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/k0kubun/pp/v3"
)

func SelectMainContent(r *bufio.Reader) *bufio.Reader {
	doc, er := goquery.NewDocumentFromReader(r)
	if er != nil {
		log.Fatal("Cannot convert to document")
	}
	content := doc.Find("div.content").Eq(1)
	main, er := goquery.OuterHtml(content)
	if er != nil {

	}
	r.Reset(strings.NewReader(main))
	return r
}

func GetMainContent(r *bufio.Reader) *goquery.Selection {
	doc, er := goquery.NewDocumentFromReader(r)
	if er != nil {
		log.Fatal("Cannot convert to document")
	}
	content := doc.Find("div.content")
	return content.First()
}

// Mark and split contents of each section, this will also add extra desciption in future is not marked by
// any h2 element at the start
func GetAllSection(content *goquery.Selection) map[string][]*goquery.Selection {
	var (
		matcher  = goquery.Single("h2[id]")
		sections = make(map[string][]*goquery.Selection)
	)

	first := content.FindMatcher(matcher)
	{
		blk := make([]*goquery.Selection, 0)
		for _, b := range content.Children().EachIter() {
			if b.IsMatcher(matcher) {
				break
			}
			blk = append(blk, b)
		}
		sections["basic-description"] = blk
	}
	first.Each(func(_ int, s *goquery.Selection) {
		val, found := s.Attr("id")
		if !found {
			log.Panicln("Not found")
		}
		blk := make([]*goquery.Selection, 0)
		for _, si := range s.NextUntilMatcher(matcher).EachIter() {
			blk = append(blk, si)
		}
		sections[val] = blk
	})

	return sections
}

func JoinBlocks(blocks []*goquery.Selection) string {
	out := make([]string, 0, len(blocks))
	for _, block := range blocks {
		htm, er := block.Html()
		if er != nil {
			log.Panicln("Failed to get html")
		}
		out = append(out, htm)
	}
	return strings.Join(out, " ")
}

// Extract key value pairs out of the table
// at the moment made with only requirements section in mind
// TODO: But should also work with tables found in some other parts
func HandleTable(table_block *goquery.Selection) (found bool, output AssociativeArray[string, string]) {
	if !table_block.Is("table") {
		found = false
		return
	}
	body := table_block.Find("tbody")
	if body.Length() == 0 {
		found = false
		return
	}

	children := body.Children()
	found = true
	for i := range children.Length() {
		table_row := children.Eq(i)
		table_data := table_row.Find("td")
		if table_data.Length() == 2 {
			key := strings.Trim(table_data.Eq(0).Text(), " \n")
			value := strings.Trim(table_data.Eq(1).Text(), " \n")
			output = append(output, KV[string, string]{Key: key, Value: value})
		} else {
			log.Panicln("Cannot operate on multiple values.")
		}
	}
	return
}

var (
	ErrNotSingleElement     = errors.New("Expect only 1 element found more than one.")
	ErrRequirementsNotFound = errors.New("Cannot find the requirements table")
)

func HandleRequriementSectionOfFunction(blocks []*goquery.Selection) (out string, err error) {
	arr, er := handleRequriementSectionOfFunction(blocks)
	if er != nil {
		err = er
		return
	}
	backingbuf := make([]byte, 0, 256)
	buf := bytes.NewBuffer(backingbuf)
	buf.WriteRune('[')
	for n, i := range arr {
		k, v := i.Key, i.Value
		buf.WriteString(`{"`)
		json.HTMLEscape(buf, []byte(k))
		buf.WriteString(`": "`)
		json.HTMLEscape(buf, []byte(v))
		buf.WriteString(`"}`)
		if n < len(arr)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteRune(']')
	out = buf.String()
	return
}
func handleRequriementSectionOfFunction(blocks []*goquery.Selection) (table AssociativeArray[string, string], err error) {
	if len(blocks) == 1 {
		rawTable := blocks[0]
		var found bool
		if found, table = HandleTable(rawTable); !found {
			err = ErrRequirementsNotFound
		}
	} else {
		for _, b := range blocks {
			if ht, er := b.Html(); er == nil {
				pp.Println(ht)
			}
		}
		err = ErrNotSingleElement
	}
	return
}
