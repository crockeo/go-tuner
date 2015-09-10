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
	ebo           uint32 // The ordering of the vertices.
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
	gl.GenBuffers(1, &renderObject.ebo)

	// Filling the RenderObject with information.
	gl.BindVertexArray(renderObject.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, renderObject.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, renderObject.ebo)
	vertOrder := []uint32{
		0, 1, 2,
		2, 3, 0,
	}
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(vertOrder)*4, gl.Ptr(vertOrder), gl.STATIC_DRAW)

	// Loading up vertex attributes.
	vertAttrib := uint32(gl.GetAttribLocation(renderObject.shaderProgram, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	// Loading up texture attributes.
	texAttrib := uint32(gl.GetAttribLocation(renderObject.shaderProgram, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texAttrib)
	gl.VertexAttribPointer(texAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))

	return renderObject
}

// Destroying the resources of a RenderObject.
func (renderObject *RenderObject) Destroy() {
	gl.DeleteVertexArrays(1, &renderObject.vao)
	gl.DeleteBuffers(1, &renderObject.vbo)
	gl.DeleteBuffers(1, &renderObject.ebo)
}

func (renderObject *RenderObject) Render() {
	gl.BindVertexArray(renderObject.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, renderObject.vbo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, renderObject.ebo)
	gl.UseProgram(renderObject.shaderProgram)

	// Binding the texture.
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, renderObject.texture)
	gl.BindFragDataLocation(renderObject.shaderProgram, 0, gl.Str("outputColor\x00"))

	// Drawing the object.
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}
