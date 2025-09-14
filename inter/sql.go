// This file contain all the interractions with sqlite especially the ones which create or insert the into the
// database or will be used in mutliple places
package inter

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/cloakwiss/ntdocs/symbols/function"
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
	Header, Name, Ttype, Url string
}

func (sym *SymbolRecord) ScrapableUrl() string {
	return fmt.Sprintf("https://learn.microsoft.com/en-us%s", sym.Url)
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

// Only for debug use, Not really useful
// func generateStatements(declaration symbols.FunctionDeclarationForInsertion, outputBuffer *bufio.Writer) {
// 	stmt1 := "INSERT OR IGNORE INTO FunctionSymbols (name, arity, return, description) VALUES ('%s', %d, '%s', '%s');\n"
// 	stmt2 := "INSERT OR IGNORE INTO FunctionParameters (function_name, srno, name, datatype, usage, documentation) VALUES ('%s', %d, '%s', '%s', '%s', '%s');\n"
// 	defer outputBuffer.Flush()

// 	fmt.Fprintf(outputBuffer, stmt1, declaration.Name, declaration.Arity, declaration.ReturnType, declaration.Description)
// 	for idx, para := range declaration.FunctionDeclaration.Parameters {
// 		joined := strings.Join(declaration.ParameterDescription[idx].Value, " ")
// 		fmt.Fprintf(outputBuffer, stmt2, declaration.Name, idx+1, para.Name, para.TypeHint, para.UsageHint, joined)
// 	}
// }

func AddToFunctionSymbol(conn *sql.DB, declaration function.FunctionDeclarationForInsertion) error {
	// Prepare statements within the transaction
	functionSymbolInsertion, err := conn.Prepare("INSERT OR IGNORE INTO FunctionSymbols (name, arity, return, description, requirements) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		return fmt.Errorf("cannot create functionSymbol insert statement: %w", err)
	}
	defer functionSymbolInsertion.Close()

	functionParameter, err := conn.Prepare("INSERT INTO FunctionParameters (function_name, srno, name, datatype, usage, documentation) VALUES (?, ?, ?, ?, ?, ?);")
	if err != nil {
		return fmt.Errorf("cannot create functionParameter insert statement: %w", err)
	}
	defer functionParameter.Close()

	// Insert function symbol
	_, err = functionSymbolInsertion.Exec(declaration.Name, declaration.Arity, declaration.ReturnType, declaration.Description, declaration.Requirements)
	if err != nil {
		return fmt.Errorf("cannot insert functionSymbol: %w", err)
	}

	// Insert parameters
	for idx, para := range declaration.FunctionDeclaration.Parameters {
		joined := strings.Join(declaration.ParameterDescription[idx].Value, " ")
		_, err = functionParameter.Exec(declaration.Name, idx+1, para.Name, para.TypeHint, para.UsageHint, joined)
		if err != nil {
			return fmt.Errorf("cannot insert functionParameter at index %d: %w", idx, err)
		}
	}

	return nil
}
