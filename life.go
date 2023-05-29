package main

import (
	"fmt"

)

func main(){
	defer func (){
		fmt.Println("Ending")
	}()
	fmt.Print("We are in")
	study()
	finances()
	marriage()
	hike()
	bike()
	read()
	travel()
}