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
type FunctionDeclarationForInsertion struct {
	FunctionDeclaration
	Description          string
	ParameterDescription utils.AssociativeArray[string, []string]
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
					parameter := splitParameter(trimmed)
					functionDeclaration.Parameters = append(functionDeclaration.Parameters, parameter)
					functionDeclaration.Arity += 1
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
	if len(blocks) == 0 {
		log.Println("Warning: empty blocks slice")
		return output
	}

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
			log.Println("No more parameters found, ending")
			break
		}

		end = l // Default to end of slice if no next parameter found
		for ; i < l; i += 1 {
			_, found := checkParameterHeader(blocks[i])
			if found {
				end = i
				i -= 1 // Step back so next iteration processes this parameter
				break
			}
		}

		if start < 0 || start > l || end < 0 || end > l || start > end {
			log.Printf("Invalid slice indices: start=%d, end=%d, len=%d", start, end, l)
			continue
		}

		if parameter != "" {
			capacity := end - start
			if capacity < 0 {
				log.Printf("Negative capacity: end=%d, start=%d", end, start)
				continue
			}

			// log.Printf("Processing parameter '%s' from blocks[%d:%d] (capacity: %d)", parameter, start, end, capacity)

			stringifiedDescription := make([]string, 0, capacity)

			parameterBlocks := blocks[start:end]
			for idx, blk := range parameterBlocks {
				if blk == nil {
					log.Printf("Warning: nil block at index %d", start+idx)
					continue
				}

				conv, er := goquery.OuterHtml(blk)
				if er != nil {
					log.Printf("Cannot convert block to html at index %d: %v", start+idx, er)
					continue
				}
				stringifiedDescription = append(stringifiedDescription, conv)
			}

			splits := strings.Fields(parameter) // Use Fields instead of Split to handle multiple spaces
			if len(splits) == 0 {
				log.Printf("Warning: empty parameter string after splitting: '%s'", parameter)
				continue
			}

			paramName := splits[len(splits)-1]
			output = append(output, utils.KV[string, []string]{Key: paramName, Value: stringifiedDescription})
		} else {
			log.Printf("Warning: empty parameter found at position %d", i)
		}

		if i >= l-1 {
			break
		}
	}

	return output
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
			UsageHint: "------",
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
