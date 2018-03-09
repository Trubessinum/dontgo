package main

import (
    "fmt"
    "strings"

    "github.com/go-gl/gl/v4.1-core/gl"
)

var vertexShaderSource string = `
    #version 410
    uniform mat4 model;

    in vec3 vp;

    void main() {
        gl_Position = model * vec4(vp, 1);
    }
`

var fragmentShaderSource string = `
    #version 410
    out vec4 frag_colour;
    void main() {
        frag_colour = vec4(0, 1, 0, 1);
    }
`

func compileShader(source string, shaderType uint32) (uint32, error) {
    source += "\x00"

    shader := gl.CreateShader(shaderType)

    csources, free := gl.Strs(source)
    gl.ShaderSource(shader, 1, csources, nil)
    free()
    gl.CompileShader(shader)

    var status int32
    gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

        return 0, fmt.Errorf("failed to compile %v: %v", source, log)
    }

    return shader, nil
}
