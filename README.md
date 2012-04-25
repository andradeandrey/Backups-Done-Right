# Backups-Done-Right

Backups-Done-Right is a P2P backup program providing easy, fast and secure encrypted off-site backups.


## Features

* file transfer is always encrypted
* very fast filesystem walker
* posibility to run more than one fs walker at once (huge speed increasing on multi HDD systems)
* simple configuration - just one config file
* allows to have one server and multiple clients
* simple installation (static linked build)
* restores with permissions, symlinks etc.
* open source


## Maintainers

Backups-Done-Right does have two project maintainers:

* Bill Broaldey   aka spikebike	<bill@broadley.org>	(english)
* Joel Bodenmann  aka Tectu	<joel@unormal.org>	(german / english)


## Technical Description

Once the filesystem walker created a database for the directories that have to be backed up, it will just update the database on every run. On each run the walker decides if the file got any changes. If yes, the file gets encrypted over AES-512 and gets uploaded to the backup server over an SSL secured TCP/IP connection. The server keeps the files encryptet - this is much safer when we have more than one client on the same backup server.
Whenever we need a backup, we can select which files need to be restored (mostly over the last-modified time). The server will send the encrypted files to the matching client over an SSL secured TCP/IP connection again. The client will then decrypt the received files and restores the complete directory tree with all the permissions, symlinks etc. 


## Build

Backups Done Right depends on 3 external go packages that need to be installed:

goconfig - to install, simply run:

	$ go get github.com/kless/goconfig/gonfig


goconfig - to install, simply run:

	$ go get github.com/mattn/go-sqlite3


go-rpcgen - to install, simply run:

	$ go get github.com/kylelemons/go-rpcgen/protoc-gen-go


Please note that Backups-Done-Right does also require sqlite3 which is not a part of go or Backups-Done-Right itself. 
Installing sqlite3 depends on your system. Therefor we cannot give you an installing howto. If you don't know how to install sqlite3 on your system, we recommend to use google to find out.


## Run

Before you can run the software the first time, you need to create a config file which fits your needs. Please copy the example config file:

	$ cd Backups-Done-Right/etc
	$ cp config.cfg.example config.cfg

Then, edit the config file to your needs.

You do also need certificates for the SSL encryption:

	$ cd Backups-Done-Right/src
	$ ./makecerts <your_email_address>


## Misc

Please see documentation/* for additional informations

