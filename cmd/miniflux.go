package cmd

import (
	"fmt"
	"strconv"

	"github.com/vlaborie/remarkable-sync/remarkable"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	miniflux "miniflux.app/client"
)

var (
	minifluxConfig *viper.Viper

	minifluxCmd = &cobra.Command{
		Use:   "miniflux",
		Short: "Sync from Miniflux service",
		Long:  `Sync from Miniflux service.`,
		Run: func(cmd *cobra.Command, args []string) {
			Remarkable := remarkable.New("/home/root/.local/share/remarkable/xochitl/")

			if err := minifluxConfig.ReadInConfig(); err == nil {
				fmt.Println("Enable Miniflux sync with config file:", minifluxConfig.ConfigFileUsed())
				Remarkable.Items = append(Remarkable.Items, Remarkable.AddDir("miniflux", "Miniflux", ""))
				Miniflux := miniflux.New("https://"+minifluxConfig.GetString("host"), minifluxConfig.GetString("token"))
				m, _ := Miniflux.CategoryEntries(5, &miniflux.Filter{})
				for _, MinifluxItem := range m.Entries {
					RemarkableItem := remarkable.RemarkableItem{
						Id:               strconv.FormatInt(MinifluxItem.ID, 10),
						Type:             "DocumentType",
						Parent:           "miniflux",
						VisibleName:      MinifluxItem.Title,
						LastModified:     "",
						Version:          1,
						Deleted:          false,
						MetadataModified: false,
						Modified:         false,
						Pinned:           MinifluxItem.Starred,
						Synced:           false,
						ContentType:      "html",
						Content:          []byte(MinifluxItem.Content),
					}
					Remarkable.Items = append(Remarkable.Items, RemarkableItem)
				}
			}

			Remarkable.Write()
		},
	}
)

func init() {
	rootCmd.AddCommand(minifluxCmd)
	cobra.OnInitialize(initMinifluxConfig)
}

func initMinifluxConfig() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	minifluxConfig = viper.New()
	minifluxConfig.SetDefault("host", "app.miniflux.net")
	minifluxConfig.AddConfigPath("/etc/remarkable-sync")
	minifluxConfig.AddConfigPath(home + "/.config/remarkable-sync")
	minifluxConfig.AddConfigPath("./config")
	minifluxConfig.SetConfigName("miniflux")
}
