package toolkit

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

func TestTools_RandomString(t *testing.T) {
	var testTools Tools

	s := testTools.RandomString(10)

	if len(s) != 10 {
		t.Error("expected 10 characters")
	}

}

var uploadTest = []struct {
	name          string
	allowedTypes  []string
	renameFile    bool
	errorExpected bool
}{
	{name: "allowed no rename", allowedTypes: []string{"image/jpg", "image/png"}, renameFile: false, errorExpected: false},
	{name: "allowed rename", allowedTypes: []string{"image/jpg", "image/png"}, renameFile: true, errorExpected: false},
	{name: "not allowed", allowedTypes: []string{"image/jpg"}, renameFile: false, errorExpected: true},
}

func TestTools_UploadFiles(t *testing.T) {

	for _, e := range uploadTest {

		//setup a pipe to avoid buffering
		pr, pw := io.Pipe()

		writer := multipart.NewWriter(pw)

		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer writer.Close()
			defer wg.Done()

			// create the form data field
			part, err := writer.CreateFormFile("file", "./testdata/img.png")
			if err != nil {
				t.Error(err)
			}

			f, err := os.Open("./testdata/img.png")

			if err != nil {
				t.Error(err)
			}

			defer f.Close()

			img, _, err := image.Decode(f)

			if err != nil {
				t.Error("error decoding image", err)
			}
			err = png.Encode(part, img)

			if err != nil {
				t.Error(err)
			}

		}()

		//read from the pipe which receives data

		request := httptest.NewRequest("POTS", "/", pr)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		var testTools Tools
		testTools.AllowedFileTypes = e.allowedTypes

		uploadedFiles, err := testTools.UploadFiles(request, "./testdata/upload/", e.renameFile)

		if err != nil && !e.errorExpected {
			t.Error(err)
		}

		if !e.errorExpected {
			if _, err := os.Stat(fmt.Sprintf("./testdata/upload/%s", uploadedFiles[0].NewFileName)); os.IsNotExist(err) {
				t.Errorf("%s: expected file to exists", e.name)
			}
			//clean up
			_ = os.Remove(fmt.Sprintf("./testdata/upload/%s", uploadedFiles[0].NewFileName))
		}

		if !e.errorExpected && err != nil {
			t.Error(err)
		}

		wg.Wait()

	}

}

func TestTools_UploadOneFile(t *testing.T) {

	//setup a pipe to avoid buffering
	pr, pw := io.Pipe()

	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		// create the form data field
		part, err := writer.CreateFormFile("file", "./testdata/img.png")
		if err != nil {
			t.Error(err)
		}

		f, err := os.Open("./testdata/img.png")

		if err != nil {
			t.Error(err)
		}

		defer f.Close()

		img, _, err := image.Decode(f)

		if err != nil {
			t.Error("error decoding image", err)
		}
		err = png.Encode(part, img)

		if err != nil {
			t.Error(err)
		}

	}()

	//read from the pipe which receives data
	request := httptest.NewRequest("POTS", "/", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	var testTools Tools

	uploadedFile, err := testTools.UploadOneFile(request, "./testdata/upload/", true)

	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(fmt.Sprintf("./testdata/upload/%s", uploadedFile.NewFileName)); os.IsNotExist(err) {
		t.Errorf("expected file to exists")
	}
	//clean up
	_ = os.Remove(fmt.Sprintf("./testdata/upload/%s", uploadedFile.NewFileName))

}
