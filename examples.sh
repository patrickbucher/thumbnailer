#!/bin/sh

# provide width, infer height automatically (conserve aspect ratio)
curl -X POST -F width=400 -F pdf=@demo.pdf localhost:1337/thumbnail > thumbnail-400x.png

# provide height, infer width automatically (conserve aspect ratio)
curl -X POST -F height=400 -F pdf=@demo.pdf localhost:1337/thumbnail > thumbnail-x400.png

# provide width and height (ignore aspect ratio)
curl -X POST -F width=400 -F height=566 -F pdf=@demo.pdf localhost:1337/thumbnail > thumbnail-400x566.png

# high quality settings
curl -X POST -F width=800 -F quality=95 -F density=300 -F pdf=@demo.pdf localhost:1337/thumbnail > thumbnail-hires.png

# JPEG
curl -X POST -F width=800 -F format=jpg -F pdf=@demo.pdf localhost:1337/thumbnail > thumbnail.jpg
