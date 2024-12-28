package main

import (
	"fmt"
	"math"
	"time"
)

type tower struct {
	position    vec2
	dmgRange    float64
	dmg         float64
	price       int
	coolDownMax int   //ms to next shoot
	lastAttack  int64 //ms where last shot happened
	texturePath string
}

func tower_init(positionIn vec2, dmgRangeIn float64, dmgIn float64, priceIn int, coolDownMaxIn int, texturePathIn string) tower {
	return tower{
		position:    positionIn,
		dmgRange:    dmgRangeIn,
		dmg:         dmgIn,
		price:       priceIn,
		coolDownMax: coolDownMaxIn,
		lastAttack:  (time.Now().UnixNano() / 1_000_000),
		texturePath: texturePathIn,
	}
}

func (towerIn *tower) dmgToETKList(etkList [](*etk)) []attack {
	attackList := []attack{}
	for i := 0; i < len(etkList); i++ {
		disX := math.Abs(towerIn.position.x - etkList[i].position.x)
		disY := math.Abs(towerIn.position.y - etkList[i].position.y)
		if towerIn.lastAttack+int64(towerIn.coolDownMax) <= time.Now().UnixNano()/1_000_000 {
			if math.Sqrt(disX*disX+disY*disY) <= towerIn.dmgRange && etkList[i].health > 0 {
				etkList[i].health -= towerIn.dmg
				fmt.Println("hit")
				towerIn.lastAttack = time.Now().UnixNano() / 1_000_000

				attackList = append(attackList, attack{
					startMS:  time.Now().UnixNano() / 1_000_000,
					duration: int64(towerIn.coolDownMax) / 3,
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
