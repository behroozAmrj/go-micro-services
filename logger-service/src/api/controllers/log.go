package logController

import (
	helper "logger-service/src/api/helpers"
	"logger-service/src/data"
	"net/http"

	//"go.mongodb.org/mongo-driver/event"
)

type JSONPayLoad struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayLoad
	_ = helper.ReadJson(w,r, &requestPayload)

	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	var logEntry data.LogEntry
	err := logEntry.Insert(event)
	if err != nil {
		helper.ErrorJSON(w, err)
		return
	}

	resp := data.JsonResponse{
		Error: false,
		Message: "logged",
	}

 	helper.WriteJson(w, http.StatusAccepted,resp)
}

