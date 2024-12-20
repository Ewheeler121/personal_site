package main

import (
    "os"
)

func getStatus() string {
    //gonna do some wonky Steam + Discord API things here
    status := os.Getenv("STATUS")
    if status == "" {
        status = "No Status Found"
    }
    return status
}
