package remarkable

import (
    "strconv"

    "github.com/vlaborie/reMarkable-sync/wallabag"
)

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
}
