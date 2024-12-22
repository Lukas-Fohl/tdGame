package main

import (
	"time"
)

type spawn struct {
	level      int
	etkToSpawn etk
	amount     int
	lastSpawn  int64
	interval   int64
	orderNum   int
	started    bool
}

func spawn_init(levelIn int, etkToSpawnIn etk, amountIn int, intervalIn int64, orderNumIn int) spawn {
	return spawn{
		level:      levelIn,
		etkToSpawn: etkToSpawnIn,
		amount:     amountIn,
		lastSpawn:  0,
		interval:   intervalIn,
		orderNum:   orderNumIn,
		started:    false,
	}
}

func (spawnIn *spawn) getSpawnEtks() [](etk) {
	spawnIn.started = true
	etkOutList := []etk{}
	if spawnIn.amount > 0 {
		if spawnIn.lastSpawn+spawnIn.interval <= time.Now().UnixNano()/1e6 {
			spawnIn.amount--
			etkOutList = append(etkOutList, spawnIn.etkToSpawn)
			spawnIn.lastSpawn = time.Now().UnixNano() / 1e6
		}
	}
	return etkOutList
}
