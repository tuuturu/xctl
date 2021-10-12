package binaries

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
)

func Download(opts DownloadOpts) (string, error) {
	binaryDir := path.Join(opts.BinariesDir, opts.Name, opts.Version)
	binaryPath := path.Join(binaryDir, opts.Name)

	_, err := opts.Fs.Stat(binaryPath)
	if err == nil {
		return binaryPath, nil
	}

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.URL.Opaque = req.URL.Path

			return nil
		},
	}

	resp, err := client.Get(opts.URL)
	if err != nil {
		return "", fmt.Errorf("fetching file: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var content io.Reader
	content = resp.Body

	for _, unpackFn := range opts.UnpackingFn {
		content, err = unpackFn(content)
		if err != nil {
			return "", fmt.Errorf("unpacking: %w", err)
		}
	}

	err = opts.Fs.MkdirAll(binaryDir, 0o700)
	if err != nil {
		return "", fmt.Errorf("preparing binary directory: %w", err)
	}

	f, err := opts.Fs.Create(binaryPath)
	if err != nil {
		return "", fmt.Errorf("creating binary file: %w", err)
	}

	err = opts.Fs.Chmod(binaryPath, 0o700)
	if err != nil {
		return "", fmt.Errorf("setting executable permission: %w", err)
	}

	_, err = io.Copy(f, content)
	if err != nil {
		return "", fmt.Errorf("storing binary file: %w", err)
	}

	return binaryPath, nil
}

func GenerateZipUnpacker(target string) UnpackingFn {
	return func(r io.Reader) (io.Reader, error) {
		raw, err := io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("reading zip contents: %w", err)
		}

		inputAt := bytes.NewReader(raw)

		result, err := zip.NewReader(inputAt, int64(len(raw)))
		if err != nil {
			return nil, fmt.Errorf("unzipping: %w", err)
		}

		buf := bytes.Buffer{}

		f, err := result.Open(target)
		if err != nil {
			return nil, fmt.Errorf("opening target: %w", err)
		}

		defer func() {
			_ = f.Close()
		}()

		_, err = io.Copy(&buf, f)
		if err != nil {
			return nil, fmt.Errorf("extracting target: %w", err)
		}

		return &buf, nil
	}
}

func GzipUnpacker(input io.Reader) (io.Reader, error) {
	ginput, err := gzip.NewReader(input)
	if err != nil {
		return nil, fmt.Errorf("reading gzip format: %w", err)
	}

	return ginput, nil
}

func GenerateTarUnpacker(target string) UnpackingFn {
	return func(r io.Reader) (io.Reader, error) {
		tr := tar.NewReader(r)

		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}

			if err != nil {
				return nil, fmt.Errorf("acquiring header: %w", err)
			}

			if strings.HasSuffix(hdr.Name, target) {
				buf := bytes.Buffer{}

				_, err = io.Copy(&buf, tr)
				if err != nil {
					return nil, fmt.Errorf("reading tar format: %w", err)
				}

				return &buf, nil
			}
		}

		return nil, fmt.Errorf("couldn't find target")
	}
}
