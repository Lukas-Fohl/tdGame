package main

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (playIn play) drawStats() {
	rl.DrawText("money: ", 20, 20, 20, rl.Black)
	rl.DrawText(strconv.Itoa(playIn.money), 120, 20, 20, rl.Black)
	rl.DrawText("level: ", 20, 40, 20, rl.Black)
	rl.DrawText(strconv.Itoa(playIn.level), 120, 40, 20, rl.Black)
	rl.DrawText("damage:", 20, 60, 20, rl.Black)
	rl.DrawText(strconv.Itoa(playIn.etkCurrentDmg)+" / "+strconv.Itoa(playIn.etkMaxDmg), 120, 60, 20, rl.Black)
}

type button struct {
	position    vec2
	size        vec2
	text        string
	texturePath string
}

func button_init(positionIn vec2, sizeIn vec2, textIn string, texturePathIn string) button {
	return button{
		position:    positionIn,
		size:        sizeIn,
		text:        textIn,
		texturePath: texturePathIn,
	}
}

func (buttonIn button) draw() {
	rl.DrawRectangle(
		int32(buttonIn.position.x),
		int32(buttonIn.position.y),
		int32(buttonIn.size.x),
		int32(buttonIn.size.y),
		rl.Gray)
	fontSize := 20
	rl.DrawText(buttonIn.text, int32(buttonIn.position.x)+int32(fontSize)/2, int32(buttonIn.position.y)+(int32(buttonIn.size.y)/2)-int32(fontSize/2), int32(fontSize), rl.Black)
}

func (buttonIn button) isClicked() bool {
	mousePoint := rl.GetMousePosition()
	buttonRectangle := rl.NewRectangle(float32(buttonIn.position.x), float32(buttonIn.position.y), float32(buttonIn.size.x), float32(buttonIn.size.y))
	if rl.CheckCollisionPointRec(mousePoint, buttonRectangle) {
		if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
			return true
		}
	}
	return false
}

type towerInfo struct {
	show  bool
	tower *tower
}

func towerInfo_init(showIn bool, towerIn *tower) towerInfo {
	return towerInfo{
		show:  showIn,
		tower: towerIn,
	}
}

func (towerInfoIn towerInfo) drawTowerInfo(textureIn rl.Texture2D) {
	if towerInfoIn.show {
		backgroundRectangle := rl.NewRectangle(760, 20, 300, 500)
		rl.DrawRectangle(
			backgroundRectangle.ToInt32().X,
			backgroundRectangle.ToInt32().Y,
			backgroundRectangle.ToInt32().Width,
			backgroundRectangle.ToInt32().Height,
			rl.Gray)
		rl.DrawText("Position:", 780, 40, 20, rl.Black)
		rl.DrawText(strconv.Itoa(int(towerInfoIn.tower.position.x))+", "+strconv.Itoa(int(towerInfoIn.tower.position.y)), 940, 40, 20, rl.Black)
		rl.DrawText("Range:", 780, 60, 20, rl.Black)
		rl.DrawText(strconv.Itoa(int(towerInfoIn.tower.dmgRange)), 940, 60, 20, rl.Black)
		rl.DrawText("Dmg:", 780, 80, 20, rl.Black)
		rl.DrawText(strconv.Itoa(int(towerInfoIn.tower.dmg)), 940, 80, 20, rl.Black)
		rl.DrawText("Cooldown (ms):", 780, 100, 20, rl.Black)
		rl.DrawText(strconv.Itoa(int(towerInfoIn.tower.coolDownMax)), 940, 100, 20, rl.Black)
		rl.DrawText("Killed:", 780, 120, 20, rl.Black)
		rl.DrawText(strconv.Itoa(int(towerInfoIn.tower.killed)), 940, 120, 20, rl.Black)
		rl.DrawText("Money made:", 780, 140, 20, rl.Black)
		rl.DrawText(strconv.Itoa(int(towerInfoIn.tower.moneyMade)), 940, 140, 20, rl.Black)
		rl.DrawTexture(textureIn, 780, 170, rl.Gray)
	}
	return
}
