package main

import (
	"flag"
	"fmt"

	"github.com/ucloud/ucloud-sdk-go/services/umem"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

var (
	pubK   = flag.String("pub", "", "PublicKey 默认为空")
	priK   = flag.String("pri", "", "PrivateKey 默认为空")
	file   = flag.String("f", "", "日期 默认为空")
	action = flag.String("a", "", "行为 默认为空")
)

func main() {

	//获取命令行参数
	flag.Parse()

	cfg := ucloud.NewConfig()
	cfg.Region = "kr-seoul"
	cfg.BaseUrl = "https://api.ucloud.cn"

	// replace the public/private key by your own
	credential := auth.NewCredential()
	credential.PublicKey = *pubK
	credential.PrivateKey = *priK

	umemClient := umem.NewClient(&cfg, &credential)

	if *action == "create" {
		res := CreateBackup(umemClient)
		fmt.Println(res)
	} else if *action == "down" {
		res := DownloadBackup(umemClient)
		fmt.Println(res)
	} else {
		fmt.Println("no such action")
	}

}

func CreateBackup(umemClient *umem.UMemClient) string {
	req := umemClient.NewCreateURedisBackupRequest()
	req.Zone = ucloud.String("kr-seoul-01")
	req.ProjectId = ucloud.String("org-4ak3mv")
	req.GroupId = ucloud.String("uredis-112q4qie")
	req.BackupName = ucloud.String(*file)

	resp, err := umemClient.CreateURedisBackup(req)
	if err != nil {
		return fmt.Sprint("[ERROR]", err)
	}

	return fmt.Sprint("[RESPONSE]", resp)
}

func DownloadBackup(umemClient *umem.UMemClient) string {
	req := umemClient.NewDescribeURedisBackupURLRequest()
	req.Zone = ucloud.String("kr-seoul-01")
	req.ProjectId = ucloud.String("org-4ak3mv")
	req.BackupId = ucloud.String(*file)
	req.GroupId = ucloud.String("uredis-112q4qie")

	resp, err := umemClient.DescribeURedisBackupURL(req)
	if err != nil {
		fmt.Println("[ERROR]", err)
		return fmt.Sprint("[ERROR]", err)
	}

	return fmt.Sprint("[RESPONSE]", resp)
}
