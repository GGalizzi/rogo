package main

import (
	"fmt"
	"os"
	"runtime"
)

func init() {
	runtime.LockOSThread()

	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		fmt.Println("Can't open log.txt -- ", err)
	}

	err = f.Truncate(0)
	if err != nil {
		fmt.Println("Can't truncate log.txt -- ", err)
	}

}

func main() {
	G := NewGame()

	G.run()
}
