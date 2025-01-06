package local_user

import (
	"database/sql"
)

// GetUserById retrieves a user from the database by their Clerk_id.
//
// db: The database connection to use.
// id: The Clerk_id of the user to retrieve.
//
// Returns:
//   - A pointer to a User struct containing the fetched user data, or nil if not found.
//     error: Any error encountered during the retrieval process.
func GetUserByClerkId(db *sql.DB, clerk_id string) (*User, error) {
	row := db.QueryRow("SELECT * FROM users WHERE Clerk_Id = ?", clerk_id)

	var user User

	err := row.Scan(&user.ID, &user.Lastname, &user.Firstname, &user.Clerk_Id, &user.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here

			return nil, nil
		}
		return nil, nil
	}

	return &user, nil
}
