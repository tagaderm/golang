
package main

import (
    "html/template"
    "log"
    "net/http"
)

func handleIt(res http.ResponseWriter, req *http.Request) {
    templatePage, err := template.ParseFiles("templates/base.html")
    if err != nil {
        log.Fatalln(err)
    }

    templatePage.Execute(res, nil)
}

func main() {
    http.HandleFunc("/", handleIt)
    http.ListenAndServe(":8080", nil)
}