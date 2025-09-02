// This file contain all the interractions with sqlite especially the ones which create or insert the into the
// database or will be used in mutliple places
package inter

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/cloakwiss/ntdocs/symbols"
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

func GenerateStatements(declaration symbols.FunctionDeclarationForInsertion, outputBuffer *bufio.Writer) {
	stmt1 := "INSERT OR IGNORE INTO FunctionSymbols (name, arity, return, description) VALUES (%s, %d, %s, %s);\n"
	stmt2 := "INSERT OR IGNORE INTO FunctionParameters (function_name, srno, name, datatype, usage, documentation) VALUES (%s, %d, %s, %s, %s, %s);\n"
	defer outputBuffer.Flush()

	fmt.Fprintf(outputBuffer, stmt1, declaration.Name, declaration.Arity, declaration.ReturnType, declaration.Description)
	for idx, para := range declaration.FunctionDeclaration.Parameters {
		joined := strings.Join(declaration.ParameterDescription[idx].Value, " ")
		fmt.Fprintf(outputBuffer, stmt2, declaration.Name, idx+1, para.Name, para.TypeHint, para.UsageHint, joined)
	}
}

func AddToFunctionSymbol(conn *sql.DB, declaration symbols.FunctionDeclarationForInsertion) error {
	// Use a transaction to reduce lock contention
	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Prepare statements within the transaction
	functionSymbolInsertion, err := tx.Prepare("INSERT OR IGNORE INTO FunctionSymbols (name, arity, return, description) VALUES (?, ?, ?, ?);")
	if err != nil {
		return fmt.Errorf("cannot create functionSymbol insert statement: %w", err)
	}
	defer functionSymbolInsertion.Close()

	functionParameter, err := tx.Prepare("INSERT OR IGNORE INTO FunctionParameters (function_name, srno, name, datatype, usage, documentation) VALUES (?, ?, ?, ?, ?, ?);")
	if err != nil {
		return fmt.Errorf("cannot create functionParameter insert statement: %w", err)
	}
	defer functionParameter.Close()

	// Insert function symbol
	_, err = functionSymbolInsertion.Exec(declaration.Name, declaration.Arity, declaration.ReturnType, declaration.Description)
	if err != nil {
		return fmt.Errorf("cannot insert functionSymbol: %w", err)
	}

	// Validate parameter lengths match
	if len(declaration.FunctionDeclaration.Parameters) != len(declaration.ParameterDescription) {
		return fmt.Errorf("parameter length mismatch: %d parameters vs %d descriptions",
			len(declaration.FunctionDeclaration.Parameters),
			len(declaration.ParameterDescription))
	}

	// Insert parameters
	for idx, para := range declaration.FunctionDeclaration.Parameters {
		joined := strings.Join(declaration.ParameterDescription[idx].Value, " ")
		_, err = functionParameter.Exec(declaration.Name, idx+1, para.Name, para.TypeHint, para.UsageHint, joined)
		if err != nil {
			return fmt.Errorf("cannot insert functionParameter at index %d: %w", idx, err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// func AddToFunctionSymbol(conn *sql.DB, declaration symbols.FunctionDeclarationForInsertion) {
// 	functionSymbolInsertion, er := conn.Prepare("INSERT OR IGNORE INTO FunctionSymbols (name, arity, return, description) VALUES (?, ?, ?, ?);")
// 	if er != nil {
// 		log.Panicln("Cannot create funtionSymbol insert statement: ", er.Error())
// 	}
// 	defer functionSymbolInsertion.Close()

// 	functionParameter, er := conn.Prepare("INSERT OR IGNORE INTO FunctionParameters (function_name, srno, name, datatype, usage, documentation) VALUES (?, ?, ?, ?, ?, ?);")
// 	if er != nil {
// 		log.Panicln("Cannot create functionParameter insert statement: ", er.Error())
// 	}
// 	defer functionParameter.Close()

// 	_, er = functionSymbolInsertion.Exec(declaration.Name, declaration.Arity, declaration.ReturnType, declaration.Description)
// 	if er != nil {
// 		log.Panicln("Cannot insert in functionSymbol :", er.Error())
// 	}

// 	if len(declaration.FunctionDeclaration.Parameters) == len(declaration.ParameterDescription) {
// 		for idx, para := range declaration.FunctionDeclaration.Parameters {
// 			joined := strings.Join(declaration.ParameterDescription[idx].Value, " ")
// 			_, er = functionParameter.Exec(declaration.Name, idx+1, para.Name, para.TypeHint, para.UsageHint, joined)
// 			if er != nil {
// 				log.Panicln("Cannot insert in functionParameter")
// 			}
// 		}
// 	} else {
// 		log.Panicln("This should not be.")
// 	}

// }
