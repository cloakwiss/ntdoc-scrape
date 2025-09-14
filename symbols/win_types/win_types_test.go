package symbols_test

import (
	"bufio"
	// "encoding/json"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/k0kubun/pp/v3"

	"github.com/cloakwiss/ntdocs/symbols"
	"github.com/cloakwiss/ntdocs/utils"
)

func GetSelectionContentAsList(content *goquery.Selection) []*goquery.Selection {
	len := content.Length()

	contentAsList := make([]*goquery.Selection, 0, len)
	content.Each(func(idx int, row *goquery.Selection) {
		contentAsList = append(contentAsList, row)
	})

	return contentAsList
}

func TestParseWinTypes(t *testing.T) {
	fd, er := os.Open("../test/windows-data-types.html")
	if er != nil {
		t.Fatal("Cannot open the file")
	}
	defer fd.Close()
	bufFile := bufio.NewReader(fd)

	sections := utils.GetMainContent(bufFile)

	tableBody := sections.Find("table").First().Find("tbody").First().Children()
	//NOTE: changed the html becuase of change in GetSelectionContentAsList's implementation
	typesInHtmlRows := GetSelectionContentAsList(tableBody)

	// pp.Println(typesInHtmlRows)

	winTypes := symbols.ParseWinTypes(typesInHtmlRows)
	pp.Println(winTypes)
	symbols.PutWinTypesinDataBase(winTypes)
}
