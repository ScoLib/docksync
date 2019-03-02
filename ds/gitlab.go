package ds

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/scolib/docksync/utils"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sync"
)

func (ds *DS) gitLabImageList() []string {
	var images []string
	imageNames := ds.gitLabImageNames()

	imageNameCh := make(chan string, 20)
	imageGetWg := new(sync.WaitGroup)
	imageGetWg.Add(len(imageNames))

	for id, imageName := range imageNames {
		tmpId := id
		tmpImageName := imageName

		go func() {
			defer imageGetWg.Done()

			for i := 1; ; i++ {
				addr := fmt.Sprintf(GitLabImageTags, tmpId, i)
				// fmt.Println(addr)
				req, err := http.NewRequest("GET", addr, nil)
				utils.CheckAndExit(err)

				resp, err := ds.httpClient.Do(req)
				utils.CheckAndExit(err)

				b, err := ioutil.ReadAll(resp.Body)
				utils.CheckAndExit(err)
				_ = resp.Body.Close()
				if len(b) < 10 {
					logrus.Infof("id %d end", tmpId)
					break
				}

				var val []struct {
					Name string // version
				}
				_ = jsoniter.Unmarshal(b, &val)

				for _, tag := range val {
					if len(tag.Name) < 20 {
						imageNameCh <- tmpImageName + ":" + tag.Name
					}
				}
			}
		}()
	}

	var imageReceiveWg sync.WaitGroup
	imageReceiveWg.Add(1)
	go func() {
		defer imageReceiveWg.Done()
		for {
			select {
			case imageName, ok := <-imageNameCh:
				if ok {
					images = append(images, imageName)
				} else {
					goto imageSetExit
				}
			}
		}
	imageSetExit:
	}()

	imageGetWg.Wait()
	close(imageNameCh)
	imageReceiveWg.Wait()

	return images
}

func (ds *DS) gitLabImageNames() map[int]string {

	images := make(map[int]string)
	var val []struct {
		Id   int
		Name string
	}

	req, _ := http.NewRequest("GET", GitLabRegistry, nil)
	resp, err := ds.httpClient.Do(req)
	utils.CheckAndExit(err)
	if resp.StatusCode != http.StatusOK {
		utils.ErrorExit("Get gitlab images failed!", 1)
	}
	b, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	utils.CheckAndExit(err)

	_ = jsoniter.Unmarshal(b, &val)

	for _, v := range val {
		if len(ds.Repositories) == 0 {
			images[v.Id] = v.Name
		} else {
			for _, img := range ds.Repositories {
				if img == v.Name {
					images[v.Id] = v.Name
				}
			}
		}
	}
	return images
}
