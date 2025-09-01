package inter

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/cloakwiss/ntdocs/symbols"
	"golang.org/x/net/html"
)

var (
	ErrCompressionFailed   = errors.New("Compression Failed")
	ErrDecompressionFailed = errors.New("Decompression Failed")
)

func GetCompressed(r *bufio.Reader) ([]byte, error) {
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

func GetDecompressed(data string) ([]byte, error) {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	reader := bytes.NewReader(decodedData)
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	decompressedData, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress data: %w", err)
	}

	return decompressedData, nil
}
