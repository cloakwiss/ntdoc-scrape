package main_test

import (
	"fmt"
	"testing"
	// "github.com/cloakwiss/ntdocs/inter"
)

func TestMain(t *testing.T) {
	nu := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	fmt.Print(nu[:25])
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
