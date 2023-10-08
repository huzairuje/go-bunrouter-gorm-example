package httplib

import (
	"encoding/json"
	"net/http"
)

type DefaultResponse struct {
	Status    string      `json:"status"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	DataError interface{} `json:"dataError"`
}

type DefaultPaginationResponse struct {
	Status     string      `json:"status"`
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Page       int         `json:"page"`
	Size       int         `json:"size"`
	TotalCount uint64      `json:"totalCount"`
	TotalPages uint64      `json:"totalPages"`
	Data       interface{} `json:"data"`
}

func SetSuccessResponse(w http.ResponseWriter, code int, message string, data interface{}) error {
	return ToJSON(w, code, DefaultResponse{
		Status:  http.StatusText(code),
		Code:    code,
		Data:    data,
		Message: message,
	})
}

func SetPaginationResponse(w http.ResponseWriter, code int, message string, data interface{}, totalCount uint64, pg *Query) error {
	return ToJSON(w, code, DefaultPaginationResponse{
		Status:     http.StatusText(code),
		Code:       code,
		Message:    message,
		Page:       pg.GetPage(),
		Size:       pg.GetSize(),
		TotalCount: totalCount,
		TotalPages: uint64(GetTotalPages(int(totalCount), pg.GetSize())),
		Data:       data,
	})
}

func SetErrorResponse(w http.ResponseWriter, code int, message string) error {
	return ToJSON(w, code, DefaultResponse{
		Status:  http.StatusText(code),
		Code:    code,
		Data:    nil,
		Message: message,
	})
}

func SetCustomResponse(w http.ResponseWriter, code int, message string, data interface{}, dataErr interface{}) error {
	return ToJSON(w, code, DefaultResponse{
		Status:    http.StatusText(code),
		Code:      code,
		Data:      data,
		Message:   message,
		DataError: dataErr,
	})
}

func ToJSON(w http.ResponseWriter, code int, value interface{}) error {
	if value == nil {
		return nil
	}

	//set header content type first
	w.Header().Set("Content-Type", "application/json")
	//and set the response code http
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	if err := enc.Encode(value); err != nil {
		return err
	}
	return nil
}
