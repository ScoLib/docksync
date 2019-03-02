package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var debug bool
var proxy, dockerUser, dockerPassword, dockerOrg, nameSpace, imagesRegistry string
var githubRepo, githubToken string
var queryLimit, processLimit, monitorCount int
var httpTimeOut, syncTimeOut time.Duration
var repositories []string

var rootCmd = &cobra.Command{
	Use:   "docksync",
	Short: "A docker image sync tool",
	Long: `
A docker image sync tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode")
	rootCmd.PersistentFlags().StringVar(&proxy, "proxy", "", "http client proxy")
	rootCmd.PersistentFlags().StringVar(&dockerUser, "user", "", "docker registry user")
	rootCmd.PersistentFlags().StringVar(&dockerPassword, "password", "", "docker registry user password")
	rootCmd.PersistentFlags().StringVar(&dockerOrg, "org", "", "docker registry user organization")
	rootCmd.PersistentFlags().StringVar(&nameSpace, "namespace", "google-containers", "google container registry namespace")
	rootCmd.PersistentFlags().IntVar(&queryLimit, "querylimit", 50, "http query limit")
	rootCmd.PersistentFlags().DurationVar(&httpTimeOut, "httptimeout", 100*time.Second, "http request timeout")
	rootCmd.PersistentFlags().DurationVar(&syncTimeOut, "synctimeout", 0, "sync timeout")
	rootCmd.PersistentFlags().IntVar(&processLimit, "processlimit", 10, "image process limit")
	rootCmd.PersistentFlags().StringVar(&githubRepo, "githubrepo", "klgd/ds-changelog", "github commit repo")
	rootCmd.PersistentFlags().StringVar(&githubToken, "githubtoken", "", "github commit token")
	rootCmd.PersistentFlags().StringVar(&imagesRegistry, "imagesregistry", "gitlab", "images registry(gitlab|quay|gcr)")
	rootCmd.PersistentFlags().StringSliceVar(&repositories, "repositories", []string{}, "images repository")
}
