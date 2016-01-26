package main

import "fmt"

func main() {
    var x int
    var y int

    fmt.Print("Enter a small number and then large number: ")
    fmt.Scanln(&x, &y)

    z := y%x

    fmt.Println(z)
}