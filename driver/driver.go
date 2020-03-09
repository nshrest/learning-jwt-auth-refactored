package driver

import (
	"database/sql"
	"log"
	"os"

	"github.com/lib/pq"
)

var db *sql.DB

// func ConnectDB connects to elephantsql and returns the db variable
// which is pointer to sql.DB instance
func ConnectDB() *sql.DB {
	//  DB connect (using builtin sql package instead documents says jet package.)
	pgURL, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	if err != nil {
		log.Fatal(err)
	}

	// Test what is inside pgURL
	// fmt.Println(pgUrl)
	// dbname=loiocvro host=raja.db.elephantsql.com password=NX6fuGUBk12YapRmI0un2Sf_TDheGsld port=5432 user=loiocvro
	db, err = sql.Open("postgres", pgURL)
	if err != nil {
		log.Fatal(err)
	}

	// Test what is inside db
	// fmt.Println(db)
	// &{0 {dbname=loiocvro host=raja.db.elephantsql.com password=NX6fuGUBk12YapRmI0un2Sf_TDheGsld port=5432 user=loiocvro 0x9f3b70} 0 {0 0} [] map[] 0 0 0xc0000560c0 0xc00001a240 false map[] map[] 0 0 0 <nil> 0 0 0 0x491ce0}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
