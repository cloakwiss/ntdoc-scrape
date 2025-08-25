package main_test

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
