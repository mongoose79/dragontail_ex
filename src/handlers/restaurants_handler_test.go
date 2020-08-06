package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"models"
	"net/http"
	"sort"
	"source_manager_service"
	"testing"
	"utils"
)

func init() {
	fileName := "C:\\workspace\\dragontail_ex\\restaurants.csv"
	sm := source_manager_service.NewSourceManagerService()
	err := sm.LoadCsvSourceFile(fileName)
	if err != nil {
		panic(err)
	}
}

func TestGetAllRestaurants(t *testing.T) {
	url := fmt.Sprintf("/restaurants")
	request, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	response := utils.InvokeRequest(request, GetAllRestaurants, "/restaurants")
	assert.Equal(t, http.StatusOK, response.Code)

	var products []models.Restaurant
	err = json.Unmarshal(response.Body.Bytes(), &products)
	assert.NoError(t, err)

	sort.Sort(models.ByName(products))

	assert.Equal(t, 4, len(products))

	assert.Equal(t, "Amora mio", products[0].Name)
	assert.Equal(t, "italian", products[0].Type)
	assert.Equal(t, "35244040", products[0].Phone)
	assert.Equal(t, "Shlomo Ibn Gabirol St 98, Tel Aviv-Yafo, Israel", products[0].Location)

	assert.Equal(t, "Giraffe", products[1].Name)
	assert.Equal(t, "Asian", products[1].Type)
	assert.Equal(t, "36916294", products[1].Phone)
	assert.Equal(t, "Shlomo Ibn Gabirol St 49, Tel Aviv-Yafo, Israel", products[1].Location)

	assert.Equal(t, "Hudson", products[2].Name)
	assert.Equal(t, "Grill", products[2].Type)
	assert.Equal(t, "036444733", products[2].Phone)
	assert.Equal(t, "HaBarzel St 27, Tel Aviv-Yafo, Israel", products[2].Location)

	assert.Equal(t, "Humongous", products[3].Name)
	assert.Equal(t, "Burger", products[3].Type)
	assert.Equal(t, "98943656", products[3].Phone)
	assert.Equal(t, "Atir Yeda St 1, Kefar Sava, Israel", products[3].Location)
}

func TestDeleteRestaurants(t *testing.T) {
	url := fmt.Sprintf("/restaurant/Giraffe")
	request, err := http.NewRequest("DELETE", url, nil)
	assert.NoError(t, err)

	response := utils.InvokeRequest(request, DeleteRestaurant, "/restaurant/{name}")
	assert.Equal(t, http.StatusOK, response.Code)

	url = fmt.Sprintf("/restaurants")
	request, err = http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	response = utils.InvokeRequest(request, GetAllRestaurants, "/restaurants")
	assert.Equal(t, http.StatusOK, response.Code)

	var products []models.Restaurant
	err = json.Unmarshal(response.Body.Bytes(), &products)
	assert.NoError(t, err)

	sort.Sort(models.ByName(products))

	assert.Equal(t, 3, len(products))

	assert.Equal(t, "Amora mio", products[0].Name)
	assert.Equal(t, "italian", products[0].Type)
	assert.Equal(t, "35244040", products[0].Phone)
	assert.Equal(t, "Shlomo Ibn Gabirol St 98, Tel Aviv-Yafo, Israel", products[0].Location)

	assert.Equal(t, "Hudson", products[1].Name)
	assert.Equal(t, "Grill", products[1].Type)
	assert.Equal(t, "036444733", products[1].Phone)
	assert.Equal(t, "HaBarzel St 27, Tel Aviv-Yafo, Israel", products[1].Location)

	assert.Equal(t, "Humongous", products[2].Name)
	assert.Equal(t, "Burger", products[2].Type)
	assert.Equal(t, "98943656", products[2].Phone)
	assert.Equal(t, "Atir Yeda St 1, Kefar Sava, Israel", products[2].Location)
}
