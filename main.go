package main

import (
	"log"
)
 
func main(){
	store,err:= NewPostgressStore()
	if err!=nil{
		log.Fatal("error")
	}
	er1,er2,er3,er4,er5:=store.init()
	if er1!=nil{
		log.Fatal(er1)
	}
	if er2!=nil{
		log.Fatal(er2)
	}
	if er3!=nil{
		log.Fatal(er3)
	}
	if er4!=nil{
		log.Fatal(er4)
	}
	if er5!=nil{
		log.Fatal(er5)
	}

	server:=NewApiServer(":8080",store)
	server.Run()
}