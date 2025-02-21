package main

import (
	"fmt"
	"github.com/RajVerma97/golang-banking-ledger/internal/api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routes.Setup(r)
	PORT := 3000
	fmt.Printf("Server Listening on Port %d", PORT)
	r.Run(fmt.Sprintf(":%d", PORT))
}
