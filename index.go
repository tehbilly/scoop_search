package scoop_search

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SearcherOpts struct {
	Verbose bool
}

type ScoopIndex struct {
	index bleve.Index
}

func NewPersistentSearcher(opts SearcherOpts) (*ScoopIndex, error) {
	home, _ := os.UserHomeDir()

	// TODO: Determine if index exists, and age? If too old recreate it?
	op := "open"
	start := time.Now()

	idx, err := bleve.Open(filepath.Join(home, "scoop", ".search-index"))
	if err != nil && errors.Is(err, bleve.ErrorIndexPathDoesNotExist) {
		op = "create"
		idx, err = bleve.New(filepath.Join(home, "scoop", ".search-index"), indexMapping())
		if err != nil {
			return nil, err
		}
	}

	if opts.Verbose {
		fmt.Printf("Time to %s index: %s\n", op, time.Since(start))
	}

	if op == "create" {
		start = time.Now()
		if err := indexManifests(idx); err != nil {
			return nil, err
		}
		if opts.Verbose {
			fmt.Printf("Time to index manifests: %s\n", time.Since(start))
		}
	}

	return &ScoopIndex{index: idx}, nil
}

func NewSearcher(opts SearcherOpts) (*ScoopIndex, error) {
	start := time.Now()
	idx, err := bleve.NewMemOnly(indexMapping())
	if err != nil {
		return nil, err
	}
	if opts.Verbose {
		fmt.Printf("Time to create mem-only index: %s\n", time.Since(start))
	}

	start = time.Now()
	if err := indexManifests(idx); err != nil {
		return nil, err
	}
	if opts.Verbose {
		fmt.Printf("Time to index manifests: %s\n", time.Since(start))
	}

	return &ScoopIndex{index: idx}, nil
}

func (i ScoopIndex) Search(opts SearchOptions) ([]*AppManifest, error) {
	req := bleve.NewSearchRequestOptions(bleve.NewMatchQuery(opts.Query), opts.Num, 0, false)
	req.Fields = []string{"name", "bucket", "version", "description"}

	results, err := i.index.Search(req)
	if err != nil {
		return nil, err
	}

	if opts.Verbose && results.Total > uint64(len(results.Hits)) {
		fmt.Printf("Displaying %d of %d results.\n", len(results.Hits), results.Total)
	}

	var manifests []*AppManifest

	for _, hit := range results.Hits {
		am, err := AppManifestFromDocumentMatch(hit)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Unable to get AppManifest from result doc: %s\n", err)
			continue
		}
		manifests = append(manifests, am)
	}

	return manifests, nil
}

func indexManifests(index bleve.Index) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	scoopPath := filepath.Join(home, "scoop")

	if _, err := os.Stat(scoopPath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Scoop does not appear to be installed, cannot find directory: %s\n", scoopPath)
		return err
	}

	buckets, err := getBucketDirs(scoopPath)
	if err != nil {
		return err
	}

	for _, bucket := range buckets {
		bucketName := filepath.Base(strings.TrimSuffix(bucket, "bucket"))

		mfs, err := readBucket(bucket)
		if err != nil {
			return err
		}

		batch := index.NewBatch()
		for _, mf := range mfs {
			if err := batch.Index(mf.Name, mf); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Unable to index %s: %s\n", mf.Name, err)
				continue
			}
		}
		if err := index.Batch(batch); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Unable to index %s batch :%s\n", bucketName, err)
			continue
		}
	}

	return nil
}

func getBucketDirs(scoopPath string) ([]string, error) {
	bucketsRoot := filepath.Join(scoopPath, "buckets")

	entries, err := os.ReadDir(bucketsRoot)
	if err != nil {
		return nil, err
	}

	var dirs []string

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		bd := filepath.Join(bucketsRoot, entry.Name(), "bucket")
		if _, err := os.Stat(bd); errors.Is(err, os.ErrNotExist) {
			continue
		}

		dirs = append(dirs, bd)
	}

	return dirs, nil
}

func readBucket(bucket string) ([]*AppManifest, error) {
	bucketName := filepath.Base(strings.TrimSuffix(bucket, "bucket"))

	entries, err := os.ReadDir(bucket)
	if err != nil {
		return nil, err
	}

	var manifests []*AppManifest

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			fmt.Printf("Skipping entry: %#v\n", entry)
			continue
		}

		fileName := filepath.Join(bucket, entry.Name())
		fb, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Printf("Unable to read '%s': %s\n", fileName, err)
			continue
		}

		var mf AppManifest
		if err := json.Unmarshal(fb, &mf); err != nil {
			fmt.Printf("Unable to unmarshal '%s': %s\n", fileName, err)
			continue
		}

		mf.Name = strings.TrimSuffix(entry.Name(), ".json")
		mf.Bucket = bucketName

		manifests = append(manifests, &mf)
	}

	return manifests, nil
}

type SearchOptions struct {
	Verbose bool
	Query   string
	Num     int
}
