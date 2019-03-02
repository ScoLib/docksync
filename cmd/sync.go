package cmd

import (
	"github.com/scolib/docksync/ds"
	"github.com/spf13/cobra"
)

var commitMsg string

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync docker images",
	Long: `
Sync docker images.`,
	Run: func(cmd *cobra.Command, args []string) {
		ds := &ds.DS{
			Proxy:          proxy,
			DockerUser:     dockerUser,
			DockerPassword: dockerPassword,
			DockerOrg:      dockerOrg,
			NameSpace:      nameSpace,
			QueryLimit:     make(chan int, queryLimit),
			ProcessLimit:   make(chan int, processLimit),
			SyncTimeOut:    syncTimeOut,
			HttpTimeOut:    httpTimeOut,
			GithubRepo:     githubRepo,
			GithubToken:    githubToken,
			CommitMsg:      commitMsg,
			Debug:          debug,
			Repositories:   repositories,
			ImagesRegistry: imagesRegistry,
		}
		ds.Init()
		ds.Sync()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.PersistentFlags().StringVar(&commitMsg, "commitmsg", "Travis CI Auto Synchronized.", "commit message")
}
