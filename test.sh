#!/bin/sh

curl -X POST -F width=200 -F format=png -F pdf=@demo.pdf localhost:1337/thumbnail > thumbnail.png
curl -X POST -F height=708 -F density=150 -F quality=95 -F format=jpg -F pdf=@demo.pdf localhost:1337/thumbnail > thumbnail.jpg
