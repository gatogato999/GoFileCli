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

func validateDir(path string) error {
	dirInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist %s", path)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("premission error \n %s", path)
		}
		// other
		return err
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("this is not a dir \n %s", path)
	}
	// can i open it ?
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can't open the dir :\n %v", err)
	}
	f.Close()
	return nil
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
		err := validateDir(originDir)
		if err != nil {
			fmt.Printf(" error: \n %v", err)
			os.Exit(1)
		}
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
