package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/vlaborie/reMarkable-sync/remarkable"

	"github.com/bmaupin/go-epub"
)

func main() {
	Remarkable := remarkable.New(".local/share/remarkable/xochitl/")

	wallabagConfig, err := os.Open(".config/reMarkable-sync/wallabag.json")
	if err == nil {
		Remarkable.Wallabag(wallabagConfig)
	}

	minifluxConfig, err := os.Open(".config/reMarkable-sync/miniflux.json")
	if err == nil {
		Remarkable.Miniflux(minifluxConfig)
	}

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
		fmt.Println("Metadata of " + RemarkableItem.Id + " updated")
	}
}
