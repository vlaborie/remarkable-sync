package main

import (
    "io/ioutil"
    "os"
    "fmt"
    "time"
    "encoding/json"

    "github.com/vlaborie/reMarkable-sync/remarkable"
    "github.com/vlaborie/reMarkable-sync/wallabag"
)

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

    var Remarkable = remarkable.New(".local/share/remarkable/xochitl/")

    dir := Remarkable.AddDir("wallabag", "Wallabag", "")
    j, _ := json.Marshal(dir)
    _ = ioutil.WriteFile(Remarkable.Dir+"wallabag.metadata", j, 0644)

    for _, WallabagItem := range Wallabag.Items {
        var RemarkableItem remarkable.RemarkableItem
        RemarkableItem.FromWallabag(WallabagItem)

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
