package main

import (
    "html/template"
    "log"
    "net/http"
)

func server(res http.ResponseWriter, req *http.Request) {
    tpl, err := template.ParseFiles("templates/base.html")
    if err != nil {
        log.Fatalln(err)
    }
    tpl.Execute(res, nil)
}

func main() {
    http.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("img"))))
    http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("css"))))
    http.HandleFunc("/", server)
    http.ListenAndServe(":8080", nil)
}