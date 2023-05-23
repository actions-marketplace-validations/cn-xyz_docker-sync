package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	url         string
	user        string
	pass        string
	hubLoginUrl string
	namespace   string
)

func init() {
	rootCmd.Flags().StringVar(&url, "config", getEnvStr("REMOTE_URL", "https://raw.githubusercontent.com/cn-xyz/image-sync/main/registry.k8s.io/containers.txt"), "Remote file URL address.")
	rootCmd.Flags().StringVar(&hubLoginUrl, "login", getEnvStr("HUB_LOGIN_URL", "registry.cn-hangzhou.aliyuncs.com"), "HUB URL address.")
	rootCmd.Flags().StringVar(&user, "user", getEnvStr("DOCKER_USER", "jinyinji_1994@163.com"), "DestHub User.")
	rootCmd.Flags().StringVar(&pass, "pass", getEnvStr("DOCKER_PASS", "skopeo10086"), "DestHub Password.")
	rootCmd.Flags().StringVar(&namespace, "namespace", getEnvStr("NAMESPACE", "cnxyz"), "DestHub Password.")
}

var rootCmd = &cobra.Command{
	Use:   "docker-sync is a synchronization tool created to solve the problem that domestic mirror repositories cannot be used.",
	Short: `Only syncing images to docker.io is currently supported`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// todo skopeo 登录
		destCmd := exec.Command("skopeo", "login", "-u", user, "-p", pass, hubLoginUrl)
		destOut, err := destCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("combined out:\n%s\n", string(destOut))
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ttt := NewDownload(url)
		ttt.GetRemoteCtx()
		ttt.CopyImage(hubLoginUrl + "/" + namespace + "/")
	},
}

func getEnvStr(env string, defaultValue string) string {
	v := os.Getenv(env)
	if v == "" {
		return defaultValue
	}
	return v
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrf("docker-sync root cmd execute: %s", err)
		os.Exit(1)
	}
}
