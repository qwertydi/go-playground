package aggregator

import (
	"github.com/qwertydi/go-challenge/models"
	"testing"
	"time"
)

// Mock implementation of DataServiceHandler for testing
type MockDataServiceHandler struct {
	data    Data
	parents map[string]string
}

func (m *MockDataServiceHandler) GetParents() map[string]string {
	return m.parents
}

func (m *MockDataServiceHandler) GetData(time time.Time) *Data {
	return &m.data
}

func (m *MockDataServiceHandler) RemoveData(time time.Time) {
	//TODO implement me
	panic("implement me")
}

func (m *MockDataServiceHandler) ProcessMessage(message []byte) {
	//TODO implement me
	panic("implement me")
}

func TestAggregateService_AggregateData(t *testing.T) {
	// Arrange
	records := make(map[string][]RecordEntry)

	records["1"] = append(make([]RecordEntry, 1), RecordEntry{destination: "2", count: 1})
	records["2"] = append(make([]RecordEntry, 1), RecordEntry{destination: "3", count: 2})
	records["3"] = append(make([]RecordEntry, 1), RecordEntry{destination: "5", count: 1})
	records["6"] = append(make([]RecordEntry, 1), RecordEntry{destination: "1", count: 1})
	records["3"] = append(make([]RecordEntry, 1), RecordEntry{destination: "4", count: 1})
	records["9"] = append(make([]RecordEntry, 1), RecordEntry{destination: "10", count: 1})

	data := Data{
		records: records,
	}

	parents := make(map[string]string)
	parents["1"] = "parent1"
	parents["2"] = "parent2"
	parents["3"] = "parent3"
	parents["4"] = "parent4"

	mockDataService := &MockDataServiceHandler{
		data:    data,
		parents: parents,
	}
	handler := AggregateServiceHandlerImpl{dataService: mockDataService}

	// Act
	result := handler.AggregateData(time.Now())

	// Assert
	expectedData := []models.AggregatedData{
		{Source: "parent3", Destination: "parent4", Count: 1},
		{Source: "parent1", Destination: "parent2", Count: 1},
		{Source: "parent2", Destination: "parent3", Count: 2},
	}

	// t.Log(result, expectedData)
	for _, r := range result {
		if !expectedDataExists(expectedData, r) {
			t.Error("Expected data not in aggregated result")
		}
	}
}

func expectedDataExists(expectedData []models.AggregatedData, r models.AggregatedData) bool {
	exists := false
	for _, e := range expectedData {

		if r.Destination == e.Destination && r.Source == e.Source && r.Count == r.Count {
			exists = true
		}
	}

	return exists
}
