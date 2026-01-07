# GoFileCli
a go cli tool to upload and download files from and into a valkey database.

## usage

- clone the repo
``` git clone ```

- install go 
```fish
sudo pacman -S go 
go version
```

- init go module 
``` go mod init [dir or repo name]```

- install valkey
```fish 
sudo pacman -S valkey 
# for me i needed to start the server manually
sudo systemctl start valkey 
# u can check it with 
valkey-cli ping 
# you will get : "PONG"
```

- install the go-redis/v9 client (thought there is a one for 
valkey called valkey-go this one is compatible with valkey too)
``` go get github.com/redis/go-redis/v9```

- uploading a file
Provide directory (folder) name in parameter for this app to let it fetches this directory contents (files) and upload them in valkey database as key-values 
``` bash
# -u : upload 
goFileCli -u originDir
```

- downloading a file
read folder from valkey database and copy them locally in different directory
``` bash
# the newFolder will be created if not exist
# -d : download 
goFileCli -d dirFromDatabase destinationDir
```

## check list of mandatory features 
- [ x ] verify dir exists on the device
- [ x ] confirm connection with the valkey database
- [ x ] check args
- [ x ] can upload all files inside dir (skip sub-dir for now)
- [ ] can download a dir and it's content from the database
- [ ] if downloading to non-existent dir create it
- [ ] permission error handling 
- [ ] modularize the tool (insdeat of a single file)

### features to consider 
- [ ] handling sub dirs
- [ ] handling different types of systems (path differ from other OS)
- [ ] preserve file attributes as possible
- [ ] consider the override of the data in the valkey db.

#### notes
 
- valkey behaves like a remote dictionary.
- each peice of data is stored under a unique key.
- values can be accessed and modified using the associated keys
