package main

import (
	"bufio"
	"bytes"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/cloakwiss/ntdocs/inter"
	"github.com/cloakwiss/ntdocs/symbols"
	_ "github.com/mattn/go-sqlite3"
)

var Pages map[SymbolType]string = map[SymbolType]string{
	Function:    "test/nf-aclapi-treeresetnamedsecurityinfow",
	Structure:   "test/ns-accctrl-actrl_access_entry_lista",
	Enumeration: "test/ne-accctrl-access_mode",
	Callback:    "test/nc-activitycoordinatortypes-activity_coordinator_callback",
	Macro:       "test/nf-amsi-amsiresultismalware",
	Union:       "test/ns-appmgmt-installspec",
	Class:       "test/nl-gdiplusimaging-bitmapdata",
}

func main() {
	// conn, closer := inter.OpenDB()
	// defer closer()
	// records := inter.ToBeAddedToRawHTML(conn)

	// tunnel := make(chan inter.RawHTMLRecord)
	// go inter.ReqWorkers(records, tunnel)

	// for r := range tunnel {
	// 	inter.AddToRawHTML(conn, r)
	// }

	db, closer := inter.OpenDB()
	defer closer()

	log.SetFlags(log.Llongfile)

	resultRows, er := db.Query("SELECT symbolName, html FROM RawHTML;")
	if er != nil {
		log.Panicf("Failed to query RawHTML table: %s\n", er)
	}
	var data, name string

	buf := bufio.NewWriter(os.Stdout)
	for resultRows.Next() {
		resultRows.Scan(&name, &data)
		// fmt.Fprintf(buf, "%s: %s\n", name, data)
		func() {
			decompressed, er := inter.GetDecompressed(data)
			if er != nil {
				log.Panicf("Failed to scan rows: %s\n", er)
			}

			// fmt.Println(string(decompressed))
			backing := bytes.NewBuffer(decompressed)
			buffer := bufio.NewReader(backing)
			allContent, er := goquery.NewDocumentFromReader(buffer)
			if er != nil {
				log.Panicln("Cannot create the document")
			}
			mainContent := allContent.Find("div.content").First()
			content := symbols.GetAllSection(symbols.GetContentAsList(mainContent))
			sig := symbols.HandleFunctionDeclarationSectionOfFunction(content["syntax"])
			// pp.Fprintf(buf, "%+v\n", sig)
			if sig.Arity > 0 {
				declar := symbols.FunctionDeclarationForInsertion{
					FunctionDeclaration:  sig,
					ParameterDescription: symbols.HandleParameterSectionOfFunction(content["parameters"]),
					Description:          symbols.JoinBlocks(content["basic-description"]),
				}
				inter.GenerateStatements(declar, buf)
				// pp.Fprintf(buf, "%+v\n", declar)
			}
		}()
		buf.Flush()
		// if er := inter.AddToFunctionSymbol(db, declar); er != nil {
		// 	log.Fatal("Error occured: ", er.Error())
		// }
		// time.Sleep(400 * time.Millisecond)
	}
	if er := resultRows.Close(); er != nil {
		log.Panicln("Cannot close Connection")
	}
}

type SymbolType int8

const (
	Function SymbolType = iota
	Structure
	Enumeration
	Callback
	Macro
	Union
	Class
)

// type SectionsName int8

// const (
// 	Syntax SectionsName = iota
// 	Parameters
// 	ReturnValue
// 	Constants
// 	Members
// 	Remarks
// 	Requirements
// 	SeeAlso
// )

// // This function will be very central to the process and also it cal be very error prone
// func getSectionList(sym SymbolType) (sections []SectionsName) {
// 	switch sym {
// 	case Function:
// 		sections = []SectionsName{Syntax, Parameters, ReturnValue, Requirements, SeeAlso}
// 	case Structure:
// 		sections = []SectionsName{Syntax, Members, Remarks, Requirements, SeeAlso}
// 	case Enumeration:
// 		sections = []SectionsName{Syntax, Constants, Requirements, SeeAlso}
// 	case Callback:
// 		sections = []SectionsName{Syntax, Parameters, ReturnValue, Requirements, SeeAlso}
// 	case Macro:
// 		sections = []SectionsName{Syntax, Requirements, SeeAlso}
// 	case Union:
// 		sections = []SectionsName{Syntax, Members, Requirements, SeeAlso}
// 	case Class:
// 		sections = []SectionsName{Requirements, SeeAlso}
// 	default:
// 		sections = nil
// 	}
// 	return
// }
