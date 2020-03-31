package main

import (
    "io/ioutil"
    "fmt"
    "encoding/json"

    "github.com/vlaborie/reMarkable-sync/remarkable"
    "github.com/vlaborie/reMarkable-sync/wallabag"

    "github.com/bmaupin/go-epub"
)

func main() {
    Wallabag := wallabag.New(".config/reMarkable-sync/wallabag.json")

    Remarkable := remarkable.New(".local/share/remarkable/xochitl/")
    dir := Remarkable.AddDir("wallabag", "Wallabag", "")
    Remarkable.Items = append(Remarkable.Items, dir)

    for _, WallabagItem := range Wallabag.Items {
        var RemarkableItem remarkable.RemarkableItem
        RemarkableItem.FromWallabag(WallabagItem)
        Remarkable.Items = append(Remarkable.Items, RemarkableItem)
    }

    for _, RemarkableItem := range Remarkable.Items {
        if RemarkableItem.ContentType == "html" {
            RemarkableItem.ContentType = "epub"
            e := epub.NewEpub(RemarkableItem.VisibleName)
            e.AddSection(string(RemarkableItem.Content), "Section 1", "", "")
            e.Write(Remarkable.Dir+RemarkableItem.Id+"."+RemarkableItem.ContentType)
            fmt.Println("EPUB of "+RemarkableItem.Id+" writed")
        }

        j, _ := json.Marshal(RemarkableItem)
        _ = ioutil.WriteFile(Remarkable.Dir+RemarkableItem.Id+".metadata", j, 0644)
        fmt.Println("Metadata of "+RemarkableItem.Id+" updated")
    }
}
