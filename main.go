package main

import (
    "fmt"
    "runtime"

    "github.com/go-gl/gl/v4.1-core/gl"
    "github.com/go-gl/glfw/v3.2/glfw"
    "github.com/go-gl/mathgl/mgl32"
)

var points = []float32{
    -0.8, 0, -0.8,
    -0.8, 0,  0.8,
     0.8, 0,  0.8,
     0.8, 0, -0.6,
    -0.6, 0, -0.6,
    -0.6, 0,  0.6,
     0.6, 0,  0.6,
     0.6, 0, -0.4,
    -0.4, 0, -0.4,
    -0.4, 0,  0.4,
     0.4, 0,  0.4,
     0.4, 0, -0.2,
    -0.2, 0, -0.2,
    -0.2, 0,  0.2,
     0.2, 0,  0.2,
     0.2, 0,  0.0,
     0.0, 0,  0.0,
}

// for ability to define points order (and/or reuse them)
var vertices = []uint32{
    // There is a functionality to reuse vertices in OpenGL with offset,
    // I just don't remember how to do it... In this case it should be eq 1.
    // TODO: get rid of this repeating mess
    0,1, 1,2, 2,3, 3,4, 4,5, 5,6, 6,7, 7,8,
    8,9, 9,10, 10,11, 11,12, 12,13, 13,14, 14,15, 15,16,
}

func main() {
    const width = 200
    const height = 200

    runtime.LockOSThread()

    window := initGlfw(width, height)
    defer glfw.Terminate()

    program := initOpenGL(vertexShader, fragmentShader);
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

    vao := makeVao(points, vertices)

    vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vp\x00")))
    gl.EnableVertexAttribArray(vertAttrib)
    gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

//---------------

    angle := 0.0
    previousTime := glfw.GetTime()

    for !window.ShouldClose() {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        time := glfw.GetTime()
        elapsed := time - previousTime
        previousTime = time

        angle += elapsed
        model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
        gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

        draw(vao, int32(len(vertices)), window, program)
    }
}

func draw(vao uint32, num int32, window *glfw.Window, program uint32) {
    gl.UseProgram(program)

    gl.BindVertexArray(vao)
    gl.DrawElements(gl.LINES, num, gl.UNSIGNED_INT, gl.PtrOffset(0))

    glfw.PollEvents()
    window.SwapBuffers()
}

func makeVao(points []float32, vertices []uint32) uint32 {
    var vao uint32
    gl.GenVertexArrays(1, &vao)
    gl.BindVertexArray(vao)

    var vbo uint32
    gl.GenBuffers(1, &vbo)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(
        gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

    var ibo uint32
    gl.GenBuffers(1, &ibo)
    gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
    gl.BufferData(
        gl.ELEMENT_ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

    return vao
}

func initGlfw(windowWidth int, windowHeight int) *glfw.Window {

    if err := glfw.Init(); err != nil {
        panic("failed to initialize glfw")
    }

    glfw.WindowHint(glfw.Resizable, glfw.False)
    glfw.WindowHint(glfw.ContextVersionMajor, 3)
    glfw.WindowHint(glfw.ContextVersionMinor, 3)
    glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
    glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

    window, err := glfw.CreateWindow(
        windowWidth, windowHeight, "Spiral", nil, nil)
    if err != nil {
        panic(err)
    }
    window.MakeContextCurrent()

    return window
}

func initOpenGL(vertexShaderSource string, fragmentShaderSource string) uint32 {
    if err := gl.Init(); err != nil {
        panic(err)
    }

    vertexShader, err := compileShader(vertexShader, gl.VERTEX_SHADER)
    if err != nil {
        panic(err)
    }
    fragmentShader, err := compileShader(fragmentShader, gl.FRAGMENT_SHADER)
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
