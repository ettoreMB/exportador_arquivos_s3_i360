package system

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
)

func ConnectToDb() (*sql.DB, error) {
	query := url.Values{}
	query.Add("database", "csl_2023")
	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword("db", "@Mepm2412"),
		Host:     fmt.Sprintf("%s:%d", "wfbdbserver.database.windows.net", 1433),
		RawQuery: query.Encode(),
	}

	//sqlserver://username:password@host/instance?param1=value&param2
	//sqlserver://db:%40Mepm2412@wfbdbserver.database.windows.net:1433?csl_2023=csl_2023
	db, err := sql.Open("sqlserver", u.String())
	if err != nil {
		err := errors.New("Error creating connection pool: " + err.Error())
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		err := errors.New("Error establishing connection: " + err.Error())
		return nil, err
	}

	fmt.Println("Connected to SQL Server")

	return db, nil
}

func GetDbData(db *sql.DB, filechannel chan []string, view string) {
	query := fmt.Sprintf("SELECT *   FROM %v", view)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("Error executing query: ", err)
	}

	defer rows.Close()

	columns, err := rows.Columns()
	//write columns name in text file
	if err != nil {
		log.Fatal("Error getting columns: " + err.Error())
	}

	filechannel <- columns

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuesPtrs := make([]interface{}, len(columns))

		for i := range values {
			valuesPtrs[i] = &values[i]
		}

		err := rows.Scan(valuesPtrs...)
		if err != nil {
			log.Fatal("Error scanning row: " + err.Error())
		}

		row := make([]string, len(columns))
		for i := range columns {
			var v interface{}
			val := values[i]

			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}

			row[i] = fmt.Sprintf("%v", v)
		}

		filechannel <- row

	}
}
