package weights

import (
	"database/sql"
	"fmt"
)

func (w *Weights) Read(db *sql.DB) (*Weights, error) {
	row := db.QueryRow("SELECT * FROM weights where ID=%s", w.ID)

	err := row.Scan(&w.ID, &w.Dcr, &w.DnrDpmo, &w.Ce, &w.Cc, &w.Dex)
	if err != nil {
		return nil, fmt.Errorf("could not scan weights: %w", err)
	}

	return w, nil
}
