package main_test

import (
	"bufio"
	"bytes"
	"log"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/cloakwiss/ntdocs/inter"
	"github.com/cloakwiss/ntdocs/symbols"
	"github.com/k0kubun/pp/v3"
	// "github.com/cloakwiss/ntdocs/inter"
)

func TestHandleSyntaxSection(t *testing.T) {
	db, closer := inter.OpenDB()
	defer closer()

	resultRows, er := db.Query("SELECT symbolName, html FROM RawHTML LIMIT 2;")
	if er != nil {
		log.Panicf("Failed to query RawHTML table: %s\n", er)
	}
	var (
		data, name string
	)

	for resultRows.Next() {
		resultRows.Scan(&name, &data)
		decompressed, er := inter.GetDecompressed(data)
		if er != nil {
			log.Panicf("Failed to scan rows: %s\n", er)
		}

		// fmt.Println(string(decompressed))
		buffer := bufio.NewReader(bytes.NewBuffer(decompressed))
		allContent, er := goquery.NewDocumentFromReader(buffer)
		if er != nil {
			log.Panicln("Cannot create the document")
		}
		mainContent := allContent.Find("div.content").First()
		content := symbols.GetAllSection(symbols.GetContentAsList(mainContent))
		// fmt.Println(len(content))
		// for k := range maps.Keys(content) {
		// 	pp.Println(k)
		// }
		pp.Println(symbols.HandleFunctionDeclarationSectionOfFunction(content["syntax"]))
	}
	if er := resultRows.Close(); er != nil {
		log.Panicln("Cannot close Connection")
	}
}

// func TestFunctionPage(t *testing.T) {
// 	sections := symbols.GetAllSection(symbols.GetMainContent(bufFile))
// 	pp.Println(symbols.HandleFunctionDeclaration(sections["syntax"]))
// 	pp.Println(symbols.HandleParameterSectionOfFunction(sections["parameters"]))
// 	// Need to remove the indexing by moving this check inside
// 	if found, table := symbols.HandleTable(sections["requirements"][0]); found {
// 		pp.Println(table)
// 	}
// 	for k := range maps.Keys(sections) {
// 		pp.Println(k)
// 	}
// }

// func TestClient(t *testing.T) {
// 	reader, er := inter.HttpClient("https://learn.microsoft.com/en-us/windows/win32/api/aclapi/nf-aclapi-treeresetnamedsecurityinfow")
// 	_ = er

// 	mainContent := symbols.GetMainContent(reader)
// 	// if htm, er := goquery.OuterHtml(mainContent[len(mainContent)-1]); er == nil {
// 	// 	fmt.Println("\n\n\nLast one:\n", htm)
// 	// }
// 	sections := symbols.GetAllSection(mainContent)
// 	pp.Println(symbols.HandleFunctionDeclarationSectionOfFunction(sections["syntax"]))
// 	pp.Println(symbols.HandleParameterSectionOfFunction(sections["parameters"]))
// 	// for k := range maps.Keys(sections) {
// 	// 	pp.Println(k)
// 	// }
// 	// // Need to remove the indexing by moving this check inside
// 	table, err := symbols.HandleRequriementSectionOfFunction(sections["requirements"])
// 	if err != nil {
// 		pp.Println(table)
// 	}
// }
