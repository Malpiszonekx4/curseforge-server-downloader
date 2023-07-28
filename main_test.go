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

func shouldDownloadCorrectFile(t *testing.T, packId int, version string, validHash string) {

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

	if hashString != validHash {
		t.Error("Hashes don't match. File is malformed.")
	}

	os.Remove("./server.zip")

}

func TestDownloadSpecifiedVersion(t *testing.T) {
	shouldDownloadCorrectFile(t, 439293, "1.5.3", "ex+q2BusZWGZhKG4a8V5yLf0i4BJGKf4XhjtiJK/46c=")
}

func TestDownloadLatest(t *testing.T) {
	// latest sky factory 4 https://www.curseforge.com/minecraft/modpacks/skyfactory-4/files/3565683
	shouldDownloadCorrectFile(t, 296062, "latest", "aovsWWh/PVVHV5WBGSCWa6tvkUAN9m2ffKkhUyCLnZY=")
}

/*
func testGenHash(t *testing.T) {
	file, _ := os.Open("./server.zip")

	hash := sha256.New()
	_, _ = io.Copy(hash, file)

	file.Close()

	hashString := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	t.Log(hashString)
}
*/
