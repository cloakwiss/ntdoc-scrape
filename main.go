package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/k0kubun/pp/v3"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	recs := GetSymbolsByGroups()
	for k, rec := range recs {
		pp.Println(k, len(rec))
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

// Working with headers' toc.json files
type HeaderRecord struct {
	header, name, symbol_type, url string
}

// Read url from db and fetches it contents
// for now we and not covering symbols of type interface
// but that logic in not encoded in scrapper at the moment
func fetchHeaderSSymbols() {

	parseTOCJson := func(input []byte) []HeaderRecord {
		var full any
		er := json.Unmarshal(input, &full)
		if er != nil {
			log.Fatal(er)
		}
		var record []HeaderRecord = make([]HeaderRecord, 0, 16)
		if unwrapped, ok := cast[map[string]any](full); ok {
			raw_title, title_present := unwrapped["toc_title"]
			if !title_present {
				log.Fatal("Header file name not present")
			}
			title, casted := cast[string](raw_title)
			if !casted {
				log.Fatal("Casting of title to string failed")
			}
			node_list, nodes_present := unwrapped["children"]
			if !nodes_present {
				log.Fatal("Children not present")
			}
			if url_list, ok := cast[[]any](node_list); ok {
				for _, raw_node := range url_list[1:] {
					if casted, ok := cast[map[string]any](raw_node); ok {
						raw_url, url_found := casted["href"]
						if !url_found {
							if _, contains_subtree := casted["children"]; contains_subtree {
								continue
							}
							log.Fatal("Value not found")
						}
						raw_key, key_found := casted["toc_title"]
						if !key_found {
							log.Fatal("Key not found or not casted")
						}
						key, key_ok := cast[string](raw_key)
						url, value_ok := cast[string](raw_url)
						if !key_ok {
							log.Fatal("Key not found or not casted")
						}
						if !value_ok {
							log.Fatal("Value not casted")
						}

						name, symbol_type, found := strings.Cut(strings.Trim(key, " "), " ")
						if !found {
							log.Fatal("Some error occured in spliting the record.")
						}
						record = append(record, HeaderRecord{header: title, name: name, symbol_type: symbol_type, url: url})

					}
				}
			}
		}
		return record
	}
	db, err := sql.Open("sqlite3", "./ntdocs.db")
	if err != nil {
		log.Fatalln("Cannot open ntdocs.db")
	}
	defer db.Close()

	res, err := db.Query("SELECT name, json_blob FROM Headers;")
	if err != nil {
		log.Fatalln("Cannot query Headers from ntdocs.db")
	}
	defer res.Close()

	var (
		name      string
		json_blob []byte
	)
	stmt, er := db.Prepare("INSERT INTO Symbol VALUES (?, ?, ?, ?)")
	if er != nil {
		log.Fatalln("Failed to prepare statement")
	}
	defer stmt.Close()

	for res.Next() {
		res.Scan(&name, &json_blob)
		records := parseTOCJson(json_blob)
		if len(records) > 0 {
			for _, record := range records {
				_, er := stmt.Exec(record.header, record.name, record.symbol_type, record.url)
				if er != nil {
					pp.Printf("Failed due to %v\nOn the write: %v\n", er, record)
					return
				}
			}
		}
	}
}

// Produced after query of database
type SymbolRecord struct {
	header, name, ttype, url string
}

func (sym *SymbolRecord) ScrapableUrl() string {
	return fmt.Sprintf("https://learn.microsoft.com/en-us%s", sym.url)
}

// This function open the db for all the records based on types
func GetSymbolsByGroups() map[string][]SymbolRecord {
	db, err := sql.Open("sqlite3", "./ntdocs.db")
	if err != nil {
		log.Fatalln("Cannot open ntdocs.db")
	}
	defer db.Close()

	listOfTypesRes, err := db.Query("SELECT type FROM Symbol GROUP BY type;")
	if err != nil {
		log.Fatalln("Cannot query Symbols types from ntdocs.db")
	}
	defer listOfTypesRes.Close()

	var (
		typeName, header, name, ttype, url string
		records                            = make(map[string][]SymbolRecord)
	)

	for listOfTypesRes.Next() {
		listOfTypesRes.Scan(&typeName)
		res, err := db.Query("SELECT * FROM Symbol WHERE type = ?;", typeName)
		if err != nil {
			log.Fatalln("Cannot query Symbols from ntdocs.db")
		}
		subRecord := make([]SymbolRecord, 0, 12000)
		for res.Next() {
			res.Scan(&header, &name, &ttype, &url)
			record := SymbolRecord{
				header: header,
				name:   name,
				ttype:  ttype,
				url:    url,
			}
			subRecord = append(subRecord, record)
		}
		records[typeName] = subRecord
	}

	return records
}

// Handle Function Page
type FunctionDeclaration struct {
	name, returnType                   string
	arity                              uint8
	usageHint, typeHint, parameterName []string
}

func HandleFunctionDeclaration(block []*goquery.Selection) (functionDeclaration FunctionDeclaration) {
	reverseSplit := func(line string) (usageHint string, typeHint string, parameter string) {
		var (
			idx    int
			marker [4]int
		)

		for ; line[idx] != '['; idx += 1 {
		}
		for ; line[idx] != ']'; idx += 1 {
		}
		idx += 1
		marker[0] = idx
		for ; line[idx] == ' '; idx += 1 {
			marker[1] = idx
		}
		marker[1] += 1
		for ; line[idx] != ' '; idx += 1 {
			marker[2] = idx
		}
		marker[2] += 1
		for ; line[idx] == ' '; idx += 1 {
			marker[3] = idx
		}
		marker[3] += 1

		usageHint, typeHint, parameter = line[:marker[0]], line[marker[1]:marker[2]], line[marker[3]:]
		return
	}

	if len(block) == 1 {
		seq := strings.SplitSeq(block[0].Text(), "\n")
		next, stop := iter.Pull(seq)
		defer stop()
		{

			if firstLine, n := next(); n {
				tokens := strings.Split(strings.Trim(firstLine, " \t("), " ")
				if len(tokens) == 2 {
					functionDeclaration.returnType = tokens[0]
					functionDeclaration.name = tokens[1]
					// pp.Println(returnType, name)
				} else {
					log.Panic("Found something strange in first line of function")
				}
			}
		}
		for {
			if line, found := next(); found {
				if trimmed := strings.TrimLeft(line, " "); trimmed != "" && trimmed != ");" {
					u, t, p := reverseSplit(trimmed)
					functionDeclaration.usageHint = append(functionDeclaration.usageHint, u)
					functionDeclaration.typeHint = append(functionDeclaration.typeHint, t)
					functionDeclaration.parameterName = append(functionDeclaration.parameterName, p)
					functionDeclaration.arity += 1
				}
			} else {
				break
			}
		}
	} else {
		log.Panicln("It have more than one block")
	}
	return
}

func HandleParameterSectionOfFunction(blocks []*goquery.Selection) (output AssociativeArray[string, []string]) {
	codeElem := goquery.Single("p code")
	checkParameterHeader := func(blk *goquery.Selection) (string, bool) {
		var (
			code  = blk.FindMatcher(codeElem)
			inner string
			found bool
		)
		switch code.Length() {
		case 0:
			inner, found = "", false
		case 1:
			inner, found = strings.Trim(code.Text(), " "), true
		default:
			log.Panic("Some new case")
		}
		return inner, found
	}

	var (
		start, end, i int
		l             = len(blocks)
	)
	for {
		var parameter string

		for ; i < l; i += 1 {
			rawParameter, found := checkParameterHeader(blocks[i])
			if found {
				parameter = rawParameter
				start = i + 1
				i += 1
				break
			}
		}
		if i >= l {
			break
		}
		for ; i < l; i += 1 {
			_, found := checkParameterHeader(blocks[i])
			if found {
				end = i
				i -= 1
				break
			}
		}
		if i >= l {
			break
		}

		if parameter != "" {
			stringifiedDescription := make([]string, 0, end-start)
			for _, blk := range blocks[start:end] {
				conv, er := goquery.OuterHtml(blk)
				if er != nil {
					log.Fatalln("Cannot convert to html")
				}
				stringifiedDescription = append(stringifiedDescription, conv)
			}
			output = append(output, KV[string, []string]{parameter, stringifiedDescription})
		} else {
			log.Fatalln("Cannot find paramter ")
		}
	}
	return
}

// ===========================================================================

// Some util function on html files
// ===========================================================================
// Find the main content
func GetMainContent(r io.Reader) []*goquery.Selection {
	doc, er := goquery.NewDocumentFromReader(r)
	if er != nil {
		log.Fatalln("Cannot convert to document")
	}
	content := doc.Find("div.content")
	mainContentRaw := content.Eq(1)
	if mainContentRaw.Nodes == nil {
		log.Fatalln("This doc does not contains the section")
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
			log.Println("This should not have orrcured.")
		}

		if i >= l {
			break
		}
	}
	return sections
}

func HandleTable(table_block *goquery.Selection) (found bool, output AssociativeArray[string, string]) {
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
			output = append(output, KV[string, string]{key, value})
		} else {
			log.Panic("Cannot operate on multiple values.")
		}
	}
	return
}

// ===========================================================================

// Generric Data Structure
type AssociativeArray[K, V any] []KV[K, V]
type KV[K, V any] struct {
	key   K
	value V
}

// Util function for casting
func cast[T any](in any) (T, bool) {
	var zero T
	switch inter := in.(type) {
	case T:
		return inter, true
	default:
		return zero, false
	}
}
