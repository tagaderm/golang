package checker

import (
    "html/template"
    "net/http"
    "google.golang.org/appengine"
    "google.golang.org/appengine/datastore"
    "google.golang.org/appengine/log"
    "io"
    "io/ioutil"
)

type Word struct {
    Name string
}

var tpl *template.Template

func init() {
    http.HandleFunc("/", index)
    http.HandleFunc("/api/check", wordCheck)

    // serve public resources
    http.Handle("/css_js/", http.StripPrefix("/css_js", http.FileServer(http.Dir("public"))))

    // parse templates
    tpl = template.Must(template.ParseGlob("*.html"))
}

func index(res http.ResponseWriter, req *http.Request) {

    if req.Method == "POST" {

                // get the word from the form's input box.
        var w Word
        w.Name = req.FormValue("word")

        ctx := appengine.NewContext(req)
        log.Infof(ctx, "WORD SUBMITTED: %v", w.Name)

                // save word into the datastore.
        key := datastore.NewKey(ctx, "Dictionary", w.Name, 0, nil)
        _, err := datastore.Put(ctx, key, &w)
        if err != nil {
            http.Error(res, err.Error(), http.StatusInternalServerError)
            return
        }
    }
    tpl.ExecuteTemplate(res, "ajax.html", nil)
}

func wordCheck(res http.ResponseWriter, req *http.Request) {

    ctx := appengine.NewContext(req)

    // retrieve the incoming word as it is typed.
    var w Word
    bs, err := ioutil.ReadAll(req.Body)
    //
    log.Infof(ctx, "Received information: %v", string(bs))
    //
    if err != nil {
        log.Infof(ctx, err.Error())
    }
    w.Name = string(bs)
    log.Infof(ctx, "ENTERED wordCheck - w.Name: %v", w.Name)

    // check the incoming word against what is currently in the datastore
    key := datastore.NewKey(ctx, "Dictionary", w.Name, 0, nil)
    err = datastore.Get(ctx, key, &w)
    if err != nil {
        io.WriteString(res, "false")
        return
    }
    io.WriteString(res, "true")
}