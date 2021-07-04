//vert:
attribute vec4 a_Vertex;	//position

uniform mat4 u_Projection;
uniform vec2 u_Resolution;

varying vec2 v_uv;

void main() {
    vec4 pos = vec4(a_Vertex.x * u_Resolution.x, a_Vertex.y * u_Resolution.y, a_Vertex.z, 1);
    gl_Position = u_Projection * pos;
    v_uv = pos.xy;
}

//frag:
precision mediump float;

varying vec2 v_uv;

float grid(vec2 st, float res) {
    vec2 grid = fract(st * res);
    return (step(res, grid.x) * step(res, grid.y));
}

void main(void) {
    float scale = 10.0;
    float resolution = 0.1;
    vec3 colour = vec3(0.25, 0.25, 0.25);

    // Scale the colours a bit for the other quads
    if (v_uv.y < 0.0)
        colour.g = 0.0;
    if (v_uv.x < 0.0)
        colour.r = 0.0;


    vec2 grid_uv = v_uv.xy * scale; // scale
    float x = grid(grid_uv, resolution); // resolution
    gl_FragColor.rgb = colour * x;
    gl_FragColor.a = 1.0;
}
