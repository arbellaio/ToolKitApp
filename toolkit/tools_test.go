package toolkit

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
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

func TestTools_CreateDirIfNotExist(t *testing.T) {
	var testTool Tools
	err := testTool.CreateDirIfNotExist("./test_data/myDir")
	if err != nil {
		t.Error(err)
	}

	err = testTool.CreateDirIfNotExist("./test_data/myDir")
	if err != nil {
		t.Error(err)
	}

	_ = os.Remove("./test_data/myDir")
}

var slugTests = []struct {
	name          string
	s             string
	expected      string
	errorExpected bool
}{
	{name: "valid string", s: "now is the time", expected: "now-is-the-time", errorExpected: false},
	{name: "empty string", s: "", expected: "", errorExpected: true},
	{name: "complex string", s: "Now is the time for all GOOD men! + fish & such &^123", expected: "now-is-the-time-for-all-good-men-fish-such-123", errorExpected: false},
	{name: "japanese string", s: "こんにちは世界", expected: "", errorExpected: true},
	{name: "japanese string and roman characters", s: "hello world こんにちは世界", expected: "hello-world", errorExpected: false},
}

func TestTools_Slugify(t *testing.T) {
	var testTools Tools

	for _, e := range slugTests {
		slug, err := testTools.Slugify(e.s)
		if err != nil && !e.errorExpected {
			t.Errorf("%s: error received when none expected: %s", e.name, e.s)
			log.Println(err)
		}

		if !e.errorExpected && slug != e.expected {
			t.Errorf("%s: wrong slug returned; expected %s but got %s", e.name, e.expected, slug)
		}

		log.Println(slug)
	}
}
