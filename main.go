// main.go
package main

import (
    "syscall/js"
)

func draw(this js.Value, p []js.Value) interface{} {
    document := js.Global().Get("document")
    
    // 创建 canvas 元素
    canvas := document.Call("createElement", "canvas")
    canvas.Set("id", "myCanvas")
    canvas.Set("width", 800)
    canvas.Set("height", 600)
    
    // 将 canvas 添加到 body
    document.Get("body").Call("appendChild", canvas)
    
    // 获取 2D 绘图上下文
    ctx := canvas.Call("getContext", "2d")
    
    // 绘制一个矩形
    ctx.Set("fillStyle", "lightblue")
    ctx.Call("fillRect", 50, 50, 200, 100)
    
    // 绘制一些文本
    ctx.Set("fillStyle", "black")
    ctx.Set("font", "20px Arial")
    ctx.Call("fillText", "Hello, WebAssembly with Go!", 60, 100)
    
    return nil
}

func main() {
    js.Global().Set("draw", js.FuncOf(draw))
    select {}
}
