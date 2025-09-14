package inter

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/andybalholm/brotli"
	"github.com/cloakwiss/ntdocs/utils"
	"golang.org/x/net/html"
)

var (
	ErrCompressionFailed   = errors.New("Compression Failed")
	ErrDecompressionFailed = errors.New("Decompression Failed")
)

func GetCompressed(r *bufio.Reader) ([]byte, error) {
	var (
		backingBuffer    = make([]byte, 0, 4<<(10*2))
		buffer           = bytes.NewBuffer(backingBuffer)
		brotliCompressor = brotli.NewWriter(buffer)

		htmlBackingBuffer = make([]byte, 0, 4<<(10*2))
		htmlBuffer        = bytes.NewBuffer(htmlBackingBuffer)
	)
	defer brotliCompressor.Close()
	main := utils.GetMainContent(r)
	for _, node := range main.Nodes {
		html.Render(htmlBuffer, node)
	}
	_, er := htmlBuffer.WriteTo(brotliCompressor)
	if er != nil {
		return nil, ErrCompressionFailed
		//TODO: some logging
	}
	if er := brotliCompressor.Close(); er != nil {
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
	brotliReader := brotli.NewReader(reader)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create brotli reader: %w", err)
	// }
	// defer brotliReader.Close()

	decompressedData, err := io.ReadAll(brotliReader)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress data: %w", err)
	}

	return decompressedData, nil
}
