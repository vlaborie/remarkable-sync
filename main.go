package main

import (
    "io/ioutil"
    "os"
    "fmt"
    "strconv"
    "time"
    "encoding/json"

    "github.com/vlaborie/reMarkable-sync/wallabag"
)

type Remarkable struct {
    Dir string
    Items []RemarkableItem
}

type RemarkableItem struct {
    Id string
    Type string `json:"type"`
    Parent string `json:"parent"`
    VisibleName string `json:"visibleName"`
    LastModified string `json:"lastModified"`
    Version int64 `json:"version"`
    Deleted bool `json:"deleted"`
    MetadataModified bool `json:"metadataModified"`
    Modified bool `json:"modified"`
    Pinned bool `json:"pinned"`
    Synced bool `json:"synced"`
    Content []byte
}

func (Remarkable *Remarkable) addDir(id string, name string, parent string) RemarkableItem {
    var RemarkableItem = RemarkableItem {
        Id: id,
        Type: "CollectionType",
        Parent: parent,
        VisibleName: name,
        LastModified: "",
        Version: 1,
        Deleted: false,
        MetadataModified: false,
        Modified: false,
        Pinned: false,
        Synced: false,
    }
    return RemarkableItem
}

func (RemarkableItem *RemarkableItem) fromWallabag(WallabagItem wallabag.WallabagItem) {
    RemarkableItem.Id = strconv.FormatInt(WallabagItem.Id, 10)
    RemarkableItem.Type = "DocumentType"
    RemarkableItem.Parent = "wallabag"
    RemarkableItem.VisibleName = WallabagItem.Title
    RemarkableItem.LastModified = WallabagItem.UpdatedAt
    RemarkableItem.Version = 1
    RemarkableItem.Deleted = false
    RemarkableItem.MetadataModified = false
    RemarkableItem.Modified = false
    RemarkableItem.Pinned = WallabagItem.IsStarred
    RemarkableItem.Synced = false
}

func indicator(channel <-chan struct{}) {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            fmt.Print(".")
        case <-channel:
            return
        }
    }
}

func main() {
    Wallabag := wallabag.New(".config/reMarkable-sync/wallabag.json")
    Wallabag.GetPages(1)

    var Remarkable = Remarkable {
        Dir: ".local/share/remarkable/xochitl/",
    }

    dir := Remarkable.addDir("wallabag", "Wallabag", "")
    j, _ := json.Marshal(dir)
    _ = ioutil.WriteFile(Remarkable.Dir+"wallabag.metadata", j, 0644)

    for _, WallabagPage := range Wallabag.Pages {
        for _, WallabagItem := range WallabagPage.Embedded.Items {
            var RemarkableItem RemarkableItem
            RemarkableItem.fromWallabag(WallabagItem)

            if _, err := os.Stat(Remarkable.Dir+RemarkableItem.Id+".epub"); os.IsNotExist(err) {
                fmt.Print("Get EPUB of Wallabag element "+RemarkableItem.Id+" => ")
                channel := make(chan struct{})
                go indicator(channel)
                RemarkableItem.Content = Wallabag.GetEpub(RemarkableItem.Id)
                close(channel)
                fmt.Print("done\n")
            }

            j, _ := json.Marshal(RemarkableItem)
            _ = ioutil.WriteFile(Remarkable.Dir+RemarkableItem.Id+".metadata", j, 0644)
            fmt.Println("Metadata of "+RemarkableItem.Id+" updated")

            if len(RemarkableItem.Content) > 0 {
                _ = ioutil.WriteFile(Remarkable.Dir+RemarkableItem.Id+".epub", RemarkableItem.Content, 0644)
                fmt.Println("EPUB of "+RemarkableItem.Id+" writed")
            }
        }
    }
}
