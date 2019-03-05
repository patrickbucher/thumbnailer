#!/bin/sh

curl -X POST -F pdf=@demo.pdf localhost:1337/thumbnail > thumbnail.png
