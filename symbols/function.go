// Contains the function to create FunctionDeclration struct
package symbols

import (
	"errors"
	"iter"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/cloakwiss/ntdocs/utils"
)

// This type will be used to match the schema of database
// it is required as I decided to throw away the part which contained header name in the page
// so now we have to query the data from DB
type FunctionDeclarationWithHeader struct {
	header string
	FunctionDeclaration
}

// This type will only data available in Function Page
type FunctionDeclaration struct {
	name, returnType string
	arity            uint8
	parameters       []Parameter
}

type Parameter struct {
	usageHint, typeHint, name string
}

func HandleFunctionDeclarationSectionOfFunction(block []*goquery.Selection) (functionDeclaration FunctionDeclaration) {
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
					log.Fatal("Found something strange in first line of function")
				}
			}
		}
		for {
			if line, found := next(); found {
				if trimmed := strings.TrimLeft(line, " "); trimmed != "" && trimmed != ");" {
					parameter := splitParameter(trimmed)
					functionDeclaration.parameters = append(functionDeclaration.parameters, parameter)
					functionDeclaration.arity += 1
				}
			} else {
				break
			}
		}
	} else {
		log.Fatal("It have more than one block")
	}
	return
}

var (
	ErrNotSingleElement     = errors.New("Expect only 1 element found more than one.")
	ErrRequirementsNotFound = errors.New("Cannot find the requirements table")
)

func HandleRequriementSectionOfFunction(blocks []*goquery.Selection) (table utils.AssociativeArray[string, string], err error) {
	if len(blocks) == 1 {
		rawTable := blocks[0]
		var found bool
		if found, table = handleTable(rawTable); !found {
			err = ErrRequirementsNotFound
		}
	} else {
		err = ErrNotSingleElement
	}
	return
}

func HandleParameterSectionOfFunction(blocks []*goquery.Selection) (output utils.AssociativeArray[string, []string]) {
	codeElem := goquery.Single("p > code")
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
			log.Fatal("Some new case")
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
			log.Fatal("Should abort at the moment")
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
					log.Fatal("Cannot convert to html")
				}
				stringifiedDescription = append(stringifiedDescription, conv)
			}
			splits := strings.Split(parameter, " ")
			output = append(output, utils.KV[string, []string]{Key: splits[len(splits)-1], Value: stringifiedDescription})
		} else {
			log.Fatal("Cannot find paramter ")
		}
	}
	return
}

func splitParameter(line string) Parameter {
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

	return Parameter{
		usageHint: line[:marker[0]],
		typeHint:  line[marker[1]:marker[2]],
		name:      line[marker[3]:],
	}
}
