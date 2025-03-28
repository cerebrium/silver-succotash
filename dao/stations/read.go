package stations

import (
	"database/sql"
	"fmt"
)

func (s *Station) Read(db *sql.DB) (*Station, error) {
	row := db.QueryRow("SELECT * FROM stations where name ilike %s", w.station)

	err := row.Scan(&s.ID, &s.station, &s.fan, &s.great, &s.fair)
	if err != nil {
		return nil, fmt.Errorf("could not scan stations: %w", err)
	}

	return s, nil
}
