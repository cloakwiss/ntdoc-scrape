// This file contain all the interractions with sqlite especially the ones which create or insert the into the
// database or will be used in mutliple places
package inter

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/cloakwiss/ntdocs/symbols/function"
	"github.com/cloakwiss/ntdocs/symbols/structure"
	"github.com/k0kubun/pp/v3"
	_ "github.com/mattn/go-sqlite3"
)

func OpenDB() (*sql.DB, func() error) {
	db, err := sql.Open("sqlite3", "./ntdocs.db")
	if err != nil {
		log.Panicf("Cannot open ntdocs.db : %s\n", err)
	}
	return db, db.Close
}

func RunQuery(dbConnection *sql.DB, query string) []SymbolRecord {
	res, err := dbConnection.Query(query)
	if err != nil {
		log.Panic("Cannot query Symbols types from ntdocs.db")
	}
	defer res.Close()

	var (
		header, name, tokentype, url string
		records                      = make([]SymbolRecord, 0, 200)
	)

	for res.Next() {
		res.Scan(&header, &name, &tokentype, &url)
		record := SymbolRecord{
			Header: header,
			Name:   name,
			Ttype:  tokentype,
			Url:    url,
		}
		records = append(records, record)
	}

	return records
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

func AddToStructSymbol(conn *sql.DB, declarations []structure.StructDeclaration, stdoutbuf *bufio.Writer) error {
	// These mirrors table schema
	type (
		structureSymbol struct {
			name                      string
			members                   int
			description, requirements string
		}
		structureMembers struct {
			structure_name string
			srno           int
			datatype, name string
		}
		structurePointer struct {
			pointer_name, structure_name string
		}
	)

	structureSymbolInsertion, er := conn.Prepare("INSERT INTO StructureSymbols(name, member_count, description, requirement) VALUES (?, ?, ?, ?);")
	if er != nil {
		return fmt.Errorf("cannot create StructureSymbol insert statement: %w", er)
	}
	defer structureSymbolInsertion.Close()

	structureMemberInsertion, er := conn.Prepare("INSERT INTO StructureMembers(structure_name, srno, datatype, name) VALUES (?, ?, ?, ?);")
	if er != nil {
		return fmt.Errorf("cannot create StructureMembers insert statement: %w", er)
	}
	defer structureMemberInsertion.Close()

	structurePointerInsertion, er := conn.Prepare("INSERT INTO StructurePointer(pointer_name, structure_name) VALUES (?,?);")
	if er != nil {
		return fmt.Errorf("cannot create StructurePointer insert statement: %w", er)
	}
	defer structurePointerInsertion.Close()

	for _, decl := range declarations {
		pp.Fprintln(stdoutbuf, decl)
		{
			value := structureSymbol{
				name:         decl.Names[0],
				members:      len(decl.Fields),
				description:  "",
				requirements: "",
			}
			_, er := structureSymbolInsertion.Exec(value.name, value.members, value.description, value.requirements)
			if er != nil {
				return fmt.Errorf("Some error in adding structureSymbol: %w", er)
			}
		}
		{
			for i := range decl.Fields {
				value := structureMembers{
					structure_name: decl.Names[0],
					srno:           i + 1,
					datatype:       decl.Fields[i].Datatype,
					name:           decl.Fields[i].Name,
				}
				_, er := structureMemberInsertion.Exec(value.structure_name, value.srno, value.datatype, value.name)
				if er != nil {
					return fmt.Errorf("Some error in adding structureMember: %w", er)
				}
			}
		}
		if len(decl.Names) > 1 {
			for _, n := range decl.Names[1:] {
				value := structurePointer{
					pointer_name:   n,
					structure_name: decl.Names[0],
				}
				_, er := structurePointerInsertion.Exec(value.pointer_name, value.structure_name)
				if er != nil {
					return fmt.Errorf("Some error in adding structurePointer: %w", er)
				}
			}
		}
	}
	return nil
}
