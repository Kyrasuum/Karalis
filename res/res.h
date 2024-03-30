#include "raylib.h"
#include "tinyobj_loader.h"

extern char* GetData(char* file, char* dir);
Model* LoadOBJ(char *fileName, char *fileText);
Model* LoadIQM(char *fileName, char *fileText);
ModelAnimation* LoadAnimIQM(char *fileName, char *fileText, int *animCount);

void ProcessMaterialsOBJ(Material *rayMaterials, tinyobj_material_t *materials, int materialCount);
extern unsigned int rlGetTextureIdDefault(void);
