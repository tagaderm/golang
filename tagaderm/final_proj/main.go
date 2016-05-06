package main

import (
    "fmt"
    "time"
    "html/template"
    "log"
    "net/http"
    // "encoding/base64"
    "encoding/json"
    "github.com/nu7hatch/gouuid"
    "github.com/boltdb/bolt"
)


type visit struct {
    IsNew  bool
}

type User struct {
    Username string
    Email string
    Password string
}

func server(res http.ResponseWriter, req *http.Request){
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

func signup(res http.ResponseWriter, req *http.Request){
    templatePage, err := template.ParseFiles("templates/sign_up_form.html")
    if err != nil {
        log.Fatalln(err)
    }
    if req.Method == "POST"{
        // templatePage, err := template.ParseFiles("templates/sign_up_form.html")
        // if err != nil {
        //     log.Fatalln(err)
        // }
        username := req.FormValue("username")
        email := req.FormValue("email")
        password := req.FormValue("password")

        user := User{
            Username: username,
            Email: email,
            Password: password,
        }

        encodeToJSon, err := json.Marshal(user)
        if err != nil {
            fmt.Printf("error: ", err)
        }

        db, err := bolt.Open("final_proj.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
        if err != nil {
            log.Fatal(err)
        }

        db.Update(func(tx *bolt.Tx) error {
            b, err := tx.CreateBucketIfNotExists([]byte("users"))
            if err != nil {
                return err
            }

            return b.Put([]byte(user.Username),[]byte(encodeToJSon))
        })
    }
    templatePage.Execute(res, nil)
}

func usernameCheck(res http.ResponseWriter, req *http.Request){
    
}

func main() {
    http.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("img"))))
    http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("css"))))
    http.HandleFunc("/", server)
    http.HandleFunc("/signup", signup)
    http.HandleFunc("/api/username_check", usernameCheck)
    http.ListenAndServe(":8080", nil)
}