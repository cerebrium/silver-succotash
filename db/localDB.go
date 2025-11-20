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
	AddLorColumnIfNotExists(db)
	AddCdfDpmoColumnIfNotExists(db)
	AddPsbColumnIfNotExists(db)
	PopulateWeights(db)
	PopulateLor(db)
	PopulateCdfDpmo(db)
	PopulatePsb(db)

	// DropStation(db)

	// DNR DPMO
	CreateStation(db)
	PopulateStations(db)

	// TIERS
	CreateTiers(db)
}

func CreateTiers(db *sql.DB) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS tiers (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT UNIQUE NOT NULL,
    FanPlus REAL NOT NULL,
    Fan REAL NOT NULL,
    Great REAL NOT NULL,
    Fair REAL NOT NULL
    );`)
	if err != nil {
		fmt.Printf("\n\n\n error creating table tiers: %s", err)
		return
	}

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM tiers`).Scan(&count)
	if err != nil {
		log.Fatal("Failed to count weights:", err)
	}

	if count > 0 {
		fmt.Println("there are tiers in the database: ", count)
		return
	}

	// If not exists, then populate the tiers
	populateTiers(db)
}

func populateTiers(db *sql.DB) {
	columns := []string{"Dcr", "DnrDpmo", "Ce", "Pod", "Cc", "Dex", "Lor", "CdfDpmo", "Psb"}

	query_str := `INSERT INTO tiers (name, FanPlus, Fan, Great, Fair) VALUES (?, ?, ?, ?, ?)`
	for _, name := range columns {
		_, err := db.Exec(query_str, name, 98, 95, 90, 80)
		if err != nil {
			fmt.Printf("Could not create tier: %s \n err: %s", name, err)
		}
	}
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

func AddLorColumnIfNotExists(db *sql.DB) {
	// SQLite doesn't have a direct "ADD COLUMN IF NOT EXISTS" syntax.
	// We need to query the table schema to check if the column exists.

	rows, err := db.Query("PRAGMA table_info(weights)")
	if err != nil {
		fmt.Println("Error querying table info:", err)
		return
	}
	defer rows.Close()

	columnExists := false
	for rows.Next() {
		var cid int
		var name string
		var typeStr string
		var notnull int
		var dfltValue *string
		var pk int

		if err := rows.Scan(&cid, &name, &typeStr, &notnull, &dfltValue, &pk); err != nil {
			fmt.Println("Error scanning table info row:", err)
			return
		}
		if name == "Lor" {
			columnExists = true
			break
		}
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating through table info:", err)
		return
	}

	if !columnExists {
		_, err := db.Exec(
			`ALTER TABLE weights
			ADD COLUMN Lor REAL;`)
		if err != nil {
			fmt.Println("\n\n\n error adding column 'Lor': %s", err)
		} else {
			fmt.Println("Column 'Lor' added successfully.")
			err = PopulateLor(db)
			if err != nil {
				fmt.Printf("error with populate: %s", err)
			}
		}
	} else {
		fmt.Println("Column 'Lor' already exists.")
	}
}

func PopulateLor(db *sql.DB) error {
	insertDefaultQuery := `
		UPDATE weights set Lor=0.1 where ID=1`
	_, err := db.Exec(insertDefaultQuery)
	if err != nil {
		log.Fatal("Failed to insert default values:", err)
		return err
	}
	return nil
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
		INSERT INTO weights (Dcr, DnrDpmo, Ce, Pod, Cc, Dex) VALUES
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

func AddCdfDpmoColumnIfNotExists(db *sql.DB) {
	rows, err := db.Query("PRAGMA table_info(weights)")
	if err != nil {
		fmt.Println("Error querying table info:", err)
		return
	}
	defer rows.Close()

	columnExists := false
	for rows.Next() {
		var cid int
		var name string
		var typeStr string
		var notnull int
		var dfltValue *string
		var pk int

		if err := rows.Scan(&cid, &name, &typeStr, &notnull, &dfltValue, &pk); err != nil {
			fmt.Println("Error scanning table info row:", err)
			return
		}
		if name == "CdfDpmo" {
			columnExists = true
			break
		}
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating through table info:", err)
		return
	}

	if !columnExists {
		_, err := db.Exec(
			`ALTER TABLE weights
			ADD COLUMN CdfDpmo REAL;`)
		if err != nil {
			fmt.Println("\n\n\n error adding column 'CdfDpmo': %s", err)
		} else {
			fmt.Println("Column 'CdfDpmo' added successfully.")
		}
	} else {
		fmt.Println("Column 'CdfDpmo' already exists.")
	}
}

func AddPsbColumnIfNotExists(db *sql.DB) {
	rows, err := db.Query("PRAGMA table_info(weights)")
	if err != nil {
		fmt.Println("Error querying table info:", err)
		return
	}
	defer rows.Close()

	columnExists := false
	for rows.Next() {
		var cid int
		var name string
		var typeStr string
		var notnull int
		var dfltValue *string
		var pk int

		if err := rows.Scan(&cid, &name, &typeStr, &notnull, &dfltValue, &pk); err != nil {
			fmt.Println("Error scanning table info row:", err)
			return
		}
		if name == "Psb" {
			columnExists = true
			break
		}
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating through table info:", err)
		return
	}

	if !columnExists {
		_, err := db.Exec(
			`ALTER TABLE weights
			ADD COLUMN Psb REAL;`)
		if err != nil {
			fmt.Println("\n\n\n error adding column 'Psb': %s", err)
		} else {
			fmt.Println("Column 'Psb' added successfully.")
		}
	} else {
		fmt.Println("Column 'Psb' already exists.")
	}
}

func PopulateCdfDpmo(db *sql.DB) error {
	insertDefaultQuery := `
		UPDATE weights set CdfDpmo=0.0 where ID=1`
	_, err := db.Exec(insertDefaultQuery)
	if err != nil {
		log.Fatal("Failed to insert default CdfDpmo values:", err)
		return err
	}
	return nil
}

func PopulatePsb(db *sql.DB) error {
	insertDefaultQuery := `
		UPDATE weights set Psb=0.0 where ID=1`
	_, err := db.Exec(insertDefaultQuery)
	if err != nil {
		log.Fatal("Failed to insert default Psb values:", err)
		return err
	}
	return nil
}
