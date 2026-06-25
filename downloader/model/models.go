// Package model
package model

import "net/url"

type DownloadData struct {
	URL              string
	Key              string
	DirectURL        string
	WithoutWaterMark string
}

type DownloadOptions struct {
	HDWithoutWaterMark url.URL
	WithoutWaterMark   url.URL
}
