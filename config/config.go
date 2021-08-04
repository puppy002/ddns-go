package config

import (
	"ddns-go/util"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// Ipv4Reg IPv4正则
const Ipv4Reg = `((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])`

// Ipv6Reg IPv6正则
const Ipv6Reg = `((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))`

// Config 配置
type Config struct {
	Ipv4 struct {
		Enable bool
		// 获取IP类型 url/netInterface
		GetType      string
		URL          string
		NetInterface string
		Domains      []string
	}
	Ipv6 struct {
		Enable bool
		// 获取IP类型 url/netInterface
		GetType      string
		URL          string
		NetInterface string
		Domains      []string
	}
	DNS DNSConfig
	User
	Webhook
	IPS IPSConfig
	// 禁止公网访问
	NotAllowWanAccess bool
}

// DNSConfig DNS配置
type DNSConfig struct {
	// 名称。如：alidns,webhook
	Name   string
	ID     string
	Secret string
}

type IPSConfig struct {
	Enable              bool
	Region              string
	AccessKeyID         string
	AccessSecret        string
	Scheme              string
	SecurityIps         string
	DBInstanceId        string
	ModifyMode          string
	SecurityIpGroupName string
}

// ConfigCache ConfigCache
type cacheType struct {
	ConfigSingle *Config
	Err          error
	Lock         sync.Mutex
}

var cache = &cacheType{}

// GetConfigCache 获得配置
func GetConfigCache() (conf Config, err error) {

	if cache.ConfigSingle != nil {
		return *cache.ConfigSingle, cache.Err
	}

	cache.Lock.Lock()
	defer cache.Lock.Unlock()

	// init config
	cache.ConfigSingle = &Config{}

	configFilePath := util.GetConfigFilePath()
	_, err = os.Stat(configFilePath)
	if err != nil {
		cache.Err = err
		return *cache.ConfigSingle, err
	}

	byt, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Println("config.yaml读取失败")
		cache.Err = err
		return *cache.ConfigSingle, err
	}

	err = yaml.Unmarshal(byt, cache.ConfigSingle)
	if err != nil {
		log.Println("反序列化配置文件失败", err)
		cache.Err = err
		return *cache.ConfigSingle, err
	}
	// remove err
	cache.Err = nil
	return *cache.ConfigSingle, err
}

// SaveConfig 保存配置
func (conf *Config) SaveConfig() (err error) {
	byt, err := yaml.Marshal(conf)
	if err != nil {
		log.Println(err)
		return err
	}

	configFilePath := util.GetConfigFilePath()
	err = ioutil.WriteFile(configFilePath, byt, 0600)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("配置文件已保存在: %s\n", configFilePath)

	// 清空配置缓存
	cache.ConfigSingle = nil

	return
}

// GetIpv4Addr 获得IPv4地址
func (conf *Config) GetIpv4Addr() (result string) {
	if conf.Ipv4.Enable {
		// 判断从哪里获取IP
		if conf.Ipv4.GetType == "netInterface" {
			// 从网卡获取IP
			ipv4, _, err := GetNetInterface()
			if err != nil {
				log.Println("从网卡获得IPv4失败!")
				return
			}

			for _, netInterface := range ipv4 {
				if netInterface.Name == conf.Ipv4.NetInterface && len(netInterface.Address) > 0 {
					return netInterface.Address[0]
				}
			}

			log.Println("从网卡中获得IPv4失败! 网卡名: ", conf.Ipv4.NetInterface)
			return
		}

		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(conf.Ipv4.URL)
		if err != nil {
			log.Println(fmt.Sprintf("未能获得IPv4地址! <a target='blank' href='%s'>点击查看接口能否返回IPv4地址</a>,", conf.Ipv4.URL))
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("读取IPv4结果失败! 查询URL: ", conf.Ipv4.URL)
			return
		}
		comp := regexp.MustCompile(Ipv4Reg)
		result = comp.FindString(string(body))
	}
	return
}

// GetIpv6Addr 获得IPv6地址
func (conf *Config) GetIpv6Addr() (result string) {
	if conf.Ipv6.Enable {
		// 判断从哪里获取IP
		if conf.Ipv6.GetType == "netInterface" {
			// 从网卡获取IP
			_, ipv6, err := GetNetInterface()
			if err != nil {
				log.Println("从网卡获得IPv6失败!")
				return
			}

			for _, netInterface := range ipv6 {
				if netInterface.Name == conf.Ipv6.NetInterface && len(netInterface.Address) > 0 {
					return netInterface.Address[0]
				}
			}

			log.Println("从网卡中获得IPv6失败! 网卡名: ", conf.Ipv6.NetInterface)
			return
		}

		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(conf.Ipv6.URL)
		if err != nil {
			log.Println(fmt.Sprintf("未能获得IPv6地址! <a target='blank' href='%s'>点击查看接口能否返回IPv6地址</a>, 官方说明:<a target='blank' href='%s'>点击访问</a> ", conf.Ipv6.URL, "https://github.com/jeessy2/ddns-go#使用ipv6"))
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("读取IPv6结果失败! 查询URL: ", conf.Ipv6.URL)
			return
		}
		comp := regexp.MustCompile(Ipv6Reg)
		result = comp.FindString(string(body))
	}
	return
}
