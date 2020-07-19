package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
)

var (
	flagRootDryRun   bool
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

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleRootCmd(cmd *cobra.Command, args []string) {
	var (
		err                error
		encURL             *url.URL
		remoteFilename     string
		targetFilename     string
		invalidFilename    = regexp.MustCompile(`[^-–.0-9A-Za-z-ÁÀȦÂÄǞǍĂĀÃÅǺǼǢĆĊĈČĎḌḐḒÉÈĖÊËĚĔĒẼE̊ẸǴĠĜǦĞG̃ĢĤḤáàȧâäǟǎăāãåǻǽǣćċĉčďḍḑḓéèėêëěĕēẽe̊ẹǵġĝǧğg̃ģĥḥÍÌİÎÏǏĬĪĨỊĴĶǨĹĻĽĿḼM̂M̄ʼNŃN̂ṄN̈ŇN̄ÑŅṊÓÒȮȰÔÖȪǑŎŌÕȬŐỌǾƠíìiîïǐĭīĩịĵķǩĺļľŀḽm̂m̄ŉńn̂ṅn̈ňn̄ñņṋóòôȯȱöȫǒŏōõȭőọǿơP̄ŔŘŖŚŜṠŠȘṢŤȚṬṰÚÙÛÜǓŬŪŨŰŮỤẂẀŴẄÝỲŶŸȲỸŹŻŽẒǮp̄ŕřŗśŝṡšşṣťțṭṱúùûüǔŭūũűůụẃẁŵẅýỳŷÿȳỹźżžẓǯßœŒçÇ]`) // https://stackoverflow.com/questions/22017723/regex-for-umlaut/56293848#56293848
		colonInTitle       = regexp.MustCompile(`\b:`)
		multipleWhitespace = regexp.MustCompile(`\s+`)
	)

	// parse supplied feed
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(args[0])
	fmt.Println(feed.Title)

	for _, f := range feed.Items {
		fmt.Println("")
		fmt.Println("Episode title:", f.Title)
		for _, e := range f.Enclosures {
			encURL, err = url.Parse(e.URL)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Remote URL:", encURL.String())

			if flagRootUseTitle {
				targetFilename = strings.ReplaceAll(f.Title, "|", "-")
				targetFilename = colonInTitle.ReplaceAllString(targetFilename, " - ")
				targetFilename = invalidFilename.ReplaceAllString(targetFilename, " ")
				targetFilename = multipleWhitespace.ReplaceAllString(targetFilename, " ")
				targetFilename = strings.TrimSpace(targetFilename)
				targetFilename = filepath.Join(args[1], targetFilename+path.Ext(encURL.Path))
			} else {
				remoteFilename = path.Base(encURL.Path)
				targetFilename = filepath.Join(args[1], remoteFilename)
			}
			fmt.Println("Local file:", targetFilename)

			// already downloaded?
			info, err := os.Stat(targetFilename)
			if !os.IsNotExist(err) && !info.IsDir() {
				continue
			}

			// actual download
			fmt.Println("Downloading episode...")
			if !flagRootDryRun {
				downloadFile(targetFilename, encURL.String())
			}
		}
	}
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
