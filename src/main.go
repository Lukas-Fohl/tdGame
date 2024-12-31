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
	for y := 0; y <= int(biggest.y); y++ {
		temp := []uint8{}
		for x := 0; x <= int(biggest.x); x++ {
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
	textureList = append(textureList, textureSave_init("cursor.png"))

	spawnList := [](spawn){}

	for i := 0; i < 30; i++ {
		spawnList = append(spawnList, spawn_init(
			i+1,
			etk_init(vec_init(0.0, 0.0), 0.0, 1.0, 3.0, 0.0, 0.1, "etk.png"),
			30,
			400,
			i))
	}

	etkList := [](etk){}

	myPlay := play{money: 100, etkMaxDmg: 100, etkCurrentDmg: 0, level: 1, levelMsStart: time.Now().UnixNano() / 1e6, attackList: []attack{}}

	towerList := [](tower){}

	towerList1 := tower_init(vec_init(8.0, 8.0), 5.0, 20.0, 10, 900, "tower.png")
	towerList2 := tower_init(vec_init(11.0, 11.0), 5.0, 20.0, 10, 900, "tower.png")

	state1 := state_init(1, "", tower_init(vec_init(8.0, 8.0), 5.0, 20.0, 10, 900, "tower.png"))
	state2 := state_init(5, "", tower_init(vec_init(8.0, 8.0), 6.0, 20.0, 10, 900, "tower.png"))
	state3 := state_init(1, "", tower_init(vec_init(8.0, 8.0), 7.0, 20.0, 10, 900, "tower.png"))

	upgrade1 := upgrade_init([]state{state1, state2, state3})

	towerList1.addUpgrade(upgrade1)

	for i := 0; i < 3; i++ {
		towerList1.levelUpUpgrade(0, &myPlay)
	}

	towerList = append(towerList, towerList1)
	towerList = append(towerList, towerList2)

	myTowerInfo := towerInfo_init(false, &towerList1)

	testButton := button_init(vec_init(940.0, 580.0), vec_init(120.0, 120.0), "next level", "")

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

		myPlay.drawStats()

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

		//tempOut := gameToScreenVec2(vec_init(float64(myPath.wayPoint[len(myPath.wayPoint)-1].x), float64(myPath.wayPoint[len(myPath.wayPoint)-1].y)), float64(floorTexture.Width), float64(floorTexture.Height), float64(screenWidth))
		//rl.DrawTexture(floorTexture, int32(tempOut.x), int32(tempOut.y), rl.White)

		//for yDraw := 0; yDraw < len(myMap); yDraw++ {
		//	vecOut := gameToScreenVec2(vec_init(float64(len(myMap[0])), float64(yDraw)), float64(blockTexture.Width), float64(blockTexture.Height), float64(screenWidth))
		//	rl.DrawTexture(blockTexture, int32(vecOut.x), int32(vecOut.y), rl.White)
		//}

		//for xDraw := 0; xDraw <= len(myMap[0]); xDraw++ {
		//	vecOut := gameToScreenVec2(vec_init(float64(xDraw), float64(len(myMap))), float64(blockTexture.Width), float64(blockTexture.Height), float64(screenWidth))
		//	rl.DrawTexture(blockTexture, int32(vecOut.x), int32(vecOut.y), rl.White)
		//}

		//draw tower
		for i := 0; i < len(towerList); i++ {
			towerTexture := findTexture(textureList, towerList[i].texturePath)
			vecOut := gameToScreenVec2(towerList[i].position, float64(towerTexture.Width), float64(towerTexture.Height), float64(screenWidth))
			rl.DrawTexture(towerTexture, int32(vecOut.x), int32(vecOut.y), rl.White)
		}

		//draw cursor
		mouseGame := screenToGameVec2(vec_init(float64(rl.GetMousePosition().X), float64(rl.GetMousePosition().Y)), float64(blockTexture.Width), float64(blockTexture.Height), float64(screenWidth))
		mouseGame.x = float64(int(mouseGame.x))
		mouseGame.y = float64(int(mouseGame.y))
		if int(mouseGame.x) <= len(myMap[0]) && int(mouseGame.y) <= len(myMap) && int(mouseGame.y) >= 0 && int(mouseGame.x) >= 0 {
			vecOut := gameToScreenVec2(mouseGame, float64(blockTexture.Width), float64(blockTexture.Height), float64(screenWidth))
			rl.DrawTexture(findTexture(textureList, "cursor.png"), int32(vecOut.x), int32(vecOut.y), rl.White)
		}

		myTowerInfo.drawTowerInfo(findTexture(textureList, myTowerInfo.tower.texturePath))
		if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
			buttonClicked := false
			for i := 0; i < len(myTowerInfo.buttons); i++ {
				buttonRectangle := rl.NewRectangle(float32(myTowerInfo.buttons[i].position.x), float32(myTowerInfo.buttons[i].position.y), float32(myTowerInfo.buttons[i].size.x), float32(myTowerInfo.buttons[i].size.y))
				if rl.CheckCollisionPointRec(rl.GetMousePosition(), buttonRectangle) {
					buttonClicked = true
				}
			}
			if !buttonClicked {
				foundTower := false
				for i := 0; i < len(towerList); i++ {
					if mouseGame == towerList[i].position {
						foundTower = true
						if myTowerInfo.tower == &towerList[i] && myTowerInfo.show {
							myTowerInfo.show = false
						} else {
							myTowerInfo.show = true
							myTowerInfo.tower = &towerList[i]
						}
					}
				}
				if !foundTower {
					myTowerInfo.show = false
				}
			}
		}

		//draw etk
		for i := 0; i < len(etkList); i++ {
			etkTexture := findTexture(textureList, etkList[i].texturePath)
			vecOut := gameToScreenVec2(etkList[i].position, float64(etkTexture.Width), float64(etkTexture.Height), float64(screenWidth))
			rl.DrawTexture(etkTexture, int32(vecOut.x), int32(vecOut.y), rl.White)
		}

		//draw attack
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

		//check if spawns for this level are left
		hasSpawnWithSameLevel := false
		for j := 0; j < len(spawnList); j++ {
			if spawnList[j].amount > 0 && spawnList[j].level == myPlay.level {
				hasSpawnWithSameLevel = true
			}
		}

		if !hasSpawnWithSameLevel && len(etkList) == 0 {
			testButton.draw()
			if testButton.isClicked() {
				myPlay.level++
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
		--> screen to game position [x]

	######################IMPORTANT######################


	build button comp [x]

	options for tower --> onclick [x]
		--> on tower click
			open side pannel
			don't draw game-cursor
			check for upgrades
			--> close button

	button for next level --> gray if not needed (don't increase level) [x]

	show numbers: money, dmg/max, level [x]

	button options for texture when is clicked
	button make textures work

	bigger textures
	scaling

	build upgrades [x]
	--> handel upgrade button + click

	place tower!!!!!!

	add etk types!!!!!!
	--> add types to attack for tower

	add map format
*/
