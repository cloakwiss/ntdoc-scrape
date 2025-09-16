package ntquery_test

import (
	"log"
	"testing"

	"github.com/cloakwiss/ntdocs/inter"
	"github.com/cloakwiss/ntdocs/ntquery"
	"github.com/k0kubun/pp/v3"
)

func TestQuery(t *testing.T) {
	connection, closer := inter.OpenDB()
	defer func() {
		if er := closer(); er != nil {
			log.Fatalln("Failed to close the db connection")
		}
	}()
	search := ntquery.NewSearch(connection, 0)
	pp.Println(search.Get("CreateIoRing"))
	pp.Println(search.Get("CreateIoRing"))
}
