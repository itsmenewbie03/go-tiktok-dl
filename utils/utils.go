package utils

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"time"
)

type IOUtils struct{}

func (s *IOUtils) getFilenameFromHeader(url string) (string, string, error) {
	// Create a list of common User-Agent strings.
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (iPad; CPU OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Windows NT 10.0; Trident/7.0; rv:11.0) like Gecko",               // IE11
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS_7.0; rv:11.0) like Gecko", // IE11 on win7
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
	}

	// Choose a random User-Agent.
	userAgent := userAgents[rand.Intn(len(userAgents))]

	// Create a new HTTP request.
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", "", err
	}

	// Set the User-Agent header.
	req.Header.Set("User-Agent", userAgent)

	// Perform the HTTP request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	contentDisp := resp.Header

	// Generate the default filename with date.
	now := time.Now().Format("20060102_150405")
	defaultFilename := fmt.Sprintf("video_%s.mp4", now)

	if contentDisp == nil {
		return defaultFilename, userAgent, nil // Default filename if header is missing
	}

	contentDisposition := resp.Header.Get("Content-Disposition")

	if contentDisposition != "" {
		re := regexp.MustCompile(`filename=(.*)`)
		matches := re.FindStringSubmatch(contentDisposition)
		if len(matches) > 1 {
			return matches[1], userAgent, nil
		}
	}

	return defaultFilename, userAgent, nil // Default filename if header is missing
}

func (s *IOUtils) DownloadFile(url string) error {
	filename, ua, err := s.getFilenameFromHeader(url)
	if err != nil {
		fmt.Printf("❌ Error getting filename: %v\n", err)
		return err
	}

	fmt.Printf("🎬 Downloading file: %v\n", filename)

	cmd := exec.Command("curl", "-A", ua, "-L", url, "-o", filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("🔥 Download failed: %v\n", err)
		return err
	}

	fmt.Printf("✅ Download complete: %v\n", filename)
	return nil
}
