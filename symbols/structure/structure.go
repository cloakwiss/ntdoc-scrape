package structure

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/cloakwiss/ntdocs/inter"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
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

func QueryStructure(dbConnection *sql.DB) []inter.SymbolRecord {
	queryRemaingSymbol := `WITH ttypes as (select datatype from FunctionParameters group by FunctionParameters.datatype),
			splits as (SELECT datatype, CASE
				WHEN like('PSS%', datatype) THEN datatype
				WHEN like('PTP%', datatype) THEN datatype
				WHEN like('LP%', datatype) THEN substr(datatype, 3)
				WHEN like('P%', datatype) THEN substr(datatype, 2)
				WHEN like('const %', datatype) THEN substr(datatype, 7)
				ELSE datatype
				END AS n
			FROM ttypes),
			completed as (select symbolName from RawHTML)
		select Symbol.* from Symbol JOIN splits ON splits.n = Symbol.name WHERE Symbol.type IS 'structure' AND Symbol.name NOT IN completed;`

	res, err := dbConnection.Query(queryRemaingSymbol)
	if err != nil {
		log.Panic("Cannot query Symbols types from ntdocs.db")
	}
	defer res.Close()

	var (
		header, name, tokentype, url string
		records                      = make([]inter.SymbolRecord, 0, 200)
	)

	for res.Next() {
		res.Scan(&header, &name, &tokentype, &url)
		record := inter.SymbolRecord{
			Header: header,
			Name:   name,
			Ttype:  tokentype,
			Url:    url,
		}
		records = append(records, record)
	}

	return records
}

type StructDeclaration struct {
	structName string
	names      []string
	fields     []DatatypeNamePair
}

type DatatypeNamePair struct {
	datatype string
	name     string
}

var (
	ErrorSomeNewNode = errors.New("Some new ndoe: ")
)

func getString(node *tree_sitter.Node, code []byte) string {
	return string(code[node.StartByte():node.EndByte()])
}

func HandleSyntaxSection(tree *tree_sitter.Tree, code []byte) (StructDeclaration, error) {
	rootNode := tree.RootNode()
	if rootNode.ChildCount() > 1 {
		log.Panicln("Unexpected case")
	}
	rootNode = rootNode.Child(0)

	var structDecl StructDeclaration

	for _, node := range rootNode.Children(rootNode.Walk()) {
		switch node.Kind() {
		// look for the name and go inside for field
		case "struct_specifier":
			inner := node.ChildByFieldName("body")
			structDecl.structName = getString(node.ChildByFieldName("name"), code)
			if er := handleStruct(inner, &structDecl, code); er != nil {
				return StructDeclaration{}, er
			}

		case "type_identifier", "pointer_declarator":
			structDecl.names = append(structDecl.names, getString(&node, code))

		case "typedef", ";", ",":
		default:
			return StructDeclaration{}, fmt.Errorf("%w : %s", ErrorSomeNewNode, node.Kind())

		}
	}
	return structDecl, nil
}

func handleStruct(node *tree_sitter.Node, structDecl *StructDeclaration, code []byte) error {
	for _, field := range node.Children(node.Walk()) {
		switch field.Kind() {
		case "field_declaration":
			structDecl.fields = append(structDecl.fields, DatatypeNamePair{
				datatype: getString(field.ChildByFieldName("type"), code),
				name:     getString(field.ChildByFieldName("declarator"), code),
			})
		case "{", "}":
		default:
			return fmt.Errorf("%w : %s", ErrorSomeNewNode, node.Kind())
		}
	}
	return nil
}
