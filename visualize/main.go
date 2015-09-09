package visualize

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"time"
)

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

	for !window.ShouldClose() {
		glfw.PollEvents()
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}
