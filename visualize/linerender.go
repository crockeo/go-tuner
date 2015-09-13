package visualize

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

// Generating the VBO from a slice of Points.
func generateVBOData(points []Point) []float32 {
	vboData := make([]float32, len(points)*2)
	for i := 0; i < len(points)*2; i += 2 {
		vboData[i] = points[i/2].X
		vboData[i+1] = points[i/2].Y
	}

	return vboData
}

// Generating the EBO from a given length.
func generateEBOData(length int) []uint32 {
	eboData := make([]uint32, length)
	for i := 0; i < length; i++ {
		eboData[i] = uint32(i)
	}

	return eboData
}

// Type Point refers to an (x, y) pair representing location in 2d space.
type Point struct {
	X float32
	Y float32
}

// Type LineRender holds all of the data necessary to render a set of lines.
type LineRender struct {
	shaderProgram uint32  // The shader program to render the line.
	color         Color   // The color of this line.
	static        bool    // Whether or not this is a GL_STATIC_DRAW or a GL_DYNAMIC_DRAW.
	weight        float32 // The proposed weight of the LineRender.

	vao uint32 // the vertex array object.
	vbo uint32 // the vertex buffer object.
	ebo uint32 // the element buffer object.

	points int32 // The number of points.
}

// Creating a LineReader with an initial list of points.
func NewLineRender(shaderProgram ShaderProgram, color Color, static bool, weight float32, points []Point) *LineRender {
	lr := new(LineRender)

	// Setting the values non-changing values.
	lr.shaderProgram = uint32(shaderProgram)
	lr.static = static
	lr.weight = weight

	// Setting up some OpenGL information.
	gl.GenVertexArrays(1, &lr.vao)
	gl.GenBuffers(1, &lr.vbo)
	gl.GenBuffers(1, &lr.ebo)

	lr.UpdatePoints(points)

	return lr
}

// Creating an empty LineReader.
func NewLineReaderEmpty(shaderProgram ShaderProgram, color Color, static bool, weight float32) *LineRender {
	return NewLineRender(shaderProgram, color, static, weight, []Point{})
}

// Cleaning up the memory from a LineReader.
func (lr *LineRender) Destroy() {
	gl.DeleteVertexArrays(1, &lr.vao)
	gl.DeleteBuffers(1, &lr.vbo)
	gl.DeleteBuffers(1, &lr.ebo)
}

// Updating the list of Points that this LineRender should be rendering.
func (lr *LineRender) UpdatePoints(points []Point) {
	// Determining the buffer mode.
	var mode uint32
	if lr.static {
		mode = gl.STATIC_DRAW
	} else {
		mode = gl.DYNAMIC_DRAW
	}

	// Setting the number of points.
	lr.points = int32(len(points))

	// Generating the buffer data.
	vboData := generateVBOData(points)
	eboData := generateEBOData(len(points))

	// Filling the buffer data.
	gl.BindVertexArray(lr.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, lr.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vboData)*4, gl.Ptr(vboData), mode)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, lr.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(eboData)*4, gl.Ptr(eboData), mode)
}

// Rendering this LineRender.
func (lr *LineRender) Render() {
	// Binding the appropriate information.
	gl.BindVertexArray(lr.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, lr.vbo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, lr.ebo)
	gl.UseProgram(lr.shaderProgram)

	// Loading up vertex attributes.
	vertAttrib := uint32(gl.GetAttribLocation(lr.shaderProgram, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	// Fragment shader color information.
	gl.Uniform4f(
		gl.GetUniformLocation(lr.shaderProgram, gl.Str("in_color\x00")),
		lr.color.Red,
		lr.color.Green,
		lr.color.Blue,
		lr.color.Alpha)

	gl.BindFragDataLocation(lr.shaderProgram, 0, gl.Str("out_color\x00"))

	// Setting the line weight.
	gl.LineWidth(lr.weight)

	// Performing the render.
	gl.DrawElements(gl.LINES, lr.points, gl.UNSIGNED_INT, nil)
}
