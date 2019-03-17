#!/bin/sh

curl -X POST -F pdf=@../demo.pdf http://thumbnailer.192.168.42.154.nip.io/thumbnail > thumbnail.png 
