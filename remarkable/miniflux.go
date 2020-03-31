package remarkable

import (
    "os"
    "io/ioutil"
    "encoding/json"
    "strconv"

    miniflux "miniflux.app/client"
)

type MinifluxConfig struct {
    Host string `json:"host"`
    Token string `json:"token"`
}

func (Remarkable *Remarkable) Miniflux(path string) {
    Remarkable.Items = append(Remarkable.Items, Remarkable.AddDir("miniflux", "Miniflux", ""))

    jsonFile, _ := os.Open(path)
    byteValue, _ := ioutil.ReadAll(jsonFile)
    var config MinifluxConfig
    json.Unmarshal(byteValue, &config)

    Miniflux := miniflux.New("https://"+config.Host, config.Token)
    m, _ := Miniflux.Entries(&miniflux.Filter {})
    for _, MinifluxItem := range m.Entries {
        var RemarkableItem RemarkableItem
        RemarkableItem.FromMiniflux(MinifluxItem)
        Remarkable.Items = append(Remarkable.Items, RemarkableItem)
    }
}

func (RemarkableItem *RemarkableItem) FromMiniflux(MinifluxItem *miniflux.Entry) {
    RemarkableItem.Id = strconv.FormatInt(MinifluxItem.ID, 10)
    RemarkableItem.Type = "DocumentType"
    RemarkableItem.Parent = "miniflux"
    RemarkableItem.VisibleName = MinifluxItem.Title
    RemarkableItem.LastModified = ""
    RemarkableItem.Version = 1
    RemarkableItem.Deleted = false
    RemarkableItem.MetadataModified = false
    RemarkableItem.Modified = false
    RemarkableItem.Pinned = MinifluxItem.Starred
    RemarkableItem.Synced = false
    RemarkableItem.ContentType = "html"
    RemarkableItem.Content = []byte(MinifluxItem.Content)
}