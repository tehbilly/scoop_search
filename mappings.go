package scoop_search

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
)

func indexMapping() *mapping.IndexMappingImpl {
	// generic english text field
	englishTextMapping := bleve.NewTextFieldMapping()
	englishTextMapping.Analyzer = en.AnalyzerName

	// non-indexed data field
	dataFieldMapping := bleve.NewTextFieldMapping()
	dataFieldMapping.Store = true
	dataFieldMapping.Index = false
	dataFieldMapping.IncludeTermVectors = false
	dataFieldMapping.IncludeInAll = false

	// keyword mapping
	keywordMapping := bleve.NewTextFieldMapping()
	keywordMapping.Analyzer = keyword.Name

	mfMapping := bleve.NewDocumentMapping()
	mfMapping.AddFieldMappingsAt("name", englishTextMapping)
	mfMapping.AddFieldMappingsAt("bucket", keywordMapping)
	mfMapping.AddFieldMappingsAt("version", dataFieldMapping)
	mfMapping.AddFieldMappingsAt("description", englishTextMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultMapping = mfMapping

	return indexMapping
}
