package cmd

import (
	"github.com/scolib/docksync/ds"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test sync",
	Long: `
Test sync.`,
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
			TestMode:       true,
			Repositories:   repositories,
			ImagesRegistry: imagesRegistry,
		}
		ds.Init()
		ds.Sync()
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
