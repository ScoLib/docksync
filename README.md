# docksync [![Build Status](https://travis-ci.com/ScoLib/docksync.svg?branch=master)](https://travis-ci.com/ScoLib/docksync)

A docker image sync tool

fork from [https://github.com/mritd/gcrsync.git](https://github.com/mritd/gcrsync.git)
## 安装

工具采用 go 编写，安装可直接从 release 页下载对应平台二进制文件即可；如预编译文件不含有您的平台，
可自行 build:

```bash
export GO111MODULE=on
go get github.com/scolib/docksync
```

## 使用

```bash
A docker image sync tool.

Usage:
  docksync [flags]
  docksync [command]

Available Commands:
  help        Help about any command
  monitor     Monitor sync images
  sync        Sync docker images
  test        Test sync

Flags:
      --debug                   debug mode
      --githubrepo string       github commit repo (default "klgd/ds-changelog")
      --githubtoken string      github commit token
  -h, --help                    help for docksync
      --httptimeout duration    http request timeout (default 1m40s)
      --imagesregistry string   images registry(gitlab|quay|gcr) (default "gitlab")
      --namespace string        google container registry namespace (default "google-containers")
      --org string              docker registry user organization
      --password string         docker registry user password
      --processlimit int        image process limit (default 10)
      --proxy string            http client proxy
      --querylimit int          http query limit (default 50)
      --repositories strings    images repository
      --synctimeout duration    sync timeout
      --user string             docker registry user

Use "docksync [command] --help" for more information about a command.
```

### monitor 命令

实时在控制台显示对比差异；一般用于监测同步进度

### sync 命令

该命令进行真正的同步操作，工作流程如下:

- 首先使用 `--githubtoken` 给定的 token 克隆 `--githubrepo` 给定的仓库到本地(这个仓库应当为元数据仓库)
- 获取 `gcr.io` 下由 `--namespace` 给定命名空间下的所有镜像
- 获取 `docker hub` 中 `--user` 指定用户下所有镜像
- 对比两者差异，得出待同步镜像
- 执行 `pull`、`tag`、`push` 操作，将其推送到由 `--user` 给定的 Docker Hub 用户仓库中
- 生成 CHANGELOG 并推送元数据仓库到远程

### test 命令

该命令与 sync 命令基本行为一致，只不过不进行真正的同步，会生成 CHANGELOG，但不会推送到远程

## 其他说明

该工具并不建议个人使用，因为同步镜像会有很多；我自己同步大约产生了 2T 左右的流量，耗时 3 天左右才算基本
同步完成；目前我已经将镜像同步到了 Docker Hub 的 `gcrxio` 用户下，可直接使用；这个工具开源目的是为了确保
`gcrxio` 用户下的镜像安全得到保证；具体更细节说明，可参考[博客文章](https://mritd.me/2018/09/17/google-container-registry-sync/)
