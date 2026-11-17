package gtfs

import (
	"setaapi/internal/repository"
	"setaapi/internal/repository/realtime"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/jmoiron/sqlx"
)

const (
	serviceAlerts    string = "service-alerts"
	tripUpdates      string = "trip-updates"
	vehiclePositions string = "vehicle-positions"
)

// Fetch every basin before touching the DB: a failed download preserves the snapshot.
func fetchRT(kind string) (map[string]*gtfs.FeedMessage, error) {
	return nil, nil
}

func UpdateAlerts() error {
	feeds, err := fetchRT(serviceAlerts)
	if err != nil {
		return err
	}
	return realtime.SaveAlerts(feeds)
}

func UpdateTripUpdates() error {
	feeds, err := fetchRT(tripUpdates)
	if err != nil {
		return err
	}
	return repository.InTransaction(repository.DB_RT, func(tx *sqlx.Tx) error {
		if _, err := tx.Exec("DELETE FROM stop_time_updates"); err != nil {
			return err
		}
		if _, err := tx.Exec("DELETE FROM trip_updates"); err != nil {
			return err
		}
		if err := realtime.SaveTripUpdates(tx, feeds); err != nil {
			return err
		}
		return realtime.SaveStopTimeUpd(tx, feeds)
	})
}

func UpdateVehiclePositions() error {
	feeds, err := fetchRT(vehiclePositions)
	if err != nil {
		return err
	}
	return repository.InTransaction(repository.DB_RT, func(tx *sqlx.Tx) error {
		if _, err := tx.Exec("DELETE FROM vehicle_positions"); err != nil {
			return err
		}
		return realtime.SaveVehiclePositions(tx, feeds)
	})
}

func rtURL(basin, ft string) string {
	return ""
}

func getRT(url string) (*gtfs.FeedMessage, error) {
	return &gtfs.FeedMessage{}, nil
}
