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

uniform sampler2D u_Sampler;
uniform vec4 u_Window;
uniform vec2 u_Dimension;
uniform vec2 u_Size;

uniform vec2 u_Offset;
uniform vec2 u_Tiling;


float map(float value, float inMin, float inMax, float outMin, float outMax) {
  return outMin + (outMax - outMin) * (value - inMin) / (inMax - inMin);
}

float window(float pointUV, float leftWindow, float rightWindow, float textureWidth, float sourceWidth, float tiling) {
    float scale = 2.0;

    float textureWidthPX        = textureWidth / scale;
    float sourceWidthPX         = sourceWidth / tiling;

    float leftWindowPX          = leftWindow;
    float leftWindowSourcePX    = leftWindow;
    float rightWindowPX         = textureWidthPX - rightWindow;
    float rightWindowSourcePX   = sourceWidthPX - rightWindow;

    float pointPX = pointUV * textureWidthPX;

    // float sourceUV = map(pointPX, 0.0, leftWindowPX, 0.0, leftWindowPX) / sourceWidth;
    float sourceUV = pointPX / sourceWidth;
    if (pointPX > leftWindowPX)
        sourceUV = map(pointPX, leftWindowPX, rightWindowPX, leftWindowSourcePX, rightWindowSourcePX) / sourceWidth;
    if (pointPX > rightWindowPX)
        sourceUV = map(pointPX, rightWindowPX, textureWidthPX, rightWindowSourcePX, sourceWidthPX) / sourceWidth;

    return sourceUV;
}

void main(void) {
   
    float x = window(v_UV.x, u_Window.y, u_Window.w, u_Size.x, u_Dimension.x, u_Tiling.x) + u_Offset.x / u_Tiling.x;
    float y = window(v_UV.y, u_Window.x, u_Window.z, u_Size.y, u_Dimension.y, u_Tiling.y) + u_Offset.y / u_Tiling.y;
    gl_FragColor = texture2D(u_Sampler, vec2(x, y));

    //sample.x = kp;
    //gl_FragColor = texture2D(u_Sampler, sample);
    

}