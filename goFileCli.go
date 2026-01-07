package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/redis/go-redis/v9"
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
			return fmt.Errorf("path does not exist : %s\n", path)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("premission error : %s\n", path)
		}
		// other
		return err
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("this is not a dir : %s\n", path)
	}
	// can i open it ?
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can't open the dir : %v\n", err)
	}
	f.Close()
	return nil
}

var cntx = context.Background()

func main() {
	// args
	if len(os.Args) < 2 {
		usageMsg()
		return
	}
	vlkClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	switch os.Args[1] {
	case "-h", "--help":
		usageMsg()
		break
	case "-u":
		if len(os.Args) < 3 {
			fmt.Printf(" -u need the directory name parm\n")
			return
		}
		originDir := os.Args[2]
		err := validateDir(originDir)
		if err != nil {
			fmt.Printf(" error: \n %v\n", err)
			os.Exit(1)
		}
		// do upload
		fmt.Printf("uploading %s\n", originDir)
		uploadDir(vlkClient, originDir)
		break
	case "-d":
		if len(os.Args) < 4 {
			fmt.Printf(" -d need used with 2 parm : originDir distDir\n")
			return
		}
		valkeyDir := os.Args[2]
		destinationDir := os.Args[3]
		// do download
		fmt.Printf("downloading from %s to %s\n", valkeyDir, destinationDir)
		break
	default:
		usageMsg()
	}
}

func uploadDir(vlkClient *redis.Client, dirPath string) {
	fullPath, err := filepath.Abs(dirPath)
	if err != nil {
		fmt.Printf("can't get the abs path : %s\n", err)
	}
	dirName := filepath.Base(fullPath)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		fmt.Printf("error reading dir : %v\n", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filepath := filepath.Join(fullPath, entry.Name())
			fileContent, err := os.ReadFile(filepath)
			if err != nil {
				fmt.Printf("can't read file %s: %v\n", entry.Name(), err)
				continue
			}
			// save key as dir:file
			key := fmt.Sprintf("%s:%s", dirName, entry.Name())

			err = vlkClient.Set(cntx, key, fileContent, 0).Err()
			if err != nil {
				fmt.Printf("can't upload file %s , %v\n", entry.Name(), err)
			} else {
				fmt.Printf("done uploading : %s\n", entry.Name())
			}

		}
	}
}
