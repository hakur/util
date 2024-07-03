package network

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHttpDownload(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*600)
	defer cancel()

	dl, err := NewHttpDownloader(&HttpDownloaderOpts{
		// ProxyAddr: "http://localhost:1082",
	})
	assert.Equal(t, nil, err)

	var checksumFile = "b:/kubernetes-node-linux-arm64.tar.gz.sha256"
	err = dl.Download(ctx, &DownloadOpts{FileURL: "https://dl.k8s.io/v1.23.3/kubernetes-node-linux-arm64.tar.gz.sha256", OutputFilename: checksumFile, Checksum: ""})
	assert.Equal(t, nil, err)

	buf, err := os.ReadFile(checksumFile)
	assert.Equal(t, nil, err)
	err = dl.Download(ctx, &DownloadOpts{FileURL: "https://dl.k8s.io/v1.23.3/kubernetes-node-linux-arm64.tar.gz", OutputFilename: "b:/kubernetes-node-linux-arm64.tar.gz", Checksum: "sha256:" + string(buf)})
	assert.Equal(t, nil, err)
}
