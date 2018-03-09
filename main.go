package main

import (
    "fmt"
    "log"
    "runtime"

    "github.com/go-gl/gl/v4.1-core/gl"
    "github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
    const width = 400
    const height = 300

    runtime.LockOSThread()

    window := initGlfw(width, height)
    defer glfw.Terminate()

    program := initOpenGL();

    for !window.ShouldClose() {
        draw(window, program)
    }
}

func draw(window *glfw.Window, program uint32) {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.UseProgram(program)

    glfw.PollEvents()
    window.SwapBuffers()
}

func initGlfw(windowWidth int, windowHeight int) *glfw.Window {

    if err := glfw.Init(); err != nil {
        log.Fatalln("failed to initialize glfw:", err)
    }

    glfw.WindowHint(glfw.Resizable, glfw.False)
    glfw.WindowHint(glfw.ContextVersionMajor, 4)
    glfw.WindowHint(glfw.ContextVersionMinor, 1)
    glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
    glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

    window, err := glfw.CreateWindow(
        windowWidth, windowHeight, "Window", nil, nil)
    if err != nil {
        panic(err)
    }
    window.MakeContextCurrent()

    return window
}

func initOpenGL() uint32 {
    if err := gl.Init(); err != nil {
        panic(err)
    }

    version := gl.GoStr(gl.GetString(gl.VERSION))
    fmt.Println("OpenGL version", version)

    prog := gl.CreateProgram()
    gl.LinkProgram(prog)

    gl.Enable(gl.DEPTH_TEST)
    gl.DepthFunc(gl.LESS)
    gl.ClearColor(0.0, 0.0, 0.4, 1.0)

    return prog
}
