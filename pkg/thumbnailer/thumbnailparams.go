package thumbnailer

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	defaultDimension = 0
	defaultDensity   = 150
	defaultQuality   = 90
	pngFormat        = "png"
	jpgFormat        = "jpg"
	defaultFormat    = pngFormat
)

var formats = []string{pngFormat, jpgFormat}

type ThumbnailParams struct {
	Format  string
	Width   uint
	Height  uint
	Density uint
	Quality uint8 // 0..100
}

func ParseParams(r *http.Request) (*ThumbnailParams, error) {
	var params ThumbnailParams

	if r.FormValue("width") != "" {
		width, err := parseUintParam(r, "width")
		if err != nil {
			return nil, err
		}
		params.Width = width
	} else {
		params.Width = defaultDimension
	}

	if r.FormValue("height") != "" {
		height, err := parseUintParam(r, "height")
		if err != nil {
			return nil, err
		}
		params.Height = height
	} else {
		params.Height = defaultDimension
	}

	if r.FormValue("density") != "" {
		density, err := parseUintParam(r, "density")
		if err != nil {
			return nil, err
		}
		params.Density = density
	} else {
		params.Density = defaultDensity
	}

	if r.FormValue("quality") != "" {
		quality, err := parseUintParamInRange(r, "quality", 0, 100)
		if err != nil {
			return nil, err
		}
		params.Quality = uint8(quality)
	} else {
		params.Quality = uint8(defaultQuality)
	}

	if format := strings.ToLower(r.FormValue("format")); format != "" {
		contained := false
		for _, f := range formats {
			if f == format {
				contained = true
				break
			}
		}
		if !contained {
			return nil, fmt.Errorf("illegal format '%s', use one of %q", format, formats)
		}
		params.Format = format
	} else {
		params.Format = defaultFormat
	}

	return &params, nil
}

func (tp ThumbnailParams) AsArgs(extraArgs ...string) []string {
	var params []string

	params = append(params, "-thumbnail") // optimized version of "-resize"
	widthStr, heightStr := "", ""
	if tp.Width > 0 {
		widthStr = strconv.Itoa(int(tp.Width))
	}
	if tp.Height > 0 {
		heightStr = strconv.Itoa(int(tp.Height))
	}
	// worst case: "x" (defaults to 1x1: pointless but legal)
	resizeParam := fmt.Sprintf("%sx%s", widthStr, heightStr)
	if widthStr != "" && heightStr != "" {
		// "!" ignores aspect ratio
		resizeParam = fmt.Sprintf("%sx%s!", widthStr, heightStr)
	}
	params = append(params, resizeParam)

	params = append(params, "-quality")
	params = append(params, strconv.Itoa(int(tp.Quality)))

	params = append(params, "-density")
	params = append(params, strconv.Itoa(int(tp.Density)))

	for _, arg := range extraArgs {
		params = append(params, arg)
	}

	return params
}

func parseUintParam(r *http.Request, name string) (uint, error) {
	var intVal uint
	stringParam := r.FormValue(name)
	if stringParam != "" {
		i, err := strconv.Atoi(stringParam)
		if err != nil {
			return 0, fmt.Errorf("'%s' param not an integer: %v", name, err)
		}
		if i < 0 {
			return 0, fmt.Errorf("'%s' param must not be negative", name)
		}
		intVal = uint(i)
	}
	return intVal, nil
}

func parseUintParamInRange(r *http.Request, name string, min, max uint) (uint, error) {
	val, err := parseUintParam(r, name)
	if err != nil {
		return 0, err
	}
	if val < min || val > max {
		return 0, fmt.Errorf("value %d out of range [%d, %d]", val, min, max)
	}
	return val, nil
}
