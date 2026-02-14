#version 430

precision highp float;
precision highp sampler2D;

// Input vertex attributes (from vertex shader)
in vec4 fragPosition;
in vec2 fragTexCoord;
in vec4 fragColor;
in vec3 fragNormal;

#ifdef PORTAL_OBJ
in vec4 fragScrPos;
#endif
#ifdef PORTAL_SCN
in float fragClip;
#endif

// Input uniform values
uniform sampler2D texture0;
uniform sampler2D texture1;
uniform sampler2D texture2;
uniform vec4 colDiffuse;

// Output fragment color
out vec4 finalColor;

// Custom uniforms

void main() {
    #ifdef PORTAL_SCN
    if (fragClip < 0.0) {
        discard;
    }
    #endif

    // Texel color fetching from texture sampler
    #ifdef PORTAL_OBJ
    vec2 uv = (fragScrPos.xy / fragScrPos.w);
    uv = uv*0.5 + 0.5;
    vec4 texelColor = texture(texture0, uv)*colDiffuse*fragColor;
    #else
    vec4 texelColor = texture(texture0, fragTexCoord)*colDiffuse*fragColor;
    #endif

    // Calculate final fragment color
    finalColor = texelColor;
}
