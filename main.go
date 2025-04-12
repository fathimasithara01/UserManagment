package main

import (
	"fmt"

	"github.com/fathima-sithara/UserManagment/controllers"
	"github.com/fathima-sithara/UserManagment/initalizeres"
	"github.com/fathima-sithara/UserManagment/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initalizeres.LoadEncVariable()
	initalizeres.ConnectToDb()
	initalizeres.Pooling()
	initalizeres.SyncDB()

}

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.LoginUser)
	r.POST("/logout", controllers.LogOut)

	r.GET("/validate", middleware.RequireAuth)

	fmt.Println("server running")
	r.Run(":8080")
}
