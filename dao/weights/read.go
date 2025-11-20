package weights

import (
	"database/sql"
	"fmt"
)

func (w *Weights) Read(db *sql.DB) (*Weights, error) {
	row := db.QueryRow("SELECT * FROM weights where ID=?", 1)

	err := row.Scan(&w.ID, &w.Dcr, &w.DnrDpmo, &w.Ce, &w.Pod, &w.Cc, &w.Dex, &w.Lor, &w.CdfDpmo, &w.Psb)
	if err != nil {
		return nil, fmt.Errorf("could not scan weights: %v", err)
	}

	return w, nil
}
