package dns

import (
	"ddns-go/config"
	"log"

	dds "github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
)

func ModifyAliyunSecurityIps(domains *config.Domains, conf *config.IPSConfig) {
	client, err := dds.NewClientWithAccessKey(conf.Region, conf.AccessKeyID, conf.AccessSecret)
	if err != nil {
		//
		log.Println(err.Error())
	}
	//recordType := "A"
	// ipAddr, _ := conf.Domains.ParseDomainResult(recordType)
	request := dds.CreateModifySecurityIpsRequest()
	request.Scheme = conf.Scheme
	request.SecurityIps = conf.SecurityIps
	request.DBInstanceId = conf.DBInstanceId
	request.ModifyMode = conf.ModifyMode
	request.SecurityIpGroupName = conf.SecurityIpGroupName
	// request.Scheme = "https"

	// request.SecurityIps = "125.122.57.167"
	// request.DBInstanceId = "dds-bp11f44078198d24"
	// request.ModifyMode = "Cover"
	// request.SecurityIpGroupName = "linping"

	//修改Ips
	response, err := client.ModifySecurityIps(request)
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("response is %#v\n", response)
}
