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

	binary, found := os.LookupEnv("BINARY")
	if !found {
		t.Fatal("BINARY env var not found")
	}

	cmd := exec.Command(binary, fmt.Sprintf("--pack=%d", packId), fmt.Sprintf("--version=%s", version))
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Error running downloader: %s", err)
	}
	_, err = os.Stat("./server.zip")

	if err != nil {
		t.Fatal("Didn't create a server.zip file")
	}

	file, err := os.Open("./server.zip")
	if err != nil {
		t.Fatal(err)
	}

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		t.Fatal(err)
		file.Close()
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
	shouldDownloadCorrectFile(t, 296062, "latest", "crG65hy9agerVdcfnhqUI5pJkvNaf59KFJ5PHuoEoWs=")
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
