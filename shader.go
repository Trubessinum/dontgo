package main

import (
    "fmt"
    "strings"

    "github.com/go-gl/gl/v4.1-core/gl"
)

var vertexShader string = `
    #version 330

    varying vec3 pos;

    uniform mat4 projection;
    uniform mat4 camera;
    uniform mat4 model;

    in vec3 vp;

    void main() {
        pos = vp;
        gl_Position = projection * camera * model * vec4(vp, 1);
    }
`

var fragmentShader string = `
    #version 330

    varying vec3 pos;

    out vec4 frag_colour;

    void main() {
        frag_colour = vec4(pos, 1);
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
