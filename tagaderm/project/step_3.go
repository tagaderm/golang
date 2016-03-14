package main

import (
    "github.com/nu7hatch/gouuid"
    "html/template"
    "log"
    "net/http"
)

func handleIt(res http.ResponseWriter, req *http.Request) {
    templatePage, err := template.ParseFiles("templates/proj_form.html")
    if err != nil {
        log.Fatalln(err)
    }
    name := req.FormValue("name")
    age := req.FormValue("age")
    sex := req.FormValue("sex")
    location := req.FormValue("location")

    cookie, err := req.Cookie("session-fino")
    if err != nil {
        id, _ := uuid.NewV4()
        cookie = &http.Cookie{
            // Secure: true,
            Name:     "session-fino",
            Value:    id.String() + "," + name + "," + age + "," + sex + "," + location,
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