package structure

import (
	"database/sql"
	"log"

	"github.com/cloakwiss/ntdocs/inter"
)

// This should be handled in structure stuff, but it was found int WinTypes
// Please add it in the structure table. Manually. Is Fine
//
// UNICODE_STRING
// A Unicode string. This type is declared in Winternl.h as follows: C++
// typedef struct _UNICODE_STRING {
//   USHORT  Length;
//   USHORT  MaximumLength;
//   PWSTR  Buffer;
// } UNICODE_STRING;
// typedef UNICODE_STRING *PUNICODE_STRING;
// typedef const UNICODE_STRING *PCUNICODE_STRING;

// This query find out all the parameter's types which are currently required
// it also contains some basic Transformations like:
//   - cutting Long pointer and pointer prefixes
//   - const modifier
//   - more to come

const query string = `WITH ttypes as (select datatype from FunctionParameters group by FunctionParameters.datatype),
	splits as (SELECT datatype, CASE
		WHEN like('PSS%', datatype) THEN datatype
		WHEN like('PTP%', datatype) THEN datatype
		WHEN like('LP%', datatype) THEN substr(datatype, 3)
		WHEN like('P%', datatype) THEN substr(datatype, 2)
		WHEN like('const %', datatype) THEN substr(datatype, 7)
		ELSE datatype
		END AS n
	FROM ttypes)
select Symbol.* from Symbol JOIN splits ON splits.n = Symbol.name WHERE Symbol.type IS 'structure';`

// This function open the db for all the records based on types
func getSymbolsByType(dbConnection *sql.DB) map[string][]inter.SymbolRecord {
	listOfTypesRes, err := dbConnection.Query(query)
	if err != nil {
		log.Panic("Cannot query Symbols types from ntdocs.db")
	}
	defer listOfTypesRes.Close()

	var (
		typeName, header, name, ttype, url string
		records                            = make(map[string][]inter.SymbolRecord)
	)

	for listOfTypesRes.Next() {
		listOfTypesRes.Scan(&typeName)
		res, err := dbConnection.Query("SELECT * FROM Symbol WHERE type = ?;", typeName)
		if err != nil {
			log.Panic("Cannot query Symbols from ntdocs.db")
		}
		subRecord := make([]inter.SymbolRecord, 0, 12000)
		for res.Next() {
			res.Scan(&header, &name, &ttype, &url)
			record := inter.SymbolRecord{
				Header: header,
				Name:   name,
				Ttype:  ttype,
				Url:    url,
			}
			subRecord = append(subRecord, record)
		}
		records[typeName] = subRecord
	}

	return records
}
