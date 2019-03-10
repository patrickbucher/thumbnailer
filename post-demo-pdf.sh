#!/bin/sh

curl -X POST -F width=500 -F height=708 -F density=300 -F pdf=@demo.pdf localhost:1337/thumbnail > thumbnail.png
