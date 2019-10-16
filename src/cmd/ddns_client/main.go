package main

import (
	"auto_ddns/src/common"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)
// Define our message object
type Message struct {
	To string `json:"to"`
	From string `json:"from"`
	Type string `json:"type"`
	Message  string `json:"message"`
}

var request_ip string
var chSend chan string

func main() {

	log.SetOutput(os.Stdout)


	request_ip=""

	setting,err:=common.ParseConfig("config.json")

	if(err!=nil){
		log.Fatal(err)
	}


	//test account
	if(setting.Cloudflare.UserName!=""){
		if err := common.TestCloudflareAccount(setting.Cloudflare.UserName, setting.Cloudflare.Token); err != nil {
			log.Println("FATAL Cloudfare Account:",err)
			return
		}
		log.Println("Cloudflare Account OK")
	}
	if(setting.Dnspod.TokenId!=""){
		if err := common.TestDnspodRequestByToken(setting.Dnspod.TokenId, setting.Dnspod.Token); err != nil {
			fmt.Println("FATAL Dnspod Account:",err)
			return
		}
		log.Println("Dnspod Account OK")
	}


	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	last_ip:="no-ip"

reconnet:
	chSend = make(chan string,10)
	u := url.URL{Scheme: "ws", Host: setting.Config.Server, Path: "/ws/chat/ddns"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("dial:", err)
		time.Sleep(3*time.Second)
		goto reconnet
	}
	defer c.Close()


	go func() {
		defer close(chSend)
		var msg Message
		for {
			err := c.ReadJSON(&msg)
			if err != nil {
				log.Println("read:", err)
				return
			}
			if msg.Type=="whoareyou"{
				ip:=strings.Split(msg.Message,":")
				log.Printf("IP: %s", ip[0])

				if last_ip==ip[0]{
					log.Println("same ip, no update")
					continue
				}
				request_ip=ip[0]
				last_ip=ip[0]

				updateddns(setting)
			}

			log.Printf("recv: %s", msg)
		}
	}()

	ticker := time.NewTicker(60*time.Second)
	defer ticker.Stop()

	//_ = c.WriteJSON(Message{ Type: "setwhoami", Message: string("@ddns-"+domain)})

	for {
		select {
		case msg,ok:=<-chSend:
			if(ok==false){
				goto reconnet
			}
			log.Println("msg received:"+msg)
			err := c.WriteJSON(Message{To:"@console-",Type:"text",Message:msg})
			if err != nil {
				log.Println("write:", err)
				c.Close()
				goto reconnet
			}

		case _ = <-ticker.C:
			//updateddns(setting)

			err := c.WriteJSON(Message{To:"@console-",Type:"text",Message:"[hb]"})
			if err != nil {
				log.Println("write:", err)
				c.Close()
				goto reconnet
			}
			log.Println("[heartbeat]")
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
			}
			c.Close()
			return

		}

	}
}

func updateddns(setting common.Setting) bool {
	if(request_ip==""){
		return false
	}
	var domain=""
	if(setting.Cloudflare.UserName!="") {
		if err := common.CloudflareRequest(setting.Cloudflare.UserName, setting.Cloudflare.Token,
			setting.Cloudflare.Domain, setting.Cloudflare.SubDomain, request_ip); err != nil {
			log.Println("error:", err)
			return false
		}
		domain+=setting.Cloudflare.SubDomain+"."+setting.Cloudflare.Domain
	}

	if(setting.Dnspod.TokenId!=""){
		if err := common.DnspodRequestByToken(setting.Dnspod.TokenId, setting.Dnspod.Token,
			setting.Dnspod.Domain, setting.Cloudflare.SubDomain, request_ip); err != nil {
			log.Println("error:", err)
			return false
		}
		domain=domain + " "+ setting.Dnspod.SubDomain+"."+setting.Dnspod.Domain
	}


	update_message:="[======> ddns update]"+domain+" "+request_ip+" "+
		runtime.GOOS+" "+runtime.GOARCH+" "+time.Now().Format("15:04:05")

	chSend <- update_message

	request_ip=""
	return true

}