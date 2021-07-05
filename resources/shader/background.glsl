// References:
// - initial reference: https://www.geeks3d.com/hacklab/20180611/demo-simple-2d-grid-in-glsl/
// - shader: https://thebookofshaders.com/edit.php#10/ikeda-simple-grid.frag
// - book: https://thebookofshaders.com/09/

//vert:
precision mediump float;
attribute vec4 a_Vertex;	//position

uniform mat4 u_Projection;
uniform vec2 u_Resolution;

varying vec2 v_uv;
varying vec2 v_pos;

void main() {
    vec4 pos = vec4(a_Vertex.x * u_Resolution.x, a_Vertex.y * u_Resolution.y, a_Vertex.z, 1);
    gl_Position = u_Projection * pos;
    v_uv = (gl_Position.xy + 1.0) / 2.0; 
    v_pos = pos.xy;
}

//frag:
precision mediump float;

varying vec2 v_uv;
varying vec2 v_pos;
uniform vec2 u_Resolution;

float grid(vec2 st, float res) {
    vec2 grid = fract(st * res);
    return (step(res, grid.x) * step(res, grid.y));
}

void main(void) {
    gl_FragColor.a = 1.0;

    vec2 iResolution = vec2(1.0 / u_Resolution.x, 1.0 / u_Resolution.y);

    if (v_uv.x < 10.0*iResolution.x) {

    } else {
        gl_FragColor.r = v_uv.x;
        gl_FragColor.g = v_uv.y;
    }
}
