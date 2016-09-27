package services

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"
)

type MockJSONStruct struct {
	JSON string `json:"field"`
}

func TestPrintJSON(t *testing.T) {

	input := &MockJSONStruct{JSON: "val"}

	w := httptest.NewRecorder()

	PrintJSON(w, input)

	retrieved := new(MockJSONStruct)
	err := json.Unmarshal([]byte(w.Body.String()), retrieved)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(input, retrieved) {
		t.Errorf("Expected: %+v ,but responsewriter contained:%+v", input, retrieved)
	}

}
