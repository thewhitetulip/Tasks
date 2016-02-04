package views

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestUploadedFileHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(UploadedFileHandler))
	defer ts.Close()
	os.Mkdir("./files", 0777)
	defer os.RemoveAll("./files")
	file, err := os.Create("./files/testfile")
	defer file.Close()
	defer os.Remove(file.Name())
	expectedContent := []byte("test content")
	file.Write(expectedContent)
	res, err := http.Get(ts.URL + "/files/testfile")
	if err != nil {
		t.Error("Error occured while getting response from test server:", err)
	}
	actualContent, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error("Error occured while reading content from response Body: ", err)
	}

	if !bytes.Equal(actualContent, expectedContent) {
		t.Errorf("Actual content (%s) did not match expected content (%s)", actualContent, expectedContent)
	}
}

//TestAddEmptyCategory tests that if the category field is empty it should do nothing
func TestAddEmptyCategory(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(AddCategoryFunc))
	defer ts.Close()
	req, err := http.NewRequest("POST", ts.URL, nil)
	req.Form, _ = url.ParseQuery("category=")
	if err != nil {
		t.Errorf("Error occured while constracting request:%s", err)
	}
	w := httptest.NewRecorder()
	AddCategoryFunc(w, req)
	body := w.Body.String()
	if len(body) != 0 {
		t.Error("Body should be empty. Instead contained data: ", body)
	}
}
