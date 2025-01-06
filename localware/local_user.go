package localware

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"hopdf.com/dao/local_user"
)

type LocalUserClerkDbContext struct {
	ClerkDbContext
	Local_User *local_user.User
}

// This translates the clerk user into the
// user for the apps db. From this point
// forward there should not be a need to
// reference clerk. We will keep the ref
// around for now in case (July 20, 2024)
//
// TODO: Decide if should remove clerk
// user ref
func AddLocalUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc, ok := c.(*ClerkDbContext)

		if !ok {
			err := errors.New("could not validate clerkDbContext")
			cc.Logger().Error(err)
			return c.HTML(http.StatusTeapot, Make404Html())
		}

		current_usr, err := local_user.GetUserByClerkId(cc.Db, cc.Clerk_User.ID)
		if err != nil {

			cc.Logger().Error(err)
			return c.HTML(http.StatusTeapot, Make404Html())
		}

		if current_usr == nil {
			email_to_write := ""
			for _, email := range cc.Clerk_User.EmailAddresses {
				if email.ID == *cc.Clerk_User.PrimaryEmailAddressID {
					email_to_write = email.EmailAddress
				}
			}
			// No user found, create user
			user_to_create := &local_user.User{ID: int(uuid.New().ID()), Clerk_Id: cc.Clerk_User.ID, Email: email_to_write}

			current_usr, err = local_user.CreateUser(cc.Db, user_to_create)
			if err != nil {
				cc.Logger().Error(err)
				// If error here then everything is broken
				return c.HTML(http.StatusUnprocessableEntity, Make404Html())
			}
		}

		new_ctx := &LocalUserClerkDbContext{*cc, current_usr}

		return next(new_ctx)
	}
}
