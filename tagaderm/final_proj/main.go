package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "github.com/nu7hatch/gouuid"
    // "os"
)


type visit struct {
    IsNew  bool
}

func server(res http.ResponseWriter, req *http.Request) {
    obj := visit{
            IsNew: false,
        }

    cookie, err := req.Cookie("session_id")

    if err != nil {
        obj.IsNew = true
        id, _ := uuid.NewV4()
        cookie = &http.Cookie{
            Name:  "session_id",
            Value: id.String(),
            // Secure: true
            HttpOnly: true,
        }
        http.SetCookie(res, cookie)
    }
    fmt.Println(cookie)

    tpl, err := template.ParseFiles("templates/base.html")
    if err != nil {
        log.Fatalln(err)
    }

    err = tpl.Execute(res, obj)
    if err != nil{
        log.Fatalln(err)
    }
}

func main() {
    http.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("img"))))
    http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("css"))))
    http.HandleFunc("/", server)
    http.ListenAndServe(":8080", nil)
}