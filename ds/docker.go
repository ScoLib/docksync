package ds

import (
	"context"
	"encoding/base64"
	"github.com/scolib/docksync/utils"
	"io"
	"io/ioutil"

	"github.com/json-iterator/go"

	"github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types"
)

func (ds *DS) Process(imageName string) {
	logrus.Infof("Process image: %s", imageName)

	ctx := context.Background()

	oldImageName := ds.getOldImageName(imageName)
	// logrus.Infof("oldImageName: %s", oldImageName)

	newImageName := "docker.io/" + ds.dockerRegistry() + "/" + imageName
	// logrus.Infof("newImageName: %s", newImageName)

	if !ds.TestMode {

		// pull image
		r, err := ds.dockerClient.ImagePull(ctx, oldImageName, types.ImagePullOptions{})
		if !utils.CheckErr(err) {
			logrus.Errorf("Failed to pull image: %s", oldImageName)
			return
		}
		_, _ = io.Copy(ioutil.Discard, r)
		_ = r.Close()
		logrus.Infof("Pull image: %s success.", oldImageName)

		// tag it
		err = ds.dockerClient.ImageTag(ctx, oldImageName, newImageName)
		if !utils.CheckErr(err) {
			logrus.Errorf("Failed to tag image [%s] ==> [%s]", oldImageName, newImageName)
			return
		}
		logrus.Infof("Tag image: %s success.", oldImageName)

		// push image
		authConfig := types.AuthConfig{
			Username: ds.DockerUser,
			Password: ds.DockerPassword,
		}
		encodedJSON, err := jsoniter.Marshal(authConfig)
		if !utils.CheckErr(err) {
			logrus.Errorln("Failed to marshal docker config")
			return
		}
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)
		r, err = ds.dockerClient.ImagePush(ctx, newImageName, types.ImagePushOptions{RegistryAuth: authStr})
		if !utils.CheckErr(err) {
			logrus.Errorf("Failed to push image: %s", newImageName)
			return
		}
		_, _ = io.Copy(ioutil.Discard, r)
		_ = r.Close()
		logrus.Infof("Push image: %s success.", newImageName)

		// clean image
		_, _ = ds.dockerClient.ImageRemove(ctx, oldImageName, types.ImageRemoveOptions{})
		_, _ = ds.dockerClient.ImageRemove(ctx, newImageName, types.ImageRemoveOptions{})
		logrus.Debugf("Remove image: %s success.", oldImageName)

	}
	ds.update <- imageName
	logrus.Debugln("Done.")

}
