package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type play struct {
	money         int
	etkMaxDmg     int
	etkCurrentDmg int
	level         int
	levelMsStart  int64
}

type vec2 struct {
	x, y float64
}

func vec_init(xIn float64, yIn float64) vec2 {
	return vec2{x: xIn, y: yIn}
}

func spawnListHandle(spawnList []spawn, etkList []etk, myPlay play) ([]spawn, []etk, play) {
	//check if spawns for this level are left
	hasSpawnWithSameLevel := false
	for j := 0; j < len(spawnList); j++ {
		if spawnList[j].amount > 0 && spawnList[j].level == myPlay.level {
			hasSpawnWithSameLevel = true
		}
	}
	if !hasSpawnWithSameLevel && len(etkList) == 0 {
		myPlay.level++
	}

	//spawn etks
	smallestSpawnByOrd := &spawn{orderNum: 1024}
	for j := 0; j < len(spawnList); j++ {
		if spawnList[j].orderNum < smallestSpawnByOrd.orderNum && spawnList[j].level == myPlay.level {
			smallestSpawnByOrd = &spawnList[j]
		}
	}

	etkList = append(etkList, smallestSpawnByOrd.getSpawnEtks()...)

	//remove all spawns with amount 0, level < level
	spawnToRemove := []int{}

	for j := len(spawnList) - 1; j >= 0; j-- {
		if (spawnList[j].amount <= 0 || spawnList[j].level < myPlay.level) && len(etkList) == 0 {
			spawnToRemove = append(spawnToRemove, j)
		}
	}

	for k := len(spawnToRemove) - 1; k >= 0; k-- {
		idx := spawnToRemove[k]
		if idx+1 >= len(spawnList) {
			spawnList = spawnList[:len(spawnList)-1]
		} else {
			spawnList = append(spawnList[:idx], spawnList[idx+1:]...)
		}
	}
	return spawnList, etkList, myPlay
}

func etkListHandle(spawnList []spawn, etkList []etk, towerList []tower, myPlay play, deltaTime float64, myPath path) ([]spawn, []etk, []tower, play) {
	//move etks
	for j := 0; j < len(etkList); j++ {
		etkList[j].wayPointPerc += 0.1 * deltaTime
		(&etkList[j]).poisitionFromPath(myPath)
		//fmt.Print(etkList[j].position.x)
		//fmt.Print(" ; ")
		//fmt.Println(etkList[j].position.y)
	}

	//dmg to etks
	etkRefList := [](*etk){}
	for j := 0; j < len(etkList); j++ {
		etkRefList = append(etkRefList, &etkList[j])
	}

	for j := 0; j < len(towerList); j++ {
		(&towerList[j]).dmgToETKList(etkRefList)
	}

	//kill etks
	listToRemove := []int{}

	for j := len(etkList) - 1; j >= 0; j-- {
		if etkList[j].health <= 0.0 || etkList[j].wayPointPerc >= 100.0 {
			listToRemove = append(listToRemove, j)
		}

		if etkList[j].health <= 0.0 {
			myPlay.money += etkList[j].reward
		}
		if etkList[j].wayPointPerc >= 100.0 {
			myPlay.etkCurrentDmg += 1
		}
	}

	//remove etk from list
	for k := len(listToRemove) - 1; k >= 0; k-- {
		idx := listToRemove[k]
		if idx+1 >= len(etkList) {
			etkList = etkList[:len(etkList)-1]
		} else {
			etkList = append(etkList[:idx], etkList[idx+1:]...)
		}
	}
	return spawnList, etkList, towerList, myPlay
}

func main() {
	screenWidth := int32(1080)
	screenHeight := int32(720)

	rl.InitWindow(screenWidth, screenHeight, "omtfyb")
	rl.SetTargetFPS(60)
	defer rl.CloseWindow()
	blockImg := rl.LoadImage("block.png")
	blockTexture := rl.LoadTextureFromImage(blockImg)
	rl.UnloadImage(blockImg)

	etkImg := rl.LoadImage("etk.png")
	etkTexture := rl.LoadTextureFromImage(etkImg)
	rl.UnloadImage(etkImg)

	towerImg := rl.LoadImage("tower.png")
	towerTexture := rl.LoadTextureFromImage(towerImg)
	rl.UnloadImage(towerImg)

	spawnList := [](spawn){}

	spawnList = append(spawnList, spawn_init(
		1,
		etk_init(vec_init(0.0, 0.0), 0.0, 1.0, 0.0, 0.0),
		30,
		750,
		0))

	etkList := [](etk){}

	myPlay := play{money: 100, etkMaxDmg: 100, etkCurrentDmg: 0, level: 1, levelMsStart: time.Now().UnixNano() / 1e6}

	//myEtk := etk_init(vec_init(0.0, 0.0), 0.0, 1.0, 0.0, 0.0)
	//myEtk2 := etk_init(vec_init(0.0, 0.0), 0.0, 1.0, 0.0, 0.0)
	//etkList = append(etkList, myEtk)
	//etkList = append(etkList, myEtk2)

	towerList := [](tower){}

	towerList1 := tower_init(vec_init(11.0, 8.0), 5.0, 20.0, 10, 2500)
	towerList2 := tower_init(vec_init(8.0, 11.0), 5.0, 20.0, 10, 2500)

	towerList = append(towerList, towerList1)
	towerList = append(towerList, towerList2)

	myPath := path_init([]vec2{
		vec_init(0.0, 0.0),
		vec_init(5.0, 5.0),
		vec_init(5.0, 10.0),
		vec_init(15.0, 15.0),
		vec_init(25.0, 20.0),
		vec_init(25.0, 25.0),
		vec_init(30.0, 30.0),
		//vec_init(35.0, 35.0),
		//vec_init(40.0, 40.0),
		//vec_init(45.0, 45.0),
	})

	time.Sleep(1 * time.Millisecond)

	//set delta-time to 0
	var deltaTime float64 = 0.0

	//etk loop - run for all level
	for !rl.WindowShouldClose() {
		//takes first time
		timeStart := time.Now().UnixNano() / 1e6

		//time.Sleep(10 * time.Millisecond)

		spawnList, etkList, myPlay = spawnListHandle(spawnList, etkList, myPlay)

		spawnList, etkList, towerList, myPlay = etkListHandle(spawnList, etkList, towerList, myPlay, deltaTime, myPath)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		for xDraw := 0; xDraw < 31; xDraw++ {
			for yDraw := 0; yDraw < 31; yDraw++ {
				vecOut := gameToScreenVec2(vec_init(float64(xDraw), float64(yDraw)), float64(blockTexture.Width), float64(blockTexture.Height), float64(screenWidth))
				rl.DrawTexture(blockTexture, int32(vecOut.x), int32(vecOut.y), rl.White)
			}
		}

		for i := 0; i < len(towerList); i++ {
			vecOut := gameToScreenVec2(vec_init(float64(towerList[i].position.x), float64(towerList[i].position.y)), float64(towerTexture.Width), float64(towerTexture.Height), float64(screenWidth))
			rl.DrawTexture(towerTexture, int32(vecOut.x), int32(vecOut.y), rl.White)
		}

		for i := 0; i < len(etkList); i++ {
			vecOut := gameToScreenVec2(vec_init(float64(etkList[i].position.x), float64(etkList[i].position.y)), float64(etkTexture.Width), float64(etkTexture.Height), float64(screenWidth))
			rl.DrawTexture(etkTexture, int32(vecOut.x), int32(vecOut.y), rl.White)
		}

		rl.EndDrawing()

		if (len(etkList) == 0 && len(spawnList) == 0) || myPlay.etkCurrentDmg >= myPlay.etkMaxDmg {
			fmt.Println("ende")
			fmt.Println(myPlay.etkCurrentDmg)
			if myPlay.etkCurrentDmg >= myPlay.etkMaxDmg {
				fmt.Println("weil du schlecht bist")
			}
			break
		}

		//take second time
		timeEnd := time.Now().UnixNano() / 1e6
		//set delta-time to (second-first)/(1/60)
		deltaTime = float64(timeEnd-timeStart) / 16.6
	}
	return
}

/*

ablauf:
	check for spawn -- TODO
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

	spwan loop
		find smalles orderNum
		check for level
			--> next level if not current
			else start
		check if stared
			--> then call getSpawnEtk
			--> append to etk list
		delete from list if 0

	idea for level spawn:
		struct with
		level
		type of etk
		amount
		interval
		last spawn
		--> how to spawn
			list of spawns
			iterate over list
				check for last spawn with interval
				spawn set last spawn
				reduce amount


TODO:
	restrict tower hit am [x]
	restrict etk movement --> delay to 60 ticks/second
	--> multiply over delta time [x]
	gameplay
		level-strcuture [x]
			check if no spawns left for this level && etk list empty --> level++
		spawn mechanic [x]
	graphic
		raylib
		isometric view
		--> game to screen position [x]
		--> screen to game position


		get hits --> save in list
		--> draw hit for x-ms --> remove from list
		draw path
		--> safe as map???

	sanity:
		inline main loop parts [bad idea but works]
		make stop between level:
			only call functions, etc
		etk speed variable
*/
