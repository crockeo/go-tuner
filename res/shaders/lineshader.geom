#version 330

// Width in pixels of the desired lines.
#define THICKNESS 5

// The size of the window relative to the viewpoint.
#define SCALE vec2(640, 480)

// Declaring the input and output layouts.
layout(lines) in;
layout(triangle_strip, max_vertices=4) out;

// Converting a vector to a magnitude / angle pair.
vec2 toMagAnglePair(vec2 vec) {
    vec2 pair;

    pair.x = sqrt(pow(vec.x, 2) + pow(vec.y, 2));
    pair.y = atan(vec.y / vec.x);

    return pair;
}

// Converting a magnitude / angle pair to a vector.
vec2 toVector(vec2 pair) {
    vec2 vec;

    vec.x = pair.x * cos(pair.y);
    vec.y = pair.x * sin(pair.y);

    return vec;
}

void main() {
    // Getting the source and destination values for the line.
    vec2 src = gl_in[0].gl_Position.xy,
         dst = gl_in[1].gl_Position.xy;

    // Determining the offset (negative or positive) for the src and dst points.
    vec2 calc = toMagAnglePair(abs(src - dst));

    calc.x = THICKNESS / 2;

    vec2 off1 = toVector(vec2(calc.x, calc.y + 90)),
         off2 = toVector(vec2(calc.x, calc.y - 90));

    off1 /= SCALE;
    off2 /= SCALE;

    // Creating the output vertices.
    vec2 verts[4];
    verts[0] = src + off1;
    verts[1] = src + off2;
    verts[2] = dst + off1;
    verts[3] = dst + off2;

    for (int i = 0; i < 4; i++) {
        gl_Position = vec4(verts[i], 0, 1);
        EmitVertex();
    }
}
