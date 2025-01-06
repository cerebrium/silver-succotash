package db

import (
	"database/sql"
	"fmt"
)

func WriteLocalDb(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE users (ID INT PRIMARY KEY,Firstname VARCHAR(255),Lastname VARCHAR(255),Clerk_id VARCHAR(255),Email VARCHAR(255));`)
	if err != nil {
		fmt.Println("\n\n\n error creating table: ", err, "\n\n")
	}
}
