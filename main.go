package main

import (
    "io/ioutil"
    "os"
    "fmt"
    "strconv"
    "time"
    "encoding/json"

    "github.com/vlaborie/reMarkable-sync/services"
)

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

func (RemarkableItem *RemarkableItem) fromWallabag(WallabagItem services.WallabagItem) {
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

const output = ".local/share/remarkable/xochitl/"

func main() {
    Wallabag := services.NewWallabag(".config/reMarkable-sync/wallabag.json")
    Wallabag.Login()
    Wallabag.GetPages(1)

    j, _ := json.Marshal(RemarkableItem {
        Id: "wallabag",
        Type: "CollectionType",
        Parent: "",
        VisibleName: "Wallabag",
        LastModified: "",
        Version: 1,
        Deleted: false,
        MetadataModified: false,
        Modified: false,
        Pinned: false,
        Synced: false,
    })
    _ = ioutil.WriteFile(output+"wallabag.metadata", j, 0644)

    for _, WallabagPage := range Wallabag.Pages {
        for _, WallabagItem := range WallabagPage.Embedded.Items {
            var RemarkableItem RemarkableItem
            RemarkableItem.fromWallabag(WallabagItem)

            if _, err := os.Stat(output+RemarkableItem.Id+".epub"); os.IsNotExist(err) {
                fmt.Print("Get EPUB of Wallabag element "+RemarkableItem.Id+" => ")
                channel := make(chan struct{})
                go indicator(channel)
                RemarkableItem.Content = Wallabag.GetEpub(RemarkableItem.Id)
                close(channel)
                fmt.Print("done\n")
            }

            j, _ := json.Marshal(RemarkableItem)
            _ = ioutil.WriteFile(output+RemarkableItem.Id+".metadata", j, 0644)
            fmt.Println("Metadata of "+RemarkableItem.Id+" updated")

            if len(RemarkableItem.Content) > 0 {
                _ = ioutil.WriteFile(output+RemarkableItem.Id+".epub", RemarkableItem.Content, 0644)
                fmt.Println("EPUB of "+RemarkableItem.Id+" writed")
            }
        }
    }
}
