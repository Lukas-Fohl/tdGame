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
	attackList    []attack
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
		//wait for button to be pressed
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
		etkList[j].wayPointPerc += etkList[j].speed * deltaTime
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
		myPlay.attackList = append(myPlay.attackList, (&towerList[j]).dmgToETKList(etkRefList)...)
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

func getMapFromPath(pathIn path) [][]uint8 {
	biggest := vec_init(0, 0)
	for i := 0; i < len(pathIn.wayPoint); i++ {
		if pathIn.wayPoint[i].x > biggest.x {
			biggest.x = pathIn.wayPoint[i].x
		}
		if pathIn.wayPoint[i].y > biggest.y {
			biggest.y = pathIn.wayPoint[i].y
		}
	}
	outMap := [][]uint8{}
	for y := 0; y < int(biggest.y); y++ {
		temp := []uint8{}
		for x := 0; x < int(biggest.x); x++ {
			temp = append(temp, uint8(0))
		}
		outMap = append(outMap, temp)
	}
	myEtk := etk_init(vec_init(0.0, 0.0), 0.0, 0.0, 0.0, 0.0, 0.0, "")
	for i := 0; i < 1000; i++ {
		(&myEtk).poisitionFromPath(pathIn)
		myEtk.wayPointPerc += 0.1
		outMap[int(myEtk.position.y)][int(myEtk.position.x)] = 1
	}
	return outMap
}

type attack struct {
	startMS  int64
	duration int64
	start    vec2
	end      vec2
}

func main() {
	screenWidth := int32(1080)
	screenHeight := int32(720)

	rl.InitWindow(screenWidth, screenHeight, "omtfyb")
	rl.SetTargetFPS(60)
	defer rl.CloseWindow()

	textureList := []textureSave{}
	textureList = append(textureList, textureSave_init("block.png"))
	textureList = append(textureList, textureSave_init("floor.png"))
	textureList = append(textureList, textureSave_init("etk.png"))
	textureList = append(textureList, textureSave_init("tower.png"))

	spawnList := [](spawn){}

	for i := 0; i < 30; i++ {
		spawnList = append(spawnList, spawn_init(
			1,
			etk_init(vec_init(0.0, 0.0), 0.0, 1.0, 0.0, 0.0, 0.1, "etk.png"),
			30,
			400,
			i))
	}

	etkList := [](etk){}

	myPlay := play{money: 100, etkMaxDmg: 100, etkCurrentDmg: 0, level: 1, levelMsStart: time.Now().UnixNano() / 1e6, attackList: []attack{}}

	//myEtk := etk_init(vec_init(0.0, 0.0), 0.0, 1.0, 0.0, 0.0)
	//myEtk2 := etk_init(vec_init(0.0, 0.0), 0.0, 1.0, 0.0, 0.0)
	//etkList = append(etkList, myEtk)
	//etkList = append(etkList, myEtk2)

	towerList := [](tower){}

	towerList1 := tower_init(vec_init(8.0, 8.0), 5.0, 20.0, 10, 900, "tower.png")
	towerList2 := tower_init(vec_init(11.0, 11.0), 5.0, 20.0, 10, 900, "tower.png")

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

	myMap := getMapFromPath(myPath)

	//etk loop - run for all level
	for !rl.WindowShouldClose() {
		//takes first time
		timeStart := time.Now().UnixNano() / 1e6

		time.Sleep(10 * time.Millisecond)

		spawnList, etkList, myPlay = spawnListHandle(spawnList, etkList, myPlay)

		spawnList, etkList, towerList, myPlay = etkListHandle(spawnList, etkList, towerList, myPlay, deltaTime, myPath)

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		blockTexture := findTexture(textureList, "block.png")
		floorTexture := findTexture(textureList, "floor.png")

		for xDraw := 0; xDraw < len(myMap[0]); xDraw++ {
			for yDraw := 0; yDraw < len(myMap); yDraw++ {
				vecOut := gameToScreenVec2(vec_init(float64(xDraw), float64(yDraw)), float64(blockTexture.Width), float64(blockTexture.Height), float64(screenWidth))
				if myMap[yDraw][xDraw] == 0 {
					rl.DrawTexture(blockTexture, int32(vecOut.x), int32(vecOut.y), rl.White)
				} else {
					rl.DrawTexture(floorTexture, int32(vecOut.x), int32(vecOut.y), rl.White)
				}
			}
		}

		for yDraw := 0; yDraw < len(myMap); yDraw++ {
			vecOut := gameToScreenVec2(vec_init(float64(len(myMap[0])), float64(yDraw)), float64(blockTexture.Width), float64(blockTexture.Height), float64(screenWidth))
			rl.DrawTexture(blockTexture, int32(vecOut.x), int32(vecOut.y), rl.White)
		}

		for xDraw := 0; xDraw <= len(myMap[0]); xDraw++ {
			vecOut := gameToScreenVec2(vec_init(float64(xDraw), float64(len(myMap))), float64(blockTexture.Width), float64(blockTexture.Height), float64(screenWidth))
			rl.DrawTexture(blockTexture, int32(vecOut.x), int32(vecOut.y), rl.White)
		}

		tempOut := gameToScreenVec2(vec_init(float64(myPath.wayPoint[len(myPath.wayPoint)-1].x), float64(myPath.wayPoint[len(myPath.wayPoint)-1].y)), float64(floorTexture.Width), float64(floorTexture.Height), float64(screenWidth))
		rl.DrawTexture(floorTexture, int32(tempOut.x), int32(tempOut.y), rl.White)

		for i := 0; i < len(towerList); i++ {
			towerTexture := findTexture(textureList, towerList[i].texturePath)
			vecOut := gameToScreenVec2(towerList[i].position, float64(towerTexture.Width), float64(towerTexture.Height), float64(screenWidth))
			rl.DrawTexture(towerTexture, int32(vecOut.x), int32(vecOut.y), rl.White)
		}

		for i := 0; i < len(etkList); i++ {
			etkTexture := findTexture(textureList, etkList[i].texturePath)
			vecOut := gameToScreenVec2(etkList[i].position, float64(etkTexture.Width), float64(etkTexture.Height), float64(screenWidth))
			rl.DrawTexture(etkTexture, int32(vecOut.x), int32(vecOut.y), rl.White)
		}

		rl.SetLineWidth(3.5)
		attackRemoveList := []int{}
		for i := 0; i < len(myPlay.attackList); i++ {
			if myPlay.attackList[i].startMS+myPlay.attackList[i].duration >= time.Now().UnixNano()/1_000_000 {
				vecStart := gameToScreenVec2(myPlay.attackList[i].start, float64(floorTexture.Width), float64(floorTexture.Height), float64(screenWidth))
				vecEnd := gameToScreenVec2(myPlay.attackList[i].end, float64(floorTexture.Width), float64(floorTexture.Height), float64(screenWidth))
				rl.DrawLine(
					int32(vecStart.x)+(int32(floorTexture.Width)/2),
					int32(vecStart.y)+(int32(floorTexture.Height)/4),
					int32(vecEnd.x)+(int32(floorTexture.Width)/2),
					int32(vecEnd.y)+(int32(floorTexture.Height)/4),
					rl.Red)
			} else {
				attackRemoveList = append(attackRemoveList, i)
			}
		}

		for k := len(attackRemoveList) - 1; k >= 0; k-- {
			idx := attackRemoveList[k]
			if idx+1 >= len(myPlay.attackList) {
				myPlay.attackList = myPlay.attackList[:len(myPlay.attackList)-1]
			} else {
				myPlay.attackList = append(myPlay.attackList[:idx], myPlay.attackList[idx+1:]...)
			}
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
		draw path [x]
		--> safe as map??? [x]
		load all textures --> saved in list
		pop function
		mid point for texture
		--> bigger textures
		scaling

	sanity:
		inline main loop parts [bad idea but works]
		make stop between level:
			only call functions, etc
		etk speed variable, texture, texture for shot, time for shot

		save in list
		--> texture-struct
			have name of texture
			and loaded texture
			-- etk has name of texture
			--> when render check list of textures load
		path speed multiply???
*/
