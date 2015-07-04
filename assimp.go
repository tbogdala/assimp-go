// Copyright 2015, Timothy Bogdala <tdb@animal-machine.com>
// See the LICENSE file for more details.

/*
Package assimp implements a basic wrapper for the ASSIMP library: http://assimp.sourceforge.net/

At present there's only a hard-coded, basic file loader that returns
a basic MeshData slice.

*/
package assimp

/*
#cgo CPPFLAGS: -I/MinGW/msys/1.0/include -std=c99
#cgo LDFLAGS: -L/MinGW/msys/1.0/lib -lassimp -lz -lstdc++

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <assimp/cimport.h>
#include <assimp/scene.h>
#include <assimp/mesh.h>
#include <assimp/cimport.h>
#include <assimp/matrix4x4.h>
#include <assimp/postprocess.h>

struct aiAnimation* animation_at(struct aiScene* s, unsigned int index)
{
	return s->mAnimations[index];
}

char* animation_name(struct aiAnimation* a)
{
	return a->mName.data;
}

struct aiMesh* mesh_at(struct aiScene* s, unsigned int index)
{
	return s->mMeshes[index];
}

struct aiVector3D* mesh_vertex_at(struct aiMesh* m, unsigned long index)
{
	return &(m->mVertices[index]);
}

struct aiVector3D* mesh_normal_at(struct aiMesh* m, unsigned long index)
{
	return &(m->mNormals[index]);
}

struct aiVector3D* mesh_tangent_at(struct aiMesh* m, unsigned long index)
{
	return &(m->mTangents[index]);
}

struct aiVector3D* mesh_uv_channel_at(struct aiMesh* m, unsigned long index)
{
	return m->mTextureCoords[index];
}

struct aiVector3D* mesh_uv_at(struct aiVector3D* uvChan, unsigned long index)
{
	return &(uvChan[index]);
}

struct aiBone* mesh_bone_at(struct aiMesh* m, unsigned long index)
{
	struct aiBone* b =  m->mBones[index];
	return b;
}

char* mesh_bone_name_at(struct aiMesh* m, unsigned long index)
{
	struct aiBone* b =  m->mBones[index];
	return b->mName.data;
}

struct aiVertexWeight* bone_vertex_weight_at(struct aiBone* b, unsigned long index)
{
	struct aiVertexWeight* w =  &b->mWeights[index];
	return w;
}

struct face {
	unsigned int x,y,z;
};

struct face mesh_face_at(struct aiMesh* m, unsigned long index)
{
	struct face f;
	struct aiFace* tempFace = &(m->mFaces[index]);
	unsigned int* tempIndices = tempFace->mIndices;
	f.x = tempIndices[0];
	f.y = tempIndices[1];
	f.z = tempIndices[2];
	return f;
}

struct aiNode* find_assimp_node(struct aiNode* node, const char* name)
{
	if (strcmp(node->mName.data, name) == 0) return node;
	for (unsigned int i=0; i<node->mNumChildren; ++i) {
		struct aiNode* p = find_assimp_node(node->mChildren[i], name);
		if (p) return p;
	}

	return NULL;
}

struct aiMatrix4x4* mesh_bone_transform(struct aiNode* root_node, struct aiMesh* m, unsigned long index)
{
	struct aiBone* b =  m->mBones[index];
	struct aiNode* n = find_assimp_node(root_node, b->mName.data);
	return &n->mTransformation;
}

struct aiMatrix4x4* mesh_bone_offset(struct aiMesh* m, unsigned long index)
{
	struct aiBone* b =  m->mBones[index];
	return &b->mOffsetMatrix;
}

struct aiNodeAnim* animation_channel_at(struct aiAnimation* a, unsigned long index)
{
	struct aiNodeAnim* na =  a->mChannels[index];
	return na;
}

char* channel_name(struct aiNodeAnim* a)
{
	return a->mNodeName.data;
}

struct aiVectorKey* animation_channel_poskey_at(struct aiNodeAnim* chan, unsigned long index)
{
	struct aiVectorKey* vk = &chan->mPositionKeys[index];
	return vk;
}

struct aiVectorKey* animation_channel_scalekey_at(struct aiNodeAnim* chan, unsigned long index)
{
	struct aiVectorKey* vk = &chan->mScalingKeys[index];
	return vk;
}

struct aiQuatKey* animation_channel_rotkey_at(struct aiNodeAnim* chan, unsigned long index)
{
	struct aiQuatKey* qk = &chan->mRotationKeys[index];
	return qk;
}


char* node_name(struct aiNode* n) {
  return n->mName.data;
}

struct aiMatrix4x4* scene_root_transform(struct aiScene* s)
{
	return &s->mRootNode->mTransformation;
}

*/
import "C"
import (
	"fmt"
	"unsafe"

	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/gombz"
)

// MatToGombzMat converts the row-major order of assimp
// to the column-major order of OpenGL
func MatToGombzMat(src *C.struct_aiMatrix4x4, dest []float32) {
	dest[0] = float32(src.a1)
	dest[1] = float32(src.b1)
	dest[2] = float32(src.c1)
	dest[3] = float32(src.d1)

	dest[4] = float32(src.a2)
	dest[5] = float32(src.b2)
	dest[6] = float32(src.c2)
	dest[7] = float32(src.d2)

	dest[8] = float32(src.a3)
	dest[9] = float32(src.b3)
	dest[10] = float32(src.c3)
	dest[11] = float32(src.d3)

	dest[12] = float32(src.a4)
	dest[13] = float32(src.b4)
	dest[14] = float32(src.c4)
	dest[15] = float32(src.d4)
}

// ParseFile loads a file at the given file path and returns all of
// the MeshData objects that get created from the meshes contained.
// err is non-nil on fail.
func ParseFile(modelFile string) (outMeshes []*gombz.Mesh, err error) {
	///////////////////////////////////////////////////////////
	// attempt to load the file
	cModelFile := C.CString(modelFile)
	defer C.free(unsafe.Pointer(cModelFile))

	cScene := C.aiImportFile(cModelFile,
		C.aiProcess_JoinIdenticalVertices|
			C.aiProcess_Triangulate|
			C.aiProcess_GenNormals|
			C.aiProcess_CalcTangentSpace|
			C.aiProcess_FindInvalidData|
			C.aiProcess_LimitBoneWeights|
			C.aiProcess_ImproveCacheLocality|
			C.aiProcess_FixInfacingNormals|
			C.aiProcess_OptimizeMeshes|
			C.aiProcess_ValidateDataStructure)

	// make sure that we got a scene back
	if uintptr(unsafe.Pointer(cScene)) == 0 {
		return nil, fmt.Errorf("Unable to load %s.\n", modelFile)
	}

	// make sure we have at least one mesh
	if cScene.mNumMeshes < 1 {
		return nil, fmt.Errorf("Unable to load %s -- no meshes were found!\n", modelFile)
	}

	// loop through each mesh
	outMeshes = make([]*gombz.Mesh, uint(cScene.mNumMeshes))
	for i := uint(0); i < uint(cScene.mNumMeshes); i++ {
		cMesh := C.mesh_at(cScene, C.uint(i))

		/*
		   // write out some information about the mesh
		   fmt.Printf("Mesh index: %d\n", i)
		   fmt.Printf("\tFace count: %d\n", cMesh.mNumFaces)
		   fmt.Printf("\tBone count: %d\n", cMesh.mNumBones)
		   fmt.Printf("\tUV component count: %d\n", cMesh.mNumUVComponents[0])
		   fmt.Printf("\tMaterial index: %d\n", cMesh.mMaterialIndex)
		   if cMesh.mTangents != nil {
		     fmt.Printf("\tHas tangents: true\n")
		   } else {
		     fmt.Printf("\tHas tangents: false\n")
		   }
		*/

		// fill up our data structure
		outMesh := new(gombz.Mesh)
		outMesh.FaceCount = uint32(cMesh.mNumFaces)
		outMesh.BoneCount = uint32(cMesh.mNumBones)
		outMesh.VertexCount = uint32(cMesh.mNumVertices)

		// copy the verts
		outMesh.Vertices = make([]mgl.Vec3, outMesh.VertexCount)
		for vi := uint32(0); vi < outMesh.VertexCount; vi++ {
			cVec3 := C.mesh_vertex_at(cMesh, C.ulong(vi))
			outMesh.Vertices[vi][0] = float32(cVec3.x)
			outMesh.Vertices[vi][1] = float32(cVec3.y)
			outMesh.Vertices[vi][2] = float32(cVec3.z)
		}

		// copy the faces
		outMesh.Faces = make([]gombz.MeshFace, outMesh.FaceCount)
		for fi := uint32(0); fi < outMesh.FaceCount; fi++ {
			cFace := C.mesh_face_at(cMesh, C.ulong(fi))
			outMesh.Faces[fi][0] = uint32(cFace.x)
			outMesh.Faces[fi][1] = uint32(cFace.y)
			outMesh.Faces[fi][2] = uint32(cFace.z)
		}

		// copy the normals
		if uintptr(unsafe.Pointer(cMesh.mNormals)) != 0 {
			outMesh.Normals = make([]mgl.Vec3, outMesh.VertexCount)
			for vi := uint32(0); vi < outMesh.VertexCount; vi++ {
				cNormal := C.mesh_normal_at(cMesh, C.ulong(vi))
				outMesh.Normals[vi][0] = float32(cNormal.x)
				outMesh.Normals[vi][1] = float32(cNormal.y)
				outMesh.Normals[vi][2] = float32(cNormal.z)
			}
		}

		// copy the tangents
		if uintptr(unsafe.Pointer(cMesh.mTangents)) != 0 {
			outMesh.Tangents = make([]mgl.Vec3, outMesh.VertexCount)
			for vi := uint32(0); vi < outMesh.VertexCount; vi++ {
				cTangent := C.mesh_tangent_at(cMesh, C.ulong(vi))
				outMesh.Tangents[vi][0] = float32(cTangent.x)
				outMesh.Tangents[vi][1] = float32(cTangent.y)
				outMesh.Tangents[vi][2] = float32(cTangent.z)
			}
		}

		// copy the UV channels
		for uvchi := uint32(0); uvchi < gombz.MaxUVChannelCount; uvchi++ {
			cUVChannel := C.mesh_uv_channel_at(cMesh, C.ulong(uvchi))
			if uintptr(unsafe.Pointer(cUVChannel)) != 0 {
				// if we have a valid UV channel, copy all of the UV's -- one per vert
				outMesh.UVChannels[uvchi] = make([]mgl.Vec2, outMesh.VertexCount)
				for vi := uint32(0); vi < outMesh.VertexCount; vi++ {
					cUV := C.mesh_uv_at(cUVChannel, C.ulong(vi))
					outMesh.UVChannels[uvchi][vi][0] = float32(cUV.x)
					outMesh.UVChannels[uvchi][vi][1] = float32(cUV.y)
				}
			}
		}

		// copy the bones
		if uintptr(unsafe.Pointer(cMesh.mBones)) != 0 {
			outMesh.Bones = make([]gombz.Bone, cMesh.mNumBones)
			outMesh.VertexWeightIds = make([]mgl.Vec4, outMesh.VertexCount)
			outMesh.VertexWeights = make([]mgl.Vec4, outMesh.VertexCount)

			for bi := uint32(0); bi < outMesh.BoneCount; bi++ {
				// setup basic bone properties
				cBone := C.mesh_bone_at(cMesh, C.ulong(bi))
				outMesh.Bones[bi].Id = int32(bi)
				outMesh.Bones[bi].Name = C.GoString(C.mesh_bone_name_at(cMesh, C.ulong(bi)))
				//fmt.Printf("\tBone #%d ; Weights=%d ; Name=%s\n", bi, cBone.mNumWeights, outMesh.Bones[bi].Name)

				// copy over the offset matrix that transforms from mesh space to
				// bone space in pose mode
				cOffsetMat4x4 := C.mesh_bone_offset(cMesh, C.ulong(bi))
				MatToGombzMat(cOffsetMat4x4, outMesh.Bones[bi].Offset[:])

				// copy over the transform matrix (relative to parent)
				cTransformMat4x4 := C.mesh_bone_transform(cScene.mRootNode, cMesh, C.ulong(bi))
				MatToGombzMat(cTransformMat4x4, outMesh.Bones[bi].Transform[:])

				/*
					printMat := func(m mgl.Mat4) {
						groggy.Logsf("DEBUG", "\t%3.1f\t%3.1f\t%3.1f\t%3.1f", m[0], m[4], m[8], m[12])
						groggy.Logsf("DEBUG", "\t%3.1f\t%3.1f\t%3.1f\t%3.1f", m[1], m[5], m[9], m[13])
						groggy.Logsf("DEBUG", "\t%3.1f\t%3.1f\t%3.1f\t%3.1f", m[2], m[6], m[10], m[14])
						groggy.Logsf("DEBUG", "\t%3.1f\t%3.1f\t%3.1f\t%3.1f", m[3], m[7], m[11], m[15])
					}
					groggy.Logsf("DEBUG", "bone.Transform")
					printMat(outMesh.Bones[bi].Transform)
					groggy.Logsf("DEBUG", "bone.Offset")
					printMat(outMesh.Bones[bi].Offset)
				*/

				// copy over the vertex weights
				for wi := C.uint(0); wi < cBone.mNumWeights; wi++ {
					cWeight := C.bone_vertex_weight_at(cBone, C.ulong(wi))
					//fmt.Printf("\t\tWeight %d ; vert=%d ; value=%f\n", wi, cWeight.mVertexId, cWeight.mWeight)

					// get the curent weights for the vertex by id
					tmpWeightVec := outMesh.VertexWeights[cWeight.mVertexId]

					// see if there's an empty spot to set a weight
					for twi := 0; twi < 4; twi++ {
						if tmpWeightVec[twi] == 0.0 {
							outMesh.VertexWeights[cWeight.mVertexId][twi] = float32(cWeight.mWeight)
							outMesh.VertexWeightIds[cWeight.mVertexId][twi] = float32(bi)
							break
						}

						// Note: DOES NOT RAISE AN ERROR IF 4 BONES ARE ALREADY ASSIGNED
						if twi == 4 {
							fmt.Printf("TOO MANY WEIGHTS: Weight %d ; vert=%d ; value=%f\n", wi, cWeight.mVertexId, cWeight.mWeight)
						}
					} // twi
				} // wi
			} // bi
		}

		// now that all bones are copied over, time to set parent id's ...
		for bi := uint32(0); bi < outMesh.BoneCount; bi++ {
			bone := &outMesh.Bones[bi]

			// start with no parent
			bone.Parent = -1

			// find the scene node for the bone
			cBoneName := C.CString(bone.Name)
			cAssimpNode := C.find_assimp_node(cScene.mRootNode, cBoneName)
			C.free(unsafe.Pointer(cBoneName))

			if uintptr(unsafe.Pointer(cAssimpNode)) != 0 {
				// get the scene node for the parent bone
				cAssimpParentNode := cAssimpNode.mParent
				if uintptr(unsafe.Pointer(cAssimpParentNode)) != 0 {
					parentName := C.GoString(C.node_name(cAssimpParentNode))

					// now loop through the bones again and find the id for the bone
					// matching the parent bone name.
					for pi := uint32(0); pi < outMesh.BoneCount; pi++ {
						parentBone := outMesh.Bones[pi]
						if parentName == parentBone.Name {
							//fmt.Printf("Bone name == %s ; parent == %s ; parentid == %d\n", bone.Name, parentBone.Name, parentBone.Id)
							// we found the parent, so set the bone's parent id now.
							bone.Parent = parentBone.Id
							break
						}
					} // pi
				}
			}
		} // bi

		// process all of the animations if they exist
		if uintptr(unsafe.Pointer(cScene.mAnimations)) != 0 {
			// create the Animations slice
			outMesh.Animations = make([]gombz.Animation, cScene.mNumAnimations)

			// loop through all of the animations
			for aniIdx := C.uint(0); aniIdx < cScene.mNumAnimations; aniIdx++ {
				cAni := C.animation_at(cScene, aniIdx)
				animation := &outMesh.Animations[aniIdx]

				// setup the animation object
				animation.Name = C.GoString(C.animation_name(cAni))
				animation.Duration = float32(cAni.mDuration)
				animation.TicksPerSecond = float32(cAni.mTicksPerSecond)

				// TODO: make this a rotation flag
				// Note: this fixes an export problem from blender where the armature doesn't seem to respect
				// the axis settings in the exporter.
				x90Q := mgl.QuatRotate(mgl.DegToRad(90.0), mgl.Vec3{1.0, 0.0, 0.0})
				x180Q := mgl.QuatRotate(mgl.DegToRad(180.0), mgl.Vec3{0.0, 1.0, 0.0})
				x90QMat := x90Q.Mul(x180Q).Mat4()

				cRootMat4x4 := C.scene_root_transform(cScene)
				MatToGombzMat(cRootMat4x4, animation.Transform[:])
				animation.Transform = animation.Transform.Inv().Mul4(x90QMat)

				// now setup all of the animation channels
				animation.Channels = make([]gombz.AnimationChannel, cAni.mNumChannels)
				for aniChI := C.uint(0); aniChI < cAni.mNumChannels; aniChI++ {
					aniChan := &animation.Channels[aniChI]
					cNodeAni := C.animation_channel_at(cAni, C.ulong(aniChI))

					// set the channel name
					aniChan.Name = C.GoString(C.channel_name(cNodeAni))

					// try to find the bone with the same name to set the BoneId.
					aniChan.BoneId = -1
					for _, bone := range outMesh.Bones {
						if bone.Name == aniChan.Name {
							// we found a match, so set the BoneId.
							aniChan.BoneId = bone.Id
							break
						}
					}

					// create the slices for the keys
					aniChan.PositionKeys = make([]gombz.AnimationVec3Key, cNodeAni.mNumPositionKeys)
					aniChan.ScaleKeys = make([]gombz.AnimationVec3Key, cNodeAni.mNumScalingKeys)
					aniChan.RotationKeys = make([]gombz.AnimationQuatKey, cNodeAni.mNumRotationKeys)

					// copy over the position keys
					for pki := C.uint(0); pki < cNodeAni.mNumPositionKeys; pki++ {
						cPosKey := C.animation_channel_poskey_at(cNodeAni, C.ulong(pki))
						aniChan.PositionKeys[pki].Time = float32(cPosKey.mTime)
						aniChan.PositionKeys[pki].Key[0] = float32(cPosKey.mValue.x)
						aniChan.PositionKeys[pki].Key[1] = float32(cPosKey.mValue.y)
						aniChan.PositionKeys[pki].Key[2] = float32(cPosKey.mValue.z)
					}

					// copy over the scale keys
					for ski := C.uint(0); ski < cNodeAni.mNumScalingKeys; ski++ {
						cScaleKey := C.animation_channel_scalekey_at(cNodeAni, C.ulong(ski))
						aniChan.ScaleKeys[ski].Time = float32(cScaleKey.mTime)
						aniChan.ScaleKeys[ski].Key[0] = float32(cScaleKey.mValue.x)
						aniChan.ScaleKeys[ski].Key[1] = float32(cScaleKey.mValue.y)
						aniChan.ScaleKeys[ski].Key[2] = float32(cScaleKey.mValue.z)
					}

					// copy over the rotation keys
					for rki := C.uint(0); rki < cNodeAni.mNumRotationKeys; rki++ {
						cRotKey := C.animation_channel_rotkey_at(cNodeAni, C.ulong(rki))
						aniChan.RotationKeys[rki].Time = float32(cRotKey.mTime)
						aniChan.RotationKeys[rki].Key.W = float32(cRotKey.mValue.w)
						aniChan.RotationKeys[rki].Key.V[0] = float32(cRotKey.mValue.x)
						aniChan.RotationKeys[rki].Key.V[1] = float32(cRotKey.mValue.y)
						aniChan.RotationKeys[rki].Key.V[2] = float32(cRotKey.mValue.z)
					}
				} // aniChI

			} // aniIdx
		}

		// add the new mesh to the slice
		outMeshes[i] = outMesh
	}

	// drop the scene now that we got our data
	C.aiReleaseImport(cScene)

	return outMeshes, nil
}
