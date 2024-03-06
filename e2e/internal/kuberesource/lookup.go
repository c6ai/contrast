package kuberesource

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

type Image struct {
	Registry string
	Path     string
	Tag      string
	Hash     string
}

type ImageLookup struct {
	m map[string]Image
}

var imageReg = regexp.MustCompile(`(?P<registry>[^/]+/)?(?P<path>[^:]+)(:(?P<tag>[^@]+))?(@(?P<hash>.*))?`)

func NewImageLookupFromFile(file io.ReadCloser) (*ImageLookup, error) {
	lookup := &ImageLookup{m: make(map[string]Image)}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		matches := imageReg.FindStringSubmatch(line)
		if matches == nil {
			return nil, fmt.Errorf("invalid image line: %s", line)
		}

		img := Image{}
		if imageReg.SubexpIndex("registry") == -1 {
			return nil, fmt.Errorf("registry not found for image line: %s", line)
		}
		img.Registry = matches[imageReg.SubexpIndex("registry")]
		if imageReg.SubexpIndex("path") == -1 {
			return nil, fmt.Errorf("path not found for image line: %s", line)
		}
		img.Path = matches[imageReg.SubexpIndex("path")]
		if imageReg.SubexpIndex("tag") != -1 {
			img.Tag = matches[imageReg.SubexpIndex("tag")]
		}
		if imageReg.SubexpIndex("hash") != -1 {
			img.Hash = matches[imageReg.SubexpIndex("hash")]
		}
		img.Hash = matches[imageReg.SubexpIndex("hash")]

		lookup.m[img.Path] = img
	}

	return lookup, nil
}

func (l *ImageLookup) Lookup(path string) (string, bool) {
	img, ok := l.m[path]
	if !ok {
		return "", false
	}

	return fmt.Sprintf("%s/%s:%s@%s", img.Registry, img.Path, img.Tag, img.Hash), true
}
