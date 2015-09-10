package visualize

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

// Type RenderObject is used to concisely represent information necessary to
// perform a 2 dimensional textured render.
type RenderObject struct {
	shaderProgram uint32 // The shader program.
	texture       uint32 // The texture.
	vao           uint32 // The Vertex Array Object.
	vbo           uint32 // The Vertex Buffer Object
}

// Creating a RenderObject with a given shaderProgram, texture, and set of
// vertices.
func CreateRenderObject(shaderProgram ShaderProgram, texture Texture, vertices []float32) *RenderObject {
	renderObject := new(RenderObject)

	// Creating the basic information.
	renderObject.shaderProgram = uint32(shaderProgram)
	renderObject.texture = uint32(texture)

	gl.GenVertexArrays(1, &renderObject.vao)
	gl.GenBuffers(1, &renderObject.vbo)

	// Filling the RenderObject with information.
	gl.BindVertexArray(renderObject.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, renderObject.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Loading up vertex attributes.
	vertAttrib := uint32(gl.GetAttribLocation(renderObject.shaderProgram, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 2, gl.FLOAT, false, 2*4, nil)

	// Loading up texture attributes.
	texAttrib := uint32(gl.GetAttribLocation(renderObject.shaderProgram, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texAttrib)
	gl.VertexAttribPointer(texAttrib, 2, gl.FLOAT, false, 2*4, nil)

	return renderObject
}

// Destroying the resources of a RenderObject.
func (renderObject *RenderObject) Destroy() {
	gl.DeleteVertexArrays(1, &renderObject.vao)
	gl.DeleteBuffers(1, &renderObject.vbo)
}

func (renderObject *RenderObject) Render() {
	gl.BindVertexArray(renderObject.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, renderObject.vbo)
	gl.UseProgram(renderObject.shaderProgram)

	// Binding the texture.
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, renderObject.texture)
	gl.Uniform1i(gl.GetUniformLocation(renderObject.shaderProgram, gl.Str("in_tex\x00")), 0)

	// Drawing the object.
	gl.DrawArrays(gl.TRIANGLES, 0, 4*3)
}
