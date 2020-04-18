package cmd

import (
	"fmt"
	"strconv"

	"github.com/vlaborie/remarkable-sync/remarkable"
	"github.com/vlaborie/remarkable-sync/wallabag"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	wallabagConfig *viper.Viper

	wallabagCmd = &cobra.Command{
		Use:   "wallabag",
		Short: "Sync from Wallabag service",
		Long:  `Sync from Wallabag service.`,
		Run: func(cmd *cobra.Command, args []string) {
			Remarkable := remarkable.New("/home/root/.local/share/remarkable/xochitl/")

			if err := wallabagConfig.ReadInConfig(); err == nil {
				fmt.Println("Enable Wallabag sync with config file:", wallabagConfig.ConfigFileUsed())
				Remarkable.Items = append(Remarkable.Items, Remarkable.AddDir("wallabag", "Wallabag", ""))
				Wallabag := wallabag.New(wallabagConfig.GetString("host"), wallabagConfig.GetString("client_id"), wallabagConfig.GetString("client_secret"), wallabagConfig.GetString("username"), wallabagConfig.GetString("password"))
				for _, WallabagItem := range Wallabag.Items {
					RemarkableItem := remarkable.RemarkableItem{
						Id:               strconv.FormatInt(WallabagItem.Id, 10),
						Type:             "DocumentType",
						Parent:           "wallabag",
						VisibleName:      WallabagItem.Title,
						LastModified:     WallabagItem.UpdatedAt,
						Version:          1,
						Deleted:          false,
						MetadataModified: false,
						Modified:         false,
						Pinned:           WallabagItem.IsStarred,
						Synced:           false,
						ContentType:      "html",
						Content:          []byte(WallabagItem.Content),
					}
					Remarkable.Items = append(Remarkable.Items, RemarkableItem)
				}
			}

			Remarkable.Write()
		},
	}
)

func init() {
	rootCmd.AddCommand(wallabagCmd)
	cobra.OnInitialize(initWallabagConfig)
}

func initWallabagConfig() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	wallabagConfig = viper.New()
	wallabagConfig.SetDefault("host", "app.wallabag.it")
	wallabagConfig.AddConfigPath("/etc/remarkable-sync")
	wallabagConfig.AddConfigPath(home + "/.config/remarkable-sync")
	wallabagConfig.AddConfigPath("./config")
	wallabagConfig.SetConfigName("wallabag")
}
