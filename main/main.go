package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"webssh/internal"
)

func main(){
	r := gin.Default()
	r.GET("/term", internal.WsSsh)
	log.Fatal(r.Run(":3000"))
}
