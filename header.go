package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"

	"github.com/k0kubun/pp/v3"

	"github.com/cloakwiss/ntdocs/utils"
)

// Working with headers' toc.json files
type HeaderRecord struct {
	header, name, symbol_type, url string
}

// Read url from db and fetches it contents
// for now we and not covering symbols of type interface
// but that logic in not encoded in scrapper at the moment
func fetchHeaderSSymbols() {
	parseTOCJson := func(input []byte) []HeaderRecord {
		var full any
		er := json.Unmarshal(input, &full)
		if er != nil {
			log.Fatal(er)
		}
		var record []HeaderRecord = make([]HeaderRecord, 0, 16)
		if unwrapped, ok := utils.Cast[map[string]any](full); ok {
			raw_title, title_present := unwrapped["toc_title"]
			if !title_present {
				log.Fatal("Header file name not present")
			}
			title, casted := utils.Cast[string](raw_title)
			if !casted {
				log.Fatal("Casting of title to string failed")
			}
			node_list, nodes_present := unwrapped["children"]
			if !nodes_present {
				log.Fatal("Children not present")
			}
			if url_list, ok := utils.Cast[[]any](node_list); ok {
				for _, raw_node := range url_list[1:] {
					if casted, ok := utils.Cast[map[string]any](raw_node); ok {
						raw_url, url_found := casted["href"]
						if !url_found {
							if _, contains_subtree := casted["children"]; contains_subtree {
								continue
							}
							log.Fatal("Value not found")
						}
						raw_key, key_found := casted["toc_title"]
						if !key_found {
							log.Fatal("Key not found or not casted")
						}
						key, key_ok := utils.Cast[string](raw_key)
						url, value_ok := utils.Cast[string](raw_url)
						if !key_ok {
							log.Fatal("Key not found or not casted")
						}
						if !value_ok {
							log.Fatal("Value not casted")
						}

						name, symbol_type, found := strings.Cut(strings.Trim(key, " "), " ")
						if !found {
							log.Fatal("Some error occured in spliting the record.")
						}
						record = append(record, HeaderRecord{header: title, name: name, symbol_type: symbol_type, url: url})

					}
				}
			}
		}
		return record
	}
	db, err := sql.Open("sqlite3", "./ntdocs.db")
	if err != nil {
		log.Fatal("Cannot open ntdocs.db")
	}
	defer db.Close()

	res, err := db.Query("SELECT name, json_blob FROM Headers;")
	if err != nil {
		log.Fatal("Cannot query Headers from ntdocs.db")
	}
	defer res.Close()

	var (
		name      string
		json_blob []byte
	)
	stmt, er := db.Prepare("INSERT INTO Symbol VALUES (?, ?, ?, ?)")
	if er != nil {
		log.Fatal("Failed to prepare statement")
	}
	defer stmt.Close()

	for res.Next() {
		res.Scan(&name, &json_blob)
		records := parseTOCJson(json_blob)
		if len(records) > 0 {
			for _, record := range records {
				_, er := stmt.Exec(record.header, record.name, record.symbol_type, record.url)
				if er != nil {
					pp.Printf("Failed due to %v\nOn the write: %v\n", er, record)
					return
				}
			}
		}
	}
}
