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

//TestAddEmptyCategory tests that if the category field is empty it should do nothing
func TestAddEmptyCategory(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(AddCategoryFunc))
	defer ts.Close()
	req, err := http.NewRequest("POST", ts.URL, nil)
	req.Form = make(map[string][]string, 0)
	req.Form.Add("category", "")
	// req.Form, _ = url.ParseQuery("category=")
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

func TestEditTaskWithWrongMethod(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(EditTaskFunc))
	defer ts.Close()
	ts.URL = ts.URL + "/edit/"
	req, err := http.NewRequest("OPTIONS", ts.URL, nil)
	if err != nil {
		t.Errorf("Error occured while constracting request:%s", err)
	}
	w := httptest.NewRecorder()
	EditTaskFunc(w, req)

	if w.Code != http.StatusFound && message != "Method not allowed" {
		t.Errorf("Message was: %s Return code was: %d. Should have been message: %s return code: %d", message, w.Code, "Method not allowed", http.StatusBadRequest)
	}
}

func TestEditTaskWrongTaskName(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(EditTaskFunc))
	defer ts.Close()
	ts.URL = ts.URL + "/edit/invalidID"
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Errorf("Error occured while constracting request:%s", err)
	}
	w := httptest.NewRecorder()
	EditTaskFunc(w, req)
	//TODO: Error should be returned as part of the string so that the message can tell
	//why there was a problem.
	if w.Code != http.StatusBadRequest {
		t.Errorf("Actual status: (%d); Expected status:(%d)", w.Code, http.StatusBadRequest)
	}
}

func TestAddCommentWithWrongMethod(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(AddCommentFunc))
	defer ts.Close()
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Errorf("Error occured while constructing request: %s", err)
	}

	w := httptest.NewRecorder()
	AddCommentFunc(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Actual status: (%d); Expected status:(%d)", w.Code, http.StatusBadRequest)
	}
}
