// Contains the function to create WinType struct
package symbols

import (
	// "errors"
	// "iter"
	// "log"
	// "strings"
	//
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
	// "github.com/cloakwiss/ntdocs/utils"
)

type WinType struct {
	name, alias_type /*typedef or define or null(for distilled)*/, alias_to, description string
	is_pointer                                                                           bool
}

func (typ WinType) PrintWinType() {
	fmt.Printf("name: %s\n", typ.name)
	fmt.Printf("alias_type: %s\n", typ.alias_type)
	fmt.Printf("alias_to: %s\n", typ.alias_to)
	fmt.Printf("is_pointer: %t\n", typ.is_pointer)
	fmt.Printf("description: %s\n", typ.description)
}

func ParseWinTypes(htmlTypeRows []*goquery.Selection) []WinType {

	len := len(htmlTypeRows)

	winTypes := make([]WinType, 0, len)
	fmt.Printf("Number of Types: %d\n", len)

	// for _, idx := range []int{0, 1, 4, 5, 17, 54, 69, 90} {
	for idx := 0; idx < len; idx += 1 {

		typ := WinType{
			name:        "null",
			alias_type:  "null",
			alias_to:    "null",
			description: "null",
		}

		// This is for not handling a single structure
		// found in this stuff Look at the comment in structure.go
		if idx == 165 {
			// winTypes = append(winTypes, typ)
			fmt.Printf("\n\n")
			fmt.Printf("At: %d\n", idx)
			typ.PrintWinType()
			continue
		}

		fmt.Printf("\n\n")
		fmt.Printf("At: %d\n", idx)

		// Name ----------------------------------------------------------------------- //
		name, e := htmlTypeRows[idx].Children().First().Find("code").First().Html()
		if e != nil {
			fmt.Println("Couldn't Find Name for windows type")
		}

		typ.name = name
		// ---------------------------------------------------------------------------- //

		// Description ---------------------------------------------------------------- //
		description := htmlTypeRows[idx].Children().Last().Text()
		code := htmlTypeRows[idx].Children().Last().Find("code").Text()

		code = strings.ReplaceAll(code, "far ", "")

		description = strings.ReplaceAll(description, code, "")
		description = strings.Join(strings.Fields(description), " ")

		description = description + "\n" + code

		typ.description = description
		// ---------------------------------------------------------------------------- //

		// Alias Handling ------------------------------------------------------------- //
		if strings.Contains(code, "#if") {
			code = strings.Split(code, "\n")[1]
		}
		code = strings.TrimSpace(code)

		alias_type := "null"
		alias_to := "null"
		if strings.Contains(code, "typedef") {
			alias_type = "typedef"

			code = strings.ReplaceAll(code, "typedef", "")
			code = strings.ReplaceAll(code, ";", "")
			code = strings.ReplaceAll(code, typ.name, "")
			code = strings.TrimSpace(code)

			alias_to = code
		} else if strings.Contains(code, "define") {
			alias_type = "define"

			code = strings.ReplaceAll(code, "#define", "")
			code = strings.ReplaceAll(code, typ.name, "")
			code = strings.TrimSpace(code)

			alias_to = code
		}

		typ.alias_type = alias_type
		typ.alias_to = alias_to
		typ.is_pointer = strings.Contains(typ.alias_to, "*")

		if typ.is_pointer {
			typ.alias_to = strings.TrimSpace(strings.ReplaceAll(typ.alias_to, "*", ""))
		}
		// ---------------------------------------------------------------------------- //

		typ.PrintWinType()
		winTypes = append(winTypes, typ)
	}

	return winTypes
}

func PutWinTypesinDataBase(winTypes []WinType) {
	db, err := sql.Open("sqlite3", "./ntdocs.db")
	if err != nil {
		log.Panicf("Cannot open ntdocs.db : %s\n", err)
	}
	defer db.Close()

	createWinTypesTableQuery := `
	CREATE TABLE IF NOT EXISTS win_type (
		name        TEXT PRIMARY KEY,
		alias_type  TEXT CHECK(alias_type IN ('typedef', 'define')) NULL,
		alias_to    TEXT NULL,
		description TEXT,
		is_pointer  BOOLEAN NOT NULL DEFAULT 0
	);`

	_, creationError := db.Exec(createWinTypesTableQuery)
	if creationError != nil {
		log.Panicf("Failed to create the Table %v\n", err)
	}

	insertQuery, stmtCreationError := db.Prepare(`
		INSERT INTO win_type (name, alias_type, alias_to, description, is_pointer)
		VALUES (?, ?, ?, ?, ?)
	`)
	if stmtCreationError != nil {
		log.Panicf("Preparation for win type insertion is fucked: %v\n", stmtCreationError)
	}

	for _, w := range winTypes {
		var (
			name        string
			alias_type  sql.NullString
			alias_to    sql.NullString
			description string
			is_pointer  bool
		)

		// Fill them ----------------------------------------------------- //
		name = w.name

		if w.alias_type != "" && w.alias_type != "null" {
			alias_type = sql.NullString{String: w.alias_type, Valid: true}
		} else {
			alias_type = sql.NullString{Valid: false}
		}

		if w.alias_to != "" {
			alias_to = sql.NullString{String: w.alias_to, Valid: true}
		} else {
			alias_to = sql.NullString{Valid: false}
		}

		description = w.description
		is_pointer = w.is_pointer
		// --------------------------------------------------------------- //

		_, err := insertQuery.Exec(name, alias_type, alias_to, description, is_pointer)
		if err != nil {
			log.Panicf("Insertion Failed for %v through this %v\n", w, err)
		}
	}
}
