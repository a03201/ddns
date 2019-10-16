package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type CloudflareConfigurationItem struct {
	UserName  string `json:"username"`
	Token     string `json:"token"`
	Domain    string `json:"domain"`
	SubDomain string `json:"sub_domain"`
}

type DnspodConfigurationItem struct {
	TokenId   string `json:"id"`
	Token     string `json:"token"`
	Domain    string `json:"domain"`
	SubDomain string `json:"sub_domain"`
}

type ConfigItem struct {
	Server string `json:"server"`
}

type Setting struct {

	Dnspod     DnspodConfigurationItem     `json:"dnspod"`
	Cloudflare CloudflareConfigurationItem `json:"cloudflare"`
	Config  ConfigItem `json:"config"`

}

func ParseConfig(configFile string)(Setting,error){
	var setting Setting

	appConf, err := os.Open(configFile)
	if err != nil {
		log.Println("opening app.conf failed:", err)
		return setting,err
	}

	defer func() {
		appConf.Close()
	}()

	b, err := ioutil.ReadAll(appConf)
	if err != nil {
		log.Println("reading app.conf failed:", err)
		return setting ,err
	}
	err = json.Unmarshal(b, &setting)
	if err != nil {
		log.Println("unmarshalling app.conf failed:", err)
		return setting,err
	}

	return setting,nil
}