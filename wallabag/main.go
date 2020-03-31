package wallabag

import (
    "bytes"
    "io/ioutil"
    "os"
    "strconv"
    "encoding/json"
    "net/http"
)

type Wallabag struct {
    Config WallabagConfig
    Token WallabagToken
    Items []WallabagItem
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
    Content string `json:"content"`
}

func New(path string) Wallabag {
    var wallabag Wallabag
    jsonFile, _ := os.Open(path)
    byteValue, _ := ioutil.ReadAll(jsonFile)
    var config WallabagConfig
    config.GrantType = "password"
    json.Unmarshal(byteValue, &config)
    wallabag.Config = config
    wallabag.login()
    wallabag.getPages(1)
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
    q.Add("detail", "full")
    q.Add("page", page)
    req.URL.RawQuery = q.Encode()
    res, _ := client.Do(req)
    body, _ := ioutil.ReadAll(res.Body)

    var wallabagPage WallabagPage
    json.Unmarshal(body, &wallabagPage)

    if wallabagPage.Page != wallabagPage.Pages {
        w.getPages(wallabagPage.Page+1)
    }
    w.Items = append(w.Items, wallabagPage.Embedded.Items...)
}
