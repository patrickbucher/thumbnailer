# Thumbnailer

Some quick and dirty thumbnailer server based on ImageMagick and hacked
in Go.

## Setup (ImageMagick)

Make sure ImageMagick's policy file (`/etc/ImageMagick-7/policy.xml` on Arch
Linux) allows reading PDFs by changing

```xml
<policy domain="coder" rights="none" pattern="{PS,PS2,PS3,EPS,PDF,XPS}" />
```

to


```xml
<policy domain="coder" rights="none" pattern="{PS,PS2,PS3,EPS,XPS}" />
<policy domain="coder" rights="read" pattern="{PDF}" />
```

by excluding the PDF from the restrictive rule and creating a new rule for
reading PDFs instead.

## Execute

Turn the first page of `demo.pdf` into `thumbnail.png`:

```bash
$ go run thumbnailer.go
$ curl -X POST -F pdf=@demo.pdf localhost:1337/thumbnail > thumbnail.png
```

Further specify `width`, `height`, and `density`:

```bash
$ curl -X POST -F pdf=@demo.pdf -F width=400 -F height=566 -F density=300 \
    localhost:1337/thumbnail > thumbnail.png
```

If both `width` and `height` are given, the aspect ratio will be ignored.
Otherwise, the missing parameter is calculated according to the input's aspect
ratio.
