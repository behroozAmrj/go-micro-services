package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	helper "shop.authentication/src/api/helpers"
	data "shop.authentication/src/models"
)

func Login() {
	fmt.Println("this is from controller")
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	payLoad := data.JsonResponse{
		Error:   false,
		Message: "Hi the broker This is From Controller",
	}

	out, _ := json.MarshalIndent(payLoad, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)

	var requestPayload struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err := helper.ReadJson(w , 
		r , 
		&requestPayload)

	if err != nil {
		helper.ErrorJSON(w, err , http.StatusBadRequest)
		return
	}

	var usr =  data.User{

	};
	user ,err := usr.GetByEmail(requestPayload.Email)
	if err != nil {
		helper.ErrorJSON(w, errors.New("Invalid Credentials") , http.StatusForbidden)
		return
	}

	valid ,err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid{
		helper.ErrorJSON(w, errors.New("Invalid Credentials") , http.StatusForbidden)
		return
	}

	payLoad = data.JsonResponse{
		Error: false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: user,
	}

	helper.WriteJson(w , 
		http.StatusAccepted , 
		payLoad)

}
