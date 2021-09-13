package scoop_search

import (
	"errors"
	"github.com/blevesearch/bleve/v2/search"
)

type AppManifest struct {
	Name        string `json:"name"`
	Bucket      string `json:"bucket"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

func (a AppManifest) Type() string {
	return "scoop_manifest"
}

func AppManifestFromDocumentMatch(r *search.DocumentMatch) (*AppManifest, error) {
	name, ok := r.Fields["name"].(string)
	if !ok {
		return nil, errors.New("missing field: name")
	}
	bucket, ok := r.Fields["bucket"].(string)
	if !ok {
		return nil, errors.New("missing field: bucket")
	}
	version, ok := r.Fields["version"].(string)
	if !ok {
		return nil, errors.New("missing field: version")
	}
	desc, ok := r.Fields["description"].(string)
	if !ok {
		return nil, errors.New("missing field: description")
	}

	return &AppManifest{Name: name, Bucket: bucket, Version: version, Description: desc}, nil
}
