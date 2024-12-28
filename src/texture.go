package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type textureSave struct {
	texture  rl.Texture2D
	filePath string
}

func textureSave_init(pathIn string) textureSave {
	tempImg := rl.LoadImage("../img/" + pathIn)
	tempTexture := rl.LoadTextureFromImage(tempImg)
	rl.UnloadImage(tempImg)
	return textureSave{
		texture:  tempTexture,
		filePath: pathIn,
	}
}

func findTexture(textureListIn []textureSave, pathIn string) rl.Texture2D {
	for i := 0; i < len(textureListIn); i++ {
		if textureListIn[i].filePath == pathIn {
			return textureListIn[i].texture
		}
	}
	return textureSave_init(pathIn).texture
}
