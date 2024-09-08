package main

import (
	_ "github.com/alioth-center/akasha-whisper/app/router"
	"github.com/alioth-center/infrastructure/exit"
)

func main() {
	exit.BlockedUntilTerminate()
}
