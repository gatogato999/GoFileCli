package main

import (
	"fmt"
	"os"
)

func usageMsg() {
	fmt.Println("tool usage")
	fmt.Println("gofilecli -h  --> show help ")
	fmt.Println("gofilecli -u origindir  --> upload the dir and its contents ")
	fmt.Println(
		"gofilecli -d valkeyDir destinationDir --> download the valkeyDir into the destinationDir",
	)
}

func main() {
	// args
	if len(os.Args) < 2 {
		usageMsg()
		return
	}
	switch os.Args[1] {
	case "-h", "--help":
		usageMsg()
		break
	case "-u":
		if len(os.Args) < 3 {
			fmt.Printf(" -u need the directory name parm")
			return
		}
		originDir := os.Args[2]
		// do upload
		fmt.Printf("upload %s", originDir)
		break
	case "-d":
		if len(os.Args) < 4 {
			fmt.Printf(" -d need used with 2 parm : originDir distDir")
			return
		}
		valkeyDir := os.Args[2]
		destinationDir := os.Args[3]
		// do download
		fmt.Printf(" from %s to %s", valkeyDir, destinationDir)
		break
	default:
		usageMsg()
	}
}
