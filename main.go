package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/cloakwiss/ntdocs/utils"
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
		if unwrapped, ok := utils.Cast[map[string]any](full); ok {
			raw_title, title_present := unwrapped["toc_title"]
			if !title_present {
				log.Fatal("Header file name not present")
			}
			title, casted := utils.Cast[string](raw_title)
			if !casted {
				log.Fatal("Casting of title to string failed")
			}
			node_list, nodes_present := unwrapped["children"]
			if !nodes_present {
				log.Fatal("Children not present")
			}
			if url_list, ok := utils.Cast[[]any](node_list); ok {
				for _, raw_node := range url_list[1:] {
					if casted, ok := utils.Cast[map[string]any](raw_node); ok {
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
						key, key_ok := utils.Cast[string](raw_key)
						url, value_ok := utils.Cast[string](raw_url)
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
		log.Fatal("Cannot open ntdocs.db")
	}
	defer db.Close()

	res, err := db.Query("SELECT name, json_blob FROM Headers;")
	if err != nil {
		log.Fatal("Cannot query Headers from ntdocs.db")
	}
	defer res.Close()

	var (
		name      string
		json_blob []byte
	)
	stmt, er := db.Prepare("INSERT INTO Symbol VALUES (?, ?, ?, ?)")
	if er != nil {
		log.Fatal("Failed to prepare statement")
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
		log.Fatal("Cannot open ntdocs.db")
	}
	defer db.Close()

	listOfTypesRes, err := db.Query("SELECT type FROM Symbol GROUP BY type;")
	if err != nil {
		log.Fatal("Cannot query Symbols types from ntdocs.db")
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
			log.Fatal("Cannot query Symbols from ntdocs.db")
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
