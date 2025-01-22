package main

import (
	"math"
	"reflect"
	"time"
)

type tower struct {
	position      vec2
	dmgRange      float64
	dmg           float64
	price         int
	coolDownMax   int   //ms to next shoot
	lastAttack    int64 //ms where last shot happened
	texturePath   string
	killed        int
	moneyMade     int
	upgrades      []upgrade
	typesToAttack []etkType
}

func tower_init(positionIn vec2, dmgRangeIn float64, dmgIn float64, priceIn int, coolDownMaxIn int, texturePathIn string, typesToAttackIn []etkType) tower {
	return tower{
		position:      positionIn,
		dmgRange:      dmgRangeIn,
		dmg:           dmgIn,
		price:         priceIn,
		coolDownMax:   coolDownMaxIn,
		lastAttack:    (time.Now().UnixNano() / 1_000_000),
		texturePath:   texturePathIn,
		killed:        0,
		moneyMade:     0,
		typesToAttack: typesToAttackIn,
	}
}

type towerPreset struct {
	tower       tower
	name        string
	description string
}

func towerPreset_init(towerIn tower, nameIn string, descriptionIn string) towerPreset {
	return towerPreset{
		tower:       towerIn,
		name:        nameIn,
		description: descriptionIn,
	}
}

func canAttack(towerTypeList []etkType, etkTypeList []etkType) bool {
	canAttack := true
	for i := 0; i < len(etkTypeList); i++ {
		isIn := false
		for j := 0; j < len(towerTypeList); j++ {
			if etkTypeList[i] == towerTypeList[j] {
				isIn = true
			}
		}
		if !isIn {
			canAttack = false
			break
		}
	}
	return canAttack
}

func (towerIn *tower) dmgToETKList(etkList [](*etk)) []attack {
	attackList := []attack{}
	for i := 0; i < len(etkList); i++ {
		disX := math.Abs(towerIn.position.x - etkList[i].position.x)
		disY := math.Abs(towerIn.position.y - etkList[i].position.y)
		if towerIn.lastAttack+int64(towerIn.coolDownMax) <= time.Now().UnixNano()/1_000_000 {
			if math.Sqrt(disX*disX+disY*disY) <= towerIn.dmgRange && etkList[i].health > 0 && canAttack(towerIn.typesToAttack, etkList[i].selfType) {
				etkList[i].health -= towerIn.dmg
				if etkList[i].health <= 0 {
					towerIn.killed++
					towerIn.moneyMade += etkList[i].reward
				}
				//fmt.Println("hit")
				towerIn.lastAttack = time.Now().UnixNano() / 1_000_000

				attackList = append(attackList, attack{
					startMS:  time.Now().UnixNano() / 1_000_000,
					duration: int64(towerIn.coolDownMax) / 4,
					start:    towerIn.position,
					end:      etkList[i].position,
				})
			} else {
				//cannot reach
			}
		} else {
			//have to wait

			//fmt.Println("#####")
			//fmt.Println(towerIn.lastAttack + int64(towerIn.coolDownMax))
			//fmt.Println(time.Now().UnixNano() / 1_000_000)
			//fmt.Println("#####")
		}
	}
	return attackList
}

type upgrade struct {
	level  int
	states []state
}

func upgrade_init(statesIn []state) upgrade {
	return upgrade{
		level:  0,
		states: statesIn,
	}
}

func (towerIn *tower) addUpgrade(upgradeIn upgrade) {
	towerIn.upgrades = append(towerIn.upgrades, upgradeIn)
	towerIn.applyUpgrade(upgradeIn)
}

func (towerIn *tower) applyUpgrade(upgradeIn upgrade) {
	if upgradeIn.states[upgradeIn.level].tower.dmgRange != 0 {
		towerIn.dmgRange = upgradeIn.states[upgradeIn.level].tower.dmgRange
	}
	if upgradeIn.states[upgradeIn.level].tower.dmg != 0 {
		towerIn.dmg = upgradeIn.states[upgradeIn.level].tower.dmg
	}
	if upgradeIn.states[upgradeIn.level].tower.coolDownMax != 0 {
		towerIn.coolDownMax = upgradeIn.states[upgradeIn.level].tower.coolDownMax
	}
	if upgradeIn.states[upgradeIn.level].tower.texturePath != "" {
		towerIn.texturePath = upgradeIn.states[upgradeIn.level].tower.texturePath
	}
	if !reflect.DeepEqual(upgradeIn.states[upgradeIn.level].tower.typesToAttack, []state{}) {
		towerIn.typesToAttack = upgradeIn.states[upgradeIn.level].tower.typesToAttack
	}
}

func (towerIn *tower) levelUpUpgrade(index int, playIn *play) {
	if towerIn.upgrades[index].level < len(towerIn.upgrades[index].states)-1 {
		if towerIn.upgrades[index].states[towerIn.upgrades[index].level+1].statePrice <= playIn.money {
			playIn.money -= towerIn.upgrades[index].states[towerIn.upgrades[index].level+1].statePrice
			towerIn.upgrades[index].level++
		}
	}
	towerIn.applyUpgrade(towerIn.upgrades[index])
}

type state struct {
	statePrice  int
	tower       tower
	description string
}

func state_init(priceIn int, descriptionIn string, towerIn tower) state {
	return state{
		statePrice:  priceIn,
		tower:       towerIn,
		description: descriptionIn,
	}
}
