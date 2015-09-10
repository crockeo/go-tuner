#version 330

uniform sampler2D in_tex;
in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
    outputColor = texture(in_tex, fragTexCoord);
}
