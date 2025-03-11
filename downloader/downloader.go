package downloader

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/itsmenewbie03/go-tiktok-dl/downloader/model"
)

type TiktokDownloader struct{}

func (s *TiktokDownloader) GetTT() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://ssstik.io/", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:136.0) Gecko/20100101 Firefox/136.0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:136.0) Gecko/20100101 Firefox/136.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("HX-Request", "true")
	req.Header.Set("HX-Trigger", "_gcaptcha_pt")
	req.Header.Set("HX-Target", "target")
	req.Header.Set("HX-Current-URL", "https://ssstik.io/")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://ssstik.io")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://ssstik.io/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Priority", "u=0")
	req.Header.Set("TE", "trailers")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	tt := s.extractTT(string(bodyText))
	if tt == nil {
		log.Fatal("failed to extract TT")
	}
	return *tt
}

func (s *TiktokDownloader) extractDownloadURL(body string) *string {
	pattern := `downloadX\('(.*)'\)"`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(body)
	if matches == nil {
		return nil
	}
	return &matches[1]
}

func (s *TiktokDownloader) extractKey(url string) *string {
	pattern := `url=(\w+)==`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(url)
	if matches == nil {
		return nil
	}
	return &matches[1]
}

func (s *TiktokDownloader) extractTT(body string) *string {
	pattern := `s_tt\s=\s'(.*)',`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(body)
	if matches == nil {
		return nil
	}
	return &matches[1]
}

func (s *TiktokDownloader) getDownloadData(link, tt string) *model.DownloadData {
	client := &http.Client{}
	// TODO: make tt query dynamic
	body := fmt.Sprintf("id=%s&locale=en&tt=%s", url.QueryEscape(link), tt)
	data := strings.NewReader(body)
	req, err := http.NewRequest("POST", "https://ssstik.io/abc?url=dl", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:136.0) Gecko/20100101 Firefox/136.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("HX-Request", "true")
	req.Header.Set("HX-Trigger", "_gcaptcha_pt")
	req.Header.Set("HX-Target", "target")
	req.Header.Set("HX-Current-URL", "https://ssstik.io/")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://ssstik.io")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://ssstik.io/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Priority", "u=0")
	req.Header.Set("TE", "trailers")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	downloadURL := s.extractDownloadURL(string(bodyText))
	if downloadURL == nil {
		return nil
	}
	key := s.extractKey(*downloadURL)
	if key == nil {
		return nil
	}
	return &model.DownloadData{
		URL: *downloadURL,
		Key: *key,
	}
}

func (s *TiktokDownloader) processDownloadData(downloadData *model.DownloadData) *string {
	client := &http.Client{}
	body := fmt.Sprintf("tt=%s", downloadData.Key)
	data := strings.NewReader(body)
	url := fmt.Sprintf("https://ssstik.io%s", downloadData.URL)
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:136.0) Gecko/20100101 Firefox/136.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("HX-Request", "true")
	req.Header.Set("HX-Trigger", "hd_download")
	req.Header.Set("HX-Target", "hd_download")
	req.Header.Set("HX-Current-URL", "https://ssstik.io/")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://ssstik.io")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://ssstik.io/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Priority", "u=0")
	req.Header.Set("TE", "trailers")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	downloadURL := resp.Header.Get("Hx-Redirect")
	if downloadURL == "" {
		return nil
	}
	return &downloadURL
}

func (s *TiktokDownloader) Download(link string) (*string, error) {
	tt := s.GetTT()
	downloadData := s.getDownloadData(link, tt)
	if downloadData == nil {
		return nil, errors.New("failed to fetch download data")
	}
	downloadURL := s.processDownloadData(downloadData)
	if downloadURL == nil {
		return nil, errors.New("failed to fetch download url")
	}
	return downloadURL, nil
}
