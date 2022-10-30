package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

//data user altında derli toplu görünsün diye envelope aldım
type envelope map[string]interface{}

// sık kullanılacağı için bir fonksiyona almakta fayda var dedim
// normalde yorum satırına ihtiyaç olan kod kötü koddur dense de o an ne düşündüğümü anlatmak adına ufak ufak notlar bırakıcam
func (app *application) readIDFromParameter(r *http.Request) (int64, error) {
	parameter := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(parameter.ByName("id"), 10, 64)

	if err != nil || id < 1 {
		return 0, errors.New("Bad request")
	}

	return id, nil
}

func (app *application) writeJsonHelper(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// test case olduğundan performansı gözetmediğimden  json.marshall yerine marshallindent kullandım
	// biraz daha uzun sürer ve memory yer ama terminalde şık bir görüntü veriyor istek atınca terminalde güzel duruyor
	jsonformat, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	jsonformat = append(jsonformat, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonformat)

	return nil
}
