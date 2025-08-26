package inter

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/cloakwiss/ntdocs/symbols"
	"golang.org/x/net/html"
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

var (
	ColorOff = "\033[0m"    // Text Reset
	BWhite   = "\033[1;37m" //Bold White
	UWhite   = "\033[4;37m" // White
)

func ReqWorkers(symbols []SymbolRecord, forCompressed chan<- RawHTMLRecord) {
	var (
		logger         = log.New(os.Stdout, "Request Worker ", log.Ltime)
		workersCounter = new(atomic.Int64)
		idx, l         = 0, len(symbols)
	)
	for idx < l {
		if workersCounter.Load() < 3 {
			workersCounter.Add(1)
			go func(name string, url string, i int) {
				defer workersCounter.Add(-1)
				buf := work(logger, url)
				logger.Printf("\tSymbols Left: %s%d%s,\tScraped:  %s%s%s\n", BWhite, l-i-1, ColorOff, UWhite, name, ColorOff)
				forCompressed <- RawHTMLRecord{name, buf}
			}(symbols[idx].Name, symbols[idx].ScrapableUrl(), idx)
			time.Sleep(2 * time.Second)
			idx += 1
		}
	}
	for workersCounter.Load() > 0 {
	}
	close(forCompressed)
}

func work(logger *log.Logger, url string) []byte {
	response, err := httpClient(url)
	if err == nil {
		buf, er := getCompressedSubSection(response)
		if er != nil {
			logger.Printf("ERROR : %s : %s", er.Error(), url)
			return nil
		}
		return buf
	} else {
		if errors.Is(err, ErrHttpGetRequestFailed) {
			logger.Printf("ERROR : %s => %s", err.Error(), url)
		} else if errors.Is(err, ErrHttpResponseReadingFailed) {
			logger.Printf("ERROR : %s => %s", err.Error(), url)
		}
		return nil
	}
}

var (
	ErrCompressionFailed = errors.New("Compression Failed")
)

func getCompressedSubSection(r *bufio.Reader) ([]byte, error) {
	var (
		backingBuffer  = make([]byte, 0, 4<<(10*2))
		buffer         = bytes.NewBuffer(backingBuffer)
		gzipCompressor = gzip.NewWriter(buffer)

		htmlBackingBuffer = make([]byte, 0, 4<<(10*2))
		htmlBuffer        = bytes.NewBuffer(htmlBackingBuffer)
	)
	defer gzipCompressor.Close()
	main := symbols.GetMainContent(r)
	for _, node := range main.Nodes {
		html.Render(htmlBuffer, node)
	}
	_, er := htmlBuffer.WriteTo(gzipCompressor)
	if er != nil {
		return nil, ErrCompressionFailed
		//TODO: some logging
	}
	if er := gzipCompressor.Close(); er != nil {
		log.Panicln("Closing the compressor failed")
	}
	// reseting htmlBuffer to prepare for reuse
	htmlBuffer.Reset()
	encoder := base64.NewEncoder(base64.StdEncoding, htmlBuffer)
	encoder.Write(buffer.Bytes())
	if er := encoder.Close(); er != nil {
		log.Panicln("Closing the writer failed")
	}
	return htmlBuffer.Bytes(), nil
}
