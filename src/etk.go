package main

type etk struct {
	position     vec2
	radius       float64
	health       float64
	reward       int
	wayPointPerc float64 //0..100
	speed        float64
	texturePath  string
}

func etk_init(positionIn vec2, radiusIn float64, healthIn float64, rewardIn int, wayPointPercIn float64, speedIn float64, texturePathIn string) etk {
	return etk{
		position:     positionIn,
		radius:       radiusIn,
		health:       healthIn,
		reward:       rewardIn,
		wayPointPerc: wayPointPercIn,
		speed:        speedIn,
		texturePath:  texturePathIn,
	}
}

type path struct {
	wayPoint []vec2
	//speed maybe
}

func path_init(vecArray []vec2) path {
	return path{wayPoint: vecArray}
}

func (etkIn *etk) poisitionFromPath(pathIn path) {
	arrayPoisition := int((etkIn.wayPointPerc / 100.0) * (float64)(len(pathIn.wayPoint)-1))
	if arrayPoisition < len(pathIn.wayPoint)-1 {
		vectorMoveX := pathIn.wayPoint[arrayPoisition+1].x - pathIn.wayPoint[arrayPoisition].x
		vectorMoveY := pathIn.wayPoint[arrayPoisition+1].y - pathIn.wayPoint[arrayPoisition].y
		vecStepInPerc := 100.0 / (float64)(len(pathIn.wayPoint)-1)
		currentPercFull := etkIn.wayPointPerc
		for currentPercFull >= 0 {
			currentPercFull -= vecStepInPerc
		}
		currentPercFull += vecStepInPerc
		etkIn.position.x = pathIn.wayPoint[arrayPoisition].x + (vectorMoveX * (currentPercFull / vecStepInPerc))
		etkIn.position.y = pathIn.wayPoint[arrayPoisition].y + (vectorMoveY * (currentPercFull / vecStepInPerc))
	}
}
