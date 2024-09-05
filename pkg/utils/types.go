package utils

import (
	"regexp"
)

// HTTP REST STRUCTS

// Response struct for HTTP Requests
type Response struct {
	Url    string
	Method string
	Status int
	Data   string
	Error  string
}

// PARSER INPUT STRUCTURES

// Config structure containing parser and metadata
type Config struct {
	Id      Tag      `json:"id"`
	Sources []string `json:"sources"`
	Parser  []Parser `json:"parser"`
}

// Parser structure contains post compiled regexp parsing template
type Parser struct {
	Name  string  `json:"name"`
	Match []Match `json:"match"`
}

// Match struct containing pre & post compiled regexp
type Match struct {
	Name  string        `json:"name"`
	Value regexp.Regexp `json:"value"`
}

// DATA OUTPUT STRUCTURES

// Final output contains a slice of key value tags and a collection of key value tags
type Collection struct {
	Tags    []Tag
	Buckets []Bucket
}

// Bucket struct contains array of key value pairs
type Bucket struct {
	Name  string
	Value [][]Tag
}

// Tag struct for key value data storage
type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Image struct {
	Make         string
	Model        string
	Lens         string
	LensModel    string
	FocalLength  string
	FileName     string
	FilePath     string
	FileNumber   string
	FileType     string
	FileSize     string
	ISO          string
	ShutterSpeed string
	Aperture     string
	Megapixels   string
	Date         string
	Year         int
	CreateDate   string
	ModifyDate   string
}
