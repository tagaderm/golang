package main

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "github.com/nu7hatch/gouuid"
    "html/template"
    "log"
    "net/http"
)

type User struct {
    Name     string
    Age      string
    Sex      string
    Location string
}

func userContrusctor(name string, age string, sex string, location string) string {
    user := User{
        Name:     name,
        Age:      age,
        Sex:      sex,
        Location: location,
    }

    encodeToJSon, err := json.Marshal(user)
    if err != nil {
        fmt.Printf("error: ", err)
    }

    finalPackage64Encode := base64.URLEncoding.EncodeToString(encodeToJSon)
    return finalPackage64Encode

}

func handleIt(res http.ResponseWriter, req *http.Request) {
    templatePage, err := template.ParseFiles("templates/proj_form.html")
    if err != nil {
        log.Fatalln(err)
    }
    name := req.FormValue("name")
    age := req.FormValue("age")
    sex := req.FormValue("sex")
    location := req.FormValue("location")

    encodedData := userContrusctor(name, age, sex, location)

    cookie, err := req.Cookie("session-fino")
    if err != nil {
        id, _ := uuid.NewV4()
        cookie = &http.Cookie{
            // Secure: true,
            Name:     "session-fino",
            Value:    id.String() + "," + name + "," + age + "," + sex + "," + location + "," + encodedData,
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