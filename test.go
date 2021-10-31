package main

import (
	"fmt"
	"github.com/rmoff/ksqldb-go"
)

func main() {

	client := ksqldb.NewClient("http://192.168.136.134:8088","","")
	err := client.Execute( "DROP STREAM 039F3E5A4F0647B9ACA3E4ACBBF08B61ORIGINAL;")
	if err!=nil{
		fmt.Println("LOI ROI")
		fmt.Println(err)
	} else {
		fmt.Println("execute")
	}
}
