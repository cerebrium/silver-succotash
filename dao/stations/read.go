package stations

import (
	"database/sql"
	"fmt"
)

func (s *Station) Read(db *sql.DB) (*Station, error) {
	row := db.QueryRow("SELECT * FROM stations where LOWER(name) LIKE LOWER(?)", s.Station)

	err := row.Scan(&s.ID, &s.Station, &s.Fan, &s.Great, &s.Fair)
	if err != nil {
		return nil, fmt.Errorf("could not scan stations: %w", err)
	}

	return s, nil
}

func ReadAll(db *sql.DB) ([]Station, error) {
	rows, err := db.Query(`SELECT * FROM stations`)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve all stations: %w", err)
	}

	defer rows.Close()

	var all_stations []Station

	for rows.Next() {
		var s Station
		err := rows.Scan(&s.ID, &s.Station, &s.Fan, &s.Great, &s.Fair)
		if err != nil {
			return nil, fmt.Errorf("could not scan station: %w", err)
		}
		all_stations = append(all_stations, s)
	}

	// Check for errors from iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return all_stations, nil
}
