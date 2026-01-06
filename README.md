# GoFileCli
a go cli tool to upload and download files from and into a valkey database.

# usage

- uploading a file
Provide directory (folder) name in parameter for this app to let it fetches this directory contents (files) and upload them in valkey database as key-values 
``` bash
# -u : upload 
gofilecli -u ~/photos
```


- downloading a file
read folder from valkey database and copy them locally in different directory
``` bash
# the newFolder will be created if not exist
# -d : download 
gofilecli -d ~/photos ~/newFolder
```

# check list of mandatory features 
- [] verify dir exists
- [] confirm connection with the valkey database
- [] check args
- [] can upload all files inside dir (skip sub-dir for now)
- [] can download a dir and it's content from the database
- [] if downloading to non-existent dir create it
- [] permission error handling 

# features to consider 
- [] handling sub dirs
- [] the order of arguments
- [] handling different types of systems (path differ from other OS)
- [] preserve file attributes as possible
