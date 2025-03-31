package stations_routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"hopdf.com/dao/stations"
	"hopdf.com/helpers"
	"hopdf.com/localware"
)

func UpdateStation(c echo.Context) error {
	cc, ok := c.(*localware.LocalUserClerkDbContext)
	if !ok {
		c.Logger().Error("could not resolve cc")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "could not resolve cc",
		})
	}

	// Create an instance of Station
	var incoming_station stations.Station

	// Bind the request body to the struct
	if err := cc.Bind(&incoming_station); err != nil {
		c.Logger().Errorf("could not read the body: %v", err)
		return cc.JSON(http.StatusBadRequest, map[string]string{"error": "could not read the body"})
	}

	res := cc.Db.QueryRow(`SELECT * FROM station WHERE ID=?`, incoming_station.ID)

	var station stations.Station

	err := res.Scan(&station.ID, &station.Station, &station.Fan, &station.Great, &station.Fair)
	if err != nil {
		c.Logger().Error("could not handle scan for station: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "could not handle scan for station",
		})
	}

	// TODO: update the data

	station.Fair = incoming_station.Fair
	station.Fan = incoming_station.Fan
	station.Great = incoming_station.Great

	err = station.Update(cc.Db)
	if err != nil {
		c.Logger().Errorf("could not update:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "could not update",
		})
	}

	return helpers.Success(c, incoming_station)
}
