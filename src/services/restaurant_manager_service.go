package services

import (
	"context"
	"github.com/palantir/stacktrace"
	"googlemaps.github.io/maps"
	"models"
	"strconv"
	"strings"
	"sync"
)

type RestaurantManagerService struct {
	ConfSrv *ConfService
	Dal     *DbService
}

var restaurantManagerServiceInstance *RestaurantManagerService
var restaurantManagerServiceOnce sync.Once

func NewRestaurantManagerService() *RestaurantManagerService {
	restaurantManagerServiceOnce.Do(func() {
		restaurantManagerServiceInstance = &RestaurantManagerService{
			ConfSrv: NewConfService(),
			Dal:     NewDbService(),
		}
	})
	return restaurantManagerServiceInstance
}

func (rms *RestaurantManagerService) UpdateRestaurant(oldRestName string, rest models.Restaurant) error {
	err := rms.Dal.UpdateRestaurant(oldRestName, rest)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to update the restaurant")
	}
	return nil
}

func (rms *RestaurantManagerService) DeleteRestaurant(name string) error {
	err := rms.Dal.DeleteRestaurant(name)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to delete the restaurant")
	}
	return nil
}

func (rms *RestaurantManagerService) GetAllRestaurants() ([]models.Restaurant, error) {
	res, err := rms.Dal.GetAllRestaurants()
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to get the restaurants")
	}

	res, err = rms.convertLocations(res)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to convert locations")
	}

	return res, nil
}

func (rms *RestaurantManagerService) convertLocations(restaurants []models.Restaurant) ([]models.Restaurant, error) {
	var res []models.Restaurant
	for _, restaurant := range restaurants {
		address, err := rms.getAddressByLocation(restaurant.Location)
		if err != nil {
			return nil, stacktrace.Propagate(err, "Failed to get the address by location %s", restaurant.Location)
		}
		rest := models.Restaurant{Name: restaurant.Name, Type: restaurant.Type,
			Phone: restaurant.Phone, Location: address}
		res = append(res, rest)
	}
	return res, nil
}

//  https://github.com/googlemaps/google-maps-services-go
func (rms *RestaurantManagerService) getAddressByLocation(location string) (string, error) {
	lat, lng, err := rms.getLocParams(location)
	if err != nil {
		return "", stacktrace.Propagate(err, "Failed to parse location parameters")
	}

	r := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{Lat: lat, Lng: lng},
	}

	client, err := maps.NewClient(maps.WithAPIKey(rms.ConfSrv.Config.GoogleMapsAPIKey))
	if err != nil {
		return "", stacktrace.Propagate(err, "Failed to init Google maps client")
	}

	resp, err := client.ReverseGeocode(context.Background(), r)
	if err != nil {
		return "", stacktrace.Propagate(err, "Failed to reverse geocode")
	}

	return resp[0].FormattedAddress, nil
}

func (rms *RestaurantManagerService) getLocParams(location string) (float64, float64, error) {
	locArr := strings.Split(location, "/")
	if len(locArr) != 2 {
		return -1, -1, stacktrace.NewError("Failed to parse restaurant's location")
	}

	var err error

	var lat float64
	if lat, err = strconv.ParseFloat(locArr[0], 64); err != nil {
		return -1, -1, stacktrace.NewError("Failed to parse lat parameter")
	}

	var lng float64
	if lng, err = strconv.ParseFloat(locArr[1], 64); err != nil {
		return -1, -1, stacktrace.NewError("Failed to parse lng parameter")
	}

	return lat, lng, nil
}
