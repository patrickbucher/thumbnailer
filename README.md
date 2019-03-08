# Thumbnailer

Some quick and dirty thumbnailer based on ImageMagick.

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
