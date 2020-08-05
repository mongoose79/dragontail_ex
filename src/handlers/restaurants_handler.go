package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/palantir/stacktrace"
	"io/ioutil"
	"log"
	"models"
	"net/http"
	"services"
	"utils"
)

func GetAllRestaurants(w http.ResponseWriter, r *http.Request) {
	log.Println("Get all restaurants request was received")
	rms := services.NewRestaurantManagerService()
	restaurants, err := rms.GetAllRestaurants()
	if err != nil {
		log.Println("Failed to retrieve the restaurants")
		utils.WriteJSON(err, w, http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(restaurants, w, http.StatusOK)
	log.Println("Get all restaurants request was completed successfully")
}

func EditRestaurant(w http.ResponseWriter, r *http.Request) {
	log.Println("Edit restaurant request was received")
	oldRestName, err := validateInputParameter(r)
	if err != nil {
		log.Println(err)
		utils.WriteJSON(err, w, http.StatusBadRequest)
		return
	}
	fmt.Println(oldRestName)

	var rest models.Restaurant
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err == nil {
		err = json.Unmarshal(bodyBytes, &rest)
		if err != nil {
			log.Println(fmt.Sprintf("Failed to unmarshal the body: %v", err))
			utils.WriteJSON(err, w, http.StatusBadRequest)
			return
		}
	}

	rms := services.NewRestaurantManagerService()
	err = rms.UpdateRestaurant(oldRestName, rest)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to update the restaurant due to internal server error: %v", err))
		utils.WriteJSON(err, w, http.StatusBadRequest)
		return
	}

	log.Println("Edit restaurant request was completed successfully")
}

func DeleteRestaurant(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete restaurant request was received")
	name, err := validateInputParameter(r)
	if err != nil {
		log.Println(err)
		utils.WriteJSON(err, w, http.StatusBadRequest)
		return
	}

	rms := services.NewRestaurantManagerService()
	err = rms.DeleteRestaurant(name)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to delete the restaurant due to internal server error: %v", err))
		utils.WriteJSON(err, w, http.StatusBadRequest)
		return
	}

	log.Println("Delete restaurant request was completed successfully")
}

func validateInputParameter(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		return "", stacktrace.NewError("Restaurant name is invalid")
	}
	return name, nil
}
