package ds

import (
	"fmt"
	"github.com/scolib/docksync/utils"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/json-iterator/go"

	"github.com/sirupsen/logrus"

)

func (ds *DS) dockerHubImages() []string {
	var images []string
	var val []struct {
		Name string
	}

	addr := fmt.Sprintf(DockerHubImage, ds.dockerRegistry())

	for {
		req, _ := http.NewRequest("GET", addr, nil)
		resp, err := ds.httpClient.Do(req)
		utils.CheckAndExit(err)
		if resp.StatusCode != http.StatusOK {
			utils.ErrorExit("Get docker hub images failed!", 1)
		}

		b, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		utils.CheckAndExit(err)

		_ = jsoniter.UnmarshalFromString(jsoniter.Get(b, "results").ToString(), &val)

		for _, v := range val {
			images = append(images, v.Name)
		}

		addr = jsoniter.Get(b, "next").ToString()
		if addr == "" {
			break
		}

	}
	return images
}

func (ds *DS) dockerHubImageList() []string {

	var images []string
	dockerHubImages := ds.dockerHubImages()

	logrus.Debugf("Number of docker hub images: %d", len(dockerHubImages))

	imgNameCh := make(chan string, 20)
	imgGetWg := new(sync.WaitGroup)
	imgGetWg.Add(len(dockerHubImages))

	for _, imageName := range dockerHubImages {

		tmpImageName := imageName

		go func() {
			defer func() {
				ds.QueryLimit <- 1
				imgGetWg.Done()
			}()

			addr := fmt.Sprintf(DockerHubTags, ds.dockerRegistry(), tmpImageName)

			select {
			case <-ds.QueryLimit:
				for {
					req, err := http.NewRequest("GET", addr, nil)
					utils.CheckAndExit(err)

					resp, err := ds.httpClient.Do(req)
					utils.CheckAndExit(err)

					b, err := ioutil.ReadAll(resp.Body)
					utils.CheckAndExit(err)
					_ = resp.Body.Close()

					var val []struct {
						Name string
					}
					_ = jsoniter.UnmarshalFromString(jsoniter.Get(b, "results").ToString(), &val)

					for _, tag := range val {
						imgNameCh <- tmpImageName + ":" + tag.Name
					}

					addr = jsoniter.Get(b, "next").ToString()
					if addr == "" {
						break
					}
				}
			}

		}()
	}

	var imgReceiveWg sync.WaitGroup
	imgReceiveWg.Add(1)
	go func() {
		defer imgReceiveWg.Done()
		for {
			select {
			case imageName, ok := <-imgNameCh:
				if ok {
					images = append(images, imageName)
				} else {
					goto imgSetExit
				}
			}
		}
	imgSetExit:
	}()

	imgGetWg.Wait()
	close(imgNameCh)
	imgReceiveWg.Wait()
	return images

}

func (ds *DS) dockerRegistry() string {
	if ds.DockerOrg == "" {
		return ds.DockerUser
	} else {
		return ds.DockerOrg
	}
}