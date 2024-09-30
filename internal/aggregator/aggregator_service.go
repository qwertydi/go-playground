package aggregator

import (
	"github.com/qwertydi/go-challenge/models"
	"strings"
	"time"
)

type AggregateServiceHandler interface {
	AggregateData(time time.Time) []models.AggregatedData
}

type AggregateServiceHandlerImpl struct {
	dataService DataServiceHandler
}

func (h *AggregateServiceHandlerImpl) AggregateData(time time.Time) []models.AggregatedData {
	data := h.dataService.GetData(time)
	parents := h.dataService.GetParents()

	records := data.records

	mappings := make(map[string]int32)

	addOrIncrement := func(key string, counter int32) {
		if _, exists := mappings[key]; exists {
			mappings[key] = mappings[key] + counter
		} else {
			mappings[key] = counter
		}
	}

	for source, destinations := range records {
		// get source parent id
		if _, sourceParentExists := parents[source]; sourceParentExists {
			// get destination parent id
			for _, destination := range destinations {
				if _, sourceDestinationExists := parents[destination.destination]; sourceDestinationExists {
					joinString := parents[source] + "_" + parents[destination.destination]
					addOrIncrement(joinString, destination.count)
				}
			}
		}
	}

	var aggregatedData []models.AggregatedData

	for k, v := range mappings {
		split := strings.Split(k, "_")
		aggregatedData = append(aggregatedData, models.AggregatedData{
			Source:      split[0],
			Destination: split[1],
			Count:       v,
		})
	}

	// remove data already processed
	h.dataService.RemoveData(time)

	return aggregatedData
}

func AggregateService(dateService DataServiceHandler) *AggregateServiceHandlerImpl {
	return &AggregateServiceHandlerImpl{
		dataService: dateService,
	}
}
