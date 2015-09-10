package visualize

import (
	"errors"
	"fmt"
	"github.com/crockeo/go-tuner/config"
	"github.com/go-gl/gl/v3.3-core/gl"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"strings"
)

// Type synonyms for more self-documenting code.
type Shader uint32
type ShaderProgram uint32
type Texture uint32

// Attempting to load a single shader.
func LoadShader(path string, shaderType uint32) (Shader, error, bool) {
	// Loading in the file contents.
	file, err := os.Open(path)
	if err != nil {
		return 0, err, false
	}
	defer file.Close()

	contentBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return 0, err, false
	}
	content := gl.Str(string(contentBytes) + "\x00")

	// Creating and compiling the shader.
	shader := gl.CreateShader(shaderType)

	gl.ShaderSource(shader, 1, &content, nil)
	gl.CompileShader(shader)

	var compiled int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &compiled)
	if compiled == gl.FALSE {
		var length int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)

		log := strings.Repeat("\x00", int(length+1))
		gl.GetShaderInfoLog(shader, length, nil, gl.Str(log))

		return 0, errors.New("Shader failed to compile: " + err.Error()), true
	}

	return Shader(shader), nil, false
}

// Handling the return value from a shader load.
func handleShaderLoad(name string, err error, major bool, loaded *int) {
	if err == nil {
		*loaded += 1
	} else if major || config.DebugMode {
		fmt.Println("Failed to load " + name + " shader: " + err.Error())
	}
}

// Attempting to load a ShaderProgram from a given location on the disk. It
// looks for file starting at the path, and then with the respective suffixes of
// .vert, .frag, and .geom.
func LoadShaderProgram(path string) (ShaderProgram, error) {
	loaded := 0

	// Loading the vertex shader.
	vert, err, major := LoadShader(path+".vert", gl.VERTEX_SHADER)
	handleShaderLoad("vert", err, major, &loaded)

	// Loading the fragment shader.
	frag, err, major := LoadShader(path+".frag", gl.FRAGMENT_SHADER)
	handleShaderLoad("frag", err, major, &loaded)

	// Loading the geometry shader.
	geom, err, major := LoadShader(path+".geom", gl.GEOMETRY_SHADER)
	handleShaderLoad("geom", err, major, &loaded)

	// Making sure we've loaded any shaders.
	if loaded == 0 {
		return 0, errors.New("Cannot make a shader program without any shaders.")
	}

	// Creating a loopable structure of the shaders.
	shaders := []uint32{
		uint32(vert),
		uint32(frag),
		uint32(geom),
	}

	program := gl.CreateProgram()
	for _, shader := range shaders {
		gl.AttachShader(program, shader)
	}
	gl.LinkProgram(program)

	// Checking that it linked correctly.
	var linked int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &linked)
	if linked == gl.FALSE {
		var length int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)

		log := strings.Repeat("\x00", int(length+1))
		gl.GetProgramInfoLog(program, length, nil, gl.Str(log))

		return 0, errors.New("Shader program failed to link: " + log)
	}

	return ShaderProgram(program), nil
}

// Attempting to load a Texture from a given location on the disk.
func LoadTexture(path string) (Texture, error) {
	// Loading the image data.
	imgFile, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, errors.New("Unsupported stride.")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	// Generating and populating the texture.
	var texture uint32

	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return Texture(texture), nil
}
