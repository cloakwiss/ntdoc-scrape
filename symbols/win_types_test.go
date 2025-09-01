package symbols_test

import (
	"bufio"
	// "encoding/json"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"github.com/cloakwiss/ntdocs/symbols"
)

func GetSelectionContentAsList(content *goquery.Selection) []*goquery.Selection {
	len := content.Length()

	contentAsList := make([]*goquery.Selection, 0, len)
	content.Each(func(idx int, row *goquery.Selection) {
		contentAsList = append(contentAsList, row)
	});

	return contentAsList
}

func TestParseWinTypes(t *testing.T) {
	fd, er := os.Open("../test/windows-data-types.html")
	if er != nil {
		t.Fatal("Cannot open the file")
	}
	defer fd.Close()
	bufFile := bufio.NewReader(fd)

	sections := symbols.GetMainContent(bufFile)

	tableBody := sections.Find("table").First().Find("tbody").First().Children()
	typesInHtmlRows := GetSelectionContentAsList(tableBody)

	// fmt.Println(typesInHtmlRows)

	winTypes := symbols.ParseWinTypes(typesInHtmlRows)
	symbols.PutWinTypesinDataBase(winTypes)
}
