package main

func screenToGameVec2(posIn vec2, blockWidth float64, blockHeight float64, screenWidth float64) vec2 {
	xVec := vec_init(1.0*(blockWidth/2.0), 0.5*(blockHeight/2.0))
	yVec := vec_init(-1.0*(blockWidth/2.0), 0.5*(blockHeight/2.0))

	// Offset screen position to reverse the translation
	xOffset := posIn.x - (screenWidth / 2.0)
	yOffset := posIn.y - 20

	// Solve for posIn.x and posIn.y using the inverse transformation
	denom := xVec.x*yVec.y - xVec.y*yVec.x
	posX := (xOffset*yVec.y - yOffset*yVec.x) / denom
	posY := (yOffset*xVec.x - xOffset*xVec.y) / denom

	return vec_init(posX, posY)
}

func gameToScreenVec2(posIn vec2, blockWidth float64, blockHeight float64, screenWidth float64) vec2 {
	xVec := vec_init(1.0*(blockWidth/2.0), 0.5*(blockHeight/2.0))
	yVec := vec_init(-1.0*(blockWidth/2.0), 0.5*(blockHeight/2.0))
	x := posIn.x*xVec.x + posIn.y*yVec.x
	y := posIn.x*xVec.y + posIn.y*yVec.y
	return vec_init(x-(blockWidth/2.0)+(screenWidth/2.0), y+20)
}
