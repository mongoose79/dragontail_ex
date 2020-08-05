package services

import (
	"github.com/stretchr/testify/assert"
	"models"
	"testing"
)

var conf models.Configuration

func init() {
	fileName := "C:\\workspace\\dragontail_ex\\src\\config\\tsconfig.json"
	conf, _ = NewConfService().ReadConfiguration(fileName)
}

func TestInitConfig(t *testing.T) {
	assert.NotNil(t, conf)
	assert.Equal(t, 8080, conf.ServicePort)
	assert.Equal(t, "restaurants.csv", conf.CSVSourceFile)
	assert.Equal(t, "localhost", conf.Host)
	assert.Equal(t, 5432, conf.DbPort)
	assert.Equal(t, "postgres", conf.User)
	assert.Equal(t, "sa", conf.Password)
	assert.Equal(t, "DragonTailDB", conf.Dbname)
}
