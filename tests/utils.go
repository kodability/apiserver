package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// makeGet makes a GET request
func makeGet(url string) (req *http.Request, write *httptest.ResponseRecorder, err error) {
	return makeRequest(http.MethodGet, url)
}

// makeDelete makes a DELETE request
func makeDelete(url string) (req *http.Request, write *httptest.ResponseRecorder, err error) {
	return makeRequest(http.MethodDelete, url)
}

// makePostJSON makes a POST request with JSON content-type
func makePostJSON(url string, body interface{}) (req *http.Request, write *httptest.ResponseRecorder, err error) {
	return makeRequestJSON(http.MethodPost, url, body)
}

// makePutJSON makes a POST request with JSON content-type
func makePutJSON(url string, body interface{}) (req *http.Request, write *httptest.ResponseRecorder, err error) {
	return makeRequestJSON(http.MethodPut, url, body)
}

// makeDelete makes a request
func makeRequest(method, url string) (req *http.Request, write *httptest.ResponseRecorder, err error) {
	r, e := http.NewRequest(method, url, nil)
	w := httptest.NewRecorder()
	return r, w, e
}

// makeRequestJSON makes a request with JSON content-type
func makeRequestJSON(method, url string, body interface{}) (req *http.Request, write *httptest.ResponseRecorder, err error) {
	jsonValue, e := json.Marshal(body)
	if e != nil {
		return nil, nil, e
	}

	r, e := http.NewRequest(method, url, bytes.NewBuffer(jsonValue))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	return r, w, e
}
