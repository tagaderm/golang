package main

import (
    "fmt"
    "io"
    "net/http"
)

func uploadWebpage(res http.ResponseWriter, req *http.Request) {
    page := `
    <!DOCTYPE html>
    <html>
      <head>
        <meta charset="utf-8">
        <title></title>
      </head>
      <body>
        <form method = "POST" enctype="multipart/form-data"> Input your name:
          <input type="file" name="name"><br>
          <input type="submit">
        </form>
      </body>
    </html>`

    io.WriteString(res, page)

    if req.Method == "POST" {
        _, src, err := req.FormFile("name")
        if err != nil {
            fmt.Println(err)
        }

        dst, err := src.Open()
        if err != nil {
            fmt.Println(err)
        }

        io.Copy(res, dst)
    }
}

func main() {
    http.HandleFunc("/", uploadWebpage)
    http.ListenAndServe(":8080", nil)
}