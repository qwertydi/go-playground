package aggregator

import (
	"log"
	"sync"
	"testing"
	"time"
)

type MockTimeServiceHandler struct {
	time time.Time
}

func (m MockTimeServiceHandler) GetCurrentTime() time.Time {
	return m.time
}

// Test the ProcessMessage method using single message
func TestDataService_ProcessMessage(t *testing.T) {
	// Initialize the struct
	mockTimeService := &MockTimeServiceHandler{
		time: time.Date(2024, 02, 20, 22, 22, 22, 22, time.UTC),
	}

	handler := DataServiceHandlerImpl{
		mu:                 sync.Mutex{},
		data:               *InitTimeData(),
		parent:             make(map[string]string),
		timeServiceHandler: mockTimeService,
	}

	// Prepare data
	byteArray := []byte("{\"source\":\"4e2ac2e6-eab6-427c-9bf0-4e816e827054\",\"destination\":\"e3682e74-556a-4ba5-92ba-afafa9dbd6fe\",\"method\":\"GET\",\"path\":\"/list\",\"httpStatus\":400}")

	// Execution
	handler.ProcessMessage(byteArray)

	// Validation
	dateTime := generateKeyByDateTime(time.Date(2024, 02, 20, 22, 22, 22, 22, time.UTC))

	data := handler.data.data[dateTime]
	if len(data.records) != 1 {
		t.Errorf("Expected len 1, but got %d", len(data.records))
	}
}

// Test the ProcessMessage method using messages
func TestDataService_ProcessMessages(t *testing.T) {
	// Initialize the struct
	mockTimeService := &MockTimeServiceHandler{
		time: time.Date(2024, 02, 20, 22, 22, 22, 22, time.UTC),
	}

	handler := DataServiceHandlerImpl{
		mu:                 sync.Mutex{},
		data:               *InitTimeData(),
		parent:             make(map[string]string),
		timeServiceHandler: mockTimeService,
	}

	// Prepare data
	for _, byteMsg := range getEvents() {
		handler.ProcessMessage(byteMsg)
	}

	// Execution
	dateTime := generateKeyByDateTime(time.Date(2024, 02, 20, 22, 22, 40, 22, time.UTC))
	data := handler.data.data[dateTime]

	log.Println(handler.data)
	// Validation
	if len(data.records) == 0 {
		t.Errorf("Expected len > 0, but got %d", len(data.records))
	} else if len(data.records) == 4 {
		t.Log("Valid number of records")
	}
	var validParents = []string{"b0cfe248-882c-4fce-9103-763eed9af39d", "5a5b510e-0586-4079-97d1-20717de9a58d", "89401e8a-254d-4e13-bb61-ce74842bcf33", "8b0f9f09-922f-4ffe-b4a9-73bb033bf28f", "45a07c5a-4b1f-418a-8702-8c9785077a0b", "87e28180-0df0-4d4e-af0b-1463ecb6a08a"}
	var validChildren = []string{"4523d645-244e-42c1-8831-420f9424755d", "ddd99459-4f13-4635-8c15-763d8128279e", "5ecbdbec-48d5-478d-a827-85c415b8e579", "bae92185-336f-4269-a393-bb52946222e5", "e5f1782a-8bb0-4b8f-9d6a-3bdcd0dc017f", "f4eca45f-dbe2-4159-a860-0753352fb754", "f4eca45f-dbe2-4159-a860-0753352fb754"}

	parent := handler.parent
	for k, v := range parent {
		if !stringExists(validParents, v) {
			t.Errorf("Expected parent %s\n", v)
		}
		if !stringExists(validChildren, k) {
			t.Errorf("Expected child %s\n", k)
		}
	}
}

func TestDataServiceHandlerImpl_GetData(t *testing.T) {
	// Initialize the struct
	mockTimeService := &MockTimeServiceHandler{
		time: time.Date(2024, 02, 20, 22, 22, 22, 22, time.UTC),
	}

	handler := DataServiceHandlerImpl{
		mu:                 sync.Mutex{},
		data:               *InitTimeData(),
		parent:             make(map[string]string),
		timeServiceHandler: mockTimeService,
	}

	// Prepare data
	for _, byteMsg := range getEvents() {
		handler.ProcessMessage(byteMsg)
	}

	// Execution
	data := handler.GetData(time.Date(2024, 02, 20, 22, 22, 22, 22, time.UTC))

	if data == nil {
		t.Errorf("Expected not null object %v\n", data)
	}

	if len(data.records) == 4 {
		t.Errorf("Expected len 4\n")
	}

}

func getEvents() [][]byte {
	var i [][]byte
	i = append(i, []byte("[{\"parent\":\"b0cfe248-882c-4fce-9103-763eed9af39d\",\"children\":\"4523d645-244e-42c1-8831-420f9424755d\"},{\"parent\":\"5a5b510e-0586-4079-97d1-20717de9a58d\",\"children\":\"ddd99459-4f13-4635-8c15-763d8128279e\"},{\"parent\":\"89401e8a-254d-4e13-bb61-ce74842bcf33\",\"children\":\"5ecbdbec-48d5-478d-a827-85c415b8e579\"},{\"parent\":\"8b0f9f09-922f-4ffe-b4a9-73bb033bf28f\",\"children\":\"bae92185-336f-4269-a393-bb52946222e5\"},{\"parent\":\"45a07c5a-4b1f-418a-8702-8c9785077a0b\",\"children\":\"e5f1782a-8bb0-4b8f-9d6a-3bdcd0dc017f\"},{\"parent\":\"87e28180-0df0-4d4e-af0b-1463ecb6a08a\",\"children\":\"f4eca45f-dbe2-4159-a860-0753352fb754\"}]\n"))
	i = append(i, []byte("\t{\"source\":\"4523d645-244e-42c1-8831-420f9424755d\",\"destination\":\"b28310ae-8dc8-4dc5-99ce-345663a9d3ff\",\"method\":\"POST\",\"path\":\"/create\",\"httpStatus\":200}\n"))
	i = append(i, []byte("\t{\"source\":\"5ecbdbec-48d5-478d-a827-85c415b8e579\",\"destination\":\"4523d645-244e-42c1-8831-420f9424755d\",\"method\":\"PATCH\",\"path\":\"/update\",\"httpStatus\":200}\n"))
	i = append(i, []byte("\t{\"source\":\"4523d645-244e-42c1-8831-420f9424755d\",\"destination\":\"ddd99459-4f13-4635-8c15-763d8128279e\",\"method\":\"GET\",\"path\":\"/list\",\"httpStatus\":400}\n"))
	i = append(i, []byte("\t{\"source\":\"f4eca45f-dbe2-4159-a860-0753352fb754\",\"destination\":\"752e2741-9bf9-4def-9b62-e71163ee9a66\",\"method\":\"POST\",\"path\":\"/create\",\"httpStatus\":200}\n"))
	return i
}

// Function to check if a string exists in an array of strings
func stringExists(array []string, target string) bool {
	for _, str := range array {
		if str == target {
			return true
		}
	}
	return false
}
