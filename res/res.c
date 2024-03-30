#include "res.h"
#include "stddef.h"
#include "string.h"
#include "stdio.h"
#include "stdlib.h"
#include "raymath.h"

#define TINYOBJ_MALLOC RL_MALLOC
#define TINYOBJ_CALLOC RL_CALLOC
#define TINYOBJ_REALLOC RL_REALLOC
#define TINYOBJ_FREE RL_FREE

char* ReadData(const char* file, const char* dir) {
    size_t dir_len = (unsigned int)strlen(dir);
    size_t fnm_len = (unsigned int)strlen(file);
    char* cfile = (char*)RL_CALLOC(fnm_len, sizeof(char));
    if (cfile == NULL) {
        return NULL;
    }
    char* cdir = (char*)RL_CALLOC(dir_len, sizeof(char));
    if (cdir == NULL) {
        free(cfile);
        return NULL;
    }
    strncpy(cfile, file, fnm_len);
    strncpy(cdir, dir, dir_len);

    return GetData(cfile, cdir);
}

Model* LoadOBJ(char *fileName, char *fileText) {
    Model* model = (Model*)RL_CALLOC(1, sizeof(Model));
    if (model == NULL) {
        return NULL;
    }

    tinyobj_attrib_t attrib = { 0 };
    tinyobj_shape_t *meshes = NULL;
    unsigned int meshCount = 0;

    tinyobj_material_t *materials = NULL;
    unsigned int materialCount = 0;

    if (fileText != NULL) {
        unsigned int dataSize = (unsigned int)strlen(fileText);
        const char *workingDir = GetDirectoryPath(fileName); // Switch to OBJ directory for material path correctness

        unsigned int flags = TINYOBJ_FLAG_TRIANGULATE;

        int ret = tinyobj_parse_obj2(&attrib, &meshes, &meshCount, &materials, &materialCount, fileText, dataSize, flags, workingDir);

        if (ret != TINYOBJ_SUCCESS){
            free(model);
            return NULL;
        }

        // WARNING: We are not splitting meshes by materials (previous implementation)
        // Depending on the provided OBJ that was not the best option and it just crashed
        // so, implementation was simplified to prioritize parsed meshes
        model->meshCount = meshCount;

        // Set number of materials available
        // NOTE: There could be more materials available than meshes but it will be resolved at
        // model.meshMaterial, just assigning the right material to corresponding mesh
        model->materialCount = materialCount;
        if (model->materialCount == 0)
            model->materialCount = 1;

        // Init model meshes and materials
        model->meshes = (Mesh *)RL_CALLOC(model->meshCount, sizeof(Mesh));
        if (model->meshes == NULL) {
            free(model);
            return NULL;
        }
        model->meshMaterial = (int *)RL_CALLOC(model->meshCount, sizeof(int)); // Material index assigned to each mesh
        if (model->meshMaterial == NULL) {
            free(model->meshes);
            free(model);
            return NULL;
        }
        model->materials = (Material *)RL_CALLOC(model->materialCount, sizeof(Material));
        if (model->materials == NULL) {
            free(model->meshes);
            free(model->meshMaterial);
            free(model);
            return NULL;
        }

        // Process each provided mesh
        for (int i = 0; i < model->meshCount; i++) {
            // WARNING: We need to calculate the mesh triangles manually using meshes[i].face_offset
            // because in case of triangulated quads, meshes[i].length actually report quads,
            // despite the triangulation that is efectively considered on attrib.num_faces
            unsigned int tris = 0;
            if (i == model->meshCount - 1)
                tris = attrib.num_faces - meshes[i].face_offset;
            else
                tris = meshes[i + 1].face_offset;

            model->meshes[i].vertexCount = tris*3;
            model->meshes[i].triangleCount = tris;   // Face count (triangulated)
            model->meshes[i].vertices = (float *)RL_CALLOC(model->meshes[i].vertexCount*3, sizeof(float));
            model->meshes[i].texcoords = (float *)RL_CALLOC(model->meshes[i].vertexCount*2, sizeof(float));
            model->meshes[i].normals = (float *)RL_CALLOC(model->meshes[i].vertexCount*3, sizeof(float));
            model->meshMaterial[i] = 0;  // By default, assign material 0 to each mesh

            // Process all mesh faces
            for (unsigned int face = 0, f = meshes[i].face_offset, v = 0, vt = 0, vn = 0; face < tris; face++, f++, v += 3, vt += 3, vn += 3) {
                // Get indices for the face
                tinyobj_vertex_index_t idx0 = attrib.faces[f*3 + 0];
                tinyobj_vertex_index_t idx1 = attrib.faces[f*3 + 1];
                tinyobj_vertex_index_t idx2 = attrib.faces[f*3 + 2];

                // Fill vertices buffer (float) using vertex index of the face
                for (int n = 0; n < 3; n++)
                    model->meshes[i].vertices[v*3 + n] = attrib.vertices[idx0.v_idx*3 + n];
                for (int n = 0; n < 3; n++)
                    model->meshes[i].vertices[(v + 1)*3 + n] = attrib.vertices[idx1.v_idx*3 + n];
                for (int n = 0; n < 3; n++)
                    model->meshes[i].vertices[(v + 2)*3 + n] = attrib.vertices[idx2.v_idx*3 + n];

                if (attrib.num_texcoords > 0) {
                    // Fill texcoords buffer (float) using vertex index of the face
                    // NOTE: Y-coordinate must be flipped upside-down
                    model->meshes[i].texcoords[vt*2 + 0] = attrib.texcoords[idx0.vt_idx*2 + 0];
                    model->meshes[i].texcoords[vt*2 + 1] = 1.0f - attrib.texcoords[idx0.vt_idx*2 + 1];

                    model->meshes[i].texcoords[(vt + 1)*2 + 0] = attrib.texcoords[idx1.vt_idx*2 + 0];
                    model->meshes[i].texcoords[(vt + 1)*2 + 1] = 1.0f - attrib.texcoords[idx1.vt_idx*2 + 1];

                    model->meshes[i].texcoords[(vt + 2)*2 + 0] = attrib.texcoords[idx2.vt_idx*2 + 0];
                    model->meshes[i].texcoords[(vt + 2)*2 + 1] = 1.0f - attrib.texcoords[idx2.vt_idx*2 + 1];
                }

                if (attrib.num_normals > 0) {
                    // Fill normals buffer (float) using vertex index of the face
                    for (int n = 0; n < 3; n++)
                        model->meshes[i].normals[vn*3 + n] = attrib.normals[idx0.vn_idx*3 + n];
                    for (int n = 0; n < 3; n++)
                        model->meshes[i].normals[(vn + 1)*3 + n] = attrib.normals[idx1.vn_idx*3 + n];
                    for (int n = 0; n < 3; n++)
                        model->meshes[i].normals[(vn + 2)*3 + n] = attrib.normals[idx2.vn_idx*3 + n];
                }
            }
        }

        // Init model materials
        if (materialCount > 0)
            ProcessMaterialsOBJ(model->materials, materials, materialCount);
        else
            model->materials[0] = LoadMaterialDefault(); // Set default material for the mesh

        tinyobj_attrib_free(&attrib);
        tinyobj_shapes_free(meshes, model->meshCount);
        tinyobj_materials_free(materials, materialCount);
    }

    // Make sure model transform is set to identity matrix!
    model->transform = MatrixIdentity();

    if ((model->meshCount != 0) && (model->meshes != NULL)) {
        // Upload vertex data to GPU (static meshes)
        for (int i = 0; i < model->meshCount; i++) UploadMesh(&model->meshes[i], false);
    }

    if (model->materialCount == 0) {
        model->materialCount = 1;
        model->materials = (Material *)RL_CALLOC(model->materialCount, sizeof(Material));
        model->materials[0] = LoadMaterialDefault();

        if (model->meshMaterial == NULL)
            model->meshMaterial = (int *)RL_CALLOC(model->meshCount, sizeof(int));
    }


    return model;
}

void ProcessMaterialsOBJ(Material *materials, tinyobj_material_t *mats, int materialCount)
{
    // Init model mats
    for (int m = 0; m < materialCount; m++)
    {
        // Init material to default
        // NOTE: Uses default shader, which only supports MATERIAL_MAP_DIFFUSE
        materials[m] = LoadMaterialDefault();

        // Get default texture, in case no texture is defined
        // NOTE: rlgl default texture is a 1x1 pixel UNCOMPRESSED_R8G8B8A8
        materials[m].maps[MATERIAL_MAP_DIFFUSE].texture = (Texture2D){ rlGetTextureIdDefault(), 1, 1, 1, PIXELFORMAT_UNCOMPRESSED_R8G8B8A8 };

        if (mats[m].diffuse_texname != NULL) materials[m].maps[MATERIAL_MAP_DIFFUSE].texture = LoadTexture(mats[m].diffuse_texname);  //char *diffuse_texname; // map_Kd
        else materials[m].maps[MATERIAL_MAP_DIFFUSE].color = (Color){ (unsigned char)(mats[m].diffuse[0]*255.0f), (unsigned char)(mats[m].diffuse[1]*255.0f), (unsigned char)(mats[m].diffuse[2] * 255.0f), 255 }; //float diffuse[3];
        materials[m].maps[MATERIAL_MAP_DIFFUSE].value = 0.0f;

        if (mats[m].specular_texname != NULL) materials[m].maps[MATERIAL_MAP_SPECULAR].texture = LoadTexture(mats[m].specular_texname);  //char *specular_texname; // map_Ks
        materials[m].maps[MATERIAL_MAP_SPECULAR].color = (Color){ (unsigned char)(mats[m].specular[0]*255.0f), (unsigned char)(mats[m].specular[1]*255.0f), (unsigned char)(mats[m].specular[2] * 255.0f), 255 }; //float specular[3];
        materials[m].maps[MATERIAL_MAP_SPECULAR].value = 0.0f;

        if (mats[m].bump_texname != NULL) materials[m].maps[MATERIAL_MAP_NORMAL].texture = LoadTexture(mats[m].bump_texname);  //char *bump_texname; // map_bump, bump
        materials[m].maps[MATERIAL_MAP_NORMAL].color = WHITE;
        materials[m].maps[MATERIAL_MAP_NORMAL].value = mats[m].shininess;

        materials[m].maps[MATERIAL_MAP_EMISSION].color = (Color){ (unsigned char)(mats[m].emission[0]*255.0f), (unsigned char)(mats[m].emission[1]*255.0f), (unsigned char)(mats[m].emission[2] * 255.0f), 255 }; //float emission[3];

        if (mats[m].displacement_texname != NULL) materials[m].maps[MATERIAL_MAP_HEIGHT].texture = LoadTexture(mats[m].displacement_texname);  //char *displacement_texname; // disp
    }
}

void BuildPoseFromParentJoints(BoneInfo *bones, int boneCount, Transform *transforms) {
    for (int i = 0; i < boneCount; i++) {
        if (bones[i].parent >= 0) {
            if (bones[i].parent > i) {
                continue;
            }
            transforms[i].rotation = QuaternionMultiply(transforms[bones[i].parent].rotation, transforms[i].rotation);
            transforms[i].translation = Vector3RotateByQuaternion(transforms[i].translation, transforms[bones[i].parent].rotation);
            transforms[i].translation = Vector3Add(transforms[i].translation, transforms[bones[i].parent].translation);
            transforms[i].scale = Vector3Multiply(transforms[i].scale, transforms[bones[i].parent].scale);
        }
    }
}

Model* LoadIQM(char *fileName, char *fileText) {
    #define IQM_MAGIC           "INTERQUAKEMODEL" // IQM file magic number
    #define IQM_VERSION         2          // only IQM version 2 supported

    #define BONE_NAME_LENGTH    32          // BoneInfo name string length
    #define MESH_NAME_LENGTH    32          // Mesh name string length
    #define MATERIAL_NAME_LENGTH 32         // Material name string length

    int dataSize = 0;
    unsigned char *fileDataPtr = fileText;

    // IQM file structs
    //-----------------------------------------------------------------------------------
    typedef struct IQMHeader {
        char magic[16];
        unsigned int version;
        unsigned int dataSize;
        unsigned int flags;
        unsigned int num_text, ofs_text;
        unsigned int num_meshes, ofs_meshes;
        unsigned int num_vertexarrays, num_vertexes, ofs_vertexarrays;
        unsigned int num_triangles, ofs_triangles, ofs_adjacency;
        unsigned int num_joints, ofs_joints;
        unsigned int num_poses, ofs_poses;
        unsigned int num_anims, ofs_anims;
        unsigned int num_frames, num_framechannels, ofs_frames, ofs_bounds;
        unsigned int num_comment, ofs_comment;
        unsigned int num_extensions, ofs_extensions;
    } IQMHeader;

    typedef struct IQMMesh {
        unsigned int name;
        unsigned int material;
        unsigned int first_vertex, num_vertexes;
        unsigned int first_triangle, num_triangles;
    } IQMMesh;

    typedef struct IQMTriangle {
        unsigned int vertex[3];
    } IQMTriangle;

    typedef struct IQMJoint {
        unsigned int name;
        int parent;
        float translate[3], rotate[4], scale[3];
    } IQMJoint;

    typedef struct IQMVertexArray {
        unsigned int type;
        unsigned int flags;
        unsigned int format;
        unsigned int size;
        unsigned int offset;
    } IQMVertexArray;

    // NOTE: Below IQM structures are not used but listed for reference
    /*
    typedef struct IQMAdjacency {
        unsigned int triangle[3];
    } IQMAdjacency;

    typedef struct IQMPose {
        int parent;
        unsigned int mask;
        float channeloffset[10];
        float channelscale[10];
    } IQMPose;

    typedef struct IQMAnim {
        unsigned int name;
        unsigned int first_frame, num_frames;
        float framerate;
        unsigned int flags;
    } IQMAnim;

    typedef struct IQMBounds {
        float bbmin[3], bbmax[3];
        float xyradius, radius;
    } IQMBounds;
    */
    //-----------------------------------------------------------------------------------

    // IQM vertex data types
    enum {
        IQM_POSITION     = 0,
        IQM_TEXCOORD     = 1,
        IQM_NORMAL       = 2,
        IQM_TANGENT      = 3,       // NOTE: Tangents unused by default
        IQM_BLENDINDEXES = 4,
        IQM_BLENDWEIGHTS = 5,
        IQM_COLOR        = 6,
        IQM_CUSTOM       = 0x10     // NOTE: Custom vertex values unused by default
    };

    Model* model = (Model*)RL_CALLOC(1, sizeof(Model));
    if (model == NULL) {
        return NULL;
    }

    IQMMesh *imesh = NULL;
    IQMTriangle *tri = NULL;
    IQMVertexArray *va = NULL;
    IQMJoint *ijoint = NULL;

    float *vertex = NULL;
    float *normal = NULL;
    float *text = NULL;
    char *blendi = NULL;
    unsigned char *blendw = NULL;
    unsigned char *color = NULL;

    // In case file can not be read, return an empty model
    if (fileDataPtr == NULL) return model;

    // Read IQM header
    IQMHeader *iqmHeader = (IQMHeader *)fileDataPtr;

    if (memcmp(iqmHeader->magic, IQM_MAGIC, sizeof(IQM_MAGIC)) != 0) {
        return model;
    }

    if (iqmHeader->version != IQM_VERSION) {
        return model;
    }

    //fileDataPtr += sizeof(IQMHeader);       // Move file data pointer

    // Meshes data processing
    imesh = RL_MALLOC(iqmHeader->num_meshes*sizeof(IQMMesh));
    //fseek(iqmFile, iqmHeader->ofs_meshes, SEEK_SET);
    //fread(imesh, sizeof(IQMMesh)*iqmHeader->num_meshes, 1, iqmFile);
    memcpy(imesh, fileDataPtr + iqmHeader->ofs_meshes, iqmHeader->num_meshes*sizeof(IQMMesh));

    model->meshCount = iqmHeader->num_meshes;
    model->meshes = RL_CALLOC(model->meshCount, sizeof(Mesh));

    model->materialCount = model->meshCount;
    model->materials = (Material *)RL_CALLOC(model->materialCount, sizeof(Material));
    model->meshMaterial = (int *)RL_CALLOC(model->meshCount, sizeof(int));

    char name[MESH_NAME_LENGTH] = { 0 };
    char material[MATERIAL_NAME_LENGTH] = { 0 };

    for (int i = 0; i < model->meshCount; i++) {
        //fseek(iqmFile, iqmHeader->ofs_text + imesh[i].name, SEEK_SET);
        //fread(name, sizeof(char), MESH_NAME_LENGTH, iqmFile);
        memcpy(name, fileDataPtr + iqmHeader->ofs_text + imesh[i].name, MESH_NAME_LENGTH*sizeof(char));

        //fseek(iqmFile, iqmHeader->ofs_text + imesh[i].material, SEEK_SET);
        //fread(material, sizeof(char), MATERIAL_NAME_LENGTH, iqmFile);
        memcpy(material, fileDataPtr + iqmHeader->ofs_text + imesh[i].material, MATERIAL_NAME_LENGTH*sizeof(char));

        model->materials[i] = LoadMaterialDefault();

        model->meshes[i].vertexCount = imesh[i].num_vertexes;

        model->meshes[i].vertices = RL_CALLOC(model->meshes[i].vertexCount*3, sizeof(float));       // Default vertex positions
        model->meshes[i].normals = RL_CALLOC(model->meshes[i].vertexCount*3, sizeof(float));        // Default vertex normals
        model->meshes[i].texcoords = RL_CALLOC(model->meshes[i].vertexCount*2, sizeof(float));      // Default vertex texcoords

        model->meshes[i].boneIds = RL_CALLOC(model->meshes[i].vertexCount*4, sizeof(unsigned char));  // Up-to 4 bones supported!
        model->meshes[i].boneWeights = RL_CALLOC(model->meshes[i].vertexCount*4, sizeof(float));      // Up-to 4 bones supported!

        model->meshes[i].triangleCount = imesh[i].num_triangles;
        model->meshes[i].indices = RL_CALLOC(model->meshes[i].triangleCount*3, sizeof(unsigned short));

        // Animated vertex data, what we actually process for rendering
        // NOTE: Animated vertex should be re-uploaded to GPU (if not using GPU skinning)
        model->meshes[i].animVertices = RL_CALLOC(model->meshes[i].vertexCount*3, sizeof(float));
        model->meshes[i].animNormals = RL_CALLOC(model->meshes[i].vertexCount*3, sizeof(float));
    }

    // Triangles data processing
    tri = RL_MALLOC(iqmHeader->num_triangles*sizeof(IQMTriangle));
    //fseek(iqmFile, iqmHeader->ofs_triangles, SEEK_SET);
    //fread(tri, sizeof(IQMTriangle), iqmHeader->num_triangles, iqmFile);
    memcpy(tri, fileDataPtr + iqmHeader->ofs_triangles, iqmHeader->num_triangles*sizeof(IQMTriangle));

    for (int m = 0; m < model->meshCount; m++) {
        int tcounter = 0;

        for (unsigned int i = imesh[m].first_triangle; i < (imesh[m].first_triangle + imesh[m].num_triangles); i++) {
            // IQM triangles indexes are stored in counter-clockwise, but raylib processes the index in linear order,
            // expecting they point to the counter-clockwise vertex triangle, so we need to reverse triangle indexes
            // NOTE: raylib renders vertex data in counter-clockwise order (standard convention) by default
            model->meshes[m].indices[tcounter + 2] = tri[i].vertex[0] - imesh[m].first_vertex;
            model->meshes[m].indices[tcounter + 1] = tri[i].vertex[1] - imesh[m].first_vertex;
            model->meshes[m].indices[tcounter] = tri[i].vertex[2] - imesh[m].first_vertex;
            tcounter += 3;
        }
    }

    // Vertex arrays data processing
    va = RL_MALLOC(iqmHeader->num_vertexarrays*sizeof(IQMVertexArray));
    //fseek(iqmFile, iqmHeader->ofs_vertexarrays, SEEK_SET);
    //fread(va, sizeof(IQMVertexArray), iqmHeader->num_vertexarrays, iqmFile);
    memcpy(va, fileDataPtr + iqmHeader->ofs_vertexarrays, iqmHeader->num_vertexarrays*sizeof(IQMVertexArray));

    for (unsigned int i = 0; i < iqmHeader->num_vertexarrays; i++) {
        switch (va[i].type) {
            case IQM_POSITION: {
                vertex = RL_MALLOC(iqmHeader->num_vertexes*3*sizeof(float));
                //fseek(iqmFile, va[i].offset, SEEK_SET);
                //fread(vertex, iqmHeader->num_vertexes*3*sizeof(float), 1, iqmFile);
                memcpy(vertex, fileDataPtr + va[i].offset, iqmHeader->num_vertexes*3*sizeof(float));

                for (unsigned int m = 0; m < iqmHeader->num_meshes; m++) {
                    int vCounter = 0;
                    for (unsigned int i = imesh[m].first_vertex*3; i < (imesh[m].first_vertex + imesh[m].num_vertexes)*3; i++) {
                        model->meshes[m].vertices[vCounter] = vertex[i];
                        model->meshes[m].animVertices[vCounter] = vertex[i];
                        vCounter++;
                    }
                }
            } break;
            case IQM_NORMAL: {
                normal = RL_MALLOC(iqmHeader->num_vertexes*3*sizeof(float));
                //fseek(iqmFile, va[i].offset, SEEK_SET);
                //fread(normal, iqmHeader->num_vertexes*3*sizeof(float), 1, iqmFile);
                memcpy(normal, fileDataPtr + va[i].offset, iqmHeader->num_vertexes*3*sizeof(float));

                for (unsigned int m = 0; m < iqmHeader->num_meshes; m++) {
                    int vCounter = 0;
                    for (unsigned int i = imesh[m].first_vertex*3; i < (imesh[m].first_vertex + imesh[m].num_vertexes)*3; i++) {
                        model->meshes[m].normals[vCounter] = normal[i];
                        model->meshes[m].animNormals[vCounter] = normal[i];
                        vCounter++;
                    }
                }
            } break;
            case IQM_TEXCOORD: {
                text = RL_MALLOC(iqmHeader->num_vertexes*2*sizeof(float));
                //fseek(iqmFile, va[i].offset, SEEK_SET);
                //fread(text, iqmHeader->num_vertexes*2*sizeof(float), 1, iqmFile);
                memcpy(text, fileDataPtr + va[i].offset, iqmHeader->num_vertexes*2*sizeof(float));

                for (unsigned int m = 0; m < iqmHeader->num_meshes; m++) {
                    int vCounter = 0;
                    for (unsigned int i = imesh[m].first_vertex*2; i < (imesh[m].first_vertex + imesh[m].num_vertexes)*2; i++) {
                        model->meshes[m].texcoords[vCounter] = text[i];
                        vCounter++;
                    }
                }
            } break;
            case IQM_BLENDINDEXES: {
                blendi = RL_MALLOC(iqmHeader->num_vertexes*4*sizeof(char));
                //fseek(iqmFile, va[i].offset, SEEK_SET);
                //fread(blendi, iqmHeader->num_vertexes*4*sizeof(char), 1, iqmFile);
                memcpy(blendi, fileDataPtr + va[i].offset, iqmHeader->num_vertexes*4*sizeof(char));

                for (unsigned int m = 0; m < iqmHeader->num_meshes; m++) {
                    int boneCounter = 0;
                    for (unsigned int i = imesh[m].first_vertex*4; i < (imesh[m].first_vertex + imesh[m].num_vertexes)*4; i++) {
                        model->meshes[m].boneIds[boneCounter] = blendi[i];
                        boneCounter++;
                    }
                }
            } break;
            case IQM_BLENDWEIGHTS: {
                blendw = RL_MALLOC(iqmHeader->num_vertexes*4*sizeof(unsigned char));
                //fseek(iqmFile, va[i].offset, SEEK_SET);
                //fread(blendw, iqmHeader->num_vertexes*4*sizeof(unsigned char), 1, iqmFile);
                memcpy(blendw, fileDataPtr + va[i].offset, iqmHeader->num_vertexes*4*sizeof(unsigned char));

                for (unsigned int m = 0; m < iqmHeader->num_meshes; m++) {
                    int boneCounter = 0;
                    for (unsigned int i = imesh[m].first_vertex*4; i < (imesh[m].first_vertex + imesh[m].num_vertexes)*4; i++) {
                        model->meshes[m].boneWeights[boneCounter] = blendw[i]/255.0f;
                        boneCounter++;
                    }
                }
            } break;
            case IQM_COLOR: {
                color = RL_MALLOC(iqmHeader->num_vertexes*4*sizeof(unsigned char));
                //fseek(iqmFile, va[i].offset, SEEK_SET);
                //fread(blendw, iqmHeader->num_vertexes*4*sizeof(unsigned char), 1, iqmFile);
                memcpy(color, fileDataPtr + va[i].offset, iqmHeader->num_vertexes*4*sizeof(unsigned char));

                for (unsigned int m = 0; m < iqmHeader->num_meshes; m++) {
                    model->meshes[m].colors = RL_CALLOC(model->meshes[m].vertexCount*4, sizeof(unsigned char));

                    int vCounter = 0;
                    for (unsigned int i = imesh[m].first_vertex*4; i < (imesh[m].first_vertex + imesh[m].num_vertexes)*4; i++) {
                        model->meshes[m].colors[vCounter] = color[i];
                        vCounter++;
                    }
                }
            } break;
        }
    }

    // Bones (joints) data processing
    ijoint = RL_MALLOC(iqmHeader->num_joints*sizeof(IQMJoint));
    //fseek(iqmFile, iqmHeader->ofs_joints, SEEK_SET);
    //fread(ijoint, sizeof(IQMJoint), iqmHeader->num_joints, iqmFile);
    memcpy(ijoint, fileDataPtr + iqmHeader->ofs_joints, iqmHeader->num_joints*sizeof(IQMJoint));

    model->boneCount = iqmHeader->num_joints;
    model->bones = RL_MALLOC(iqmHeader->num_joints*sizeof(BoneInfo));
    model->bindPose = RL_MALLOC(iqmHeader->num_joints*sizeof(Transform));

    for (unsigned int i = 0; i < iqmHeader->num_joints; i++) {
        // Bones
        model->bones[i].parent = ijoint[i].parent;
        //fseek(iqmFile, iqmHeader->ofs_text + ijoint[i].name, SEEK_SET);
        //fread(model.bones[i].name, sizeof(char), BONE_NAME_LENGTH, iqmFile);
        memcpy(model->bones[i].name, fileDataPtr + iqmHeader->ofs_text + ijoint[i].name, BONE_NAME_LENGTH*sizeof(char));

        // Bind pose (base pose)
        model->bindPose[i].translation.x = ijoint[i].translate[0];
        model->bindPose[i].translation.y = ijoint[i].translate[1];
        model->bindPose[i].translation.z = ijoint[i].translate[2];

        model->bindPose[i].rotation.x = ijoint[i].rotate[0];
        model->bindPose[i].rotation.y = ijoint[i].rotate[1];
        model->bindPose[i].rotation.z = ijoint[i].rotate[2];
        model->bindPose[i].rotation.w = ijoint[i].rotate[3];

        model->bindPose[i].scale.x = ijoint[i].scale[0];
        model->bindPose[i].scale.y = ijoint[i].scale[1];
        model->bindPose[i].scale.z = ijoint[i].scale[2];
    }

    BuildPoseFromParentJoints(model->bones, model->boneCount, model->bindPose);

    RL_FREE(imesh);
    RL_FREE(tri);
    RL_FREE(va);
    RL_FREE(vertex);
    RL_FREE(normal);
    RL_FREE(text);
    RL_FREE(blendi);
    RL_FREE(blendw);
    RL_FREE(ijoint);
    RL_FREE(color);

    return model;
}

ModelAnimation* LoadAnimIQM(char *fileName, char *fileText, int *animCount) {
    #define IQM_MAGIC       "INTERQUAKEMODEL"   // IQM file magic number
    #define IQM_VERSION     2                   // only IQM version 2 supported

    #define BONE_NAME_LENGTH    32          // BoneInfo name string length
    #define MESH_NAME_LENGTH    32          // Mesh name string length
    #define MATERIAL_NAME_LENGTH 32         // Material name string length
            
    int dataSize = 0;
    unsigned char *fileDataPtr = fileText;

    typedef struct IQMHeader {
        char magic[16];
        unsigned int version;
        unsigned int dataSize;
        unsigned int flags;
        unsigned int num_text, ofs_text;
        unsigned int num_meshes, ofs_meshes;
        unsigned int num_vertexarrays, num_vertexes, ofs_vertexarrays;
        unsigned int num_triangles, ofs_triangles, ofs_adjacency;
        unsigned int num_joints, ofs_joints;
        unsigned int num_poses, ofs_poses;
        unsigned int num_anims, ofs_anims;
        unsigned int num_frames, num_framechannels, ofs_frames, ofs_bounds;
        unsigned int num_comment, ofs_comment;
        unsigned int num_extensions, ofs_extensions;
    } IQMHeader;

    typedef struct IQMJoint {
        unsigned int name;
        int parent;
        float translate[3], rotate[4], scale[3];
    } IQMJoint;

    typedef struct IQMPose {
        int parent;
        unsigned int mask;
        float channeloffset[10];
        float channelscale[10];
    } IQMPose;

    typedef struct IQMAnim {
        unsigned int name;
        unsigned int first_frame, num_frames;
        float framerate;
        unsigned int flags;
    } IQMAnim;

    // In case file can not be read, return an empty model
    if (fileDataPtr == NULL) return NULL;

    // Read IQM header
    IQMHeader *iqmHeader = (IQMHeader *)fileDataPtr;

    if (memcmp(iqmHeader->magic, IQM_MAGIC, sizeof(IQM_MAGIC)) != 0) {
        return NULL;
    }

    if (iqmHeader->version != IQM_VERSION) {
        return NULL;
    }

    // Get bones data
    IQMPose *poses = RL_MALLOC(iqmHeader->num_poses*sizeof(IQMPose));
    //fseek(iqmFile, iqmHeader->ofs_poses, SEEK_SET);
    //fread(poses, sizeof(IQMPose), iqmHeader->num_poses, iqmFile);
    memcpy(poses, fileDataPtr + iqmHeader->ofs_poses, iqmHeader->num_poses*sizeof(IQMPose));

    // Get animations data
    *animCount = iqmHeader->num_anims;
    IQMAnim *anim = RL_MALLOC(iqmHeader->num_anims*sizeof(IQMAnim));
    //fseek(iqmFile, iqmHeader->ofs_anims, SEEK_SET);
    //fread(anim, sizeof(IQMAnim), iqmHeader->num_anims, iqmFile);
    memcpy(anim, fileDataPtr + iqmHeader->ofs_anims, iqmHeader->num_anims*sizeof(IQMAnim));

    ModelAnimation *animations = RL_MALLOC(iqmHeader->num_anims*sizeof(ModelAnimation));

    // frameposes
    unsigned short *framedata = RL_MALLOC(iqmHeader->num_frames*iqmHeader->num_framechannels*sizeof(unsigned short));
    //fseek(iqmFile, iqmHeader->ofs_frames, SEEK_SET);
    //fread(framedata, sizeof(unsigned short), iqmHeader->num_frames*iqmHeader->num_framechannels, iqmFile);
    memcpy(framedata, fileDataPtr + iqmHeader->ofs_frames, iqmHeader->num_frames*iqmHeader->num_framechannels*sizeof(unsigned short));

    // joints
    IQMJoint *joints = RL_MALLOC(iqmHeader->num_joints*sizeof(IQMJoint));
    memcpy(joints, fileDataPtr + iqmHeader->ofs_joints, iqmHeader->num_joints*sizeof(IQMJoint));

    for (unsigned int a = 0; a < iqmHeader->num_anims; a++)
    {
        animations[a].frameCount = anim[a].num_frames;
        animations[a].boneCount = iqmHeader->num_poses;
        animations[a].bones = RL_MALLOC(iqmHeader->num_poses*sizeof(BoneInfo));
        animations[a].framePoses = RL_MALLOC(anim[a].num_frames*sizeof(Transform *));
        // animations[a].framerate = anim.framerate;     // TODO: Use animation framerate data?

        for (unsigned int j = 0; j < iqmHeader->num_poses; j++)
        {
            // If animations and skeleton are in the same file, copy bone names to anim
            if (iqmHeader->num_joints > 0)
                memcpy(animations[a].bones[j].name, fileDataPtr + iqmHeader->ofs_text + joints[j].name, BONE_NAME_LENGTH*sizeof(char));
            else
                strcpy(animations[a].bones[j].name, "ANIMJOINTNAME"); // default bone name otherwise
            animations[a].bones[j].parent = poses[j].parent;
        }

        for (unsigned int j = 0; j < anim[a].num_frames; j++) animations[a].framePoses[j] = RL_MALLOC(iqmHeader->num_poses*sizeof(Transform));

        int dcounter = anim[a].first_frame*iqmHeader->num_framechannels;

        for (unsigned int frame = 0; frame < anim[a].num_frames; frame++)
        {
            for (unsigned int i = 0; i < iqmHeader->num_poses; i++)
            {
                animations[a].framePoses[frame][i].translation.x = poses[i].channeloffset[0];

                if (poses[i].mask & 0x01)
                {
                    animations[a].framePoses[frame][i].translation.x += framedata[dcounter]*poses[i].channelscale[0];
                    dcounter++;
                }

                animations[a].framePoses[frame][i].translation.y = poses[i].channeloffset[1];

                if (poses[i].mask & 0x02)
                {
                    animations[a].framePoses[frame][i].translation.y += framedata[dcounter]*poses[i].channelscale[1];
                    dcounter++;
                }

                animations[a].framePoses[frame][i].translation.z = poses[i].channeloffset[2];

                if (poses[i].mask & 0x04)
                {
                    animations[a].framePoses[frame][i].translation.z += framedata[dcounter]*poses[i].channelscale[2];
                    dcounter++;
                }

                animations[a].framePoses[frame][i].rotation.x = poses[i].channeloffset[3];

                if (poses[i].mask & 0x08)
                {
                    animations[a].framePoses[frame][i].rotation.x += framedata[dcounter]*poses[i].channelscale[3];
                    dcounter++;
                }

                animations[a].framePoses[frame][i].rotation.y = poses[i].channeloffset[4];

                if (poses[i].mask & 0x10)
                {
                    animations[a].framePoses[frame][i].rotation.y += framedata[dcounter]*poses[i].channelscale[4];
                    dcounter++;
                }

                animations[a].framePoses[frame][i].rotation.z = poses[i].channeloffset[5];

                if (poses[i].mask & 0x20)
                {
                    animations[a].framePoses[frame][i].rotation.z += framedata[dcounter]*poses[i].channelscale[5];
                    dcounter++;
                }

                animations[a].framePoses[frame][i].rotation.w = poses[i].channeloffset[6];

                if (poses[i].mask & 0x40)
                {
                    animations[a].framePoses[frame][i].rotation.w += framedata[dcounter]*poses[i].channelscale[6];
                    dcounter++;
                }

                animations[a].framePoses[frame][i].scale.x = poses[i].channeloffset[7];

                if (poses[i].mask & 0x80)
                {
                    animations[a].framePoses[frame][i].scale.x += framedata[dcounter]*poses[i].channelscale[7];
                    dcounter++;
                }

                animations[a].framePoses[frame][i].scale.y = poses[i].channeloffset[8];

                if (poses[i].mask & 0x100)
                {
                    animations[a].framePoses[frame][i].scale.y += framedata[dcounter]*poses[i].channelscale[8];
                    dcounter++;
                }

                animations[a].framePoses[frame][i].scale.z = poses[i].channeloffset[9];

                if (poses[i].mask & 0x200)
                {
                    animations[a].framePoses[frame][i].scale.z += framedata[dcounter]*poses[i].channelscale[9];
                    dcounter++;
                }

                animations[a].framePoses[frame][i].rotation = QuaternionNormalize(animations[a].framePoses[frame][i].rotation);
            }
        }

        // Build frameposes
        for (unsigned int frame = 0; frame < anim[a].num_frames; frame++)
        {
            for (int i = 0; i < animations[a].boneCount; i++)
            {
                if (animations[a].bones[i].parent >= 0)
                {
                    animations[a].framePoses[frame][i].rotation = QuaternionMultiply(animations[a].framePoses[frame][animations[a].bones[i].parent].rotation, animations[a].framePoses[frame][i].rotation);
                    animations[a].framePoses[frame][i].translation = Vector3RotateByQuaternion(animations[a].framePoses[frame][i].translation, animations[a].framePoses[frame][animations[a].bones[i].parent].rotation);
                    animations[a].framePoses[frame][i].translation = Vector3Add(animations[a].framePoses[frame][i].translation, animations[a].framePoses[frame][animations[a].bones[i].parent].translation);
                    animations[a].framePoses[frame][i].scale = Vector3Multiply(animations[a].framePoses[frame][i].scale, animations[a].framePoses[frame][animations[a].bones[i].parent].scale);
                }
            }
        }
    }

    RL_FREE(joints);
    RL_FREE(framedata);
    RL_FREE(poses);
    RL_FREE(anim);

    return animations;
}
