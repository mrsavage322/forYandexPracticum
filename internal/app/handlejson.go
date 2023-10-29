package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "reflect"
	_ "strings"
	_ "time"
)

//const (
//	targetField = "User" // имя поля, о котором нужно получить информацию
//	targetTag   = "json" // тег, значение которого нужно получить
//)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

//type JSONResponse struct {
//	Result string `json:"result"`
//}

func HandleJSON(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req Request
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//var jsonResponse JSONResponse
	//err = json.Unmarshal([]byte(req.Result), &jsonResponse)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}

	link := req.URL

	//obj := Response{}
	//objType := reflect.TypeOf(req)
	//field, ok := objType.FieldByName(targetField)
	//if !ok {
	//	panic(fmt.Errorf("field (%s): not found", targetField))
	//}
	//tagValue, ok := field.Tag.Lookup(targetTag)
	//if !ok {
	//	panic(fmt.Errorf("tag (%s) for field (%s): not found", targetTag, targetField))
	//}

	//link := strings.Split(tagValue, ",")
	id := GenerateRandomID(5)
	shortURL := fmt.Sprintf("%s/%s", BaseURL, id)
	URLMap.Set(id, link)

	resp := Response{Result: shortURL}
	responseData, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}
