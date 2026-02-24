#version 430

precision highp float;
precision highp sampler2D;

layout (triangles) in;
layout (triangle_strip) out;
layout (max_vertices = 256) out;

in vec4 vertPosition[];
in vec2 vertTexCoord[];
in vec4 vertColor[];
in vec3 vertNormal[];

#ifdef PORTAL_OBJ
in vec4 vertScrPos[];
#endif
#ifdef PORTAL_SCN
in float vertClip[];
#endif

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

// Input uniform values
uniform mat4 mvp;
uniform mat4 matModel;
uniform mat4 matView;
uniform mat4 matProjection;

void main() {
	int i;
	int j;

	for (i = 0; i+2 < gl_in.length(); i+=3) {
		for (j = 0; j < 4; j++) {
			vec3 randomOffset = vec3(
				fract(sin(dot(vec3(j, i, 0), vec3(12.9898, 78.233, 54.53))) * 43758.5453),
				0.0,
				fract(sin(dot(vec3(j + 1, i + 1, 0), vec3(12.9898, 78.233, 54.53))) * 43758.5453)
			);
			vec4 grassPos = mix(mix(gl_in[i].gl_Position, gl_in[i+1].gl_Position, randomOffset.x), gl_in[i+2].gl_Position, randomOffset.z);

			gl_Position = grassPos + matProjection * (matView * matModel * vec4(0,0,0,0) + vec4(0.02, 0, 0, 0));
			fragColor = vec4(0,0.3,0,1) * vertColor[i];
			fragTexCoord = vertTexCoord[i];
			fragNormal = vertNormal[i];
			EmitVertex();
			gl_Position = grassPos + mvp * vec4(0, 0.2, 0, 0) +  + matProjection * (matView * matModel * vec4(0,0,0,0) + vec4(0, 0.1, 0, 0));
			fragColor = vec4(0,0.5,0,1) * vertColor[i];
			fragTexCoord = vertTexCoord[i];
			fragNormal = vertNormal[i];
			EmitVertex();
			gl_Position = grassPos + matProjection * (matView * matModel * vec4(0,0,0,0) + vec4(-0.02, 0, 0, 0));
			fragColor = vec4(0,0.3,0,1) * vertColor[i];
			fragTexCoord = vertTexCoord[i];
			fragNormal = vertNormal[i];
			EmitVertex();
			EndPrimitive();
		}
	}
}
