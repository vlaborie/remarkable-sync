package cmd

import (
	"fmt"

	"github.com/vlaborie/remarkable-sync/remarkable"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	wallabagConfig *viper.Viper
	minifluxConfig *viper.Viper

	rootCmd = &cobra.Command{
		Use:   "remarkable-sync",
		Short: "Sync tool for reMarkable paper tablet",
		Long: `Remarkable-sync is a Go applications for syncing external
services to reMarkable paper table, like Wallabag or Miniflux.`,
		Run: func(cmd *cobra.Command, args []string) {
			Remarkable := remarkable.New("/home/root/.local/share/remarkable/xochitl/")

			if err := wallabagConfig.ReadInConfig(); err == nil {
				Remarkable.Wallabag(wallabagConfig.GetString("host"), wallabagConfig.GetString("client_id"), wallabagConfig.GetString("client_secret"), wallabagConfig.GetString("username"), wallabagConfig.GetString("password"))
			}

			if err := minifluxConfig.ReadInConfig(); err == nil {
				Remarkable.Miniflux(minifluxConfig.GetString("host"), minifluxConfig.GetString("token"))
			}

			Remarkable.Write()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
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

	if err := wallabagConfig.ReadInConfig(); err == nil {
		fmt.Println("Enable Wallabag sync with config file:", wallabagConfig.ConfigFileUsed())
	}

	minifluxConfig = viper.New()
	minifluxConfig.SetDefault("host", "app.miniflux.net")
	minifluxConfig.AddConfigPath("/etc/remarkable-sync")
	minifluxConfig.AddConfigPath(home + "/.config/remarkable-sync")
	minifluxConfig.AddConfigPath("./config")
	minifluxConfig.SetConfigName("miniflux")

	if err := minifluxConfig.ReadInConfig(); err == nil {
		fmt.Println("Enable Miniflux sync with config file:", minifluxConfig.ConfigFileUsed())
	}
}
