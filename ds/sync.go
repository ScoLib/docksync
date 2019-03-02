package ds

import (
	"context"
	"fmt"
	"github.com/scolib/docksync/utils"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	ChangeLog         = "CHANGELOG-%s.md"
	GcrRegistryTpl    = "gcr.io/%s/%s"
	GitLabRegistryTpl = "registry.gitlab.com/gitlab-org/build/cng/%s"
	QuayRegistryTpl   = "quay.io/%s/%s"
	GcrImages         = "https://gcr.io/v2/%s/tags/list"
	GcrImageTags      = "https://gcr.io/v2/%s/%s/tags/list"
	DockerHubImage    = "https://hub.docker.com/v2/repositories/%s/?page_size=100"
	DockerHubTags     = "https://hub.docker.com/v2/repositories/%s/%s/tags/?page_size=100"

	GitLabRegistry  = "https://gitlab.com/gitlab-org/build/CNG/container_registry.json"
	GitLabImageTags = "https://gitlab.com/gitlab-org/build/CNG/registry/repository/%d/tags?format=json&page=%d&per_page=100"

	QuayImageTags = "https://quay.io/api/v1/repository/%s/%s/tag/?limit=100&page=%d&onlyActiveTags=true"
)

func (ds *DS) Sync() {
	images := ds.getImageList()
	dockerHubImages := ds.dockerHubImageList()
	needSyncImages := utils.SliceDiff(images, dockerHubImages)

	logrus.Infof("%s registry images total: %d", ds.ImagesRegistry, len(images))
	logrus.Infof("Docker hub images total: %d", len(dockerHubImages))
	logrus.Infof("Number of images waiting to be processed: %d", len(needSyncImages))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if ds.SyncTimeOut != 0 {
			select {
			case <-time.After(ds.SyncTimeOut):
				cancel()
			}
		}
	}()

	processWg := new(sync.WaitGroup)
	processWg.Add(len(needSyncImages))

	for _, imageName := range needSyncImages {
		tmpImageName := imageName
		go func() {
			defer func() {
				ds.ProcessLimit <- 1
				processWg.Done()
			}()
			select {
			case <-ds.ProcessLimit:
				ds.Process(tmpImageName)
			case <-ctx.Done():
			}
		}()
	}

	// doc gen
	chgWg := new(sync.WaitGroup)
	chgWg.Add(1)
	go func() {
		defer chgWg.Done()

		var images []string
		for {
			select {
			case imageName, ok := <-ds.update:
				if ok {
					images = append(images, imageName)
				} else {
					goto ChangeLogDone
				}
			case <-ctx.Done():
				goto ChangeLogDone
			}
		}
	ChangeLogDone:
		if len(images) > 0 {
			ds.Commit(images)
		}
	}()

	processWg.Wait()
	close(ds.update)
	chgWg.Wait()

}

func (ds *DS) Monitor() {

	if ds.MonitorCount == -1 {
		for {
			select {
			case <-time.Tick(5 * time.Second):
				gcrImages := ds.gcrImageList()
				dockerHubImages := ds.dockerHubImageList()
				needSyncImages := utils.SliceDiff(gcrImages, dockerHubImages)
				logrus.Infof("Gcr images: %d | Docker hub images: %d | Waiting process: %d", len(gcrImages), len(dockerHubImages), len(needSyncImages))
			}
		}
	} else {
		for i := 0; i < ds.MonitorCount; i++ {
			select {
			case <-time.Tick(5 * time.Second):
				gcrImages := ds.gcrImageList()
				dockerHubImages := ds.dockerHubImageList()
				needSyncImages := utils.SliceDiff(gcrImages, dockerHubImages)
				logrus.Infof("Gcr images: %d | Docker hub images: %d | Waiting process: %d", len(gcrImages), len(dockerHubImages), len(needSyncImages))
			}
		}
	}

}

func (ds *DS) getImageList() []string {
	var images []string
	switch ds.ImagesRegistry {
	case "gitlab":
		images = ds.gitLabImageList()
	case "gcr":
		images = ds.gcrImageList()
	case "quay":
		images = ds.quayImageList()
	default:
		utils.ErrorExit(fmt.Sprintf("not support %s", ds.ImagesRegistry), 1)
	}
	return images
}

func (ds *DS) getOldImageName(imageName string) string {
	var name string
	switch ds.ImagesRegistry {
	case "gitlab":
		name = fmt.Sprintf(GitLabRegistryTpl, imageName)
	case "gcr":
		name = fmt.Sprintf(GcrRegistryTpl, ds.NameSpace, imageName)
	case "quay":
		name = fmt.Sprintf(QuayRegistryTpl, ds.NameSpace, imageName)
	default:
		utils.ErrorExit(fmt.Sprintf("not support %s", ds.ImagesRegistry), 1)
	}
	return name
}
