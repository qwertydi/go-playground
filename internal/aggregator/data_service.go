package aggregator

import (
	"encoding/json"
	"fmt"
	"github.com/qwertydi/go-challenge/internal/util"
	"github.com/qwertydi/go-challenge/models"
	"sync"
	"time"
)

type RecordEntry struct {
	destination string
	count       int32
}

// Data struct with a map field
type Data struct {
	records map[string][]RecordEntry
}

type TimeData struct {
	data map[string]*Data
}

type DataServiceHandler interface {
	GetParents() map[string]string
	GetData(time time.Time) *Data
	RemoveData(time time.Time)
	ProcessMessage(message []byte)
}

type DataServiceHandlerImpl struct {
	mu                 sync.Mutex
	data               TimeData
	parent             map[string]string
	timeServiceHandler util.TimeServiceHandler
}

func InitTimeData() *TimeData {
	return &TimeData{
		data: make(map[string]*Data),
	}
}

func InitData() *Data {
	return &Data{
		records: make(map[string][]RecordEntry),
	}
}

func generateKeyByDateTime(t time.Time) string {
	// Format the time to generate a key like "YYYY-MM-DD HH:MM"
	return t.Format("2006-01-02_15:04")
}

func (h *DataServiceHandlerImpl) ProcessMessage(message []byte) {
	var result interface{}
	err := json.Unmarshal(message, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	dateTime := generateKeyByDateTime(h.timeServiceHandler.GetCurrentTime())
	date := h.data.data[dateTime]
	if date == nil {
		h.data.data[dateTime] = InitData()
	}

	switch result.(type) {
	case []interface{}:
		var listParentData []models.ParentData
		err := json.Unmarshal(message, &listParentData)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
		h.insertAssociationData(listParentData)
	case map[string]interface{}:
		var requestData models.RequestData
		err := json.Unmarshal(message, &requestData)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
		h.insertRequestData(requestData, dateTime)
	default:
		fmt.Println("Unknown JSON type.")
		return
	}
}

func (h *DataServiceHandlerImpl) insertAssociationData(association []models.ParentData) {
	h.mu.Lock()
	defer h.mu.Unlock()
	// Iterate over the slice using range
	for _, data := range association {
		h.parent[data.Children] = data.Parent
	}
}

func (h *DataServiceHandlerImpl) insertRequestData(data models.RequestData, date string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	listByDate := h.data.data[date]
	key := listByDate.records[data.Source]
	getListOrDefault := func(key []RecordEntry) []RecordEntry {
		// Check if the key exists
		if len(key) == 0 {
			return []RecordEntry{}
		}
		return key
	}

	recordEntries := getListOrDefault(key)
	if contains(recordEntries, data.Destination) {
		increment(recordEntries, data.Destination)
	} else {
		recordEntries = append(recordEntries, RecordEntry{
			data.Destination,
			1,
		})
	}
	for _, recordEntry := range recordEntries {
		if recordEntry.destination == data.Destination {
			recordEntry.count++
		} else {
			recordEntry.count = 1
			recordEntry.destination = data.Destination
		}
	}

	h.data.data[date].records[data.Source] = recordEntries
}

func contains(slice []RecordEntry, item string) bool {
	for _, s := range slice {
		if s.destination == item {
			return true
		}
	}
	return false
}

func increment(slice []RecordEntry, item string) {
	for _, s := range slice {
		if s.destination == item {
			s.count++
		}
	}
}

func (h *DataServiceHandlerImpl) GetData(time time.Time) *Data {
	dateTime := generateKeyByDateTime(time)
	return h.data.data[dateTime]
}

func (h *DataServiceHandlerImpl) GetParents() map[string]string {
	return h.parent
}

func (h *DataServiceHandlerImpl) RemoveData(time time.Time) {
	dateTime := generateKeyByDateTime(time)
	delete(h.data.data, dateTime)
}

// DataService initializes a new instance of Service with a map
func DataService(timeServiceHandler util.TimeServiceHandler) *DataServiceHandlerImpl {
	return &DataServiceHandlerImpl{
		mu:                 sync.Mutex{},
		data:               *InitTimeData(),
		parent:             make(map[string]string),
		timeServiceHandler: timeServiceHandler,
	}
}
