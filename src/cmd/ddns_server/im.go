package main

import (
	"encoding/json"
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"text/template"
)

// Define our message object
type Message struct {
	To string `json:"to"`
	From string `json:"from"`
	Type string `json:"type"`
	Message  string `json:"message"`
}
// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type IMServer struct {
	clients sync.Map
	//clients map[string]*websocket.Conn
	msgqueue chan Message
}

var  chatTemplate *template.Template

func StartIMServer(r *mux.Router,templateBox *rice.Box) *IMServer {

	im:=IMServer{msgqueue:make(chan Message,100)}

	templateChat, err := templateBox.String("chat.html")
	if err != nil {
		log.Fatal(err)
	}

	chatTemplate,err=template.New("home").Parse(templateChat)
	if err != nil {
		log.Fatal(err)
	}
	r.HandleFunc("/chat",ChatHandler)
	r.HandleFunc("/chat/{room}",ChatHandler)
	r.HandleFunc("/ws/chat/{room}",im.WebsocketHandler)
	r.HandleFunc("/ws/chat",im.WebsocketHandler)
	r.HandleFunc("/ws",im.WebsocketHandler)
	r.HandleFunc("/im/command/", im.IMCommandHandler).Methods("POST")

	go im.deliverMessage()

	return &im
}

func (i*IMServer)BroadcastMessage(m Message){
	select {
	case i.msgqueue<-m:
		break
	default:
		fmt.Println("---fail to add message----")
	}
}

func (i*IMServer)sendMessage(m Message){

	select {
	case i.msgqueue<-m:
		break
	default:
		fmt.Println("---fail to add message----")
	}
}
func (i*IMServer)BroadcastText(log []byte){
	select {
	case i.msgqueue<-Message{Type:"text",Message:string(log)}:
		break
	default:
		fmt.Println("---fail to add message----")
	}


}
func (i*IMServer)WebsocketHandler(w http.ResponseWriter, r *http.Request) {

	params:=mux.Vars(r)
	room:=params["room"]


	//fmt.Println("ws client connected ",myaddr,room)

	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("upgrade fail:%v\n",err)
		return
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	myaddr:=r.RemoteAddr
	// Register our new client
	if(room==""){
		i.clients.Store(myaddr,ws)

	}else{
		myaddr="@"+room+"-"+myaddr
		i.clients.Store(myaddr,ws)
	}
	i.msgqueue<-Message{To:myaddr,Type:"whoareyou",Message:r.RemoteAddr}
	for {
		var msg Message


		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			i.clients.Delete(myaddr)
			break
		}
		if(msg.Type=="setwhoami"){

			i.clients.Store(msg.Message,ws)
			i.clients.Delete(myaddr)
			myaddr=msg.Message
			i.msgqueue<-Message{Type:"result",Message:"ok"}

			continue
		}else if(msg.To=="")&&(room!=""){
			msg.To="@"+room+"-"
		}else if(msg.Type=="hb"){
			continue
		}
		// Send the newly received message to the broadcast channel
		i.msgqueue <- msg
	}
}
func (i*IMServer)deliverMessage() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-i.msgqueue

		if(msg.To==""){
			// Send it out to every client(not room numbers) that is currently connected
			i.clients.Range(func(names, conn interface{}) bool {
				client,_:=conn.(*websocket.Conn)
				name,_:=names.(string)
				if strings.HasPrefix(name,"@")==false{
					err := client.WriteJSON(msg)
					if err != nil {
						client.Close()
						i.clients.Delete(name)
					}
				}

				return true
			})
		}else if(strings.HasPrefix(msg.To,"@")&&strings.HasSuffix(msg.To,"-")) {
			// Send it to room numbers
			i.clients.Range(func(names, conn interface{}) bool {
				client,_:=conn.(*websocket.Conn)
				name,_:=names.(string)
				if(strings.HasPrefix(name,msg.To)){
					err := client.WriteJSON(msg)
					if err != nil {
						client.Close()
						i.clients.Delete(name)
					}
				}

				return true
			})
		}else{
			tow,_ :=i.clients.Load(msg.To)
			tows,_:=tow.(*websocket.Conn)
			if(tows ==nil){
				fmt.Printf("no dist found:%v\n",msg.To)
			}else{
				err := tows.WriteJSON(msg)
				if err != nil {
					tows.Close()
					i.clients.Delete(msg.To)
				}
			}
		}
	}


}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	chatTemplate.Execute(w, nil)
}
func (i*IMServer)IMCommandHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	command := r.FormValue("command")
	switch command {
	case "alluser":
		type user struct {
			Name string
		}

		ul:=[]user{}

		i.clients.Range(func(client, value interface{}) bool {
			name,_:=client.(string)
			ul=append(ul,user{Name:name})
			return true
		})

		sort.Slice(ul, func(i, j int) bool { return ul[i].Name < ul[j].Name })
		json.NewEncoder(w).Encode(ul)

		break

	case "send":
		to:=r.FormValue("to")
		msg:=r.FormValue("message")
		i.sendMessage(Message{To:to,Message:msg})
		fmt.Fprintf(w,"{\"result\": \"ok\"}")
		break;

	case "check":
		to:=r.FormValue("name")
		if(to==""){
			fmt.Fprintf(w,"{\"result\": \"fail\"}")
		}

		ws,_:=i.clients.Load(to)
		if(ws!=nil){
			fmt.Fprintf(w,"{\"result\": \"connected\"}")
		}else {
			fmt.Fprintf(w,"{\"result\": \"notconnected\"}")
		}

		break;
	default:
		fmt.Printf("im command NOT processed:%v\n",command)
		break;

	}
}
