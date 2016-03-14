package main

import (
    "fmt"
    "net/http"
)

func Inputpage(res http.ResponseWriter, req *http.Request) {
    fmt.Fprint(res, `<!DOCTYPE html>
                    <html>
                    <body>
                        <form>
                            Input your name:<br>
                            <input type="text" name="inputField"><br>
                        </form>
                    </body>
                    </html>`)
    fmt.Fprint(res, "<font style=\" font-size: 30px;\">"+req.FormValue("inputField")+"</font>")
}

func main() {
    http.HandleFunc("/", Inputpage)
    http.ListenAndServe(":8080", nil)
}