package prepareTestWeb

import (
	"encoding/gob"
	db "github.com/richkule/prepareTestWeb/DBWorker"
	hd "github.com/richkule/prepareTestWeb/handlers"
	in "github.com/richkule/prepareTestWeb/init"
	"log"
	"net/http"
)

// Устанавливает все handler для сервера
func setHandler() {
	http.HandleFunc("/", hd.MakeHandler(hd.Index))
	http.HandleFunc("/reg", hd.MakeHandler(hd.Reg))
	http.HandleFunc("/auto", hd.MakeHandler(hd.Auto))
	http.HandleFunc("/login", hd.MakeHandler(hd.Login))
	http.HandleFunc("/logout", hd.MakeHandler(hd.Logout))
	http.HandleFunc("/tests", hd.MakeHandler(hd.Test))
	http.HandleFunc("/create", hd.MakeHandler(hd.Create))
	http.HandleFunc("/newTest", hd.MakeHandler(hd.NewTest))
	http.HandleFunc("/edit/", hd.MakeHandler(hd.Edit))
	http.HandleFunc("/newTopic", hd.MakeHandler(hd.NewTopic))
	http.Handle("/static/", http.FileServer(http.Dir("")))
}

func Init() {
	gob.Register(in.SessUs{})
	setHandler()
	if err := db.DbConnect(); err != nil {
		log.Fatal(err.Error())
	}
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
