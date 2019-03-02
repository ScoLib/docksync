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

type Image struct {
	Name string
	Tags []string
}


func (ds *DS) gcrImageList() []string {

	var images, publicImageNames []string

	if len(ds.Repositories) == 0 {
		publicImageNames = ds.gcrPublicImageNames()
	} else {
		publicImageNames = ds.Repositories
	}

	logrus.Debugf("Number of gcr images: %d", len(publicImageNames))

	imgNameCh := make(chan string, 20)
	imgGetWg := new(sync.WaitGroup)
	imgGetWg.Add(len(publicImageNames))

	for _, imageName := range publicImageNames {

		tmpImageName := imageName

		go func() {
			defer func() {
				ds.QueryLimit <- 1
				imgGetWg.Done()
			}()

			select {
			case <-ds.QueryLimit:
				req, err := http.NewRequest("GET", fmt.Sprintf(GcrImageTags, ds.NameSpace, tmpImageName), nil)
				utils.CheckAndExit(err)

				resp, err := ds.httpClient.Do(req)
				utils.CheckAndExit(err)

				b, err := ioutil.ReadAll(resp.Body)
				utils.CheckAndExit(err)
				_ = resp.Body.Close()

				var tags []string
				_ = jsoniter.UnmarshalFromString(jsoniter.Get(b, "tags").ToString(), &tags)

				for _, tag := range tags {
					imgNameCh <- tmpImageName + ":" + tag
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

func (ds *DS) gcrPublicImageNames() []string {

	req, err := http.NewRequest("GET", fmt.Sprintf(GcrImages, ds.NameSpace), nil)
	utils.CheckAndExit(err)

	resp, err := ds.httpClient.Do(req)
	utils.CheckAndExit(err)
	defer func() {
		_ = resp.Body.Close()
	}()

	b, err := ioutil.ReadAll(resp.Body)
	utils.CheckAndExit(err)

	var imageNames []string
	_ = jsoniter.UnmarshalFromString(jsoniter.Get(b, "child").ToString(), &imageNames)
	return imageNames
}
