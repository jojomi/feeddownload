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
	if dryRun {
		fmt.Print("[DRY-RUN] ")
	}
	fmt.Printf("Downloading %s to %s...\n", remoteURL, targetFilename)
	if !dryRun {
		return downloadFile(targetFilename, remoteURL)
	}
	return nil
}
