package inter

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"log"

	"github.com/cloakwiss/ntdocs/symbols"
	"golang.org/x/net/html"
)

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
