package feeddownload

import (
	"fmt"
	"os"
)

func HandleFile(remoteURL, targetFilename string, dryRun bool) error {
	// already downloaded? -> stop processing
	info, err := os.Stat(targetFilename)
	if !os.IsNotExist(err) && !info.IsDir() {
		return nil
	}

	// actual download
	fmt.Println("Downloading episode...")
	if !dryRun {
		return downloadFile(targetFilename, remoteURL)
	}
	return nil
}
