package symbols_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"github.com/cloakwiss/ntdocs/symbols"
)

func TestHandleRequriementSectionOfFunction(t *testing.T) {
	fd, er := os.Open("../test/nf-aclapi-treeresetnamedsecurityinfow")
	if er != nil {
		t.Fatal("Cannot open the file")
	}
	defer fd.Close()
	bufFile := bufio.NewReader(fd)

	sections := symbols.GetAllSection(symbols.GetMainContentAsList(bufFile))
	goquery.OuterHtml(sections["requirements"][0])
	table, er := symbols.HandleRequriementSectionOfFunction(sections["requirements"])
	if er != nil {
		t.Fatalf("%s", er.Error())
	}
	mar, er := json.MarshalIndent(table, "", "  ")
	if er != nil {
		t.Fatal("Marshalling failed")
	}
	fmt.Println(string(mar))
}
