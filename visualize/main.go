package visualize

import (
	"fmt"
	"github.com/crockeo/go-tuner/config"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"math"
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

	if config.DebugMode {
		version := gl.GoStr(gl.GetString(gl.VERSION))
		fmt.Println("Running on OpenGL: " + version)
	}

	// Testing rendering of objects.
	shaderProgram, err := LoadShaderProgram("res/shaders/texrenderer")
	if err != nil {
		fmt.Println("LSP: " + err.Error())
	}
	defer DestroyShaderProgram(shaderProgram)

	texture, err := LoadTexture("res/textures/texture.jpg")
	if err != nil {
		fmt.Println("LT: " + err.Error())
	}
	defer DestroyTexture(texture)

	renderObject := CreateRenderObject(shaderProgram, texture, []float32{
		-1.0, -1.0, 0.0, 0.0,
		1.0, -1.0, 1.0, 0.0,
		1.0, 1.0, 1.0, 1.0,
		-1.0, 1.0, 0.0, 1.0,
	})
	defer renderObject.Destroy()

	// Testing a LineRender.
	lineShader, err := LoadShaderProgram("res/shaders/lineshader")
	if err != nil {
		fmt.Println("LSP: " + err.Error())
	}
	defer DestroyShaderProgram(lineShader)

	lineRenders := []*LineRender{
		NewLineRender(lineShader, Color{0.0, 0.0, 0.0, 1.0}, true, 3.0, []Point{
			Point{-1.0, 0},
			Point{1.0, 0},
		}),

		NewLineRender(lineShader, Color{1.0, 0.0, 1.0, 1.0}, true, 8.0, []Point{
			Point{-0.5, -0.5},
			Point{0.5, 0.5},
			Point{1, 0},
		}),
	}
	defer func() {
		for _, lr := range lineRenders {
			lr.Destroy()
		}
	}()

	// Creating a 2nd LineRender from DefaultGenerateSinePoints.
	var phase float32 = 0
	lineRender2 := NewLineRender(
		lineShader,
		Color{1.0, 0.0, 0.0, 1.0},
		false,
		4.0,
		DefaultGenerateSinePoints(440, phase))
	defer lineRender2.Destroy()

	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Rendering the outputs.
		renderObject.Render()
		for _, lr := range lineRenders {
			lr.Render()
		}
		lineRender2.Render()

		// Updating render for lineRender2.
		phase += 2 * math.Pi * (440.0 / 44100.0)
		if phase > 2*math.Pi {
			phase -= 2 * math.Pi
		}
		lineRender2.UpdatePoints(DefaultGenerateSinePoints(440, phase))

		if config.DebugMode {
			// Reporting OpenGL errors.
			glErr := gl.GetError()
			if glErr != gl.NO_ERROR {
				fmt.Printf("OpenGL error: %d\n", glErr)
			}
		}

		window.SwapBuffers()
		glfw.PollEvents()
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}
