package main

import (
    "fmt"
    "kite-admin/backend/config"
    "kite-admin/backend/models"
)

func main() {
    config.ConnectDatabase()
    err := config.DB.Create(&models.Attachment{FileName: "test", OriginalName: "test"}).Error
    if err != nil {
        fmt.Println("ERROR:", err)
    } else {
        fmt.Println("SUCCESS")
    }
}
