package visualize

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"runtime"
	"time"
)

// Initializing the source file.
func init() {
	runtime.LockOSThread()
}

// Testing stuff to do with opening an OpenGL context.
func Testing() error {
	if err := glfw.Init(); err != nil {
		return err
	}
	defer glfw.Terminate()

	// Setting up a bunch of windowing stuff with GLFW.
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Creating GLFW.
	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		return err
	}

	// Creating OpenGL.
	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		return err
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("Running on OpenGL: " + version)

	// Testing loading ShaderPrograms and Textures.
	_, err = LoadShaderProgram("res/shaders/texrenderer")
	if err != nil {
		fmt.Println("LSP: " + err.Error())
	}

	_, err = LoadTexture("res/textures/texture.jpg")
	if err != nil {
		fmt.Println("LT: " + err.Error())
	}

	for !window.ShouldClose() {
		glfw.PollEvents()
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}
