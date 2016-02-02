package views

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
