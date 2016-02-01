package main

import (
    "log"
    "os"
    "text/template"
)

type entity struct {
    Name string
    IsAlive  bool
}

func main() {
    obj := entity{
        Name: "Number 5",
        IsAlive: true,
    }

    tpl, err := template.ParseFiles("cond_logic.gohtml")
    if err != nil {
        log.Fatalln(err)
    }

    err = tpl.Execute(os.Stdout, obj)
    if err != nil {
        log.Fatalln(err)
    }
}