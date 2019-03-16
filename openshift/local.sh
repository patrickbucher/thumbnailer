#!/bin/sh

docker build -t thumbnailer-base - <../Dockerfile 
docker build . -t thumbnailer-s2i
s2i build https://github.com/patrickbucher/thumbnailer.git thumbnailer-s2i thumbnailer
docker run -p 1337:1337 -dit thumbnailer
