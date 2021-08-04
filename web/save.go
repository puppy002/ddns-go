package web

import (
	"ddns-go/config"
	"ddns-go/dns"
	"net/http"
	"strings"
)

// Save 保存
func Save(writer http.ResponseWriter, request *http.Request) {

	conf, _ := config.GetConfigCache()

	idNew := request.FormValue("DnsID")
	secretNew := request.FormValue("DnsSecret")
	accessKeyID := request.FormValue("AccessKeyID")
	accessSecret := request.FormValue("AccessSecret")

	idHide, secretHide, aidHide, asecretHide := getHideIDSecret(&conf)

	if idNew != idHide {
		conf.DNS.ID = idNew
	}
	if secretNew != secretHide {
		conf.DNS.Secret = secretNew
	}
	if accessKeyID != aidHide {
		conf.IPS.AccessKeyID = accessKeyID
	}
	if accessSecret != asecretHide {
		conf.IPS.AccessSecret = accessSecret
	}

	// 覆盖以前的配置
	conf.DNS.Name = request.FormValue("DnsName")

	conf.Ipv4.Enable = request.FormValue("Ipv4Enable") == "on"
	conf.Ipv4.URL = strings.TrimSpace(request.FormValue("Ipv4Url"))
	conf.Ipv4.GetType = request.FormValue("Ipv4GetType")
	conf.Ipv4.NetInterface = request.FormValue("Ipv4NetInterface")
	conf.Ipv4.Domains = strings.Split(request.FormValue("Ipv4Domains"), "\r\n")

	conf.Ipv6.Enable = request.FormValue("Ipv6Enable") == "on"
	conf.Ipv6.GetType = request.FormValue("Ipv6GetType")
	conf.Ipv6.NetInterface = request.FormValue("Ipv6NetInterface")
	conf.Ipv6.URL = strings.TrimSpace(request.FormValue("Ipv6Url"))
	conf.Ipv6.Domains = strings.Split(request.FormValue("Ipv6Domains"), "\r\n")

	conf.Username = strings.TrimSpace(request.FormValue("Username"))
	conf.Password = request.FormValue("Password")

	conf.IPS.Region = request.FormValue("Region")
	conf.IPS.Scheme = request.FormValue("Scheme")
	conf.IPS.DBInstanceId = request.FormValue("DBInstanceId")
	conf.IPS.ModifyMode = request.FormValue("ModifyMode")
	conf.IPS.SecurityIpGroupName = request.FormValue("SecurityIpGroupName")

	conf.WebhookURL = strings.TrimSpace(request.FormValue("WebhookURL"))
	conf.WebhookRequestBody = strings.TrimSpace(request.FormValue("WebhookRequestBody"))

	conf.NotAllowWanAccess = request.FormValue("NotAllowWanAccess") == "on"

	// 保存到用户目录
	err := conf.SaveConfig()

	// 只运行一次
	go dns.RunOnce()

	// 回写错误信息
	if err == nil {
		writer.Write([]byte("ok"))
	} else {
		writer.Write([]byte(err.Error()))
	}

}
