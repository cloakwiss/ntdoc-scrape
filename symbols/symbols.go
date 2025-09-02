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

func GetContentAsList(content *goquery.Selection) []*goquery.Selection {
	firstNode := content.Children().First()
	lastNode := content.Children().Last()
	len := content.Children().Length()

	contentAsList := make([]*goquery.Selection, 0, len)
	currentNode := firstNode
	for !currentNode.IsSelection(lastNode) {
		// fmt.Println(currentNode.Html())
		contentAsList = append(contentAsList, currentNode)
		currentNode = currentNode.Next()
	}
	contentAsList = append(contentAsList, lastNode)

	return contentAsList
}

// Finds the main content of the documentation page which is in 2nd div container of the page
func GetMainContentAsList(r *bufio.Reader) []*goquery.Selection {
	mainContentRaw := GetMainContent(r)
	if mainContentRaw.Length() == 0 {
		log.Fatal("This doc does not contain div.content section")
	}

	return GetContentAsList(mainContentRaw)
}

// Mark and split contents of each section, this will also add extra desciption in future is not marked by
// any h2 element at the start
func GetAllSection(content []*goquery.Selection) map[string][]*goquery.Selection {
	var (
		first, start, end int
		firstSet          bool
		i, l              int = 0, len(content)
		matcher               = goquery.Single("h2")
		sections              = make(map[string][]*goquery.Selection)
	)

	for {
		for ; i < l && !content[i].IsMatcher(matcher); i += 1 {
		}
		if !firstSet {
			first = i
			firstSet = true
		}
		// if i >= l {
		// 	log.Println("This should not have orrcured.")
		// 	break
		// }
		sectionName, found := content[i].Attr("id")
		start = i + 1

		if i >= l {
			log.Panicln("This should not occur.")
			break
		}

		for i = start + 1; i < l && !content[i].IsMatcher(matcher); i += 1 {
		}
		end = i

		sectionContent := content[start:end]
		if found {
			sections[sectionName] = sectionContent
		} else {
			log.Print("This should not have orrcured.\n")
		}

		if i == l {
			break
		} else if i > l {
			log.Fatal("Should be Unreaachable, as should `i` sholuld not be greater than `l`")
		}
	}
	sections["basic-description"] = content[:first]
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
