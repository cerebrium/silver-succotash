package db

import (
	"database/sql"
	"fmt"
	"log"
)

func WriteLocalDb(db *sql.DB) {
	// USERS
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users(ID INT PRIMARY KEY,Firstname VARCHAR(255),Lastname VARCHAR(255),Clerk_id VARCHAR(255),Email VARCHAR(255));`)
	if err != nil {
		fmt.Println("\n\n\n error creating table: ", err, "\n\n")
	}

	// WEIGHTS
	CreateWeights(db)
	PopulateWeights(db)

	// DropStation(db)

	// DNR DPMO
	CreateStation(db)
	PopulateStations(db)
}

func CreateWeights(db *sql.DB) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS weights (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Dcr REAL UNIQUE NOT NULL,
    DnrDpmo REAL UNIQUE NOT NULL,
    Ce REAL UNIQUE NOT NULL,
    Pod REAL UNIQUE NOT NULL,
    Cc REAL UNIQUE NOT NULL,
    Dex REAL UNIQUE NOT NULL);`)
	if err != nil {
		fmt.Println("\n\n\n error creating table: ", err, "\n\n")
	}
}

func PopulateWeights(db *sql.DB) {
	// Check if weights table has any entries
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM weights`).Scan(&count)
	if err != nil {
		log.Fatal("Failed to count weights:", err)
	}

	if count > 0 {
		fmt.Println("there are weights in the database: ", count)
		return
	}

	insertDefaultQuery := `
		INSERT INTO weights (dcr, dnrDpmo, ce, pod, cc, dex) VALUES
		(0.35, 0.35, 0.075, 0.075, 0.075, 0.075);`
	_, err = db.Exec(insertDefaultQuery)
	if err != nil {
		log.Fatal("Failed to insert default values:", err)
	}
}

func DropStation(db *sql.DB) {
	_, err := db.Exec(`drop table station`)
	if err != nil {
		log.Fatal("Error creating the dnr dpmo: ", err)
	}
}

// dnr
func CreateStation(db *sql.DB) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS station (
      ID INTEGER PRIMARY KEY AUTOINCREMENT,
      station TEXT UNIQUE NOT NULL,
      fan INTEGER NOT NULL,
      great INTEGER NOT NULL,
      fair INTEGER NOT NULL);`,
	)
	if err != nil {
		log.Fatal("Error creating the dnr dpmo: ", err)
	}
}

type StationDnr struct {
	station string
	fan     int
	great   int
	fair    int
}

func PopulateStations(db *sql.DB) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM station`).Scan(&count)
	if err != nil {
		log.Fatal("Error counting stations", err)
	}

	if count > 0 {
		fmt.Println("there are stations: ", count)
		return
	}

	stationDnrs := []StationDnr{
		{"DRG2", 1000, 1250, 1550},
		{"DSN1", 900, 1100, 1400},
		{"DBS3", 1150, 1450, 1800},
		{"DBS2", 1550, 1950, 2450},
		{"DEX2", 1050, 1300, 1600},
		{"DCF1", 1250, 1550, 1950},
		{"DSA1", 1100, 1350, 1700},
		{"DPO1", 1200, 1500, 1900},
		{"DOX2", 1100, 1400, 1750},
	}

	for _, station := range stationDnrs {
		_, err := db.Exec(`INSERT INTO station (station, fan, great, fair) VALUES (?, ?, ?, ?)`, station.station, station.fan, station.great, station.fair)
		if err != nil {
			log.Fatal("could not insert into db the station: ", err)
		}
	}
}
