package main

import (
	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)
var homeTemplate *template.Template
func main() {
	// find a rice.Box
	templateBox:= rice.MustFindBox("../../public")

	// get file contents as string
	templateIndex, err := templateBox.String("index.html")
	if err != nil {
		log.Fatal(err)
	}

	homeTemplate,err=template.New("home").Parse(templateIndex)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()


	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(templateBox.HTTPBox())))

	r.HandleFunc("/", HomeHandler).Methods("GET")

	im = StartIMServer(r, templateBox)

	log.SetOutput(new(wsout))



	http.Handle("/", r)

	log.Println("running server on 8888")
	errh:=http.ListenAndServe(":8888", r)
	log.Print(errh)

}
// HomeHandler will be rendering the index.html template, it is written in Vue
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, nil)
}


var im *IMServer
type wsout struct {
	io.Writer
}
func (do *wsout)Write(b[]byte) (int, error){


	os.Stdout.Write(b)
	im.BroadcastText(b)

	return 0,nil

}