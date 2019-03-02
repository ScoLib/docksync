package ds

import (
	"github.com/docker/docker/client"
	"github.com/scolib/docksync/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

type DS struct {
	Proxy          string
	DockerUser     string
	DockerPassword string
	DockerOrg      string
	GithubToken    string
	GithubRepo     string
	CommitMsg      string
	MonitorCount   int
	TestMode       bool
	Debug          bool
	SyncTimeOut    time.Duration
	QueryLimit     chan int
	ProcessLimit   chan int
	HttpTimeOut    time.Duration
	httpClient     *http.Client
	dockerClient   *client.Client
	dockerHubToken string
	update         chan string
	commitURL      string
	Repositories   []string
	NameSpace      string
	ImagesRegistry string
}

func (ds *DS) Init() {

	if ds.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Infoln("Init http client.")
	ds.httpClient = &http.Client{
		Timeout: ds.HttpTimeOut,
	}
	if ds.Proxy != "" {
		p := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(ds.Proxy)
		}
		ds.httpClient.Transport = &http.Transport{Proxy: p}
	}

	logrus.Infoln("Init docker client.")
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.39"))
	utils.CheckAndExit(err)
	ds.dockerClient = dockerClient

	logrus.Infoln("Init limit channel.")
	for i := 0; i < cap(ds.QueryLimit); i++ {
		ds.QueryLimit <- 1
	}
	for i := 0; i < cap(ds.ProcessLimit); i++ {
		ds.ProcessLimit <- 1
	}

	logrus.Infoln("Init update channel.")
	ds.update = make(chan string, 20)

	// logrus.Infoln("Init commit repo.")
	// if ds.GithubToken == "" {
	// 	utils.ErrorExit("Github Token is blank!", 1)
	// }
	// ds.commitURL = "https://" + ds.GithubToken + "@github.com/" + ds.GithubRepo + ".git"
	// ds.Clone()

	logrus.Infoln("Init success...")
}
