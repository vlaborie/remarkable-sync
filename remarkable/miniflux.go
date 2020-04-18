package remarkable

import (
	"strconv"

	miniflux "miniflux.app/client"
)

func (Remarkable *Remarkable) Miniflux(host string, token string) {
	Remarkable.Items = append(Remarkable.Items, Remarkable.AddDir("miniflux", "Miniflux", ""))

	Miniflux := miniflux.New("https://"+host, token)
	m, _ := Miniflux.Entries(&miniflux.Filter{})
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
