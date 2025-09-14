package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/cloakwiss/ntdocs/inter"
	"github.com/cloakwiss/ntdocs/symbols/function"
	"github.com/cloakwiss/ntdocs/utils"
	_ "github.com/mattn/go-sqlite3"
)

type Command uint8

const (
	SCRAPE_Structure Command = iota + 1
	FILL_FunctionRecord
	FILL_StructureRecord
)

var usageHint = []struct{ name, description string }{
	{"scrape-structure", "Scrape only the structs which are required"},
	{"fill-function-record", "Read scraped data and fill the Function"},
	{"fill-structure-record", "Read scraped data and fill the Structure Table"},
}

func matchFlag(flag string) (Command, bool) {
	var (
		found bool
		cmd   Command
	)

	for i := range usageHint {
		k := usageHint[i].name
		if strings.Compare(flag, k) == 0 {
			found = true
			cmd = Command(uint(i + 1))
			break
		}
	}
	return cmd, found
}

func usage(out io.Writer) {
	fmt.Fprintln(out, "All available flags in the command")
	for i := range usageHint {
		k, v := usageHint[i].name, usageHint[i].description
		fmt.Fprintf(out, "\t--%s\t\t%s\n", k, v)
	}
}

func main() {
	stdout := os.Stdout
	defer stdout.Close()

	out := bufio.NewWriter(stdout)
	defer out.Flush()

	if len(os.Args) != 2 {
		fmt.Fprintln(out, "Need only 1 flag.")
		usage(out)
		return
	}

	args := os.Args[1:]

	flag, found := strings.CutPrefix(args[1], "--")
	if !found {
		fmt.Fprintf(out, "Wrong flag: '%s'\n", args[1])
		usage(out)
		return
	}
	cmd, found := matchFlag(flag)
	if !found {
		fmt.Fprintf(out, "Wrong flag: '%s'\n", args[1])
		usage(out)
		return
	}

	run(cmd, out)
}

func run(cmd Command, stdout *bufio.Writer) {
	db, closer := inter.OpenDB()
	defer closer()

	log.SetFlags(log.Llongfile)

	switch cmd {
	case SCRAPE_Structure:
		//todo
	case FILL_FunctionRecord:
		fillFunctionRecords(db, stdout)
	case FILL_StructureRecord:
		fillStructureRecords(db, stdout)

	default:
		log.Fatal("Some unknown command found")

	}
}
func fillStructureRecords(db *sql.DB, stdoutbuf *bufio.Writer) {

}

func fillFunctionRecords(db *sql.DB, stdoutbuf *bufio.Writer) {
	_ = stdoutbuf
	resultRows, er := db.Query("SELECT symbolName, html FROM RawHTML;")
	if er != nil {
		log.Panicf("Failed to query RawHTML table: %s\n", er)
	}
	var data, name string

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
			mainContent := utils.GetMainContent(buffer)
			content := utils.GetAllSection(mainContent)
			sig := function.HandleFunctionDeclarationSectionOfFunction(content["syntax"])
			if sig.Arity > 0 {
				paras, er := function.HandleParameterSectionOfFunction(content["parameters"])
				if len(paras) != int(sig.Arity) {
					log.Println("Parameter parse failed by ", int(sig.Arity)-len(paras), ": ", sig)
					return
				}
				if er == function.ErrNewCase || er == function.ErrRangingProblem {
					log.Println("Left: ", sig)
					return
				}
				if er != nil {
					log.Panicln(er)
				}
				req, er := function.HandleRequriementSectionOfFunction(content["requirements"])
				if er != nil {
					if er == function.ErrNotSingleElement {
						log.Println("Left: ", sig)
						return
					} else {
						log.Panicf("Requirements genearation of %+v failed due to: %s\n", sig, er)
					}
				}
				declar := function.FunctionDeclarationForInsertion{
					FunctionDeclaration:  sig,
					ParameterDescription: paras,
					Description:          utils.JoinBlocks(content["basic-description"]),
					Requirements:         req,
				}
				_ = declar
				if er := inter.AddToFunctionSymbol(db, declar); er != nil {
					log.Panicln("Some error in db: ", er)
				}
				fmt.Println(sig)
				// inter.GenerateStatements(declar, buf)
			}
			// stdoutbuf.Flush()
		}()
	}
	if er := resultRows.Close(); er != nil {
		log.Panicln("Cannot close Connection")
	}
}

// var Pages map[SymbolType]string = map[SymbolType]string{
// 	Function:    "test/nf-aclapi-treeresetnamedsecurityinfow",
// 	Structure:   "test/ns-accctrl-actrl_access_entry_lista",
// 	Enumeration: "test/ne-accctrl-access_mode",
// 	Callback:    "test/nc-activitycoordinatortypes-activity_coordinator_callback",
// 	Macro:       "test/nf-amsi-amsiresultismalware",
// 	Union:       "test/ns-appmgmt-installspec",
// 	Class:       "test/nl-gdiplusimaging-bitmapdata",
// }

// type SymbolType int8

// const (
// 	Function SymbolType = iota
// 	Structure
// 	Enumeration
// 	Callback
// 	Macro
// 	Union
// 	Class
// )

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
