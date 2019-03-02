package ds

import (
	"fmt"
	"github.com/scolib/docksync/utils"
	"os"
	"path/filepath"
	"strings"
	"time"

)

func (ds *DS) Commit(images []string) {

	loc, _ := time.LoadLocation("Asia/Shanghai")

	repoDir := filepath.Join(strings.Split(ds.GithubRepo, "/")[1], "changelog")
	if _, err := os.Stat(repoDir); err != nil {
		_ = os.MkdirAll(repoDir, 0755)
	}

	repoChangeLog := filepath.Join(repoDir, fmt.Sprintf(ChangeLog, time.Now().In(loc).Format("2006-01-02")))

	var content []byte
	chgLog, err := os.OpenFile(repoChangeLog, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	utils.CheckAndExit(err)
	defer func() {
		_ = chgLog.Close()
	}()

	updateInfo := fmt.Sprintf("### %s Update:\n\n", time.Now().In(loc).Format("2006-01-02 15:04:05"))
	for _, imageName := range images {
		updateInfo += "- " + ds.getOldImageName(imageName) + "\n"
	}
	_, _ = chgLog.WriteString(updateInfo + string(content))

	utils.GitCmd(repoDir, "config", "--global", "push.default", "simple")
	utils.GitCmd(repoDir, "config", "--global", "user.email", "slice1213@gmail.com")
	utils.GitCmd(repoDir, "config", "--global", "user.name", "klgd")
	utils.GitCmd(repoDir, "add", ".")
	utils.GitCmd(repoDir, "commit", "-m", ds.CommitMsg)
	if !ds.TestMode {
		utils.GitCmd(repoDir, "push", "--force", ds.commitURL, "master")
	}

}

func (ds *DS) Clone() {
	_ = os.RemoveAll(strings.Split(ds.GithubRepo, "/")[1])
	utils.GitCmd("", "clone", ds.commitURL)
}
