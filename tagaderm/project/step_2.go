package main

import (
    "github.com/nu7hatch/gouuid"
    "html/template"
    "log"
    "net/http"
)

func handleIt(res http.ResponseWriter, req *http.Request) {
    templatePage, err := template.ParseFiles("templates/base.html")
    if err != nil {
        log.Fatalln(err)
    }

    cookie, err := req.Cookie("session-fino")
    if err != nil {
        id, _ := uuid.NewV4()
        cookie = &http.Cookie{
            // Secure: true,
            Name:     "session-fino",
            Value:    id.String(),
            HttpOnly: true,
        }
    }

    http.SetCookie(res, cookie)
    templatePage.Execute(res, nil)
}

func main() {
    http.HandleFunc("/", handleIt)
    http.ListenAndServe(":8080", nil)
}