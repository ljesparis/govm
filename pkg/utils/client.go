package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"sync"

	"github.com/dustin/go-humanize"
	"golang.org/x/sync/errgroup"
)

var (
	// ErrSourceDoesNotExists represents the error when the remote file does not exists
	ErrSourceDoesNotExists = errors.New("remote file does not exists")
)

// Stats represents the download stats
type Stats struct {
	BytesDownloaded      string
	BytesToDownload      string
	DownloadedPercentage string
}

// callback used to notify others
// download progress
type dCallback func(*Stats)

// tracker example was taken
// from https://golangcode.com/download-a-file-with-progress/
type tracker struct {
	// remote file size
	t int64

	// bytes downloaded
	b int64

	cb dCallback
}

func (wc *tracker) Write(p []byte) (int, error) {
	n := len(p)
	wc.b += int64(n)

	if wc.cb != nil {
		wc.cb(&Stats{
			DownloadedPercentage: fmt.Sprintf("%d",
				int(math.Floor((float64(wc.b)*100)/(float64(wc.t))))) + "%",
			BytesToDownload: humanize.IBytes(uint64(wc.t)),
			BytesDownloaded: humanize.IBytes(uint64(wc.b)),
		})
	}

	return n, nil
}

// Client represents the client used to download files
type Client struct {
	// remote file
	srcFile string

	// local file
	dstFile string

	m      *sync.Mutex
	client *http.Client
}

// DefaultClient will create the default client to download files,
// only receive source file and destiny file.
func DefaultClient(srcFile, dstFile string) *Client {
	return &Client{
		m:       &sync.Mutex{},
		srcFile: srcFile,
		dstFile: dstFile,
		client:  http.DefaultClient,
	}
}

// Download will download files and show progress
func (d *Client) Download(cb dCallback) error {
	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		d.m.Lock()
		defer d.m.Unlock()

		response, err := d.client.Get(d.srcFile)
		if err != nil {
			return err
		}

		if response.StatusCode == http.StatusNotFound {
			return ErrSourceDoesNotExists
		}

		buffer := bytes.NewBuffer([]byte{})
		tracker := &tracker{
			cb: cb,
			t:  response.ContentLength,
		}

		if _, err = io.Copy(buffer, io.TeeReader(response.Body, tracker)); err != nil {
			return err
		}

		err = ioutil.WriteFile(d.dstFile, buffer.Bytes(), 0777)
		if err != nil {
			return err
		}

		return response.Body.Close()
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
