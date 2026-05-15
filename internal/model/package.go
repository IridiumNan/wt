// Package model: write all structs for client and server usage
// includes
// Package : meta data for package
// ClientConfig : store and handle the client config
// ServerConfig : store and handle the server config
package model

import "time"

// Package : meta data for package which store in memory for fast check
type Package struct {
	// package name
	Name string `json:"string"`

	// package tag
	Tag string `json:"tag"`

	// package size
	Size int64 `json:"size"`

	// modify time
	ModTime time.Time `json:"mod_time"`
}
