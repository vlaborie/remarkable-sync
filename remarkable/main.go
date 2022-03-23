package remarkable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/bmaupin/go-epub"
)

type Remarkable struct {
	Dir   string
	Items []RemarkableItem
}

type RemarkableItem struct {
	Id               string `json:"-"`
	Type             string `json:"type"`
	Parent           string `json:"parent"`
	VisibleName      string `json:"visibleName"`
	LastModified     string `json:"lastModified"`
	Version          int64  `json:"version"`
	Deleted          bool   `json:"deleted"`
	MetadataModified bool   `json:"metadataModified"`
	Modified         bool   `json:"modified"`
	Pinned           bool   `json:"pinned"`
	Synced           bool   `json:"synced"`
	ContentType      string `json:"-"`
	Content          []byte `json:"-"`
}

func New(dir string) Remarkable {
	var remarkable Remarkable
	remarkable.Dir = dir
	return remarkable
}

func (Remarkable *Remarkable) AddDir(id string, name string, parent string) RemarkableItem {
	var RemarkableItem = RemarkableItem{
		Id:               id,
		Type:             "CollectionType",
		Parent:           parent,
		VisibleName:      name,
		LastModified:     "",
		Version:          1,
		Deleted:          false,
		MetadataModified: false,
		Modified:         false,
		Pinned:           false,
		Synced:           false,
	}
	return RemarkableItem
}

func (Remarkable *Remarkable) Write() {
	for _, RemarkableItem := range Remarkable.Items {
		if RemarkableItem.ContentType == "html" {
			RemarkableItem.ContentType = "epub"
			e := epub.NewEpub(RemarkableItem.VisibleName)
			e.AddSection(string(RemarkableItem.Content), "Section 1", "", "")
			e.Write(Remarkable.Dir + RemarkableItem.Id + "." + RemarkableItem.ContentType)
			fmt.Println("EPUB of " + RemarkableItem.Id + " writed")
		}

		j, _ := json.Marshal(RemarkableItem)
		_ = ioutil.WriteFile(Remarkable.Dir+RemarkableItem.Id+".metadata", j, 0644)
		_ = ioutil.WriteFile(Remarkable.Dir+RemarkableItem.Id+".content", []byte("{}"), 0644)
		fmt.Println("Metadata of " + RemarkableItem.Id + " updated")
	}
}
