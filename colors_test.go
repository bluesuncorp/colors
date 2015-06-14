package colors

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"testing"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called
//

func IsEqual(t *testing.T, val1, val2 interface{}) bool {
	v1 := reflect.ValueOf(val1)
	v2 := reflect.ValueOf(val2)

	if v1.Kind() == reflect.Ptr {
		v1 = v1.Elem()
	}

	if v2.Kind() == reflect.Ptr {
		v2 = v2.Elem()
	}

	if !v1.IsValid() && !v2.IsValid() {
		return true
	}

	v1Underlying := reflect.Zero(reflect.TypeOf(v1)).Interface()
	v2Underlying := reflect.Zero(reflect.TypeOf(v2)).Interface()

	if v1 == v1Underlying {
		if v2 == v2Underlying {
			goto CASE4
		} else {
			goto CASE3
		}
	} else {
		if v2 == v2Underlying {
			goto CASE2
		} else {
			goto CASE1
		}
	}

CASE1:
	return reflect.DeepEqual(v1.Interface(), v2.Interface())

CASE2:
	return reflect.DeepEqual(v1.Interface(), v2)
CASE3:
	return reflect.DeepEqual(v1, v2.Interface())
CASE4:
	return reflect.DeepEqual(v1, v2)
}

func Equal(t *testing.T, val1, val2 interface{}) {
	EqualSkip(t, 2, val1, val2)
}

func EqualSkip(t *testing.T, skip int, val1, val2 interface{}) {

	if !IsEqual(t, val1, val2) {

		_, file, line, _ := runtime.Caller(skip)
		fmt.Printf("%s:%d %v does not equal %v\n", path.Base(file), line, val1, val2)
		t.FailNow()
	}
}

func NotEqual(t *testing.T, val1, val2 interface{}) {
	NotEqualSkip(t, 2, val1, val2)
}

func NotEqualSkip(t *testing.T, skip int, val1, val2 interface{}) {

	if IsEqual(t, val1, val2) {
		_, file, line, _ := runtime.Caller(skip)
		fmt.Printf("%s:%d %v should not be equal %v\n", path.Base(file), line, val1, val2)
		t.FailNow()
	}
}

func PanicMatches(t *testing.T, fn func(), matches string) {
	PanicMatchesSkip(t, 2, fn, matches)
}

func PanicMatchesSkip(t *testing.T, skip int, fn func(), matches string) {

	_, file, line, _ := runtime.Caller(skip)

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("%s", r)

			if err != matches {
				fmt.Printf("%s:%d Panic...  expected [%s] received [%s]", path.Base(file), line, matches, err)
				t.FailNow()
			}
		}
	}()

	fn()
}

func TestColorConversionFromHEX(t *testing.T) {

	hex, _ := ParseHEX("#5f55f5")

	Equal(t, hex.ToHEX().String(), "#5f55f5")
	Equal(t, hex.ToRGB().String(), "rgb(95,85,245)")
	Equal(t, hex.ToRGBA().String(), "rgba(95,85,245,1)")

	hex, _ = ParseHEX("#5f5")
	Equal(t, hex.ToRGB().String(), "rgb(85,255,85)")

	hex, _ = ParseHEX("Bad Hex color!")
	Equal(t, hex, nil)
}

func TestColorConversionFromRGB(t *testing.T) {

	rgb, _ := ParseRGB("rgb(95%,85%,50%)")

	Equal(t, rgb.ToRGB().String(), "rgb(242,217,128)")
	Equal(t, rgb.ToRGBA().String(), "rgba(242,217,128,1)")
	Equal(t, rgb.ToHEX().String(), "#f2d980")

	rgb, _ = ParseRGB("rgb(95,85,245)")
	Equal(t, rgb.ToRGB().String(), "rgb(95,85,245)")
	Equal(t, rgb.ToRGBA().String(), "rgba(95,85,245,1)")
	Equal(t, rgb.ToHEX().String(), "#5f55f5")

	rgb, _ = RGB(95, 85, 245)
	Equal(t, rgb.ToRGB().String(), "rgb(95,85,245)")
	Equal(t, rgb.ToRGBA().String(), "rgba(95,85,245,1)")
	Equal(t, rgb.ToHEX().String(), "#5f55f5")

	rgb, _ = ParseRGB("BAD RGB COLOR")
	Equal(t, rgb, nil)

	rgb, _ = ParseRGB("rgb(95%,85%,245)")
	Equal(t, rgb, nil)
}

func TestColorConversionFromRGBA(t *testing.T) {

	rgba, _ := ParseRGBA("rgba(95%,85%,50%,1)")

	Equal(t, rgba.ToRGB().String(), "rgb(242,217,128)")
	Equal(t, rgba.ToRGBA().String(), "rgba(242,217,128,1)")
	Equal(t, rgba.ToHEX().String(), "#f2d980")

	rgba, _ = ParseRGBA("rgba(95,85,245,1)")
	Equal(t, rgba.ToRGB().String(), "rgb(95,85,245)")
	Equal(t, rgba.ToRGBA().String(), "rgba(95,85,245,1)")
	Equal(t, rgba.ToHEX().String(), "#5f55f5")

	rgba, _ = RGBA(95, 85, 245, 1)
	Equal(t, rgba.ToRGB().String(), "rgb(95,85,245)")
	Equal(t, rgba.ToRGBA().String(), "rgba(95,85,245,1)")
	Equal(t, rgba.ToHEX().String(), "#5f55f5")

	rgba, _ = RGBA(95, 85, 245, 6)
	Equal(t, rgba, nil)

	rgba, _ = RGBA(95, 85, 245, -1)
	Equal(t, rgba, nil)

	rgba, _ = ParseRGBA("BAD RGBA COLOR")
	Equal(t, rgba, nil)

	rgba, _ = ParseRGBA("rgba(95%,85%,245,1)")
	Equal(t, rgba, nil)
}

func TestParseColor(t *testing.T) {

	color, _ := Parse("#FFF")
	NotEqual(t, color, nil)
	Equal(t, reflect.TypeOf(color), reflect.TypeOf(&HEXColor{}))

	color, _ = Parse("rgb(95,85,245)")
	NotEqual(t, color, nil)
	Equal(t, reflect.TypeOf(color), reflect.TypeOf(&RGBColor{}))

	color, _ = Parse("rgba(95,85,245,1)")
	NotEqual(t, color, nil)
	Equal(t, reflect.TypeOf(color), reflect.TypeOf(&RGBAColor{}))
}

func TestInterfaceTypes(t *testing.T) {

	fn := func(c Color) string {

		if c == nil {
			return ""
		}

		return c.String()
	}

	hex, _ := ParseHEX("#FFF")
	rgb, _ := ParseRGB("rgb(95,85,245)")
	rgba, _ := ParseRGBA("rgba(95,85,245,1)")

	fn(hex)
	fn(rgb)
	fn(rgba)
}

func BenchmarkSpeed(b *testing.B) {

	for n := 0; n < b.N; n++ {
		h, _ := ParseHEX("#FFFFFF")
		h.ToRGBA()
	}
}
