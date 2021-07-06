//vert:
attribute vec4 a_Vertex;	//position
attribute vec2 a_UV;        //UV

uniform mat4 u_Projection;
uniform vec4 u_Window;
uniform vec2 u_Dimension;
uniform vec2 u_Size;

varying vec2 v_UV;

void main() { 
	gl_Position =  u_Projection * a_Vertex;
    v_UV = a_UV;
}

//frag:
precision mediump float;
varying highp vec2 v_UV;
uniform vec4 u_Window;
uniform vec2 u_Dimension;
uniform vec2 u_Size;
uniform sampler2D u_Sampler;


float map(float value, float inMin, float inMax, float outMin, float outMax) {
  return outMin + (outMax - outMin) * (value - inMin) / (inMax - inMin);
}

void main(void) {
    float scale = 10.0;
    vec2 sample = v_UV;

    float top = u_Window.x;
    float left = u_Window.y;
    float bottom = u_Window.z;
    float right = u_Window.w;

    float width = u_Dimension.x;
    float height = u_Dimension.y;

    float w = u_Dimension.x;
    float W = u_Size.x / scale;

    float a = left;
    float b = w - right;

    float A = left;
    float B = W - right;

    float Ap = (A / W);
    float Bp = (B / W);
    float Kp = v_UV.x;
    float K = Kp * W;
    
    float kp = sample.x;
    //if (Kp < Ap) {
        kp = map(K, 0.0, A, 0.0, a) / w;
        //gl_FragColor = vec4(1.0, 0.0, 0.0, 1.0);
    //}
    
    if (Kp > Ap && Kp < Bp) {
        kp = map(K, A, B, a, b) / w;
        //gl_FragColor = vec4(0.0, 1.0, 0.0, 1.0);
    }
    if (Kp > Bp) {
        kp = map(K, B, W, b, w) / w;
        //gl_FragColor = vec4(0.0, 0.0, 1.0, 1.0);
    }
    
    gl_FragColor = texture2D(u_Sampler, vec2(kp, sample.y));

    //sample.x = kp;
    //gl_FragColor = texture2D(u_Sampler, sample);
    

}