package main

import (
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
	upPict    *ebiten.Image
	downPict  *ebiten.Image
	leftPict  *ebiten.Image
	rightPict *ebiten.Image
	xLoc      int
	yLoc      int
	dx        int
	dy        int
	width     float64
	height    float64
}

type Game struct {
	playerSprite           Sprite
	coinSprite             Sprite
	heartSprite1           Sprite
	heartSprite2           Sprite
	heartSprite3           Sprite
	bottomWallSprite       Sprite
	topWallSprite          Sprite
	leftWallSprite         Sprite
	rightWallSprite        Sprite
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
	goldWidth, goldHeight := gold.upPict.Size()
	playerWidth, playerHeight := player.upPict.Size()
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
	if game.mostRecentKeyUp == true {
		screen.DrawImage(game.playerSprite.upPict, &game.drawOps)
	} else if game.mostRecentKeyDown == true {
		screen.DrawImage(game.playerSprite.downPict, &game.drawOps)
	} else if game.mostRecentKeyRight == true {
		screen.DrawImage(game.playerSprite.rightPict, &game.drawOps)
	} else if game.mostRecentKeyLeft == true {
		screen.DrawImage(game.playerSprite.leftPict, &game.drawOps)
	} else {
		screen.DrawImage(game.playerSprite.upPict, &game.drawOps)
	}
	//screen.DrawImage(game.playerSprite.upPict, &game.drawOps)

	game.drawOps.GeoM.Reset()
	game.drawOps.GeoM.Translate(float64(game.bottomWallSprite.xLoc), float64(game.bottomWallSprite.yLoc))
	screen.DrawImage(game.bottomWallSprite.upPict, &game.drawOps)

	game.drawOps.GeoM.Reset()
	game.drawOps.GeoM.Translate(float64(game.topWallSprite.xLoc), float64(game.topWallSprite.yLoc))
	screen.DrawImage(game.topWallSprite.upPict, &game.drawOps)

	game.drawOps.GeoM.Reset()
	game.drawOps.GeoM.Rotate(180 * math.Pi / 360)
	game.drawOps.GeoM.Translate(float64(game.leftWallSprite.xLoc), float64(game.leftWallSprite.yLoc))
	screen.DrawImage(game.leftWallSprite.upPict, &game.drawOps)

	game.drawOps.GeoM.Reset()
	game.drawOps.GeoM.Rotate(180 * math.Pi / 360)
	game.drawOps.GeoM.Translate(float64(game.rightWallSprite.xLoc), float64(game.rightWallSprite.yLoc))
	screen.DrawImage(game.rightWallSprite.upPict, &game.drawOps)

	if game.deathCounter == 0 {
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.heartSprite1.xLoc), float64(game.heartSprite1.yLoc))
		screen.DrawImage(game.heartSprite1.upPict, &game.drawOps)
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.heartSprite2.xLoc), float64(game.heartSprite2.yLoc))
		screen.DrawImage(game.heartSprite2.upPict, &game.drawOps)
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.heartSprite3.xLoc), float64(game.heartSprite3.yLoc))
		screen.DrawImage(game.heartSprite3.upPict, &game.drawOps)
	} else if game.deathCounter == 1 {
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.heartSprite1.xLoc), float64(game.heartSprite1.yLoc))
		screen.DrawImage(game.heartSprite1.upPict, &game.drawOps)
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.heartSprite2.xLoc), float64(game.heartSprite2.yLoc))
		screen.DrawImage(game.heartSprite2.upPict, &game.drawOps)
	} else if game.deathCounter == 2 {
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.heartSprite1.xLoc), float64(game.heartSprite1.yLoc))
		screen.DrawImage(game.heartSprite1.upPict, &game.drawOps)
	} else if game.deathCounter > 2 {
		game.gameOver = true
	}

	if game.collectedGold == false {
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.coinSprite.xLoc), float64(game.coinSprite.yLoc))
		screen.DrawImage(game.coinSprite.upPict, &game.drawOps)
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
	playerWidth, _ := gameObject.playerSprite.upPict.Size()
	gameObject.playerSprite.xLoc = playerWidth
	gameObject.playerSprite.yLoc = ScreenHeight / 2

	coinWidth, coinHeight := gameObject.coinSprite.upPict.Size()
	rand.Seed(int64(time.Now().Second()))
	gameObject.coinSprite.xLoc = rand.Intn(ScreenWidth - coinWidth)
	gameObject.coinSprite.yLoc = rand.Intn(ScreenHeight - coinHeight)

	_, wallHeight := gameObject.leftWallSprite.upPict.Size()
	gameObject.topWallSprite.yLoc = 0
	gameObject.bottomWallSprite.yLoc = ScreenHeight - wallHeight
	gameObject.leftWallSprite.xLoc = wallHeight
	gameObject.rightWallSprite.xLoc = ScreenWidth

	heartWidth, heartHeight := gameObject.heartSprite1.upPict.Size()
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
	upPlayer, _, err := ebitenutil.NewImageFromFile("tankFilledTopSquare.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	downPlayer, _, err := ebitenutil.NewImageFromFile("tankFilledTopSquareDown.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	leftPlayer, _, err := ebitenutil.NewImageFromFile("tankFilledTopSquareLeft.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	rightPlayer, _, err := ebitenutil.NewImageFromFile("tankFilledTopSquareRight.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.playerSprite.upPict = upPlayer
	game.playerSprite.downPict = downPlayer
	game.playerSprite.leftPict = leftPlayer
	game.playerSprite.rightPict = rightPlayer

	coins, _, err := ebitenutil.NewImageFromFile("gold-coins-large.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.coinSprite.upPict = coins

	boundaryWall, _, err := ebitenutil.NewImageFromFile("blueLine.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.bottomWallSprite.upPict = boundaryWall
	game.topWallSprite.upPict = boundaryWall
	game.rightWallSprite.upPict = boundaryWall
	game.leftWallSprite.upPict = boundaryWall

	heart, _, err := ebitenutil.NewImageFromFile("heartScaled.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.heartSprite1.upPict = heart
	game.heartSprite2.upPict = heart
	game.heartSprite3.upPict = heart
}
