package tiers

import (
	"database/sql"
	"fmt"
)

func UpdateTiers(db *sql.DB, incomingTiers []Tiers) error {
	queryString := `
    UPDATE tiers 
    SET FanPlus = ?, Fan = ?, Great = ?, Fair = ? 
    WHERE ID = ?;
  `

	for _, tier := range incomingTiers {
		result, err := db.Exec(queryString, tier.FanPlus, tier.Fan, tier.Great, tier.Fair, tier.ID)
		if err != nil {
			return fmt.Errorf("failed to update weights: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get affected rows: %w", err)
		}

		if rowsAffected == 0 {
			return fmt.Errorf("no rows updated, ID %d not found", tier.ID)
		}
	}

	return nil
}
