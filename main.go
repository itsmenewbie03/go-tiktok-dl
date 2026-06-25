package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
	"github.com/itsmenewbie03/go-tiktok-dl/downloader"
	"github.com/itsmenewbie03/go-tiktok-dl/utils"
	"github.com/spf13/cobra"
)

var watch bool

var rootCmd = &cobra.Command{
	Use:   "ttdl [TikTok URL]",
	Short: "Download TitTok videos via the CLI",
	Long:  `ttdl is a CLI tool for downloading and watching tiktok videos.`,
	Args:  cobra.ArbitraryArgs,
	Run:   ttdl,
}

func init() {
	rootCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch the video instead of downloading")
}

func ttdl(cmd *cobra.Command, args []string) {
	var ttURL string

	if len(args) <= 0 {
		cmd.Help()
		return
	}

	ttURL = args[0]

	if ttURL == "" {
		fmt.Println("No tiktok url provided")
		return
	}

	client := downloader.TiktokDownloader{}
	ioutil := utils.IOUtils{}

	downloadOpts, err := client.Download(ttURL)
	if err != nil {
		fmt.Printf("❌ Download failed: %s\n", err.Error())
		os.Exit(-1)
	}

	if watch {
		streamURL := downloadOpts.WithoutWaterMark
		mpv := exec.Command("mpv", streamURL.String())
		mpv.Stderr = os.Stderr
		mpv.Stdout = os.Stdout
		if err := mpv.Run(); err != nil {
			log.Fatal(err)
		}
		log.Info("ByeBye!")
		return
	}

	downloadURL := downloadOpts.HDWithoutWaterMark
	fmt.Println("📥 Downloading video...")

	if err := ioutil.DownloadFile(downloadURL.String()); err != nil {
		fmt.Printf("❌ Download failed: %s\n", err.Error())
		os.Exit(-1)
	}

	fmt.Println("✅ Download Success! 🎉 Enjoy your video! 🎬")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
