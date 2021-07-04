//vert:
attribute vec4 a_Vertex;	//position

uniform mat4 u_Projection;

void main() { 
	gl_Position =  u_Projection * a_Vertex;
}

//frag:
precision mediump float;

void main(void) {
    gl_FragColor = vec4(1, 0.3, 0.0, 1.0);
}