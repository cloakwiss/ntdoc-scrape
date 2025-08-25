package symbols

import (
	"io"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/cloakwiss/ntdocs/utils"
)

// Find the main content
func GetMainContent(r io.Reader) []*goquery.Selection {
	doc, er := goquery.NewDocumentFromReader(r)
	if er != nil {
		log.Fatal("Cannot convert to document")
	}
	content := doc.Find("div.content")
	mainContentRaw := content.Eq(1)
	if mainContentRaw.Nodes == nil {
		log.Fatal("This doc does not contains the section")
	}

	firstNode := mainContentRaw.Children().First()
	lastNode := mainContentRaw.Children().Last()
	len := mainContentRaw.Children().Length()

	mainContent := make([]*goquery.Selection, 0, len)
	currentNode := firstNode
	for !currentNode.IsSelection(lastNode) {
		mainContent = append(mainContent, currentNode)
		currentNode = currentNode.Next()
	}
	mainContent = append(mainContent, lastNode)

	return mainContent
}

func GetAllSection(content []*goquery.Selection) map[string][]*goquery.Selection {
	var (
		start, end int
		i, l       int = 0, len(content)
		matcher        = goquery.Single("h2")
		sections       = make(map[string][]*goquery.Selection)
	)

	for {
		for ; i < l && !content[i].IsMatcher(matcher); i += 1 {
		}
		// if i >= l {
		// 	log.Println("This should not have orrcured.")
		// 	break
		// }
		sectionName, found := content[i].Attr("id")
		start = i + 1

		for i = start + 1; i < l && !content[i].IsMatcher(matcher); i += 1 {
		}
		end = i

		sectionContent := content[start:end]
		if found {
			sections[sectionName] = sectionContent
		} else {
			log.Print("This should not have orrcured.\n")
		}

		if i >= l {
			break
		}
	}
	return sections
}

func HandleTable(table_block *goquery.Selection) (found bool, output utils.AssociativeArray[string, string]) {
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
			log.Fatal("Cannot operate on multiple values.")
		}
	}
	return
}

// ===========================================================================
