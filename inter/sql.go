package inter

import (
	"database/sql"
	"fmt"
	"log"
)

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

func AddFunctionSymbol() {

}

// This function will be used to interact with the `FunctionSymbols`
func addRecordToFunctionSymbols() {

}
