package main

import (
    "fmt"
    "time"
    "html/template"
    "log"
    "net/http"
    "io"
    "io/ioutil"
    "encoding/json"
    "strconv"
    "github.com/nu7hatch/gouuid"
    "github.com/boltdb/bolt"
    "github.com/bradfitz/gomemcache/memcache"
)


type visit struct {
    IsNew  bool
}

type User struct {
    Username string
    Email string
    Password string
}

type ZipCode struct {
    Code string `json:"zip_code"`
    Distance float64 `json:"distance"`
    City string `json:"city"`
    State string `json:"state"`
}

type ExpectedJSON struct {
    ZipCodes []ZipCode `json:"zip_codes"`
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

        username := req.FormValue("username")
        email := req.FormValue("email")
        password := req.FormValue("password")

        mc := memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212")
        mc.Set(&memcache.Item{Key: "username", Value: []byte(username)})
        mc.Set(&memcache.Item{Key: "email", Value: []byte(email)})
        mc.Set(&memcache.Item{Key: "username", Value: []byte(password)})

        db, err := bolt.Open("final_proj.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
        if err != nil {
            log.Fatal(err)
        }

        var exists []byte
        db.View(func(tx *bolt.Tx) error {
            b := tx.Bucket([]byte("users"))
            v := b.Get([]byte(username))
            exists = v
            return nil
        })

        db.Close()

        if exists == nil {
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
            db.Close()
            fmt.Println("user created")
            obj := struct {
                    Creadted bool
                } {
                    true,
                }
            templatePage.Execute(res, obj)
        } else {
            fmt.Println("user not created")
            obj := struct {
                    NotCreated bool
                } {
                    true,
                }
            templatePage.Execute(res, obj)
        }
    } else if req.Method == "GET" {
        templatePage.Execute(res, nil)
    }
}

func usernameCheck(res http.ResponseWriter, req *http.Request){
    // acquire the incoming word
    var w string
    bs, err := ioutil.ReadAll(req.Body)
    if err != nil {
        log.Fatal(err)
    }
    w = string(bs)

    // check the incoming username against the db
    db, err := bolt.Open("final_proj.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
    if err != nil {
        log.Fatal(err)
    }
    var exists []byte
    db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("users"))
        v := b.Get([]byte(w))
        exists = v
        return nil
    })
    db.Close()

    if exists == nil {
        io.WriteString(res, "false")
        return
    }
    io.WriteString(res, "true")
}

func update(res http.ResponseWriter, req *http.Request) {
    templatePage, err := template.ParseFiles("templates/update_form.html")
    if err != nil {
        log.Fatalln(err)
    }
    if req.Method == "GET" {
        fmt.Println("in get")
        obj := struct {
                Display bool
            } {
                true,
            }
        templatePage.Execute(res, obj)
        } else if req.Method == "POST"{
            var user User
            username := req.FormValue("check_username")
            if username != "" {

                mc := memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212")
                it, err := mc.Get("username")
                if err != nil {
                    log.Fatalln(err)
                }
                fmt.Println(it)
                if /*username == it*/ true {
                    email, err := mc.Get("username")
                    if err != nil {
                        log.Fatalln(err)
                    }
                    fmt.Println(email)
                } else { //not in memcache
                    var exists []byte
                    db, err := bolt.Open("final_proj.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
                    if err != nil {
                        log.Fatal(err)
                    }
                    db.View(func(tx *bolt.Tx) error {
                        b := tx.Bucket([]byte("users"))
                        v := b.Get([]byte(username))
                        exists = v
                        return nil
                    })
                    db.Close()
                    var m User
                    if err := json.Unmarshal(exists, &m); err != nil {
                        log.Fatal(err)
                    }
                    user = m
                    fmt.Println(m.Email)
                }

                fmt.Println("in post")
                obj := struct {
                        Display bool
                    } {
                        false,
                    }
                templatePage.Execute(res, obj)    
            }
            obj := user
            templatePage.Execute(res, obj)
        }    
}

func upload(res http.ResponseWriter, req *http.Request) {
    tpl, err := template.ParseFiles("templates/upload.html")
    if err != nil {
        log.Fatalln(err)
    }

    tpl.Execute(res, nil)

    // client, err := storage.NewClient(ctx)
    // if err != nil {
    //     log.Fatal(err)
    // }

    // err = tpl.Execute(res, obj)
    // if err != nil{
    //     log.Fatalln(err)
    // }
}

func external(res http.ResponseWriter, req *http.Request) {
    tpl, err := template.ParseFiles("templates/external.html")
    if err != nil {
        log.Fatalln(err)
    }
    if req.Method == "GET" {
        tpl.Execute(res, nil)
    } else if req.Method == "POST" {
        zip := req.FormValue("zip")
        url := "https://www.zipcodeapi.com/rest/JyCX6C8IlSVxTcSVmad12a43G0M5bckUKXSRz0AUTvsMIlI4vW5x6aANanTmzdhk/radius.json/"+zip+"/5/mile"
        resp, err := http.Get(url)
            if err != nil {
            log.Fatalln(err)
        }
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)

        var input ExpectedJSON
        json.Unmarshal(body, &input)
        var table_items string
        table_items += "<tr><th>Zip Code</th><th>Distance</th><th>City</th><th>State</th></tr>"
        for _, element := range input.ZipCodes {
            table_items += "<tr>"
            table_items += "<td>"
            table_items += element.Code
            table_items += "</td>"
            table_items += "<td>"
            table_items += strconv.FormatFloat(element.Distance , 'f', 4, 32)
            table_items += "</td>"
            table_items += "<td>"
            table_items += element.City
            table_items += "</td>"
            table_items += "<td>"
            table_items += element.State
            table_items += "</td>"
            table_items += "</tr>"
        }

        fmt.Fprint(res, "<!DOCTYPE html><html><body><table border=\"1\" cellpadding=\"5\" cellspacing=\"5\">"+table_items+"</table><a href=\"http://localhost:8080/external\">Try Another</a></body></html>")
    }

}

func main() {
    http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
    http.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("img"))))
    http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("css"))))
    http.HandleFunc("/", server)
    http.HandleFunc("/signup", signup)
    http.HandleFunc("/update", update)
    http.HandleFunc("/upload", upload)
    http.HandleFunc("/external", external)
    http.HandleFunc("/api/username_check", usernameCheck)
    http.ListenAndServe(":8080", nil)
}