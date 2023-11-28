package helper

import (
	"encoding/json"
	
	
	"net/http"

	model "shop.service/src/models"
)

func ReadJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	//err = dec.Decode(&struct{}{})
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return errors.New("body must have onyl a single JSON value sss:" + err.Error())
//
	//}

	return nil
}

func WriteJson(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, " ", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload model.JsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return WriteJson(w,
		statusCode,
		payload)
}
