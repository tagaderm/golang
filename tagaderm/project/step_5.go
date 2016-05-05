package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "fmt"
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

func repackJSON(jsonValues User) string {
    b, _ := json.Marshal(jsonValues)
    return base64.StdEncoding.EncodeToString(b)
}

func unpackJSON(cookie *http.Cookie) (User, bool) {
    decode, _ := base64.StdEncoding.DecodeString(cookie.Value)
    var jsonValues User
    json.Unmarshal(decode, &jsonValues)
    if hmac.Equal([]byte(jsonValues.Hmac), []byte(getCode(jsonValues.Uuid+jsonValues.Name+jsonValues.Age))) {
        return jsonValues, true
    }
    return jsonValues, false
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

func HMACFunction(data string) string {
    h := hmac.New(sha256.New, []byte(data+"securekey"))
    return string(h.Sum(nil))
}

func updateCookie(cookie *http.Cookie, req *http.Request) string {
    jsonValues, _ := unpackJSON(cookie)
    jsonValues.Name = req.FormValue("name")
    jsonValues.Age = req.FormValue("age")
    jsonValues.Hmac = getCode(jsonValues.Uuid + jsonValues.Name + jsonValues.Age)
    return repackJSON(jsonValues)
}

func handleIt(res http.ResponseWriter, req *http.Request) {
    templatePage, err := template.ParseFiles("proj_form.html")
    if err != nil {
        log.Fatalln(err)
    }

    cookie, err := req.Cookie("session-fino")
    if err != nil {
        cookie = defaultCookie()
        http.SetCookie(res, cookie)
    }
    if req.Method == "POST" {
        cookie.Value = updateCookie(cookie, req)
    }
    obj, valid := unpackJSON(cookie)
    if valid {
        t, _ := template.New("name").Parse(file1)
        t.Execute(res, obj)
    }
}

func main() {
    http.HandleFunc("/", handleIt)
    http.ListenAndServe(":8080", nil)
}