// File contains Generic function which will be shared in parsing of all symbols types
package symbols

import (
	"bufio"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/cloakwiss/ntdocs/utils"
)

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
func handleTable(table_block *goquery.Selection) (found bool, output utils.AssociativeArray[string, string]) {
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
		table_row := children.Eq(i) //.Children()
		table_data := table_row.Find("td")
		if table_data.Length() == 2 {
			key := strings.Trim(table_data.Eq(0).Text(), " \n")
			value := strings.Trim(table_data.Eq(1).Text(), " \n")
			output = append(output, utils.KV[string, string]{Key: key, Value: value})
		} else {
			log.Panicln("Cannot operate on multiple values.")
		}
	}
	return
}
