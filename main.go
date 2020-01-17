package main

import (
    "bytes"
    "io/ioutil"
    "os"
    "fmt"
    "strconv"
    "time"
    "encoding/json"
    "net/http"
)

type Wallabag struct {
    Config WallabagConfig
    Token WallabagToken
    Pages []WallabagPage
}

type WallabagConfig struct {
    GrantType string `json:"grant_type"`
    Host string `json:"host"`
    ClientId string `json:"client_id"`
    ClientSecret string `json:"client_secret"`
    Username string `json:"username"`
    Password string `json:"password"`
}

type WallabagToken struct {
    AccessToken string `json:"access_token"`
    ExpiresIn int64 `json:"expires_in"`
    RefreshToken string `json:"refresh_token"`
    Scope string `json:"scope"`
    TokenType string `json:"token_type"`
}

type WallabagPage struct {
    Page int64 `json:"page"`
    Limit int64 `json:"limit"`
    Pages int64 `json:"pages"`
    Total int64 `json:"total"`
    Embedded WallabagItems `json:"_embedded"`
}

type WallabagItems struct {
    Items []WallabagItem `json:"items"`
}

type WallabagItem struct {
    Id int64 `json:"id"`
    Title string `json:"title"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
    IsStarred bool `json:"is_starred"`
    IsArchived bool `json:is_archived"`
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

const output = ".local/share/remarkable/xochitl/"

func NewWallabag(path string) Wallabag {
    var wallabag Wallabag
    jsonFile, _ := os.Open(path)
    byteValue, _ := ioutil.ReadAll(jsonFile)
    var config WallabagConfig
    config.GrantType = "password"
    json.Unmarshal(byteValue, &config)
    wallabag.Config = config
    return wallabag
}

func (w *Wallabag) login() {
    req, _ := json.Marshal(w.Config)
    resp, _ := http.Post("https://"+w.Config.Host+"/oauth/v2/token", "application/json", bytes.NewBuffer(req))
    body, _ := ioutil.ReadAll(resp.Body)
    var token WallabagToken
    json.Unmarshal(body, &token)
    w.Token = token
}

func (w *Wallabag) getPages(p int64) {
    page := strconv.FormatInt(p, 10)

    client := &http.Client{}

    req, _ := http.NewRequest("GET", "https://"+w.Config.Host+"/api/entries.json", nil)
    req.Header.Set("Authorization", "Bearer "+w.Token.AccessToken)
    q := req.URL.Query()
    q.Add("detail", "metadata")
    q.Add("page", page)
    req.URL.RawQuery = q.Encode()
    res, _ := client.Do(req)
    body, _ := ioutil.ReadAll(res.Body)

    var wallabagPage WallabagPage
    json.Unmarshal(body, &wallabagPage)

    if wallabagPage.Page != wallabagPage.Pages {
        w.getPages(wallabagPage.Page+1)
    }
    w.Pages = append(w.Pages, wallabagPage)
}

func (w *Wallabag) getEpub(id string) []byte {
    client := &http.Client{}
    req, _ := http.NewRequest("GET", "https://"+w.Config.Host+"/api/entries/"+id+"/export.epub", nil)
    req.Header.Set("Authorization", "Bearer "+w.Token.AccessToken)
    resp, _ := client.Do(req)
    epub, _ := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    return epub
}

func (i *WallabagItem) toRemarkableItem() RemarkableItem {
    var metadata RemarkableItem
    metadata.Id = strconv.FormatInt(i.Id, 10)
    metadata.Type = "DocumentType"
    metadata.Parent = "wallabag"
    metadata.VisibleName = i.Title
    metadata.LastModified = i.UpdatedAt
    metadata.Version = 1
    metadata.Deleted = false
    metadata.MetadataModified = false
    metadata.Modified = false
    metadata.Pinned = i.IsStarred
    metadata.Synced = false
    return metadata
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
    Wallabag := NewWallabag(".config/reMarkable-sync/wallabag.json")
    Wallabag.login()
    Wallabag.getPages(1)

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
            RemarkableItem := WallabagItem.toRemarkableItem()

            if _, err := os.Stat(output+RemarkableItem.Id+".epub"); os.IsNotExist(err) {
                fmt.Print("Get EPUB of Wallabag element "+RemarkableItem.Id+" => ")
                channel := make(chan struct{})
                go indicator(channel)
                RemarkableItem.Content = Wallabag.getEpub(RemarkableItem.Id)
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
