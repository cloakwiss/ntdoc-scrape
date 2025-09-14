// Contains the function to create FunctionDeclration struct
package function

import (
	"encoding/json"
	"errors"
	"iter"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/k0kubun/pp/v3"

	"github.com/cloakwiss/ntdocs/utils"
)

// This type will be used to match the schema of database
// it is required as I decided to throw away the part which contained header name in the page
// so now we have to query the data from DB
type FunctionDeclarationForInsertion struct {
	FunctionDeclaration
	Description, Requirements string
	ParameterDescription      utils.AssociativeArray[string, []string]
}

// This type will only data available in Function Page
type FunctionDeclaration struct {
	Name, ReturnType string
	Arity            uint8
	Parameters       []Parameter
}

type Parameter struct {
	UsageHint, TypeHint, Name string
}

// This function does not handle function with no parameter well
func HandleFunctionDeclarationSectionOfFunction(block []*goquery.Selection) (functionDeclaration FunctionDeclaration) {
	if len(block) == 1 {
		seq := strings.SplitSeq(block[0].Text(), "\n")
		next, stop := iter.Pull(seq)
		defer stop()
		{

			if firstLine, n := next(); n {
				tokens := strings.Split(strings.Trim(firstLine, " \t("), " ")
				if len(tokens) == 2 {
					functionDeclaration.ReturnType = tokens[0]
					functionDeclaration.Name = tokens[1]
				} else if len(tokens) > 2 {
					functionDeclaration.ReturnType = strings.Join(tokens[:len(tokens)-1], " ")
					functionDeclaration.Name = tokens[len(tokens)-1]
					// pp.Println(returnType, name)
				} else {
					st, er := goquery.OuterHtml(block[0])
					if er != nil {
						log.Panic("Some other error occured while this....")
					} else {
						log.Panicf("Found something strange in first line of function: %s", st)
					}

				}
			}
		}
		for {
			if line, found := next(); found {
				if trimmed := strings.TrimLeft(line, " "); trimmed != "" && trimmed != ");" {
					if !strings.HasPrefix(trimmed, "\t)") {
						parameter := splitParameter(trimmed)
						functionDeclaration.Parameters = append(functionDeclaration.Parameters, parameter)
						functionDeclaration.Arity += 1
					}
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

func HandleRequriementSectionOfFunction(blocks []*goquery.Selection) (out string, err error) {
	arr, er := handleRequriementSectionOfFunction(blocks)
	if er == nil {
		mar, er := json.MarshalIndent(arr, "", "  ")
		if er == nil {
			out, er = string(mar), nil
		} else {
			out, err = "", er
		}
	} else {
		out, err = "", er
	}
	return
}
func handleRequriementSectionOfFunction(blocks []*goquery.Selection) (table utils.AssociativeArray[string, string], err error) {
	if len(blocks) == 1 {
		rawTable := blocks[0]
		var found bool
		if found, table = utils.HandleTable(rawTable); !found {
			err = ErrRequirementsNotFound
		}
	} else {
		for _, b := range blocks {
			if ht, er := b.Html(); er == nil {
				pp.Println(ht)
			}
		}
		err = ErrNotSingleElement
	}
	return
}

var (
	ErrMissing        = errors.New("Something is missing")
	ErrNewCase        = errors.New("Some new case")
	ErrRangingProblem = errors.New("Ranging Problem")
)

func HandleParameterSectionOfFunction(blocks []*goquery.Selection) (output utils.AssociativeArray[string, []string], err error) {
	if len(blocks) == 0 {
		log.Println("Warning: empty blocks slice")
		return
	}

	codeElem := goquery.Single("p > code")
	checkParameterHeader := func(blk *goquery.Selection) (bool, error) {
		var (
			code  = blk.FindMatcher(codeElem)
			found bool
			err   error
		)
		switch code.Length() {
		case 0:
			found = false
		case 1:
			found = true
		default:
			if htm, er := blk.Html(); er == nil {
				pp.Println(htm)
			} else {
				log.Panicln("Comer other error")
			}
			err = ErrNewCase
		}
		return found, err
	}

	var markings = make([]int, 0)

	for i, blk := range blocks {
		found, er := checkParameterHeader(blk)
		if er == nil {
			if found {
				markings = append(markings, i)
			}
		} else {
			err = er
			return
		}
	}
	markings = append(markings, len(blocks))

	if len(markings) > 0 && markings[0] != 0 {
		log.Panicln("Cannot find zero at start")
	}

	var markers = make([][]int, 0)

	if len(markings) > 1 {
		for i := range markings[:len(markings)-1] {
			markers = append(markers, []int{markings[i], markings[i+1]})
		}
	}

	if len(markers) == 0 {
	} else if len(markers) == 1 {
		marker := markers[0]
		header := blocks[marker[0]]
		content := blocks[marker[0]+1:]

		// pp.Println(header.Text())
		var stringified = make([]string, 0, 4)
		for _, elem := range content {
			text, er := elem.Html()
			if er != nil {
				err = er
			}
			// pp.Println(text)
			stringified = append(stringified, text)
		}

		output = append(output, utils.KV[string, []string]{
			Key:   strings.Trim(header.Text(), " "),
			Value: stringified,
		})

	} else {
		for _, marker := range markers {
			// pp.Println(marker)
			header := blocks[marker[0]]
			if marker[0]+1 <= marker[1] && marker[1] <= len(blocks) {
				content := blocks[marker[0]+1 : marker[1]]

				// pp.Println(header.Text())
				var stringified = make([]string, 0, 4)
				for _, elem := range content {
					text, er := elem.Html()
					if er != nil {
						err = er
					}
					// pp.Println(text)
					stringified = append(stringified, text)
				}

				output = append(output, utils.KV[string, []string]{
					Key:   strings.Trim(header.Text(), " "),
					Value: stringified,
				})
			} else {
				err = ErrRangingProblem
			}
		}
	}
	return
}

func splitParameter(line string) Parameter {
	var idx, l int = 0, len(line)

	// usage hint
	for ; idx < l && line[idx] != '['; idx += 1 {
	}
	for ; idx < l && line[idx] != ']'; idx += 1 {
	}
	idx += 1

	if idx >= l {
		sig := strings.SplitAfter(line, " ")

		return Parameter{
			UsageHint: "",
			TypeHint:  strings.Trim(strings.Join(sig[:len(sig)-1], " "), " "),
			Name:      strings.TrimRight(sig[len(sig)-1], ", "),
		}
	} else {
		sig := strings.SplitAfter(line[idx:], " ")

		return Parameter{
			UsageHint: strings.Trim(line[:idx], "[]"),
			TypeHint:  strings.Trim(strings.Join(sig[:len(sig)-1], " "), " "),
			Name:      strings.TrimRight(sig[len(sig)-1], ", "),
		}
	}

}
