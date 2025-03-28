package stations

import (
	"database/sql"
	"fmt"
)

func (s *Station) Update(db *sql.DB) error {
	query := `
		UPDATE weights 
		SET dcr = ?, dnrDpmo = ?, ce = ?, pod = ?, cc = ?, dex = ? 
		WHERE ID = ?;
	`
	result, err := db.Exec(query, s.ID, s.Station, s.Fan, s.Great, s.Fair)
	if err != nil {
		return fmt.Errorf("failed to update weights: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated, ID %d not found", s.ID)
	}

	return nil
}
