package main_test

import (
	"bufio"
	"iter"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	main "github.com/cloakwiss/ntdocs"
	"github.com/k0kubun/pp/v3"
)

var pages map[main.SymbolType]string = map[main.SymbolType]string{
	main.Function:    "test/nf-aclapi-buildexplicitaccesswithnamea",
	main.Structure:   "test/ns-accctrl-actrl_access_entry_lista",
	main.Enumeration: "test/ne-accctrl-access_mode",
	main.Callback:    "test/nc-activitycoordinatortypes-activity_coordinator_callback",
	main.Macro:       "test/nf-amsi-amsiresultismalware",
	main.Union:       "test/ns-appmgmt-installspec",
	main.Class:       "test/nl-gdiplusimaging-bitmapdata",
}

func TestFunctionPage(t *testing.T) {
	fd, er := os.Open(pages[main.Function])
	if er != nil {
		t.Fatal("Cannot open the file")
	}
	defer fd.Close()
	bufFile := bufio.NewReader(fd)
	sections := main.GetAllSection(main.GetMainContent(bufFile))
	pp.Println(HandleFunctionDeclaration(sections["syntax"]))
	pp.Println(HandleParameterSectionOfFunction(sections["parameters"]))
	// Need to remove the indexing by moving this check inside
	if found, table := HandleTable(sections["requirements"][0]); found {
		pp.Println(table)
	}
	// debug print
	// codeElem := goquery.Single("code")
	// for _, blk := range sections["requirements"] {
	// 	all, er := goquery.OuterHtml(blk)
	// 	if er == nil {
	// 		fmt.Println("``````", all, "``````")
	// 	}
	// }
}

type FunctionDeclaration struct {
	name, returnType                   string
	arity                              uint8
	usageHint, typeHint, parameterName []string
}

func HandleFunctionDeclaration(block []*goquery.Selection) (functionDeclaration FunctionDeclaration) {
	reverseSplit := func(line string) (usageHint string, typeHint string, parameter string) {
		var (
			idx    int
			marker [4]int
		)

		for ; line[idx] != '['; idx += 1 {
		}
		for ; line[idx] != ']'; idx += 1 {
		}
		idx += 1
		marker[0] = idx
		for ; line[idx] == ' '; idx += 1 {
			marker[1] = idx
		}
		marker[1] += 1
		for ; line[idx] != ' '; idx += 1 {
			marker[2] = idx
		}
		marker[2] += 1
		for ; line[idx] == ' '; idx += 1 {
			marker[3] = idx
		}
		marker[3] += 1

		usageHint, typeHint, parameter = line[:marker[0]], line[marker[1]:marker[2]], line[marker[3]:]
		return
	}

	if len(block) == 1 {
		seq := strings.SplitSeq(block[0].Text(), "\n")
		next, stop := iter.Pull(seq)
		defer stop()
		{

			if firstLine, n := next(); n {
				tokens := strings.Split(strings.Trim(firstLine, " \t("), " ")
				if len(tokens) == 2 {
					functionDeclaration.returnType = tokens[0]
					functionDeclaration.name = tokens[1]
					// pp.Println(returnType, name)
				} else {
					log.Panic("Found something strange in first line of function")
				}
			}
		}
		for {
			if line, found := next(); found {
				if trimmed := strings.TrimLeft(line, " "); trimmed != "" && trimmed != ");" {
					u, t, p := reverseSplit(trimmed)
					functionDeclaration.usageHint = append(functionDeclaration.usageHint, u)
					functionDeclaration.typeHint = append(functionDeclaration.typeHint, t)
					functionDeclaration.parameterName = append(functionDeclaration.parameterName, p)
					functionDeclaration.arity += 1
				}
			} else {
				break
			}
		}
	} else {
		log.Panicln("It have more than one block")
	}
	return
}

func HandleParameterSectionOfFunction(blocks []*goquery.Selection) (output AssociativeArray[string, []string]) {
	codeElem := goquery.Single("code")
	checkParameterHeader := func(blk *goquery.Selection) (string, bool) {
		var (
			code  = blk.FindMatcher(codeElem)
			inner string
			found bool
		)
		switch code.Length() {
		case 0:
			inner, found = "", false
		case 1:
			inner, found = strings.Trim(code.Text(), " "), true
		default:
			log.Panic("Some new case")
		}
		return inner, found
	}

	var (
		start, end, i int
		l             = len(blocks)
	)
	for {
		var parameter string

		for ; i < l; i += 1 {
			rawParameter, found := checkParameterHeader(blocks[i])
			if found {
				parameter = rawParameter
				start = i + 1
				i += 1
				break
			}
		}
		if i >= l {
			break
		}
		for ; i < l; i += 1 {
			_, found := checkParameterHeader(blocks[i])
			if found {
				end = i
				i -= 1
				break
			}
		}
		if i >= l {
			break
		}

		if parameter != "" {
			stringifiedDescription := make([]string, 0, end-start)
			for _, blk := range blocks[start:end] {
				conv, er := goquery.OuterHtml(blk)
				if er != nil {
					log.Fatalln("Cannot convert to html")
				}
				stringifiedDescription = append(stringifiedDescription, conv)
			}
			output.key = append(output.key, parameter)
			output.value = append(output.value, stringifiedDescription)
		} else {
			log.Fatalln("Cannot find paramter ")
		}
	}
	return
}

type AssociativeArray[K, V any] struct {
	key   []K
	value []V
}

func HandleTable(table_block *goquery.Selection) (found bool, output AssociativeArray[string, string]) {
	if !table_block.Is("table") {
		found = false
		return
	}
	body := table_block.Find("tbody")
	if body.Length() == 0 {
		found = false
		return
	}

	children := body.Children()
	found = true
	for i := range children.Length() {
		table_row := children.Eq(i) //.Children()
		table_data := table_row.Find("td")
		if table_data.Length() == 2 {
			key := strings.Trim(table_data.Eq(0).Text(), " \n")
			value := strings.Trim(table_data.Eq(1).Text(), " \n")
			output.key = append(output.key, key)
			output.value = append(output.value, value)
		} else {
			log.Panic("Cannot operate on multiple values.")
		}
	}
	return
}
