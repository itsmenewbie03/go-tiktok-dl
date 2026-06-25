package downloader

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
	"github.com/itsmenewbie03/go-tiktok-dl/downloader/model"
	"github.com/itsmenewbie03/go-tiktok-dl/utils"
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

func (s *TiktokDownloader) extractDirectURL(body string) *string {
	// safer regex
	pattern := `data-directurl="([^"]+)"`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(body)
	if len(matches) < 2 {
		return nil
	}

	return &matches[1]
}

func (s *TiktokDownloader) extractLowResURL(body []byte) *string {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	dlURL, ok := doc.Find(".download_link.without_watermark").Attr("href")
	if !ok {
		return nil
	}
	return &dlURL
}

func (s *TiktokDownloader) extractDownloadURL(body string) *string {
	pattern := `data-directurl="(.*)"`
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
	pattern := `s_tt\s=\s'(\w+)',`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(body)
	if matches == nil {
		return nil
	}
	return &matches[1]
}

func (s *TiktokDownloader) getDownloadData(link, tt string) *model.DownloadData {
	client := &http.Client{}
	dbg := fmt.Sprintf("ab=1&loc=US&ip=%s", utils.RandomUSIP())
	// TODO: make tt query dynamic
	body := fmt.Sprintf("id=%s&locale=en&tt=%s&debug=%s", url.QueryEscape(link), tt, url.QueryEscape(dbg))
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
	withoutWaterMarkURL := s.extractLowResURL(bodyText)
	if withoutWaterMarkURL == nil {
		log.Fatal("failed to get url low res url")
	}
	directURL := s.extractDirectURL(string(bodyText))
	if directURL != nil {
		return &model.DownloadData{
			DirectURL:        *directURL,
			WithoutWaterMark: *withoutWaterMarkURL,
		}
	}
	downloadURL := s.extractDownloadURL(string(bodyText))
	if downloadURL == nil {
		log.Print("downloadURL not found")
		return nil
	}
	key := s.extractKey(*downloadURL)
	if key == nil {
		log.Print("key not found")
		return nil
	}
	return &model.DownloadData{
		URL: *downloadURL,
		Key: *key,
	}
}

func (s *TiktokDownloader) extractPrefix(url string) *string {
	// look for `url=` and stop at the first "=="
	pattern := `url=([^=]+)==`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return nil
	}
	return &matches[1]
}

func (s *TiktokDownloader) processDownloadData(downloadData *model.DownloadData) *string {
	client := &http.Client{}
	var url string
	var data *strings.Reader
	if downloadData.DirectURL != "" {
		key := s.extractPrefix(downloadData.DirectURL)
		body := fmt.Sprintf("tt=%s", *key)
		data = strings.NewReader(body)
		url = fmt.Sprintf("https://ssstik.io%s", downloadData.DirectURL)
	} else {
		log.Warn("DirectURL not found, falling back (this seems be an outdated branch, keeping for compatibility for now)")
		body := fmt.Sprintf("tt=%s", downloadData.Key)
		data = strings.NewReader(body)
		url = fmt.Sprintf("https://ssstik.io%s", downloadData.URL)
	}
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

func (s *TiktokDownloader) Download(link string) (*model.DownloadOptions, error) {
	tt := s.GetTT()
	downloadData := s.getDownloadData(link, tt)
	if downloadData == nil {
		return nil, errors.New("failed to fetch download data")
	}
	downloadURL := s.processDownloadData(downloadData)
	if downloadURL == nil {
		return nil, errors.New("failed to fetch download url")
	}
	dlURL, err := url.Parse(*downloadURL)
	if err != nil {
		return nil, err
	}
	lowResURL, err := url.Parse(downloadData.WithoutWaterMark)
	if err != nil {
		return nil, err
	}
	return &model.DownloadOptions{
		HDWithoutWaterMark: *dlURL,
		WithoutWaterMark:   *lowResURL,
	}, nil
}
