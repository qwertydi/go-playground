package aggregator

import (
	"github.com/qwertydi/go-challenge/models"
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

	var aggregatedData []models.AggregatedData

	for source, destinations := range records {
		// get source parent id
		if _, sourceParentExists := parents[source]; sourceParentExists {
			// get destination parent id
			// log.Printf("\nSource Parent exists for: %s, parentId: %s", source, parents[source])
			for _, destination := range destinations {
				if _, sourceDestinationExists := parents[destination.destination]; sourceDestinationExists {
					// log.Printf("\nDestination Parent exists for: %s, parentId: %s", destination.destination, parents[destination.destination])
					// todo aggregate by count by parent
					aggregatedData = append(aggregatedData, models.AggregatedData{
						Source:      parents[source],
						Destination: parents[destination.destination],
						Count:       destination.count,
					})
				}
			}
		}
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
