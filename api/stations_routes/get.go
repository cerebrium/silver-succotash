package stations_routes

import (
	"os"

	"github.com/labstack/echo/v4"
	"hopdf.com/dao/stations"
	"hopdf.com/helpers"
	"hopdf.com/localware"
)

// This handler expects the body of th incoming request
// to have a pdf in it. The pdf will have a table that
// needs to be parsed.
func GetStations(c echo.Context) error {
	cc, ok := c.(*localware.LocalUserClerkDbContext)
	if !ok {
		c.Logger().Error("could not resolve cc")
		os.Exit(1)
	}

	all_stations, err := stations.ReadAll(cc.DbContext.Db)
	if err != nil {
		c.Logger().Errorf("could not read all stations: ", err)
	}

	return helpers.Success(c, all_stations)
}
