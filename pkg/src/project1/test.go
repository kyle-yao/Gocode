package main

import (
	"fmt"
	"reflect"
)

type computer struct {
	UI     string `info:"GUI"`
	System string `info:"System"`
}

func (this *computer) GetUI() string {
	fmt.Println("UI = ", this.UI)
	return this.UI
}

func (this *computer) SetSys(Systemname string) {
	this.System = Systemname

}
func findTag(str interface{}) {
	t := reflect.TypeOf(str).Elem()
	for i := 0; i < t.NumField(); i++ {
		temp := t.Field(i).Tag.Get("info")
		fmt.Println("info=", temp)

	}
}

func main() {

	var com computer
	findTag(&com)
}

//var book1 Book
//book1.title = "Golang"
//book1.auth = "gg"
//
//fmt.Printf("%v",book1)

//var myMap1 map[string]string
//if myMap1 ==nil{
//	fmt.Println("myMap1 is empty")
//}
////
//myMap2:= make(map[string]string, 10)
//myMap1= make(map[string]string, 10)
//
//myMap1["one"] = "java"
//myMap1["two"] = "c++"
//myMap1["three"] = "python"
//
//fmt.Println(myMap1)
//fmt.Println(myMap2)
