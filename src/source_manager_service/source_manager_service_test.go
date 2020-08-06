package source_manager_service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadCsvSourceFile(t *testing.T) {
	fileName := "C:\\workspace\\dragontail_ex\\restaurants.csv"
	sm := NewSourceManagerService()
	err := sm.LoadCsvSourceFile(fileName)
	assert.NoError(t, err)
}
