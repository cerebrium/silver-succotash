package weights

import (
	"database/sql"
	"fmt"
)

func (w Weights) Update(db *sql.DB) error {
	query := `
		UPDATE weights
		SET dcr = ?, DnrDpmo = ?, ce = ?, pod = ?, cc = ?, dex = ?, lor = ?, CdfDpmo = ?, Psb = ?
		WHERE ID = ?;
	`
	result, err := db.Exec(query, w.Dcr, w.DnrDpmo, w.Ce, w.Pod, w.Cc, w.Dex, w.Lor, w.CdfDpmo, w.Psb, w.ID)
	if err != nil {
		return fmt.Errorf("failed to update weights: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated, ID %d not found", w.ID)
	}

	return nil
}
