package brokerCTRL

import (
	//"encoding/json"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	//"shop.service/src/models"

	helper "shop.service/src/api/helpers"
	model "shop.service/src/models"
)

type AuthPayLoad struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MailPayLoad struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type RequestPayLoad struct {
	Action string      `json:"action"`
	Auth   AuthPayLoad `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayLoad `json:"mail,omitempty"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func BrokersHandler(w http.ResponseWriter, r *http.Request) {
	payLoad := model.JsonResponse{
		Error:   false,
		Message: "Hi the broker This is From Controller",
	}

	//out, _ := json.MarshalIndent(payLoad, "", "\t")
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusAccepted)
	//w.Write(out)

	_ = helper.WriteJson(w, http.StatusOK, payLoad)
}

func HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayLoad

	err := helper.ReadJson(w,
		r,
		&requestPayload)

	if err != nil {
		helper.ErrorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		authenticate(w, requestPayload.Auth)
	case "log":
		logItem(w, requestPayload.Log)
	case "mail":
		sendMail(w, requestPayload.Mail)
	default:
		helper.ErrorJSON(w, errors.New("Unknown action"))
	}
}

func authenticate(w http.ResponseWriter, authPayload AuthPayLoad) {
	//create some json we`ll send to the auth microservice
	jsonData, _ := json.MarshalIndent(authPayload, "", "\t")

	//call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		helper.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		helper.ErrorJSON(w, err)
		return
	}

	defer response.Body.Close()
	//make sure we get back the correct status code

	if response.StatusCode == http.StatusUnauthorized {
		helper.ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		helper.ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	//create a variable we`ll read response.Body info
	var jsonFromService model.JsonResponse //data.JsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		helper.ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		helper.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payLoad model.JsonResponse
	payLoad.Error = false
	payLoad.Message = "Atuhenticated!"
	payLoad.Data = jsonFromService.Data

	helper.WriteJson(w, http.StatusAccepted, payLoad)
}

func logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		helper.ErrorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		helper.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		helper.ErrorJSON(w, err)
		return
	}

	var payload model.JsonResponse
	payload.Error = false
	payload.Message = "logged"

	helper.WriteJson(w, http.StatusAccepted, payload)

}

func sendMail(w http.ResponseWriter, msg MailPayLoad) {
	jsonData, err := json.MarshalIndent(msg, "", "\t")
	//call the mail service
	mailServiceURL := "http://mailer-service/send"

	//post to mail service
	request, err := http.NewRequest("POST",
		mailServiceURL,
		bytes.NewBuffer(jsonData))
	if err != nil {
		helper.ErrorJSON(w, err)
		return
	}

	request.Header.Set("Context-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		helper.ErrorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		helper.ErrorJSON(w, errors.New("error calling mail service "))
		return
	}

	var payLoad model.JsonResponse
	payLoad.Error = false
	payLoad.Message = "Message sent to " + msg.To

	helper.WriteJson(w,
		http.StatusAccepted,
		payLoad)

}
