package main

import (
	"bufio"
	"bytes"
	"fmt"
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
			content := symbols.GetAllSection(mainContent)
			sig := symbols.HandleFunctionDeclarationSectionOfFunction(content["syntax"])
			if sig.Arity > 0 {
				paras, er := symbols.HandleParameterSectionOfFunction(content["parameters"])
				if len(paras) != int(sig.Arity) {
					log.Println("Parameter parse failed by ", int(sig.Arity)-len(paras), ": ", sig)
					return
				}
				if er == symbols.ErrNewCase || er == symbols.ErrRangingProblem {
					log.Println("Left: ", sig)
					return
				}
				if er != nil {
					log.Panicln(er)
				}
				req, er := symbols.HandleRequriementSectionOfFunction(content["requirements"])
				if er != nil {
					if er == symbols.ErrNotSingleElement {
						log.Println("Left: ", sig)
						return
					} else {
						log.Panicf("Requirements genearation of %+v failed due to: %s\n", sig, er)
					}
				}
				declar := symbols.FunctionDeclarationForInsertion{
					FunctionDeclaration:  sig,
					ParameterDescription: paras,
					Description:          symbols.JoinBlocks(content["basic-description"]),
					Requirements:         req,
				}
				_ = declar
				if er := inter.AddToFunctionSymbol(db, declar); er != nil {
					log.Panicln("Some error in db: ", er)
				}
				fmt.Println(sig)
				// inter.GenerateStatements(declar, buf)
			}
			buf.Flush()
		}()
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
