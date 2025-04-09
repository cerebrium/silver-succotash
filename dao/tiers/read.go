package tiers

import (
	"database/sql"
)

func ReadTiers(db *sql.DB) ([]*Tiers, error) {
	rows, err := db.Query("SELECT * FROM tiers")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tiers := []*Tiers{}
	for rows.Next() {
		tier := &Tiers{}
		err := rows.Scan(&tier.ID, &tier.Name, &tier.FanPlus, &tier.Fan, &tier.Great, &tier.Fair)
		if err != nil {
			return nil, err
		}
		tiers = append(tiers, tier)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tiers, nil
}
