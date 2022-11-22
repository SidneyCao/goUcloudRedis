package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ucloud/ucloud-sdk-go/services/umem"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

var (
	pubK   = flag.String("pub", "", "PublicKey 默认为空")
	priK   = flag.String("pri", "", "PrivateKey 默认为空")
	name   = flag.String("n", "", "备份名 默认为空")
	action = flag.String("a", "", "行为 默认为空")
	rid    = flag.String("i", "", "redis id 默认为空")
)

const filePath string = "./lastRedisBackupID-"

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
	req.GroupId = ucloud.String(*rid)
	req.BackupName = ucloud.String(*name)

	nFile := filePath + *rid + ".txt"
	file, err := os.OpenFile(nFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	resp, err := umemClient.CreateURedisBackup(req)
	if err != nil {
		write := bufio.NewWriter(file)
		write.WriteString("backup error!")
		write.Flush()
		log.Panic("[ERROR]", err)
	}

	s := fmt.Sprint("[RESPONSE]", resp)
	backupID := strings.Split(s, "}")[1]
	backupID = strings.TrimSpace(backupID)

	write := bufio.NewWriter(file)
	write.WriteString(backupID)
	write.Flush()

	return backupID
}

func DownloadBackup(umemClient *umem.UMemClient) string {

	nFile := filePath + *rid + ".txt"
	backupID, err := os.ReadFile(nFile)
	if err != nil {
		log.Panic(err)
	}

	req := umemClient.NewDescribeURedisBackupURLRequest()
	req.Zone = ucloud.String("kr-seoul-01")
	req.ProjectId = ucloud.String("org-4ak3mv")
	req.BackupId = ucloud.String(string(backupID))
	req.GroupId = ucloud.String(*rid)

	resp, err := umemClient.DescribeURedisBackupURL(req)
	if err != nil {
		log.Panic("[ERROR]", err)
	}

	s := fmt.Sprint("[RESPONSE]", resp)
	url := strings.Split(s, " ")[5]
	url = strings.TrimSpace(url)
	return url
}
