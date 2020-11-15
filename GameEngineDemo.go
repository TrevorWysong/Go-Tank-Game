package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "image/png"
	"log"
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
	collision bool
}

type Game struct {
	playerSprite               Sprite
	tankTopper                 Sprite
	fireball                   Sprite
	coinSprite                 Sprite
	heartSprite1               Sprite
	heartSprite2               Sprite
	heartSprite3               Sprite
	firstMap                   Sprite
	secondMap                  Sprite
	thirdMap                   Sprite
	drawOps                    ebiten.DrawImageOptions
	collectedGold              bool
	playerAndWallCollision     bool
	projectileAndWallCollision bool
	mostRecentKeyLeft          bool
	mostRecentKeyRight         bool
	mostRecentKeyDown          bool
	mostRecentKeyUp            bool
	mostRecentKeyA             bool
	mostRecentKeyS             bool
	mostRecentKeyD             bool
	mostRecentKeyW             bool
	deathCounter               int
	projectileList             []Sprite
	projectileHold             bool
	playerFireballCounter      int
	playerFireballTimer        int
	gameOver                   bool
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

func wallCollisionCheckFirstLevel(anySprite Sprite, spriteWidth int) bool {
	boundaryWidth := 25
	if anySprite.xLoc < 0+boundaryWidth || anySprite.xLoc > ScreenWidth-boundaryWidth-spriteWidth ||
		anySprite.yLoc > ScreenHeight-boundaryWidth-spriteWidth || anySprite.yLoc < 0+boundaryWidth ||
		anySprite.xLoc > 200-spriteWidth && anySprite.xLoc < 275 && anySprite.yLoc < 250 ||
		anySprite.xLoc > 275-spriteWidth && anySprite.xLoc < 475 && anySprite.yLoc < 250 && anySprite.yLoc > 175-spriteWidth ||
		anySprite.xLoc > 175-spriteWidth && anySprite.xLoc < 275 && anySprite.yLoc < 475 && anySprite.yLoc > 400-spriteWidth ||
		anySprite.xLoc > 550-spriteWidth && anySprite.xLoc < 625 && anySprite.yLoc < 575 && anySprite.yLoc > 350-spriteWidth ||
		anySprite.xLoc > 475-spriteWidth && anySprite.xLoc < 550 && anySprite.yLoc < 575 && anySprite.yLoc > 500-spriteWidth {
		return true
	}
	return false
}

func (game *Game) shootFireball() []Sprite {
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) && game.projectileHold == false {
		game.projectileAndWallCollision = false
		game.playerFireballCounter += 1
		tempFireball := game.fireball
		tempFireball.xLoc = game.playerSprite.xLoc
		tempFireball.yLoc = game.playerSprite.yLoc
		tempFireball.dx = 5
		tempFireball.dy = 0

		game.projectileList = append(game.projectileList, tempFireball)
	}

	return game.projectileList
}

func (game *Game) changeTankDirection() {
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
}

func (game *Game) changeTankTopperDirection() {
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		game.mostRecentKeyA = false
		game.mostRecentKeyS = false
		game.mostRecentKeyD = false
		game.mostRecentKeyW = true
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		game.mostRecentKeyA = false
		game.mostRecentKeyS = true
		game.mostRecentKeyD = false
		game.mostRecentKeyW = false
	} else if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		game.mostRecentKeyA = true
		game.mostRecentKeyS = false
		game.mostRecentKeyD = false
		game.mostRecentKeyW = false
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		game.mostRecentKeyA = false
		game.mostRecentKeyS = false
		game.mostRecentKeyD = true
		game.mostRecentKeyW = false
	}
}
func (game *Game) manageTankTopperOffset() {
	if game.mostRecentKeyA == true {
		game.tankTopper.xLoc = game.playerSprite.xLoc - 7
		game.tankTopper.yLoc = game.playerSprite.yLoc + 20
	} else if game.mostRecentKeyD == true {
		game.tankTopper.xLoc = game.playerSprite.xLoc + 20
		game.tankTopper.yLoc = game.playerSprite.yLoc + 20
	} else if game.mostRecentKeyS == true {
		game.tankTopper.xLoc = game.playerSprite.xLoc + 20
		game.tankTopper.yLoc = game.playerSprite.yLoc + 20
	} else if game.mostRecentKeyW == true {
		game.tankTopper.xLoc = game.playerSprite.xLoc + 20
		game.tankTopper.yLoc = game.playerSprite.yLoc - 7
	} else {
		game.tankTopper.xLoc = game.playerSprite.xLoc + 20
		game.tankTopper.yLoc = game.playerSprite.yLoc - 7
	}
}

func (game *Game) manageLevel1CollisionDetection() {
	if game.collectedGold == false {
		game.collectedGold = gotGold(game.playerSprite, game.coinSprite)
	}

	if game.playerAndWallCollision == false {
		game.playerAndWallCollision = wallCollisionCheckFirstLevel(game.playerSprite, 61)
	} else {
		game.playerSprite.yLoc = ScreenHeight / 2
		game.playerSprite.xLoc = 74 //player width
		game.playerAndWallCollision = false
		game.deathCounter += 1
	}
	if len(game.projectileList) > 0 {
		for i := 0; i < len(game.projectileList); i++ {
			if game.projectileList[i].collision == false {
				game.projectileList[i].xLoc += game.projectileList[i].dx
				game.projectileList[i].yLoc += game.projectileList[i].dy
				game.projectileList[i].collision = wallCollisionCheckFirstLevel(game.projectileList[i], 20)
			}
		}
	}

	//if game.projectileAndWallCollision == false {
	//	game.projectileAndWallCollision = wallCollisionCheckFirstLevel(game.fireball, 20)
	//}
}

func (game *Game) Update() error {
	game.changeTankDirection()
	game.changeTankTopperDirection()
	game.shootFireball()
	game.manageTankTopperOffset()
	game.manageLevel1CollisionDetection()

	return nil
}

func (game Game) Draw(screen *ebiten.Image) {
	//screen.Fill(colornames.Chocolate)
	game.drawOps.GeoM.Reset()
	game.drawOps.GeoM.Translate(float64(game.firstMap.xLoc), float64(game.firstMap.yLoc))
	screen.DrawImage(game.firstMap.upPict, &game.drawOps)

	if len(game.projectileList) > 0 {
		for i := 0; i < len(game.projectileList); i++ {
			if game.projectileList[i].collision == false {
				game.drawOps.GeoM.Reset()
				game.drawOps.GeoM.Translate(float64(game.projectileList[i].xLoc), float64(game.projectileList[i].yLoc))
				screen.DrawImage(game.projectileList[i].upPict, &game.drawOps)
			}
		}
	}

	//if game.projectileAndWallCollision == false {
	//	game.drawOps.GeoM.Reset()
	//	game.drawOps.GeoM.Translate(float64(game.fireball.xLoc), float64(game.fireball.yLoc))
	//	screen.DrawImage(game.fireball.upPict, &game.drawOps)
	//}

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

	game.drawOps.GeoM.Reset()
	game.drawOps.GeoM.Translate(float64(game.tankTopper.xLoc), float64(game.tankTopper.yLoc))
	if game.mostRecentKeyW == true {
		screen.DrawImage(game.tankTopper.upPict, &game.drawOps)
	} else if game.mostRecentKeyS == true {
		screen.DrawImage(game.tankTopper.downPict, &game.drawOps)
	} else if game.mostRecentKeyD == true {
		screen.DrawImage(game.tankTopper.rightPict, &game.drawOps)
	} else if game.mostRecentKeyA == true {
		screen.DrawImage(game.tankTopper.leftPict, &game.drawOps)
	} else {
		screen.DrawImage(game.tankTopper.upPict, &game.drawOps)
	}

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
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Comp510 First Graphics")
	gameObject := Game{}
	loadImage(&gameObject)

	gameObject.tankTopper.xLoc = gameObject.playerSprite.xLoc
	gameObject.tankTopper.yLoc = gameObject.playerSprite.yLoc

	playerWidth, _ := gameObject.playerSprite.upPict.Size()
	gameObject.playerSprite.xLoc = playerWidth
	gameObject.playerSprite.yLoc = ScreenHeight / 2

	coinWidth, coinHeight := gameObject.coinSprite.upPict.Size()
	rand.Seed(int64(time.Now().Second()))
	gameObject.coinSprite.xLoc = rand.Intn(ScreenWidth - coinWidth)
	gameObject.coinSprite.yLoc = rand.Intn(ScreenHeight - coinHeight)

	boundaryWidth := 25
	heartWidth, heartHeight := gameObject.heartSprite1.upPict.Size()
	gameObject.heartSprite1.yLoc = ScreenHeight - (boundaryWidth * 2) - (heartHeight / 2)
	gameObject.heartSprite1.xLoc = boundaryWidth + 16
	gameObject.heartSprite2.yLoc = ScreenHeight - (boundaryWidth * 2) - (heartHeight / 2)
	gameObject.heartSprite2.xLoc = (boundaryWidth + 20) + (heartWidth)
	gameObject.heartSprite3.yLoc = ScreenHeight - (boundaryWidth * 2) - (heartHeight / 2)
	gameObject.heartSprite3.xLoc = (boundaryWidth + 24) + (heartWidth * 2)

	if err := ebiten.RunGame(&gameObject); err != nil {
		log.Fatal("Oh no! something terrible happened", err)
	}
}

func loadImage(game *Game) {
	firstMap, _, err := ebitenutil.NewImageFromFile("Level1Correct.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.firstMap.upPict = firstMap

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

	tankTopperUp, _, err := ebitenutil.NewImageFromFile("tankTopper.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	tankTopperDown, _, err := ebitenutil.NewImageFromFile("tankTopperDown.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	tankTopperLeft, _, err := ebitenutil.NewImageFromFile("tankTopperLeft.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	tankTopperRight, _, err := ebitenutil.NewImageFromFile("tankTopperRight.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.tankTopper.upPict = tankTopperUp
	game.tankTopper.downPict = tankTopperDown
	game.tankTopper.leftPict = tankTopperLeft
	game.tankTopper.rightPict = tankTopperRight

	fireball, _, err := ebitenutil.NewImageFromFile("fireball.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.fireball.upPict = fireball

	coins, _, err := ebitenutil.NewImageFromFile("gold-coins-large.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.coinSprite.upPict = coins

	heart, _, err := ebitenutil.NewImageFromFile("heartScaled.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.heartSprite1.upPict = heart
	game.heartSprite2.upPict = heart
	game.heartSprite3.upPict = heart
}
