package remarkable

import (
	"strconv"

	"github.com/vlaborie/remarkable-sync/wallabag"
)

func (Remarkable *Remarkable) Wallabag(host string, ID string, secret string, username string, password string) {
	Remarkable.Items = append(Remarkable.Items, Remarkable.AddDir("wallabag", "Wallabag", ""))
	Wallabag := wallabag.New(host, ID, secret, username, password)
	for _, wallabagItem := range Wallabag.Items {
		var RemarkableItem RemarkableItem
		RemarkableItem.FromWallabag(wallabagItem)
		Remarkable.Items = append(Remarkable.Items, RemarkableItem)
	}
}

func (RemarkableItem *RemarkableItem) FromWallabag(WallabagItem wallabag.WallabagItem) {
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
	RemarkableItem.ContentType = "html"
	RemarkableItem.Content = []byte(WallabagItem.Content)
}
