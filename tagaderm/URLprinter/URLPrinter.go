package main

import (
    "fmt"
    "net/http"
)

func urlName(res http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(res, "%v", req.URL.Path)
}

func main() {
    http.HandleFunc("/", urlName)
    http.ListenAndServe(":8080", nil)
}