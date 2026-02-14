#version 430

// Input vertex attributes
in vec3 vertexPosition;
in vec2 vertexTexCoord;
in vec3 vertexNormal;
in vec4 vertexColor;
in vec4 vertexTangent;
in vec4 vertexTexCoord2;

#ifdef GRASS_OBJ
layout(std430, binding = 1) buffer InstanceBuf {
    struct GrassInstance {
        vec4 pos_height;
        vec4 n_rot;
    } instances[];
};
#endif

// Input uniform values
uniform mat4 mvp;
uniform mat4 matModel;
uniform mat4 matView;
uniform mat4 matProjection;
uniform mat4 matNormal;

// Output vertex attributes (to fragment shader)
out vec4 fragPosition;
out vec2 fragTexCoord;
out vec4 fragColor;
out vec3 fragNormal;

#ifdef PORTAL_OBJ
out vec4 fragScrPos;
#endif
#ifdef PORTAL_SCN
out float fragClip;
#endif

// Custom Uniforms
uniform vec3 portalPos = vec3(0,0,0);
uniform vec3 portalNorm = vec3(0,0,0);
uniform mat4 portalMat = mat4(1.0);

mat2 rot2(float a){
    float s = sin(a), c = cos(a);
    return mat2(c,-s,s,c);
}

void main() {
    // grass based positioning
    #ifdef GRASS_OBJ
    uint id = uint(gl_InstanceID);
    vec3 basePos = instances[id].pos_height.xyz;
    float h      = instances[id].pos_height.w;
    float rot    = instances[id].n_rot.w;

    // rotate around Y in local grass space
    local.xz = rot2(rot) * vertexPosition.xz;
    // scale height
    local.y *= h;
    vertexPosition = basePos + vertexPosition
    #endif

    // Calculate final vertex position
    gl_Position = mvp*portalMat*portalMat*vec4(vertexPosition, 1.0);

    // Send vertex attributes to fragment shader
    fragPosition = matModel*vec4(vertexPosition, 1.0);
    fragTexCoord = vertexTexCoord;
    fragColor = vertexColor;
    fragNormal = normalize(vec3(matNormal*vec4(vertexNormal, 1.0)));

    #ifdef PORTAL_OBJ
    fragScrPos = gl_Position;
    #endif
    #ifdef PORTAL_SCN
    fragClip = dot((vec3(fragPosition) - portalPos), portalNorm);
    #endif
}
