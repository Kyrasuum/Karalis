#include "raylib.h"
#include "rlgl.h"
#include "shader.h"

void SetLocs(int* locs) {
	// Set default shader locations: attributes locations
	locs[RL_SHADER_LOC_VERTEX_POSITION] = getloc("vertexPosition");
	locs[RL_SHADER_LOC_VERTEX_TEXCOORD01] = getloc("vertexTexCoord");
	locs[RL_SHADER_LOC_VERTEX_COLOR] = getloc("vertexColor");

	// Set default shader locations: uniform locations
	locs[RL_SHADER_LOC_MATRIX_MVP] = getloc("mvp");
	locs[RL_SHADER_LOC_COLOR_DIFFUSE] = getloc("colDiffuse");
	locs[RL_SHADER_LOC_MAP_DIFFUSE] = getloc("texture0");
}
