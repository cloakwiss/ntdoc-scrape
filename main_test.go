package main_test

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"math/rand"
	"net/http"
	"os"
	"sync/atomic"
	"testing"
	"time"

	main "github.com/cloakwiss/ntdocs"
	"github.com/k0kubun/pp/v3"
)

var pages map[main.SymbolType]string = map[main.SymbolType]string{
	main.Function:    "test/nf-aclapi-treeresetnamedsecurityinfow",
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
	pp.Println(main.HandleFunctionDeclaration(sections["syntax"]))
	pp.Println(main.HandleParameterSectionOfFunction(sections["parameters"]))
	// Need to remove the indexing by moving this check inside
	if found, table := main.HandleTable(sections["requirements"][0]); found {
		pp.Println(table)
	}
	for k := range maps.Keys(sections) {
		pp.Println(k)
	}
}

// func TestReqPool(t *testing.T) {
// 	_ = t
// 	allSymbols := main.GetSymbolsByGroups()
// 	if funcs, found := allSymbols["function"]; found {
// 		threadActiveCounter := new(atomic.Int64)
// 		l := len(funcs) / 120
// 		idx := 0
// 		for idx < l {
// 			if threadActiveCounter.Load() < 4 {
// 				threadActiveCounter.Add(1)
// 				go runThread(threadActiveCounter, funcs[idx].ScrapableUrl())
// 				time.Sleep(4 * time.Second)
// 				idx += 1
// 			}
// 		}
// 	}
// }

func runThread(counter *atomic.Int64, url string) bool {
	defer counter.Add(-1)
	log.Println(counter.Load(), url)
	time.Sleep(time.Duration(rand.Int63n(14500)) * time.Millisecond)
	return true
}

// func TestClient(t *testing.T) {
// 	httpClient("https://learn.microsoft.com/en-us/windows/win32/api/aclapi/nf-aclapi-treeresetnamedsecurityinfow")
// }

func httpClient(url string) {
	client := new(http.Client)
	resp, er := client.Get(url)
	if er != nil {
		log.Fatalf("Cannot fetch: %s\n", url)
	}
	defer resp.Body.Close()
	buffer := make([]byte, 1024*1024*20)
	l, er := resp.Body.Read(buffer)
	if er != nil {
		log.Fatalln("Cannot read from Response.")
	}
	fmt.Println(string(buffer[:l]))
}
