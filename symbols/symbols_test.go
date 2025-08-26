package symbols_test

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/cloakwiss/ntdocs/symbols"
)

func TestGetAllSection(t *testing.T) {
	fd, er := os.Open("../test/nf-aclapi-treeresetnamedsecurityinfow")
	if er != nil {
		t.Fatal("Cannot open the file")
	}
	defer fd.Close()
	bufFile := bufio.NewReader(fd)
	sections := symbols.GetAllSection(symbols.GetMainContentAsList(bufFile))
	for k, _ := range sections {
		fmt.Println(k)
	}
	desc, found := sections["basic-description"]
	if found {
		for _, l := range desc {
			fmt.Println(goquery.OuterHtml(l))
		}
	}

}
