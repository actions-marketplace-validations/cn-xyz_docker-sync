package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
)

type Download struct {
	RemoteStr  []string // 远程文件内容
	RemoteTag  []string // 仓库Tag
	url        string
	pushImages []string
}

// RemoteImageTags 远程仓库中的镜像Tag
type RemoteImageTags struct {
	Repository string   `json:"Repository"`
	Tags       []string `json:"Tags"`
}

// DestImageTags 远程仓库中的镜像Tag
type DestImageTags struct {
	Repository string   `json:"Repository"`
	Tags       []string `json:"Tags"`
}

func NewDownload(url string) *Download {
	s := &Download{
		url: url,
	}
	return s
}

// GetRemoteCtx 获取远程地址文件中存放的需要同步的镜像名称内容
func (d *Download) GetRemoteCtx() {
	response, err := http.Get(d.url)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(response.Body)
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 获取远程文件内容进行换行处理
	data := strings.NewReader(string(responseData))
	readerLine := bufio.NewReader(data)

	for {
		readLineStr, err := readerLine.ReadString('\n')

		if len(readLineStr) > 2 {
			newText := strings.Trim(readLineStr, "\n")
			d.RemoteStr = append(d.RemoteStr, newText)
		}

		if err == io.EOF {
			break
		}
	}
}

func (d *Download) CopyImage(destHub string) {
	for _, v := range d.RemoteStr {
		// step 1 远程
		remoteCmd := exec.Command("skopeo", "list-tags", "docker://"+v)
		remoteOut, err := remoteCmd.CombinedOutput()

		if err != nil {
			fmt.Printf("combined out:\n%s\n", string(remoteOut))
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}

		var remoteTags RemoteImageTags
		err = json.Unmarshal(remoteOut, &remoteTags)
		if err != nil {
			log.Fatalf("json.Unmarshal.Tags failed with %s\n", err)
		}

		next := filepath.Base(remoteTags.Repository)

		// step 2 检查 本地仓库是否有当前镜像
		destCmd := exec.Command("skopeo", "list-tags", "docker://"+destHub+next)
		destOut, err := destCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("combined out:\n%s\n", string(destOut))
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}

		var destImageTags DestImageTags
		err = json.Unmarshal(destOut, &destImageTags)
		if err != nil {
			log.Fatalf("json.Unmarshal.Tags failed with %s\n", err)
		}

		// 对比tag差异
		diffCtx := difference(remoteTags.Tags, destImageTags.Tags)

		// 地址组合推送镜像
		name := filepath.Base(remoteTags.Repository)

		// 复制镜像
		for _, v := range diffCtx {
			pushCmd := exec.Command("skopeo",
				"copy",
				"--insecure-policy",
				"--src-tls-verify=false",
				"--dest-tls-verify=false",
				"-q",
				"docker://"+remoteTags.Repository+":"+v,
				"docker://"+destHub+name+":"+v)

			pushOut, err := pushCmd.CombinedOutput()
			if err != nil {
				fmt.Printf("combined out:\n%s\n", string(pushOut))
				log.Fatalf("cmd.Run() failed with %s\n", err)
			}

			log.Println(string(pushOut))
		}
	}
}

// 求交集
func intersect(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// 求差集 slice1-并集
func difference(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}

// skopeo copy --insecure-policy --src-tls-verify=false --dest-tls-verify=false -q docker://quay.io/prometheus/alertmanager:latest docker://docker.io/***/alertmanager:latest
// https://registry.hub.docker.com/v2/repositories/cnxyz/alertmanager/tags?n=10
// curl https://hub.docker.com/v2/namespaces/{namespace}/repository/{repository}/images/{digest}/tags
// skopeo login -u cnxyz -p pojqu0pudcAcsyxkin docker.io

// export http_proxy='http://192.168.6.151:10811'
// export https_proxy='http://192.168.6.151:10811'

// skopeo login -u jinyinji_1994@163.com -p skopeo10086 registry.cn-hangzhou.aliyuncs.com
// docker login --username=jinyin_1994@163.com registry.cn-hangzhou.aliyuncs.com
