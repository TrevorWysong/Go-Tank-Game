package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 700
)

type Sprite struct {
	pict   *ebiten.Image
	xLoc   int
	yLoc   int
	dx     int
	dy     int
	width  float64
	height float64
}

type WallsSprite struct {
	pict   *ebiten.Image
	xLoc   int
	yLoc   int
	angle  int
	width  int
	height int
}

type Game struct {
	playerSprite           Sprite
	coinSprite             Sprite
	heartSprite1           Sprite
	heartSprite2           Sprite
	heartSprite3           Sprite
	bottomWallSprite       WallsSprite
	topWallSprite          WallsSprite
	leftWallSprite         WallsSprite
	rightWallSprite        WallsSprite
	drawOps                ebiten.DrawImageOptions
	collectedGold          bool
	playerAndWallCollision bool
	mostRecentKeyLeft      bool
	mostRecentKeyRight     bool
	mostRecentKeyDown      bool
	mostRecentKeyUp        bool
	deathCounter           int
	gameOver               bool
}

func gotGold(player, gold Sprite) bool {
	goldWidth, goldHeight := gold.pict.Size()
	playerWidth, playerHeight := player.pict.Size()
	if player.xLoc < gold.xLoc+goldWidth &&
		player.xLoc+playerWidth > gold.xLoc &&
		player.yLoc < gold.yLoc+goldHeight &&
		player.yLoc+playerHeight > gold.yLoc {
		return true
	}
	return false
}

func boundaryWallCollision(player Sprite) bool {
	wallHeight := 8
	// 74 x 80 (width x height)
	playerWidth, playerHeight := 74, 80
	if player.xLoc < 0+wallHeight || player.xLoc > ScreenWidth-wallHeight-playerWidth ||
		player.yLoc > ScreenHeight-wallHeight-playerHeight || player.yLoc < 0+wallHeight {
		return true
	}
	return false
}

func (game *Game) Update() error {
	fmt.Println(game.deathCounter)
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		game.playerSprite.dx = -3
		game.mostRecentKeyLeft = true
		game.mostRecentKeyDown = false
		game.mostRecentKeyRight = false
		game.mostRecentKeyUp = false
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		game.playerSprite.dx = 3
		game.mostRecentKeyLeft = false
		game.mostRecentKeyDown = false
		game.mostRecentKeyRight = true
		game.mostRecentKeyUp = false
	} else if inpututil.IsKeyJustReleased(ebiten.KeyRight) || inpututil.IsKeyJustReleased(ebiten.KeyLeft) {
		game.playerSprite.dx = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		game.playerSprite.dy = -3
		game.mostRecentKeyLeft = false
		game.mostRecentKeyDown = false
		game.mostRecentKeyRight = false
		game.mostRecentKeyUp = true
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		game.playerSprite.dy = 3
		game.mostRecentKeyLeft = false
		game.mostRecentKeyDown = true
		game.mostRecentKeyRight = false
		game.mostRecentKeyUp = false
	} else if inpututil.IsKeyJustReleased(ebiten.KeyUp) || inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		game.playerSprite.dy = 0
	}
	game.playerSprite.yLoc += game.playerSprite.dy
	game.playerSprite.xLoc += game.playerSprite.dx
	if game.collectedGold == false {
		game.collectedGold = gotGold(game.playerSprite, game.coinSprite)
	}
	if game.playerAndWallCollision == false {
		game.playerAndWallCollision = boundaryWallCollision(game.playerSprite)
	} else {
		game.playerSprite.yLoc = ScreenHeight / 2
		game.playerSprite.xLoc = 74 //player width
		game.playerAndWallCollision = false
		game.deathCounter += 1
	}

	return nil
}

func (game Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Chocolate)
	game.drawOps.GeoM.Reset()
	game.drawOps.GeoM.Translate(float64(game.playerSprite.xLoc), float64(game.playerSprite.yLoc))
	//game.drawOps.GeoM.Scale(0.10, 0.10)
	if game.mostRecentKeyUp == true {

	} else if game.mostRecentKeyDown == true {

	} else if game.mostRecentKeyRight == true {
		//game.drawOps.GeoM.Rotate(225)
	} else if game.mostRecentKeyLeft == true {

	} else {

	}
	//game.drawOps.GeoM.Rotate(180 * math.Pi / 360)
	screen.DrawImage(game.playerSprite.pict, &game.drawOps)

	game.drawOps.GeoM.Reset()
	game.drawOps.GeoM.Translate(float64(game.bottomWallSprite.xLoc), float64(game.bottomWallSprite.yLoc))
	screen.DrawImage(game.bottomWallSprite.pict, &game.drawOps)

	game.drawOps.GeoM.Reset()
	game.drawOps.GeoM.Translate(float64(game.topWallSprite.xLoc), float64(game.topWallSprite.yLoc))
	screen.DrawImage(game.topWallSprite.pict, &game.drawOps)

	game.drawOps.GeoM.Reset()
	game.drawOps.GeoM.Rotate(180 * math.Pi / 360)
	game.drawOps.GeoM.Translate(float64(game.leftWallSprite.xLoc), float64(game.leftWallSprite.yLoc))
	screen.DrawImage(game.leftWallSprite.pict, &game.drawOps)

	game.drawOps.GeoM.Reset()
	game.drawOps.GeoM.Rotate(180 * math.Pi / 360)
	game.drawOps.GeoM.Translate(float64(game.rightWallSprite.xLoc), float64(game.rightWallSprite.yLoc))
	screen.DrawImage(game.rightWallSprite.pict, &game.drawOps)
	if game.deathCounter == 0 {
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.heartSprite1.xLoc), float64(game.heartSprite1.yLoc))
		screen.DrawImage(game.heartSprite1.pict, &game.drawOps)
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.heartSprite2.xLoc), float64(game.heartSprite2.yLoc))
		screen.DrawImage(game.heartSprite2.pict, &game.drawOps)
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.heartSprite3.xLoc), float64(game.heartSprite3.yLoc))
		screen.DrawImage(game.heartSprite3.pict, &game.drawOps)
	} else if game.deathCounter == 1 {
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.heartSprite1.xLoc), float64(game.heartSprite1.yLoc))
		screen.DrawImage(game.heartSprite1.pict, &game.drawOps)
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.heartSprite2.xLoc), float64(game.heartSprite2.yLoc))
		screen.DrawImage(game.heartSprite2.pict, &game.drawOps)
	} else if game.deathCounter == 2 {
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.heartSprite1.xLoc), float64(game.heartSprite1.yLoc))
		screen.DrawImage(game.heartSprite1.pict, &game.drawOps)
	} else if game.deathCounter > 2 {
		game.gameOver = true
	}

	if game.collectedGold == false {
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.coinSprite.xLoc), float64(game.coinSprite.yLoc))
		screen.DrawImage(game.coinSprite.pict, &game.drawOps)
	}
	//if game.playerAndWallCollision == false{
	//	game.drawOps.GeoM.Reset()
	//	game.drawOps.GeoM.Translate(float64(game.coinSprite.xLoc), float64(game.coinSprite.yLoc))
	//	screen.DrawImage(game.coinSprite.pict, &game.drawOps)
	//}
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Comp510 First Graphics")
	gameObject := Game{}
	loadImage(&gameObject)
	playerWidth, _ := gameObject.playerSprite.pict.Size()
	gameObject.playerSprite.xLoc = playerWidth
	gameObject.playerSprite.yLoc = ScreenHeight / 2

	coinWidth, coinHeight := gameObject.coinSprite.pict.Size()
	rand.Seed(int64(time.Now().Second()))
	gameObject.coinSprite.xLoc = rand.Intn(ScreenWidth - coinWidth)
	gameObject.coinSprite.yLoc = rand.Intn(ScreenHeight - coinHeight)

	_, wallHeight := gameObject.leftWallSprite.pict.Size()
	gameObject.topWallSprite.yLoc = 0
	gameObject.bottomWallSprite.yLoc = ScreenHeight - wallHeight
	gameObject.leftWallSprite.xLoc = wallHeight
	gameObject.rightWallSprite.xLoc = ScreenWidth

	heartWidth, heartHeight := gameObject.heartSprite1.pict.Size()
	gameObject.heartSprite1.yLoc = ScreenHeight - (wallHeight * 2) - heartHeight
	gameObject.heartSprite1.xLoc = wallHeight * 2
	gameObject.heartSprite2.yLoc = ScreenHeight - (wallHeight * 2) - heartHeight
	gameObject.heartSprite2.xLoc = (wallHeight * 3) + (heartWidth)
	gameObject.heartSprite3.yLoc = ScreenHeight - (wallHeight * 2) - heartHeight
	gameObject.heartSprite3.xLoc = (wallHeight * 4) + (heartWidth * 2)

	//gameObject.wallsSprite.xLoc = ScreenWidth/2

	if err := ebiten.RunGame(&gameObject); err != nil {
		log.Fatal("Oh no! something terrible happened", err)
	}
}

func loadImage(game *Game) {
	pict, _, err := ebitenutil.NewImageFromFile("tankfilledtop2.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.playerSprite.pict = pict

	coins, _, err := ebitenutil.NewImageFromFile("gold-coins-large.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.coinSprite.pict = coins

	boundaryWall, _, err := ebitenutil.NewImageFromFile("blueLine.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.bottomWallSprite.pict = boundaryWall
	game.topWallSprite.pict = boundaryWall
	game.rightWallSprite.pict = boundaryWall
	game.leftWallSprite.pict = boundaryWall

	heart, _, err := ebitenutil.NewImageFromFile("heartScaled.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.heartSprite1.pict = heart
	game.heartSprite2.pict = heart
	game.heartSprite3.pict = heart
}
