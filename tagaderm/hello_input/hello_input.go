package main

import "fmt"

func main() {
    fmt.Print("Enter your name: ")
    var x string
    fmt.Scanln(&x)
    fmt.Printf("hello, "+x+"\n")
}