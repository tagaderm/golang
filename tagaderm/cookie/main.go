package main

import (
    "fmt"
    "net/http"
    "github.com/nu7hatch/gouuid"
)

func main() {
    http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
        cook, err := req.Cookie("my")
        if err != nil {
            id, _ := uuid.NewV4()
            cook = &http.Cookie{
                Name:  "the-session",
                Value: id.String(),

                HttpOnly: true,
            }
            http.SetCookie(res, cook)
        }
        fmt.Printf("%v", cook)
    })
    http.ListenAndServe(":8080", nil)
}