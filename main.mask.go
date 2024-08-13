package main

import (
	"syscall/js"
)

func main() {

	jsProcess := js.FuncOf(Process)

	js.Global().Set("process", jsProcess)

	defer jsProcess.Release()

	select {}
}

type process struct {
	data js.Value
	img  js.Value
}

func Process(this js.Value, args []js.Value) interface{} {
	imgSrc := args[0]
	storkeWidth := args[1]
	jsCallback := args[2]

	if storkeWidth.Int() == 0 {
		return js.ValueOf(process{data: js.ValueOf(nil), img: imgSrc})
	}

	document := js.Global().Get("document")

	canvas := document.Call("createElement", "canvas")
	ctx := canvas.Call("getContext", "2d")
	fillColor := js.ValueOf("rgba(116, 90, 241, .5)")

	img := document.Call("createElement", "img")
	img.Set("crossorigin", "anonymous")
	img.Set("src", imgSrc)

	var onloadCallback js.Func

	onloadCallback = js.FuncOf(func(this js.Value, _ []js.Value) interface{} {
		canvas.Set("width", img.Get("width"))
		canvas.Set("height", img.Get("height"))

		ctx.Call("drawImage", img, 0, 0)

		imageData := ctx.Call("getImageData", 0, 0, img.Get("width").Int(), img.Get("height").Int())
		data := imageData.Get("data")

		boundaries := []int{}
		var i int = 0
		len := data.Get("length").Int()

		for i < len {
			if data.Index(i+3).Int() != 0 {
				boundaries = append(boundaries, i)
			}
			i += 4
		}

		minX := imageData.Get("width").Int()
		maxX := 0
		maxY := 0
		var w = img.Get("width").Int()


		for _, v := range boundaries {
			x := (v / 4) % w
			y := (v / 4 / w)
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}
			if y > maxY {
				maxY = y
			}

			FillCriclePath(ctx, js.ValueOf(x), js.ValueOf(y), storkeWidth, fillColor)
		}

		newImageData := ctx.Call("getImageData", 0, 0, canvas.Get("width").Int(), canvas.Get("height").Int())

		res := process{data: newImageData, img: canvas.Call("toDataURL")}

		resMap := map[string]interface{}{
			"data": res.data,
			"img":  res.img,
		}

		jsCallback.Invoke(js.ValueOf(resMap))
		
		onloadCallback.Release()
		return js.ValueOf(nil)
	})

	img.Set("onload", onloadCallback)

	return js.ValueOf(nil)
}

func FillCriclePath(ctx js.Value, x js.Value, y js.Value, N js.Value, color js.Value) interface{} {
	path := GetCriclePath(x, y, N)
	ctx.Set("fillStyle", color)
	ctx.Call("fill", path)
	return nil
}

func GetCriclePath(x js.Value, y js.Value, radius js.Value) interface{} {
	path := js.Global().Get("Path2D").New()
	path.Call("arc", x, y, radius, js.ValueOf(0), js.Global().Get("Math").Get("PI").Float()*2)
	return path
}
