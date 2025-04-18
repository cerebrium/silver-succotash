package stations

import (
	"database/sql"
	"fmt"
)

func (s *Station) Update(db *sql.DB) error {
	query := `
		UPDATE station
		SET fan = ?, great = ?, fair = ?
		WHERE ID = ?;
	`
	result, err := db.Exec(query, s.Fan, s.Great, s.Fair, s.ID)
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
