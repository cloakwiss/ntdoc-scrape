package inter

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/cloakwiss/ntdocs/symbols"
)

var (
	ErrHttpGetRequestFailed      = errors.New("HTTP GET request failed.")
	ErrHttpResponseReadingFailed = errors.New("Cannot read the response for GET request.")
)

func httpClient(url string) (*bufio.Reader, error) {
	client := new(http.Client)
	resp, er := client.Get(url)
	if er != nil {
		return nil, ErrHttpGetRequestFailed
	}
	defer resp.Body.Close()
	buffer, er := io.ReadAll(resp.Body)
	if er != nil {
		return nil, ErrHttpResponseReadingFailed
	}
	// fmt.Println(string(buffer))
	reader := bytes.NewReader(buffer)
	return bufio.NewReader(reader), nil
}

func ReqWorkers(symbols []SymbolRecord) {
	logger := log.New(os.Stdout, "Request Worker ", log.Ltime|log.Lshortfile)
	for _, sym := range symbols {
		workersCounter := new(atomic.Int64)
		idx := 0
		if workersCounter.Load() < 4 {
			workersCounter.Add(1)
			go runThread(logger, workersCounter, sym.ScrapableUrl())
			time.Sleep(4 * time.Second)
			idx += 1
		}
	}
}

func runThread(logger *log.Logger, counter *atomic.Int64, url string) bool {
	defer counter.Add(-1)
	reader, err := httpClient(url)
	if err == nil {
		logger.Printf("INFO : Scrapped %s", url)
		mainSections := symbols.GetAllSection(symbols.GetMainContent(reader))
		_ = mainSections
	} else {
		if errors.Is(err, ErrHttpGetRequestFailed) {
			logger.Printf("ERROR : %s => %s", err.Error(), url)
		} else if errors.Is(err, ErrHttpResponseReadingFailed) {
			logger.Printf("ERROR : %s => %s", err.Error(), url)
		}
		return false
	}
	return true
}
