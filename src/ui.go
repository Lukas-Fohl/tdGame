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
	color       rl.Color
}

func button_init(positionIn vec2, sizeIn vec2, textIn string, texturePathIn string, colorIn rl.Color) button {
	return button{
		position:    positionIn,
		size:        sizeIn,
		text:        textIn,
		texturePath: texturePathIn,
		color:       colorIn,
	}
}

func (buttonIn button) draw() {
	rl.DrawRectangle(
		int32(buttonIn.position.x),
		int32(buttonIn.position.y),
		int32(buttonIn.size.x),
		int32(buttonIn.size.y),
		buttonIn.color)
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
	show    bool
	tower   *tower
	buttons []button
}

func towerInfo_init(showIn bool, towerIn *tower) towerInfo {
	return towerInfo{
		show:    showIn,
		tower:   towerIn,
		buttons: []button{},
	}
}

func (towerInfoIn *towerInfo) drawTowerInfo(textureIn rl.Texture2D, playIn *play) {
	fontSize := 20
	if towerInfoIn.show {
		backgroundRectangle := rl.NewRectangle(760, 20, 300, 500)
		rl.DrawRectangle(
			backgroundRectangle.ToInt32().X,
			backgroundRectangle.ToInt32().Y,
			backgroundRectangle.ToInt32().Width,
			backgroundRectangle.ToInt32().Height,
			rl.Gray)
		rl.DrawText("Position:", 780, 40, int32(fontSize), rl.Black)
		rl.DrawText(strconv.Itoa(int(towerInfoIn.tower.position.x))+", "+strconv.Itoa(int(towerInfoIn.tower.position.y)), 940, 40, int32(fontSize), rl.Black)
		rl.DrawText("Range:", 780, 60, int32(fontSize), rl.Black)
		rl.DrawText(strconv.Itoa(int(towerInfoIn.tower.dmgRange)), 940, 60, int32(fontSize), rl.Black)
		rl.DrawText("Dmg:", 780, 80, int32(fontSize), rl.Black)
		rl.DrawText(strconv.Itoa(int(towerInfoIn.tower.dmg)), 940, 80, int32(fontSize), rl.Black)
		rl.DrawText("Cooldown (ms):", 780, 100, int32(fontSize), rl.Black)
		rl.DrawText(strconv.Itoa(int(towerInfoIn.tower.coolDownMax)), 940, 100, int32(fontSize), rl.Black)
		rl.DrawText("Killed:", 780, 120, int32(fontSize), rl.Black)
		rl.DrawText(strconv.Itoa(int(towerInfoIn.tower.killed)), 940, 120, int32(fontSize), rl.Black)
		rl.DrawText("Money made:", 780, 140, int32(fontSize), rl.Black)
		rl.DrawText(strconv.Itoa(int(towerInfoIn.tower.moneyMade)), 940, 140, int32(fontSize), rl.Black)
		rl.DrawTexture(textureIn, 780, 170, rl.Gray)
		for i := 0; i < len(towerInfoIn.buttons); i++ {
			towerInfoIn.buttons[i].draw()
		}
		//draw range of buttons
		//check for click --> levelUpUpgrade
		fontSize = 15
		towerInfoIn.buttons = []button{}
		offset := 200
		buttonWidth := 50
		buttonDis := 70
		for i := 0; i < len(towerInfoIn.tower.upgrades); i++ {
			displayText := strconv.Itoa(towerInfoIn.tower.upgrades[i].level) + "/" + strconv.Itoa(len(towerInfoIn.tower.upgrades[i].states)-1)
			tempButton := button_init(vec_init(780.0, float64(buttonDis*i+offset)), vec_init(float64(buttonWidth), float64(buttonWidth)), "", "", rl.White)
			rl.DrawText(displayText, 850, int32(buttonDis*i+(offset+(buttonDis/2)-fontSize)), int32(fontSize), rl.Black)
			towerInfoIn.buttons = append(towerInfoIn.buttons, tempButton)
			tempButton.draw()
			if towerInfoIn.tower.upgrades[i].level < len(towerInfoIn.tower.upgrades[i].states)-1 {
				rl.DrawText(towerInfoIn.tower.upgrades[i].states[towerInfoIn.tower.upgrades[i].level+1].description, 900, int32(buttonDis*i+(offset+(buttonDis/2)-fontSize)), int32(fontSize), rl.Black)
			}
			if tempButton.isClicked() {
				towerInfoIn.tower.levelUpUpgrade(i, playIn)
			}
		}
	}
	return
}
