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
		t.Error("wrong length random string returned")
	}
}

var uploadTests = []struct {
	name          string
	allowedTypes  []string
	renameFile    bool
	errorExpected bool
}{
	{
		name:          "allowed no rename",
		allowedTypes:  []string{"image/jpeg", "image/png"},
		renameFile:    false,
		errorExpected: false,
	},
	{
		name:          "allowed rename",
		allowedTypes:  []string{"image/jpeg", "image/png"},
		renameFile:    true,
		errorExpected: false,
	},
	{
		name:          "not allowed",
		allowedTypes:  []string{"image/jpeg"},
		renameFile:    false,
		errorExpected: true,
	},
}

func TestTools_UploadFile(t *testing.T) {
	//set up a pipe to avoid buffering
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer writer.Close()
		defer wg.Done()

		// create the form data field 'file'

		part, err := writer.CreateFormFile("file", "./test_data/img.png")

		if err != nil {
			t.Error(err)
		}
		f, err := os.Open("./test_data/img.png")
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
	request := httptest.NewRequest("POST", "/", pr)
	request.Header.Add("content-type", writer.FormDataContentType())

	var testtools Tools
	uploadedFiles, err := testtools.UploadFile(request, "./test_data/uploads/", true)

	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(fmt.Sprintf("./test_data/uploads/%s", uploadedFiles.NewFileName)); os.IsNotExist(err) {
		t.Errorf("expected file to exist : %s", err.Error())
	}

	_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles.NewFileName))

	wg.Wait()

}

func TestTools_UploadFiles(t *testing.T) {
	for _, e := range uploadTests {
		//set up a pipe to avoid buffering
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer writer.Close()
			defer wg.Done()

			// create the form data field 'file'

			part, err := writer.CreateFormFile("file", "./test_data/img.png")

			if err != nil {
				t.Error(err)
			}
			f, err := os.Open("./test_data/img.png")
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
		request := httptest.NewRequest("POST", "/", pr)
		request.Header.Add("content-type", writer.FormDataContentType())

		var testtools Tools
		testtools.AllowedFileTypes = e.allowedTypes
		uploadedFiles, err := testtools.UploadFiles(request, "./test_data/uploads/", e.renameFile)

		if err != nil {
			t.Error(err)
		}

		if !e.errorExpected {
			if _, err := os.Stat(fmt.Sprintf("./test_data/uploads/%s", uploadedFiles[0].NewFileName)); os.IsNotExist(err) {
				t.Errorf("%s : expected file to exist : %s", e.name, err.Error())
			}

			_ = os.Remove(fmt.Sprintf("./test_data/uploads/%s", uploadedFiles[0].NewFileName))
		}

		if !e.errorExpected && err != nil {
			t.Errorf("%s : error expected but none received", e.name)
		}

		wg.Wait()

	}
}
