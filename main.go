package main

import (
	"github.com/k0kubun/pp/v3"
	_ "github.com/mattn/go-sqlite3"

	"github.com/cloakwiss/ntdocs/inter"
)

var Pages map[SymbolType]string = map[SymbolType]string{
	Function:    "test/nf-aclapi-treeresetnamedsecurityinfow",
	Structure:   "test/ns-accctrl-actrl_access_entry_lista",
	Enumeration: "test/ne-accctrl-access_mode",
	Callback:    "test/nc-activitycoordinatortypes-activity_coordinator_callback",
	Macro:       "test/nf-amsi-amsiresultismalware",
	Union:       "test/ns-appmgmt-installspec",
	Class:       "test/nl-gdiplusimaging-bitmapdata",
}

func main() {
	recs := inter.GetSymbolsByGroups()
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
