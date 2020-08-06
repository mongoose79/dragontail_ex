package services

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/palantir/stacktrace"
	"log"
	"models"
	"sync"
)

type DbService struct {
	ConfSrv *ConfService
}

var dbServiceInstance *DbService
var dbServiceOnce sync.Once

func NewDbService() *DbService {
	dbServiceOnce.Do(func() {
		dbServiceInstance = &DbService{
			ConfSrv: NewConfService(),
		}

		err := dbServiceInstance.executeSQLQuery(createTables)
		if err != nil {
			log.Fatal(err)
		}
	})
	return dbServiceInstance
}

const createTables = "CREATE TABLE IF NOT EXISTS Restaurants (" +
	"Id SERIAL PRIMARY KEY," +
	"Name VARCHAR(50)," +
	"Type VARCHAR(50)," +
	"Phone VARCHAR(50)," +
	"Location VARCHAR(50)" +
	");"

const updateRestaurantByName = "UPDATE Restaurants SET Name=$1, Type=$2, Phone=$3, Location=$4" +
	"WHERE Name like $5"

const deleteRestaurantByName = "DELETE FROM Restaurants WHERE Name like $1"

const deleteAllRestaurants = "TRUNCATE TABLE Restaurants"

const insertRestaurant = "INSERT INTO Restaurants (Name, Type, Phone, Location)" +
	"VALUES ($1, $2, $3, $4);"

const selectAllRestaurants = "SELECT Name, Type, Phone, Location FROM Restaurants"

func (dbs *DbService) GetAllRestaurants() ([]models.Restaurant, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbs.ConfSrv.Config.Host, dbs.ConfSrv.Config.DbPort, dbs.ConfSrv.Config.User,
		dbs.ConfSrv.Config.Password, dbs.ConfSrv.Config.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to open Postgresql connection")
	}
	defer db.Close()

	rows, err := db.Query(selectAllRestaurants)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to get the restaurants")
	}
	defer rows.Close()

	var restaurants []models.Restaurant
	for rows.Next() {
		var name string
		var rType string
		var phone string
		var location string
		if err := rows.Scan(&name, &rType, &phone, &location); err != nil {
			return nil, stacktrace.Propagate(err, "Failed to get the restaurant from the database")
		}
		restaurant := models.Restaurant{Name: name, Type: rType, Phone: phone, Location: location}
		restaurants = append(restaurants, restaurant)
	}
	if err := rows.Err(); err != nil {
		return nil, stacktrace.Propagate(err, "Failed to get the restaurant")
	}

	return restaurants, nil
}

func (dbs *DbService) InsertRestaurants(restaurants []models.Restaurant) error {
	err := dbs.executeSQLQuery(deleteAllRestaurants)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to delete the restaurants")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbs.ConfSrv.Config.Host, dbs.ConfSrv.Config.DbPort, dbs.ConfSrv.Config.User,
		dbs.ConfSrv.Config.Password, dbs.ConfSrv.Config.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to open Postgresql connection")
	}
	defer db.Close()

	for _, restaurant := range restaurants {
		stmt, err := db.Prepare(insertRestaurant)
		if err != nil {
			return stacktrace.Propagate(err, "Failed to prepare the insert query")
		}

		result, err := stmt.Exec(restaurant.Name, restaurant.Type, restaurant.Phone, restaurant.Location)
		if err != nil {
			return stacktrace.Propagate(err, "Failed to insert the restaurant")
		}
		stmt.Close()

		res, err := result.RowsAffected()
		if err != nil {
			return stacktrace.Propagate(err, "Failed to retrieve effected rows after the insert")
		}
		if res != 1 {
			return stacktrace.Propagate(err, "Failed to insert the restaurant")
		}
	}
	return nil
}

func (dbs *DbService) UpdateRestaurant(oldRestName string, rest models.Restaurant) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbs.ConfSrv.Config.Host, dbs.ConfSrv.Config.DbPort, dbs.ConfSrv.Config.User,
		dbs.ConfSrv.Config.Password, dbs.ConfSrv.Config.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to open Postgresql connection")
	}
	defer db.Close()

	stmt, err := db.Prepare(updateRestaurantByName)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to prepare the delete query")
	}

	result, err := stmt.Exec(rest.Name, rest.Type, rest.Phone, rest.Location, oldRestName)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to delete the restaurant")
	}
	stmt.Close()

	res, err := result.RowsAffected()
	if err != nil {
		return stacktrace.Propagate(err, "Failed to retrieve effected rows after the delete")
	}
	if res != 1 {
		return stacktrace.Propagate(err, "Failed to delete the restaurant")
	}

	return nil
}

func (dbs *DbService) DeleteRestaurant(name string) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbs.ConfSrv.Config.Host, dbs.ConfSrv.Config.DbPort, dbs.ConfSrv.Config.User,
		dbs.ConfSrv.Config.Password, dbs.ConfSrv.Config.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to open Postgresql connection")
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteRestaurantByName)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to prepare the delete query")
	}

	result, err := stmt.Exec(name)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to delete the restaurant")
	}
	stmt.Close()

	res, err := result.RowsAffected()
	if err != nil {
		return stacktrace.Propagate(err, "Failed to retrieve effected rows after the delete")
	}
	if res != 1 {
		return stacktrace.Propagate(err, "Failed to delete the restaurant")
	}

	return nil
}

func (dbs *DbService) executeSQLQuery(query string) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbs.ConfSrv.Config.Host, dbs.ConfSrv.Config.DbPort, dbs.ConfSrv.Config.User,
		dbs.ConfSrv.Config.Password, dbs.ConfSrv.Config.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to open Postgresql connection")
	}
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to prepare the query")
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return stacktrace.Propagate(err, "Failed to execute the query")
	}
	return nil
}
