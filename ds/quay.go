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

func (ds *DS) quayImageList() []string {
	var images []string
	imageNames := ds.Repositories
	imageNameCh := make(chan string, 20)
	imageGetWg := new(sync.WaitGroup)
	imageGetWg.Add(len(imageNames))

	for _, imageName := range imageNames {
		tmpImageName := imageName

		go func() {
			defer imageGetWg.Done()

			for i := 1; ; i++ {
				addr := fmt.Sprintf(QuayImageTags, ds.NameSpace, imageName, i)
				// fmt.Println(addr)
				req, err := http.NewRequest("GET", addr, nil)
				utils.CheckAndExit(err)

				resp, err := ds.httpClient.Do(req)
				utils.CheckAndExit(err)

				b, err := ioutil.ReadAll(resp.Body)
				utils.CheckAndExit(err)
				_ = resp.Body.Close()

				var val []struct {
					Name string // version
				}
				_ = jsoniter.UnmarshalFromString(jsoniter.Get(b, "tags").ToString(), &val)

				if len(val) == 0 {
					logrus.Infof("quay image %s end", imageName)
					break
				}

				for _, tag := range val {
					imageNameCh <- tmpImageName + ":" + tag.Name
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
