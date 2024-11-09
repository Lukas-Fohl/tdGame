package main

import (
	"fmt"
	"math"
)

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
	position vec2
	dmgRange float64
	dmg      float64
	price    int
}

func tower_init(positionIn vec2, dmgRangeIn float64, dmgIn float64, priceIn int) tower {
	return tower{
		position: positionIn,
		dmgRange: dmgRangeIn,
		dmg:      dmgIn,
		price:    priceIn,
	}
}

func tower_dmgToETKList(towerIn tower, etkList [](*etk)) {
	for i := 0; i < len(etkList); i++ {
		disX := math.Abs(towerIn.position.x - etkList[i].position.x)
		disY := math.Abs(towerIn.position.y - etkList[i].position.y)
		if math.Sqrt(disX*disX+disY*disY) <= towerIn.dmgRange {
			etkList[i].health -= towerIn.dmg
			fmt.Println("hit")
		}
	}
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
	}
}

func main() {
	etkList := [](*etk){}

	myEtk := etk_init(vec_init(0.0, 0.0), 0.0, 1.0, 0.0, 0.0)
	myEtk2 := etk_init(vec_init(0.0, 0.0), 0.0, 1.0, 0.0, 0.0)
	etkList = append(etkList, &myEtk)
	etkList = append(etkList, &myEtk2)

	towerList := [](tower){}

	towerList1 := tower_init(vec_init(10.0, 9.0), 5.0, 20.0, 10)
	towerList2 := tower_init(vec_init(10.0, 9.0), 5.0, 20.0, 10)
	towerList = append(towerList, towerList1)
	towerList = append(towerList, towerList2)

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
		for j := 0; j < len(etkList); j++ {
			etkList[j].wayPointPerc += 0.1
			etk_poisitionFromPath(etkList[j], myPath)
			//fmt.Print(etkList[j].position.x)
			//fmt.Print(" ; ")
			//fmt.Println(etkList[j].position.y)
		}
		for j := 0; j < len(towerList); j++ {
			tower_dmgToETKList(towerList[j], etkList)
		}

		listToRemove := []int{}

		for j := len(etkList) - 1; j >= 0; j-- {
			if etkList[j].health <= 0.0 || etkList[j].wayPointPerc >= 100.0 {
				listToRemove = append(listToRemove, j)
			}
		}
		for k := len(listToRemove) - 1; k >= 0; k-- {
			idx := listToRemove[k]
			if idx+1 >= len(etkList) {
				etkList = etkList[:len(etkList)-1]
			} else {
				etkList = append(etkList[:idx], etkList[idx+1:]...)
			}
		}
		//remove etk from list
		//-> give mone
	}
	return
}

/*

ablauf:
	etk
		tick
			move
			check for full path -> kys
			check for health -> kys + money
	tower
		check for all in range
		-> apply dmg
	repead on emtpy etk list

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

TODO:
	restrict tower hit am
	restrict etk movement
*/
