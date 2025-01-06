package local_user

import (
	"database/sql"
)

func CreateUser(db *sql.DB, usr *User) (*User, error) {
	// We pass in the ID because auto incrementing id
	// is not as ideal as uuid's. To my knowledge there
	// is no way to automatically create them in sqlite.
	_, err := db.Exec("INSERT INTO users (ID, Firstname, Lastname, Clerk_Id, Email) VALUES (?, ?, ?, ?, ?)",
		usr.ID,
		usr.Firstname,
		usr.Lastname,
		usr.Clerk_Id,
		usr.Email)
	if err != nil {
		return nil, err
	}
	new_usr, err := GetUserByClerkId(db, usr.Clerk_Id)
	if err != nil {
		return nil, err
	}

	return new_usr, nil
}
