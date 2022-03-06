package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/jojomi/feeddownload/feeddownload"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
)

var (
	flagRootDryRun   bool
	flagRootVerbose  bool
	flagRootUseTitle bool
)

func main() {
	rootCmd := &cobra.Command{
		Use: "feeddownload",
		Run: handleRootCmd,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("feed url and target folder arguments are required")
			}
			if len(args) < 2 {
				return errors.New("a target folder argument is required")
			}
			return nil
		},
	}

	pFlags := rootCmd.PersistentFlags()
	pFlags.BoolVarP(&flagRootUseTitle, "use-title", "t", true, "use episode title for local filename")
	pFlags.BoolVarP(&flagRootDryRun, "dry-run", "d", false, "just simulate, no downloads are executed")
	pFlags.BoolVarP(&flagRootVerbose, "verbose", "v", false, "more detailed output")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleRootCmd(_ *cobra.Command, args []string) {
	var (
		err            error
		encURL         *url.URL
		remoteFilename string
		targetFilename string
	)

	// parse supplied feed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(args[0])
	if err != nil {
		log.Fatal(err)
	}
	if flagRootVerbose {
		fmt.Println(feed.Title)
	}

	for _, f := range feed.Items {
		if flagRootVerbose {
			fmt.Println("")
			fmt.Println("Episode title:", f.Title)
		}
		for _, e := range f.Enclosures {
			encURL, err = url.Parse(e.URL)
			if err != nil {
				log.Fatal(err)
			}
			if flagRootVerbose {
				fmt.Println("Remote URL:", encURL.String())
			}

			if flagRootUseTitle {
				targetFilename = filepath.Join(args[1], feeddownload.FilenameFromTitle(f.Title)+path.Ext(encURL.Path))
			} else {
				remoteFilename = path.Base(encURL.Path)
				targetFilename = filepath.Join(args[1], remoteFilename)
			}
			if flagRootVerbose {
				fmt.Println("Local file:", targetFilename)
			}

			err = feeddownload.HandleFile(encURL.String(), targetFilename, flagRootDryRun)
			if err != nil {
				panic(err)
			}
		}
	}

	fmt.Printf("checked %d feed items\n", len(feed.Items))
}
