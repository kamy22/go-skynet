package tests

import (
	"bytes"
	"io/ioutil"
	"path"
	"strings"
	"testing"

	skynet "github.com/NebulousLabs/go-skynet"
	"gopkg.in/h2non/gock.v1"
)

// TestDownloadFile tests downloading a single file.
func TestDownloadFile(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	const srcFile = "../testdata/file1.txt"
	const skylink = "testskynet"
	const sialink = skynet.URISkynetPrefix + skylink

	file, err := ioutil.TempFile("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	dstFile := file.Name()

	// Download file request.
	//
	// Match against the full URL, including the skylink.
	opts := skynet.DefaultDownloadOptions
	urlpath := strings.TrimRight(opts.PortalDownloadPath, "/") + "/" + skylink
	gock.New(skynet.DefaultPortalURL).
		Get(urlpath).
		Reply(200).
		BodyString("test\n")

	// Pass the full sialink to verify that the prefix is trimmed.
	err = skynet.DownloadFile(dstFile, sialink, opts)
	if err != nil {
		t.Fatal(err)
	}

	// Check file equality.
	f1, err1 := ioutil.ReadFile(srcFile)
	if err1 != nil {
		t.Fatal(err1)
	}
	f2, err2 := ioutil.ReadFile(path.Clean(dstFile))
	if err2 != nil {
		t.Fatal(err2)
	}
	if !bytes.Equal(f1, f2) {
		t.Fatalf("Downloaded file at %v did not equal uploaded file %v", dstFile, srcFile)
	}

	// Verify we don't have pending mocks.
	if !gock.IsDone() {
		t.Fatal("test finished with pending mocks")
	}
}