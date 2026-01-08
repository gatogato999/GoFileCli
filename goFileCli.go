package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	// check the -r flag and then rmove it from the args array
	// isRecursive := false

	toolArgs := os.Args[1:]
	if len(toolArgs) < 1 {
		usageMsg()
		return
	}

	isRecursive := false
	for index, arg := range toolArgs {
		if arg == "-r" {
			isRecursive = true
			toolArgs = append(toolArgs[:index], toolArgs[index+1:]...)
			break
		}
	}
	fmt.Println(isRecursive)
	vlkClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	switch toolArgs[0] {
	case "-h", "--help":
		usageMsg()
	case "-u":
		if len(toolArgs) < 2 {
			fmt.Printf(" -u need the directory name parm\n")
			return
		}
		originDir := toolArgs[1]
		err := validateDir(originDir)
		if err != nil {
			fmt.Printf(" error: \n %v\n", err)
			os.Exit(1)
		}
		// do upload
		fmt.Printf("uploading %s\n", originDir)
		uploadDir(vlkClient, originDir, isRecursive)
	case "-d":
		if len(toolArgs) < 3 {
			fmt.Printf(" -d need used with 2 parm : originDir distDir\n")
			return
		}
		valkeyDir := toolArgs[1]
		destinationDir := toolArgs[2]
		// do download
		fmt.Printf("downloading from %s to %s\n", valkeyDir, destinationDir)
		downloadDir(vlkClient, valkeyDir, destinationDir, isRecursive)
	default:
		usageMsg()
	}
}

func uploadDir(vlkClient *redis.Client, inputDir string, isRecursive bool) {
	fullPath := filepath.Clean(inputDir)
	info, err := os.Stat(fullPath)
	if err != nil || !info.IsDir() {
		fmt.Printf("error : %s is not a dir\n", fullPath)
		return
	}
	// process all files
	processFile := func(path string, dir os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if dir.IsDir() {

			if !isRecursive && path != fullPath {
				return filepath.SkipDir
			}
			return nil
		}

		key := strings.ReplaceAll(path, string(os.PathSeparator), "/")

		fileContent, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf(" can't read %s \n %v\n", path, err)
			return nil
		}
		err = vlkClient.Set(cntx, key, fileContent, 0).Err()
		if err != nil {
			fmt.Printf("\ncan't upload %s \n", err)
		} else {
			fmt.Printf("\ndone uploading : %s\n", key)
		}
		return nil
	}
	if isRecursive {
		filepath.WalkDir(fullPath, processFile)
	} else {
		dirContents, _ := os.ReadDir(fullPath)
		for _, subDir := range dirContents {
			if !subDir.IsDir() {
				dirPath := filepath.Join(fullPath, subDir.Name())
				processFile(dirPath, subDir, nil)
			}
		}
	}
}

func downloadDir(
	vlkClient *redis.Client,
	valkeyDir string,
	destinationDir string,
	isRecursive bool,
) {
	// search the database
	pattern := valkeyDir + ":*"
	var cursor uint64
	isFirstRun := true
	for {
		// scan doesn't block the db
		// cursor used by valkey to remeber where it left off
		keys, nextCursor, err := vlkClient.Scan(cntx, cursor, pattern, 10).Result()
		if err != nil {
			fmt.Printf("error scanning db : %v\n", err)
			return
		}

		if isFirstRun && len(keys) == 0 {
			fmt.Printf("error : directory '%s' not found id valkey \n", valkeyDir)
			return
		}

		if isFirstRun {
			if err := os.MkdirAll(destinationDir, 0o755); err != nil {
				fmt.Printf("can't create dir %s : %v\n", destinationDir, err)
			}
			isFirstRun = false

		}

		for _, key := range keys {
			// file name
			keyName := strings.SplitN(key, ":", 2)
			if len(keyName) < 2 {
				continue
			}
			fileName := keyName[1]

			fileContents, err := vlkClient.Get(cntx, key).Bytes()
			if err != nil {
				fmt.Printf("error : can't get %s : %v\n", key, err)
				continue
			}

			destinationFile := filepath.Join(destinationDir, fileName)
			err = os.WriteFile(destinationFile, fileContents, 0o644) // rw- r-- r--
			if err != nil {
				fmt.Printf("error : failed to write %s: %v \n ", destinationFile, err)
			} else {
				fmt.Printf("done downloading : %s ---> %s\n", key, destinationFile)
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			// scan is finished
			break
		}
	}
	fmt.Printf("\n finished downloading '%s' ---> '%s", valkeyDir, destinationDir)
}
