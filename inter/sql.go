// This file contain all the interractions with sqlite especially the ones which create or insert the into the
// database or will be used in mutliple places
package inter

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB() (*sql.DB, func() error) {
	db, err := sql.Open("sqlite3", "./ntdocs.db")
	if err != nil {
		log.Panicf("Cannot open ntdocs.db : %s\n", err)
	}
	return db, db.Close
}

// Produced after query of table `Headers`
type SymbolRecord struct {
	header, Name, ttype, url string
}

func (sym *SymbolRecord) ScrapableUrl() string {
	return fmt.Sprintf("https://learn.microsoft.com/en-us%s", sym.url)
}

// This function open the db for all the records based on types
func GetSymbolsByGroups(dbConnection *sql.DB) map[string][]SymbolRecord {
	listOfTypesRes, err := dbConnection.Query("SELECT type FROM Symbol GROUP BY type;")
	if err != nil {
		log.Panic("Cannot query Symbols types from ntdocs.db")
	}
	defer listOfTypesRes.Close()

	var (
		typeName, header, name, ttype, url string
		records                            = make(map[string][]SymbolRecord)
	)

	for listOfTypesRes.Next() {
		listOfTypesRes.Scan(&typeName)
		res, err := dbConnection.Query("SELECT * FROM Symbol WHERE type = ?;", typeName)
		if err != nil {
			log.Panic("Cannot query Symbols from ntdocs.db")
		}
		subRecord := make([]SymbolRecord, 0, 12000)
		for res.Next() {
			res.Scan(&header, &name, &ttype, &url)
			record := SymbolRecord{
				header: header,
				Name:   name,
				ttype:  ttype,
				url:    url,
			}
			subRecord = append(subRecord, record)
		}
		records[typeName] = subRecord
	}

	return records
}

func ToBeAddedToRawHTML(dbConnection *sql.DB) []SymbolRecord {
	query := `WITH
		Kernel AS (SELECT name FROM Kernel32Function),
		Scraped AS (SELECT symbolName FROM RawHTML),
		Remaining AS (SELECT name FROM Kernel WHERE name NOT IN Scraped)
	SELECT * FROM Symbol WHERE name IN Remaining;`

	result, err := dbConnection.Query(query)
	log.SetPrefix("RawHTML : ")
	if err != nil {
		log.Panic("Cannot perform the query")
	}
	defer result.Close()

	var (
		header, name, ttype, url string
	)
	records := make([]SymbolRecord, 0, 12000)
	for result.Next() {
		result.Scan(&header, &name, &ttype, &url)
		record := SymbolRecord{
			header: header,
			Name:   name,
			ttype:  ttype,
			url:    url,
		}
		records = append(records, record)
	}
	return records
}

type RawHTMLRecord struct {
	SymbolName string
	HtmlBlob   []byte
}

func AddToRawHTML(conn *sql.DB, rec RawHTMLRecord) {
	stmt, er := conn.Prepare("INSERT INTO RawHTML (symbolName, html) VALUES (?, ?);")
	if er != nil {
		log.Panic("Failed to prepare the insert statement")
	}
	if _, er := stmt.Exec(rec.SymbolName, rec.HtmlBlob); er != nil {
		log.Panic("Insert failed")
	}
}
