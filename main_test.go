package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
)

func runTest(t *testing.T, packId int, version string, validHash string) {

	cmd := exec.Command("./curseforge-server-downloader", fmt.Sprintf("--pack=%d", packId), fmt.Sprintf("--version=%s", version))
	err := cmd.Run()
	if err != nil {
		t.Error("Error running downloader")
	}
	_, err = os.Stat("./server.zip")

	if err != nil {
		t.Error("Didn't create a server.zip file")
	}

	file, err := os.Open("./server.zip")
	if err != nil {
		t.Fatal(err)
	}

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		t.Fatal(err)
	}

	file.Close()

	hashString := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	if hashString != "ex+q2BusZWGZhKG4a8V5yLf0i4BJGKf4XhjtiJK/46c=" {
		t.Error("Hashes don't match. File is malformed.")
	}

	os.Remove("./server.zip")

}

func Test1(t *testing.T) {
	runTest(t, 439293, "1.5.3", "ex+q2BusZWGZhKG4a8V5yLf0i4BJGKf4XhjtiJK/46c=")
}

func Test2(t *testing.T) {
	runTest(t, 439293, "latest", "ex+q2BusZWGZhKG4a8V5yLf0i4BJGKf4XhjtiJK/46c=")
}
