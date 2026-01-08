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

- recursive flag
when applying the `-r` flag for either `-u` or `-d` the tool do the operation (upload or download) for the sub-dir of the specified dir.
``` bash 
# -r can be any where
# in normal upload mode 
goFileCli -u originDir # will skip the sub-dir inside of 'origindir'
# using -r flag 
goFileCli -u originDir -r # will upload the direct files and sub-dir's files
# origindir:directFile.txt 
# originDir/subDir:directFile.txt
goFileCli -d originDir dist # originDir has sub dir 
# only direct files in originDir will be inside dist
goFileCli -d origindir dist -r
# all files and sub dirs in origindir will be in dist (full copy)
```

## check list of mandatory features 
- [x] verify dir exists on the device
- [x] confirm connection with the valkey database
- [x] check args
- [x] can upload all files inside dir (skip sub-dir for now)
- [x] can download a dir and it's content from the database
- [x] if downloading to non-existent dir create it
- [x] add a `-r` flag that imply recursive operations of subDirs.
- [ ] permission error handling 
- [ ] modularize the tool (insdeat of a single file)
- [ ] refactor the error handling (too many `if != err`)

### things to consider 
- [ ] the override of the data if the given input is the same.
- [x] handling sub dirs
- [x] handling different types of systems (path differ from other OS)
- [ ] preserve file attributes as possible

#### notes
 
- valkey behaves like a remote dictionary.
- each peice of data is stored under a unique key.
- values can be accessed and modified using the associated keys
- SCAN is better than KEYS (non-block)
