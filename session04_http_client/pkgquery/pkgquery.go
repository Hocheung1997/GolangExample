package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type pkgData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func fetchPackageData(url string) ([]pkgData, error) {
	var packages []pkgData
	packagesPtr := &packages
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.Header.Get("Content-Type") != "application/json" {
		return packages, nil
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return packages, err
	}
	err = json.Unmarshal(data, packagesPtr)
	if err != nil {
		return packages, err
	}
	return packages, nil
}
