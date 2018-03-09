package main

import (
    "fmt"
    "log"
    "runtime"

    "github.com/go-gl/gl/v4.1-core/gl"
    "github.com/go-gl/glfw/v3.2/glfw"
    "github.com/go-gl/mathgl/mgl32"
)

var pyramidVertices = []float32{
    // front
    0, 0.7, 0,
    -0.7, -0.7, -0.7,
    0.7, -0.7, -0.7,

    // back
    0, 0.7, 0,
    0.7, -0.7, 0.7,
    -0.7, -0.7, 0.7,

    // left
    0, 0.7, 0,
    -0.7, -0.7, 0.7,
    -0.7, -0.7, -0.7,

    // right
    0, 0.7, 0,
    0.7, -0.7, -0.7,
    0.7, -0.7, 0.7,
}

func main() {
    const width = 200
    const height = 200

    runtime.LockOSThread()

    window := initGlfw(width, height)
    defer glfw.Terminate()

    program := initOpenGL();
    gl.UseProgram(program)

//---------------

    projection := mgl32.Perspective(
        mgl32.DegToRad(45.0), float32(width)/height, 0.1, 10.0)
    projectionUniform := gl.GetUniformLocation(
        program, gl.Str("projection\x00"))
    gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

    camera := mgl32.LookAtV(
        mgl32.Vec3{2, 1.5, 2}, mgl32.Vec3{0, -0.25, 0}, mgl32.Vec3{0, 1, 0})
    cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
    gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

    model := mgl32.Ident4()
    modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
    gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

//---------------

    vao := makeVao(pyramidVertices)
    num := int32(len(pyramidVertices) / 3)

    angle := 0.0
    previousTime := glfw.GetTime()

    for !window.ShouldClose() {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        // Update
        time := glfw.GetTime()
        elapsed := time - previousTime
        previousTime = time

        angle += elapsed
        model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
        gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

        draw(vao, num, window, program)
    }
}

func draw(vao uint32, num int32, window *glfw.Window, program uint32) {
    gl.UseProgram(program)

    gl.BindVertexArray(vao)
    gl.DrawArrays(gl.TRIANGLES, 0, num)

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
        windowWidth, windowHeight, "Pyramid", nil, nil)
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

    vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
    if err != nil {
        panic(err)
    }
    fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
    if err != nil {
        panic(err)
    }

    version := gl.GoStr(gl.GetString(gl.VERSION))
    fmt.Println("OpenGL version", version)

    prog := gl.CreateProgram()
    gl.AttachShader(prog, vertexShader)
    gl.AttachShader(prog, fragmentShader)
    gl.LinkProgram(prog)

    gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

    gl.Enable(gl.DEPTH_TEST)
    gl.DepthFunc(gl.LESS)
    gl.ClearColor(0.0, 0.0, 0.4, 1.0)

    return prog
}

func makeVao(points []float32) uint32 {
    var vbo uint32
    gl.GenBuffers(1, &vbo)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

    var vao uint32
    gl.GenVertexArrays(1, &vao)
    gl.BindVertexArray(vao)
    gl.EnableVertexAttribArray(0)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

    return vao
}
