package main

import (
    "io"
    "net/http"
    "strconv"
)

func main() {
    http.HandleFunc("/", page)
    http.ListenAndServe(":8080", nil)
}

func page(res http.ResponseWriter, req *http.Request) {
    if req.URL.Path != "/" {
        http.NotFound(res, req)
        return
    }

    cookie, err := req.Cookie("pecan_sandy")

    if err == http.ErrNoCookie {
        cookie = &http.Cookie{
            Name:  "pecan_sandy",
            Value: "0",
        }
    }
    visit_count, _ := strconv.Atoi(cookie.Value)
    visit_count++
    cookie.Value = strconv.Itoa(visit_count)
    http.SetCookie(res, cookie)
    io.WriteString(res, cookie.Value)
}