package dns

import (
	"ddns-go/config"
	"time"
)

// DNS interface
type DNS interface {
	Init(conf *config.Config)
	// 添加或更新IPv4/IPv6记录
	AddUpdateDomainRecords() (domains config.Domains)
}

// RunTimer 定时运行
func RunTimer(firstDelay time.Duration, delay time.Duration) {
	time.Sleep(firstDelay)
	for {
		RunOnce()
		time.Sleep(delay)
	}
}

// RunOnce RunOnce
func RunOnce() {
	conf, err := config.GetConfigCache()
	if err != nil {
		return
	}

	var dnsSelected DNS
	switch conf.DNS.Name {
	case "alidns":
		dnsSelected = &Alidns{}
	case "dnspod":
		dnsSelected = &Dnspod{}
	case "cloudflare":
		dnsSelected = &Cloudflare{}
	case "huaweicloud":
		dnsSelected = &Huaweicloud{}
	default:
		dnsSelected = &Alidns{}
	}
	dnsSelected.Init(&conf)

	domains := dnsSelected.AddUpdateDomainRecords()

	//云数据库Mongodb安全组白名单
	if conf.IPS.Enable {
		ModifyAliyunSecurityIps(&domains, &conf.IPS)
	}

	config.ExecWebhook(&domains, &conf)
}
