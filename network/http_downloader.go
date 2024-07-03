package network

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var ErrUnknownChecksumType = errors.New("unknown checksum type")

func NewHttpDownloader(opts *HttpDownloaderOpts) (t *HttpDownloader, err error) {
	var transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if opts.ProxyAddr != "" {
		proxyURL, err := url.Parse(opts.ProxyAddr)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	client := &http.Client{
		Transport: transport,
	}
	t = new(HttpDownloader)
	t.Opts = opts
	t.Client = client
	return
}

type HttpDownloaderOpts struct {
	ProxyAddr string
}

type DownloadOpts struct {
	FileURL        string
	OutputFilename string
	Checksum       string
	TotalSizeChan  chan int64
	Cookie         string
}

type HttpDownloader struct {
	Opts           *HttpDownloaderOpts
	Client         *http.Client
	CopiedCallback func(bytesCount int)
	// WorkerCount    int
}

func (t *HttpDownloader) Download(ctx context.Context, opts *DownloadOpts) (err error) {
	f, fi, err := t.prepareOutputFile(opts.OutputFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	if opts.Checksum != "" {
		verifyOk, err := t.verifyChecksum(f, opts.Checksum)
		if verifyOk {
			return nil
		} else if errors.Is(err, ErrUnknownChecksumType) {
			return err
		}
	}

	req, err := t.prepareReuqest(ctx, opts.FileURL, fi.Size(), true)
	if err != nil {
		return
	}

	resp, err := t.Client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode == 501 || resp.StatusCode != 206 { // 501 Unsupported client range
		resp.Body.Close()
		if err = f.Truncate(0); err != nil {
			return err
		}

		req, err = t.prepareReuqest(ctx, opts.FileURL, fi.Size(), false)
		if err != nil {
			return
		}
		resp, err = t.Client.Do(req)
		if err != nil {
			return
		}
	}
	defer resp.Body.Close()

	contentLength, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if opts.TotalSizeChan != nil {
		opts.TotalSizeChan <- contentLength
	}

	errChan := make(chan error)
	doneChan := make(chan struct{})

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				if t.CopiedCallback != nil {
					t.CopiedCallback(n)
				}

				if _, err = f.Write(buf[:n]); err != nil {
					errChan <- err
					return
				}
			}

			if err != nil {
				if err != io.EOF {
					errChan <- err
					return
				} else {
					doneChan <- struct{}{}
					return
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err = <-errChan:
		return err
	case <-doneChan:
		return
	}
}

func (t *HttpDownloader) prepareOutputFile(filename string) (f *os.File, fi os.FileInfo, err error) {
	f, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return
	}

	fi, err = f.Stat()
	if err != nil {
		return
	}
	return
}

func (t *HttpDownloader) prepareReuqest(ctx context.Context, fileURL string, fileSize int64, useRange bool) (req *http.Request, err error) {
	req, err = http.NewRequestWithContext(ctx, "GET", fileURL, nil)
	if err != nil {
		return
	}
	if useRange {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d", fileSize))
	}
	return
}

func (t *HttpDownloader) verifyChecksum(r io.Reader, cheksum string) (ok bool, err error) {
	hashType := strings.Split(cheksum, ":")[0]
	var hasher hash.Hash
	switch hashType {
	case "sha256":
		hasher = sha256.New()
	default:
		return false, fmt.Errorf("%w %s", ErrUnknownChecksumType, hashType)
	}

	if _, err := io.Copy(hasher, r); err != nil {
		return false, err
	}

	sum := hasher.Sum(nil)

	readerChecksum := hashType + ":" + hex.EncodeToString(sum)
	return readerChecksum == cheksum, nil
}
