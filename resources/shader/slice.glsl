//vert:
attribute vec4 a_Vertex;	//position
attribute vec2 a_UV;        //UV

uniform mat4 u_Projection;

varying vec2 v_UV;

void main() { 
	gl_Position =  u_Projection * a_Vertex;
    v_UV = a_UV;
}

//frag:
precision mediump float;
varying highp vec2 v_UV;
uniform sampler2D u_Sampler;

void main(void) {
    gl_FragColor = vec4(v_UV.x, v_UV.y, 1.0, 1.0);
    gl_FragColor = texture2D(u_Sampler, v_UV);
}