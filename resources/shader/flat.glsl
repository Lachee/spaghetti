//vert:
attribute vec4 a_Vertex;	//position
attribute vec4 a_Color;     //color

uniform mat4 u_Projection;

varying vec4 v_Color;

void main() {
    gl_Position = u_Projection * a_Vertex;
    v_Color = a_Color;
}

//frag:
precision mediump float;

varying vec4 v_Color;

void main(void) {
    gl_FragColor = v_Color;
}
