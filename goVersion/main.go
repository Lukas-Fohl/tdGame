package main

import "fmt"

type play struct {
	money      int
	etkMax     int
	etkCurrent int
}

type vec2 struct {
	x, y float64
}

func vec_init(xIn float64, yIn float64) vec2 {
	return vec2{x: xIn, y: yIn}
}

type etk struct {
	position     vec2
	radius       float64
	health       float64
	reward       int
	wayPointPerc float64 //0..100
}

func etk_init(positionIn vec2, radiusIn float64, healthIn float64, rewardIn int, wayPointPercIn float64) etk {
	return etk{
		position:     positionIn,
		radius:       radiusIn,
		health:       healthIn,
		reward:       rewardIn,
		wayPointPerc: wayPointPercIn,
	}
}

type tower struct {
	positionX float64
	positionY float64
	dmgRange  float64
	dmg       float64
	price     int
}

type path struct {
	wayPoint []vec2
}

func path_init(vecArray []vec2) path {
	return path{wayPoint: vecArray}
}

func etk_poisitionFromPath(etkIn *etk, pathIn path) {
	arrayPoisition := int((etkIn.wayPointPerc / 100.0) * (float64)(len(pathIn.wayPoint)-1))
	if arrayPoisition < len(pathIn.wayPoint) {
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
		//currentVecPerc := (arrayPoisition / len(pathIn.wayPoint))
	}
}

func main() {
	myEtk := etk_init(vec_init(0.0, 0.0), 0.0, 1.0, 0.0, 0.0)
	myPath := path_init([]vec2{
		vec_init(0.0, 0.0),
		vec_init(5.0, 5.0),
		vec_init(10.0, 10.0),
		vec_init(15.0, 15.0),
		vec_init(20.0, 20.0),
		vec_init(25.0, 25.0),
		vec_init(30.0, 30.0),
		vec_init(35.0, 35.0),
		vec_init(40.0, 40.0),
		vec_init(45.0, 45.0),
	})

	for i := 0; i < 1000; i++ {
		myEtk.wayPointPerc += 0.1
		etk_poisitionFromPath(&myEtk, myPath)
		fmt.Print(myEtk.position.x)
		fmt.Print(" ; ")
		fmt.Println(myEtk.position.y)
	}
	return
}

/*
Idee:
	liste von etks:
		helth
		position
	liste von toweren:
		dmg
		wdh
		range
		attack:
			raycast --> kill
			play animation
		check for obst???
	path:
		liste an punkten
		--> interpolate

*/
