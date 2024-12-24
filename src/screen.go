package main

func screenToGameVec2(posIn vec2) vec2 {
	return posIn
}

func gameToScreenVec2(posIn vec2, blockWidth float64, blockHeight float64, screenWidth float64) vec2 {
	xVec := vec_init(1.0*(blockWidth/2.0), 0.5*(blockHeight/2.0))
	yVec := vec_init(-1.0*(blockWidth/2.0), 0.5*(blockHeight/2.0))
	x := posIn.x*xVec.x + posIn.y*yVec.x
	y := posIn.x*xVec.y + posIn.y*yVec.y
	return vec_init(x-(blockWidth/2.0)+(screenWidth/2.0), y+20)
}
