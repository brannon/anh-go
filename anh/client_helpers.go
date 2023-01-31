package anh

import (
	"encoding/xml"
	"io"
	"net/url"
)

func appendQueryString(u *url.URL, name string, value string) {
	q := u.Query()
	q.Add(name, value)
	u.RawQuery = q.Encode()
}

func parseAtomEntry(r io.Reader) (*atomEntry, error) {
	decoder := xml.NewDecoder(r)

	var entry atomEntry
	err := decoder.Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func parseAtomFeed(r io.Reader) (*atomFeed, error) {
	decoder := xml.NewDecoder(r)

	var feed atomFeed
	err := decoder.Decode(&feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}
