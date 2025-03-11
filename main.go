package main

import (
	"fmt"
	"os"

	"github.com/itsmenewbie03/go-tiktok-dl/downloader"
	"github.com/itsmenewbie03/go-tiktok-dl/utils"
)

func main() {
	appName := os.Args[0]

	if len(os.Args) < 2 {
		fmt.Printf("❌ Usage: %s <TikTok_URL>\n", appName)
		os.Exit(1)
	}

	tiktokURL := os.Args[1]

	client := downloader.TiktokDownloader{}
	ioutil := utils.IOUtils{}

	downloadURL, err := client.Download(tiktokURL)
	if err != nil {
		fmt.Printf("❌ Download failed: %s\n", err.Error())
		os.Exit(-1)
	}

	fmt.Println("📥 Downloading video...")

	if err := ioutil.DownloadFile(*downloadURL); err != nil {
		fmt.Printf("❌ Download failed: %s\n", err.Error())
		os.Exit(-1)
	}

	fmt.Println("✅ Download Success! 🎉 Enjoy your video! 🎬")
}
