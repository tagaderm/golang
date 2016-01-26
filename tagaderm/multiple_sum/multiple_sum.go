package main

import "fmt"

func main() {
    var sum int
    for x:=0; x<1000; x++ {
        if x%3 == 0 {
            sum += x
        }else if x%5 == 0 {
            sum += x
        }
    }
    
    fmt.Println(sum)
}