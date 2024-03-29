package main

import (
	"bytes"
	"database/sql"
	_ "fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	raudio "github.com/hajimehoshi/ebiten/v2/examples/resources/audio"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	_ "golang.org/x/image/font/opentype"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 700
	sampleRate   = 44100
)

var (
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func init() {
	var err error
	// Initialize audio context.
	g.audioContext = audio.NewContext(sampleRate)

	playerDeathSound, err := os.Open("sounds/death.wav")
	if err != nil {
		log.Fatal(err)
	}
	playerDeathSoundDecoded, err := wav.Decode(g.audioContext, playerDeathSound)

	// Create an audio.Player that has one stream.
	g.playerDeathAudioPlayer, err = audio.NewPlayer(g.audioContext, playerDeathSoundDecoded)
	if err != nil {
		log.Fatal(err)
	}

	winSound, err := os.Open("sounds/win.wav")
	if err != nil {
		log.Fatal(err)
	}
	winSoundDecoded, err := wav.Decode(g.audioContext, winSound)

	// Create an audio.Player that has one stream.
	g.winAudioPlayer, err = audio.NewPlayer(g.audioContext, winSoundDecoded)
	if err != nil {
		log.Fatal(err)
	}

	pickedUpBonusSound, err := os.Open("sounds/spawn.wav")
	if err != nil {
		log.Fatal(err)
	}
	pickedUpBonusSoundDecoded, err := wav.Decode(g.audioContext, pickedUpBonusSound)

	// Create an audio.Player that has one stream.
	g.pickedUpBonusAudioPlayer, err = audio.NewPlayer(g.audioContext, pickedUpBonusSoundDecoded)
	if err != nil {
		log.Fatal(err)
	}

	extraLifeSound, err := os.Open("sounds/health.wav")
	if err != nil {
		log.Fatal(err)
	}
	extraLifeSoundDecoded, err := wav.Decode(g.audioContext, extraLifeSound)

	// Create an audio.Player that has one stream.
	g.extraLifeAudioPlayer, err = audio.NewPlayer(g.audioContext, extraLifeSoundDecoded)
	if err != nil {
		log.Fatal(err)
	}

	loseSound, err := os.Open("sounds/game over.wav")
	if err != nil {
		log.Fatal(err)
	}
	loseSoundDecoded, err := wav.Decode(g.audioContext, loseSound)

	// Create an audio.Player that has one stream.
	g.loseAudioPlayer, err = audio.NewPlayer(g.audioContext, loseSoundDecoded)
	if err != nil {
		log.Fatal(err)
	}

	humanEnemyDeathSound, err := os.Open("sounds/human groan.wav")
	if err != nil {
		log.Fatal(err)
	}
	humanEnemyDeathSoundDecoded, err := wav.Decode(g.audioContext, humanEnemyDeathSound)

	// Create an audio.Player that has one stream.
	g.humanEnemyDeathAudioPlayer, err = audio.NewPlayer(g.audioContext, humanEnemyDeathSoundDecoded)
	if err != nil {
		log.Fatal(err)
	}

	monsterEnemyDeathSound, err := os.Open("sounds/monster death.wav")
	if err != nil {
		log.Fatal(err)
	}
	monsterEnemyDeathSoundDecoded, err := wav.Decode(g.audioContext, monsterEnemyDeathSound)

	// Create an audio.Player that has one stream.
	g.monsterEnemyDeathAudioPlayer, err = audio.NewPlayer(g.audioContext, monsterEnemyDeathSoundDecoded)
	if err != nil {
		log.Fatal(err)
	}

	monsterEnemyDamagedSound, err := os.Open("sounds/first monster hit.wav")
	if err != nil {
		log.Fatal(err)
	}
	monsterEnemyDamagedSoundDecoded, err := wav.Decode(g.audioContext, monsterEnemyDamagedSound)

	// Create an audio.Player that has one stream.
	g.monsterEnemyDamagedAudioPlayer, err = audio.NewPlayer(g.audioContext, monsterEnemyDamagedSoundDecoded)
	if err != nil {
		log.Fatal(err)
	}

	playerProjectileSound, err := os.Open("sounds/projectile.wav")
	if err != nil {
		log.Fatal(err)
	}
	playerProjectileSoundDecoded, err := wav.Decode(g.audioContext, playerProjectileSound)

	// Create an audio.Player that has one stream.
	g.playerShootsProjectileAudioPlayer, err = audio.NewPlayer(g.audioContext, playerProjectileSoundDecoded)
	if err != nil {
		log.Fatal(err)
	}

	enemyProjectileSound, err := os.Open("sounds/enemy projectile.wav")
	if err != nil {
		log.Fatal(err)
	}
	enemyProjectileSoundDecoded, err := wav.Decode(g.audioContext, enemyProjectileSound)

	// Create an audio.Player that has one stream.
	g.enemyShootsProjectileAudioPlayer, err = audio.NewPlayer(g.audioContext, enemyProjectileSoundDecoded)
	if err != nil {
		log.Fatal(err)
	}

	// Decode wav-formatted data and retrieve decoded PCM stream.
	d, err := wav.Decode(g.audioContext, bytes.NewReader(raudio.Jab_wav))
	if err != nil {
		log.Fatal(err)
	}

	// Create an audio.Player that has one stream.
	g.enemyAndPlayerCollisionAudioPlayer, err = audio.NewPlayer(g.audioContext, d)
	if err != nil {
		log.Fatal(err)
	}
}

type Sprite struct {
	upPict              *ebiten.Image
	downPict            *ebiten.Image
	leftPict            *ebiten.Image
	rightPict           *ebiten.Image
	xLoc                int
	yLoc                int
	dx                  int
	dy                  int
	width               float64
	height              float64
	collision           bool
	direction           string
	health              int
	inPlayerProximity   bool
	projectileHold      bool
	enemyProjectileList []Sprite
}

type Game struct {
	playerSprite                       Sprite
	personEnemy                        Sprite
	monsterEnemy                       Sprite
	tankTopper                         Sprite
	fireball                           Sprite
	coinSprite                         Sprite
	heartSprite1                       Sprite
	heartSprite2                       Sprite
	heartSprite3                       Sprite
	titleScreenBackground              Sprite
	firstMap                           Sprite
	secondMap                          Sprite
	thirdMap                           Sprite
	winnerScreen                       Sprite
	loserScreen                        Sprite
	drawOps                            ebiten.DrawImageOptions
	collectedGold                      bool
	playerAndWallCollision             bool
	projectileAndWallCollision         bool
	mostRecentKeyLeft                  bool
	mostRecentKeyRight                 bool
	mostRecentKeyDown                  bool
	mostRecentKeyUp                    bool
	mostRecentKeyA                     bool
	mostRecentKeyS                     bool
	mostRecentKeyD                     bool
	mostRecentKeyW                     bool
	deathCounter                       int
	score                              int
	projectileList                     []Sprite
	levelOneEnemyList                  []Sprite
	levelTwoEnemyList                  []Sprite
	levelThreeEnemyList                []Sprite
	levelOneIsActive                   bool
	levelTwoIsActive                   bool
	levelThreeIsActive                 bool
	spawnedLevel1Enemies               bool
	spawnedLevel2Enemies               bool
	spawnedLevel3Enemies               bool
	gameOver                           bool
	gameWon                            bool
	userNameList                       []string
	userName                           string
	startGame                          bool
	dbEntryComplete                    bool
	processedDB                        bool
	audioContext                       *audio.Context
	playerDeathAudioPlayer             *audio.Player
	humanEnemyDeathAudioPlayer         *audio.Player
	monsterEnemyDeathAudioPlayer       *audio.Player
	monsterEnemyDamagedAudioPlayer     *audio.Player
	winAudioPlayer                     *audio.Player
	loseAudioPlayer                    *audio.Player
	playerShootsProjectileAudioPlayer  *audio.Player
	enemyShootsProjectileAudioPlayer   *audio.Player
	enemyAndPlayerCollisionAudioPlayer *audio.Player
	pickedUpBonusAudioPlayer           *audio.Player
	extraLifeAudioPlayer               *audio.Player
	playedWinSound                     bool
	playedLoseSound                    bool
	extraLifeAwarded                   bool
	allScores                          bool
	playerScores                       bool
	currentPlayerAndScoreLeaderboard   bool
	playerRespawnInvincibility         bool
}

var g Game
var userNameMap = make(map[int][]string)
var scoreMap = make(map[int][]int)
var currentPlayerMap = make(map[int][]string)
var currentPlayerScoreMap = make(map[int][]int)
var dbUserNameList []string
var dbUserNameListSorted []string
var dbScoreList []int
var dbScoreListSorted []int

func (game *Game) gotGold(player, gold Sprite) bool {
	goldWidth, goldHeight := gold.upPict.Size()
	playerWidth, playerHeight := player.upPict.Size()
	if player.xLoc < gold.xLoc+goldWidth &&
		player.xLoc+playerWidth > gold.xLoc &&
		player.yLoc < gold.yLoc+goldHeight &&
		player.yLoc+playerHeight > gold.yLoc {
		if game.score%100 == 0 {
			g.pickedUpBonusAudioPlayer.Rewind()
			g.pickedUpBonusAudioPlayer.Play()
		} else if game.score%100 != 0 {
			g.extraLifeAudioPlayer.Rewind()
			g.extraLifeAudioPlayer.Play()
			if game.deathCounter > 0 {
				game.deathCounter -= 1
				game.extraLifeAwarded = true
			}
		}
		game.score += 25
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

func (game *Game) iterateAndStoreUserName() {
	for i := 0; i < len(game.userNameList); i++ {
		game.userName += game.userNameList[i]
	}
}

func (game *Game) getLeaderBoardFormat() {
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		if game.allScores == true && game.playerScores == false {
			game.allScores = false
			game.playerScores = true
		} else if game.allScores == false && game.playerScores == true {
			game.allScores = true
			game.playerScores = false
		} else {
			game.allScores = true
			game.playerScores = false
		}
		game.currentPlayerAndScoreLeaderboard = false
	}
}

func (game *Game) getUserName() {
	if inpututil.IsKeyJustReleased(ebiten.KeyBackspace) && len(game.userNameList) > 0 {
		game.userNameList = game.userNameList[:len(game.userNameList)-1]
	} else if inpututil.IsKeyJustReleased(ebiten.KeyA) {
		game.userNameList = append(game.userNameList, "A")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyB) {
		game.userNameList = append(game.userNameList, "B")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyC) {
		game.userNameList = append(game.userNameList, "C")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyD) {
		game.userNameList = append(game.userNameList, "D")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyE) {
		game.userNameList = append(game.userNameList, "E")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyF) {
		game.userNameList = append(game.userNameList, "F")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyG) {
		game.userNameList = append(game.userNameList, "G")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyH) {
		game.userNameList = append(game.userNameList, "H")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyI) {
		game.userNameList = append(game.userNameList, "I")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyJ) {
		game.userNameList = append(game.userNameList, "J")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyK) {
		game.userNameList = append(game.userNameList, "K")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyL) {
		game.userNameList = append(game.userNameList, "L")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyM) {
		game.userNameList = append(game.userNameList, "M")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyN) {
		game.userNameList = append(game.userNameList, "N")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyO) {
		game.userNameList = append(game.userNameList, "O")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyP) {
		game.userNameList = append(game.userNameList, "P")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyQ) {
		game.userNameList = append(game.userNameList, "Q")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		game.userNameList = append(game.userNameList, "R")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyS) {
		game.userNameList = append(game.userNameList, "S")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyT) {
		game.userNameList = append(game.userNameList, "T")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyU) {
		game.userNameList = append(game.userNameList, "U")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyV) {
		game.userNameList = append(game.userNameList, "V")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyW) {
		game.userNameList = append(game.userNameList, "W")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyX) {
		game.userNameList = append(game.userNameList, "X")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyY) {
		game.userNameList = append(game.userNameList, "Y")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyZ) {
		game.userNameList = append(game.userNameList, "Z")
	} else if inpututil.IsKeyJustReleased(ebiten.Key0) {
		game.userNameList = append(game.userNameList, "0")
	} else if inpututil.IsKeyJustReleased(ebiten.Key1) {
		game.userNameList = append(game.userNameList, "1")
	} else if inpututil.IsKeyJustReleased(ebiten.Key2) {
		game.userNameList = append(game.userNameList, "2")
	} else if inpututil.IsKeyJustReleased(ebiten.Key3) {
		game.userNameList = append(game.userNameList, "3")
	} else if inpututil.IsKeyJustReleased(ebiten.Key4) {
		game.userNameList = append(game.userNameList, "4")
	} else if inpututil.IsKeyJustReleased(ebiten.Key5) {
		game.userNameList = append(game.userNameList, "5")
	} else if inpututil.IsKeyJustReleased(ebiten.Key6) {
		game.userNameList = append(game.userNameList, "6")
	} else if inpututil.IsKeyJustReleased(ebiten.Key7) {
		game.userNameList = append(game.userNameList, "7")
	} else if inpututil.IsKeyJustReleased(ebiten.Key8) {
		game.userNameList = append(game.userNameList, "8")
	} else if inpututil.IsKeyJustReleased(ebiten.Key9) {
		game.userNameList = append(game.userNameList, "9")
	} else if inpututil.IsKeyJustReleased(ebiten.KeyEnter) == true && len(game.userNameList) > 0 {
		for i := 0; i < len(game.userNameList); i++ {
			game.userName += game.userNameList[i]
		}
		game.startGame = true
	}
}

func wallCollisionCheckSecondLevel(anySprite Sprite, spriteWidth int) bool {
	boundaryWidth := 25
	if anySprite.xLoc < 0+boundaryWidth || anySprite.xLoc > ScreenWidth-boundaryWidth-spriteWidth ||
		anySprite.yLoc > ScreenHeight-boundaryWidth-spriteWidth || anySprite.yLoc < 0+boundaryWidth ||
		anySprite.xLoc > 200-spriteWidth && anySprite.xLoc < 325 && anySprite.yLoc < 525 ||
		anySprite.xLoc > 500-spriteWidth && anySprite.xLoc < 600 && anySprite.yLoc > 200-spriteWidth {
		return true
	}
	return false
}

func wallCollisionCheckThirdLevel(anySprite Sprite, spriteWidth int) bool {
	boundaryWidth := 25
	if anySprite.xLoc < 0+boundaryWidth || anySprite.xLoc > ScreenWidth-boundaryWidth-spriteWidth ||
		anySprite.yLoc > ScreenHeight-boundaryWidth-spriteWidth || anySprite.yLoc < 0+boundaryWidth ||
		(anySprite.xLoc > 200-spriteWidth && anySprite.yLoc > 175-spriteWidth && anySprite.yLoc < 275) ||
		(anySprite.xLoc > 200-spriteWidth && anySprite.xLoc < 325 && anySprite.yLoc > 275-spriteWidth && anySprite.yLoc < 425) ||
		(anySprite.xLoc > 200-spriteWidth && anySprite.xLoc < 625 && anySprite.yLoc > 425-spriteWidth && anySprite.yLoc < 525) {
		return true
	}
	return false
}

func projectileCollisionWithEnemy(anyEnemy Sprite, anyProjectileSprite Sprite, enemyWidth int, projectileWidth int) (bool, bool, int, int) {
	if (anyProjectileSprite.xLoc < anyEnemy.xLoc+enemyWidth &&
		anyProjectileSprite.xLoc+projectileWidth > anyEnemy.xLoc &&
		anyProjectileSprite.yLoc < anyEnemy.yLoc+enemyWidth &&
		anyProjectileSprite.yLoc+projectileWidth > anyEnemy.yLoc) && (anyEnemy.health == 1) {
		if enemyWidth == 50 {
			g.monsterEnemyDeathAudioPlayer.Rewind()
			g.monsterEnemyDeathAudioPlayer.Play()
		} else {
			g.humanEnemyDeathAudioPlayer.Rewind()
			g.humanEnemyDeathAudioPlayer.Play()
		}
		anyEnemy.health -= 1
		additionalScore := 200
		return true, true, anyEnemy.health, additionalScore
	} else if (anyProjectileSprite.xLoc < anyEnemy.xLoc+enemyWidth &&
		anyProjectileSprite.xLoc+projectileWidth > anyEnemy.xLoc &&
		anyProjectileSprite.yLoc < anyEnemy.yLoc+enemyWidth &&
		anyProjectileSprite.yLoc+projectileWidth > anyEnemy.yLoc) && (anyEnemy.health == 2) {
		g.monsterEnemyDamagedAudioPlayer.Rewind()
		g.monsterEnemyDamagedAudioPlayer.Play()
		anyEnemy.health -= 1
		additionalScore := 100
		return false, true, anyEnemy.health, additionalScore
	}
	additionalScore := 0
	return false, false, anyEnemy.health, additionalScore
}

func projectileCollisionWithPlayer(player Sprite, anyProjectileSprite Sprite, playerWidth int, projectileWidth int) (bool, int) {
	if anyProjectileSprite.xLoc < player.xLoc+playerWidth &&
		anyProjectileSprite.xLoc+projectileWidth > player.xLoc &&
		anyProjectileSprite.yLoc < player.yLoc+playerWidth &&
		anyProjectileSprite.yLoc+projectileWidth > player.yLoc {
		return true, 1
	}
	return false, 0
}

func playerCollisionWithEnemy(anyEnemy Sprite, player Sprite, enemyWidth int, playerWidth int) int {
	if player.xLoc < anyEnemy.xLoc+enemyWidth &&
		player.xLoc+playerWidth > anyEnemy.xLoc &&
		player.yLoc < anyEnemy.yLoc+enemyWidth &&
		player.yLoc+playerWidth > anyEnemy.yLoc {
		death := 1
		return death
	}
	death := 0
	return death
}

func (game *Game) playerShootFireball() []Sprite {
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) && game.playerSprite.projectileHold == false {
		g.playerShootsProjectileAudioPlayer.Rewind()
		g.playerShootsProjectileAudioPlayer.Play()
		game.playerSprite.projectileHold = true
		go func() {
			<-time.After(500 * time.Millisecond)
			game.playerSprite.projectileHold = false
		}()
		game.projectileAndWallCollision = false
		tempFireball := game.fireball

		if game.mostRecentKeyW == true {
			tempFireball.xLoc = game.playerSprite.xLoc + 20
			tempFireball.yLoc = game.playerSprite.yLoc - 18
			tempFireball.dx = 0
			tempFireball.dy = -10
			game.projectileList = append(game.projectileList, tempFireball)
		} else if game.mostRecentKeyS == true {
			tempFireball.xLoc = game.playerSprite.xLoc + 20
			tempFireball.yLoc = game.playerSprite.yLoc + 55
			tempFireball.dx = 0
			tempFireball.dy = 10
			game.projectileList = append(game.projectileList, tempFireball)
		} else if game.mostRecentKeyA == true {
			tempFireball.xLoc = game.playerSprite.xLoc - 15
			tempFireball.yLoc = game.playerSprite.yLoc + 18
			tempFireball.dx = -10
			tempFireball.dy = 0
			game.projectileList = append(game.projectileList, tempFireball)
		} else if game.mostRecentKeyD == true {
			tempFireball.xLoc = game.playerSprite.xLoc + 55
			tempFireball.yLoc = game.playerSprite.yLoc + 18
			tempFireball.dx = 10
			tempFireball.dy = 0
			game.projectileList = append(game.projectileList, tempFireball)
		} else {
			tempFireball.xLoc = game.playerSprite.xLoc + 20
			tempFireball.yLoc = game.playerSprite.yLoc - 18
			tempFireball.dx = 0
			tempFireball.dy = -10
			game.projectileList = append(game.projectileList, tempFireball)
		}
	}
	return game.projectileList
}

func (game *Game) enemyShootFireball(i int) []Sprite {
	if game.levelOneIsActive {
		if game.levelOneEnemyList[i].projectileHold == false && game.levelOneEnemyList[i].collision == false {
			game.levelOneEnemyList[i].projectileHold = true
			g.enemyShootsProjectileAudioPlayer.Rewind()
			g.enemyShootsProjectileAudioPlayer.Play()

			go func() {
				<-time.After(3000 * time.Millisecond)
				game.levelOneEnemyList[i].projectileHold = false
			}()
			game.projectileAndWallCollision = false

			tempFireball := game.fireball

			if game.levelOneEnemyList[i].direction == "up" {
				tempFireball.xLoc = game.levelOneEnemyList[i].xLoc + 20
				tempFireball.yLoc = game.levelOneEnemyList[i].yLoc - 18
				tempFireball.dx = 0
				tempFireball.dy = -3
				game.levelOneEnemyList[i].enemyProjectileList = append(game.levelOneEnemyList[i].enemyProjectileList, tempFireball)
			} else if game.levelOneEnemyList[i].direction == "down" {
				tempFireball.xLoc = game.levelOneEnemyList[i].xLoc + 20
				tempFireball.yLoc = game.levelOneEnemyList[i].yLoc + 55
				tempFireball.dx = 0
				tempFireball.dy = 3
				game.levelOneEnemyList[i].enemyProjectileList = append(game.levelOneEnemyList[i].enemyProjectileList, tempFireball)
			} else if game.levelOneEnemyList[i].direction == "left" {
				tempFireball.xLoc = game.levelOneEnemyList[i].xLoc - 15
				tempFireball.yLoc = game.levelOneEnemyList[i].yLoc + 18
				tempFireball.dx = -3
				tempFireball.dy = 0
				game.levelOneEnemyList[i].enemyProjectileList = append(game.levelOneEnemyList[i].enemyProjectileList, tempFireball)
			} else if game.levelOneEnemyList[i].direction == "right" {
				tempFireball.xLoc = game.levelOneEnemyList[i].xLoc + 55
				tempFireball.yLoc = game.levelOneEnemyList[i].yLoc + 18
				tempFireball.dx = 3
				tempFireball.dy = 0
				game.levelOneEnemyList[i].enemyProjectileList = append(game.levelOneEnemyList[i].enemyProjectileList, tempFireball)
			} else {
				tempFireball.xLoc = game.levelOneEnemyList[i].xLoc + 20
				tempFireball.yLoc = game.levelOneEnemyList[i].yLoc - 18
				tempFireball.dx = 0
				tempFireball.dy = -3
				game.levelOneEnemyList[i].enemyProjectileList = append(game.levelOneEnemyList[i].enemyProjectileList, tempFireball)
			}
		}
		return game.levelOneEnemyList[i].enemyProjectileList

	} else if game.levelTwoIsActive {
		if game.levelTwoEnemyList[i].projectileHold == false && game.levelTwoEnemyList[i].collision == false {
			game.levelTwoEnemyList[i].projectileHold = true
			g.enemyShootsProjectileAudioPlayer.Rewind()
			g.enemyShootsProjectileAudioPlayer.Play()

			go func() {
				<-time.After(3000 * time.Millisecond)
				game.levelTwoEnemyList[i].projectileHold = false
			}()
			game.projectileAndWallCollision = false

			tempFireball := game.fireball

			if game.levelTwoEnemyList[i].direction == "up" {
				tempFireball.xLoc = game.levelTwoEnemyList[i].xLoc + 20
				tempFireball.yLoc = game.levelTwoEnemyList[i].yLoc - 18
				tempFireball.dx = 0
				tempFireball.dy = -3
				game.levelTwoEnemyList[i].enemyProjectileList = append(game.levelTwoEnemyList[i].enemyProjectileList, tempFireball)
			} else if game.levelTwoEnemyList[i].direction == "down" {
				tempFireball.xLoc = game.levelTwoEnemyList[i].xLoc + 20
				tempFireball.yLoc = game.levelTwoEnemyList[i].yLoc + 55
				tempFireball.dx = 0
				tempFireball.dy = 3
				game.levelTwoEnemyList[i].enemyProjectileList = append(game.levelTwoEnemyList[i].enemyProjectileList, tempFireball)
			} else if game.levelTwoEnemyList[i].direction == "left" {
				tempFireball.xLoc = game.levelTwoEnemyList[i].xLoc - 15
				tempFireball.yLoc = game.levelTwoEnemyList[i].yLoc + 18
				tempFireball.dx = -3
				tempFireball.dy = 0
				game.levelTwoEnemyList[i].enemyProjectileList = append(game.levelTwoEnemyList[i].enemyProjectileList, tempFireball)
			} else if game.levelTwoEnemyList[i].direction == "right" {
				tempFireball.xLoc = game.levelTwoEnemyList[i].xLoc + 55
				tempFireball.yLoc = game.levelTwoEnemyList[i].yLoc + 18
				tempFireball.dx = 3
				tempFireball.dy = 0
				game.levelTwoEnemyList[i].enemyProjectileList = append(game.levelTwoEnemyList[i].enemyProjectileList, tempFireball)
			} else {
				tempFireball.xLoc = game.levelTwoEnemyList[i].xLoc + 20
				tempFireball.yLoc = game.levelTwoEnemyList[i].yLoc - 18
				tempFireball.dx = 0
				tempFireball.dy = -3
				game.levelTwoEnemyList[i].enemyProjectileList = append(game.levelTwoEnemyList[i].enemyProjectileList, tempFireball)
			}
		}
		return game.levelTwoEnemyList[i].enemyProjectileList

	} else if game.levelThreeIsActive {
		if game.levelThreeEnemyList[i].projectileHold == false && game.levelThreeEnemyList[i].collision == false {
			game.levelThreeEnemyList[i].projectileHold = true
			g.enemyShootsProjectileAudioPlayer.Rewind()
			g.enemyShootsProjectileAudioPlayer.Play()

			go func() {
				<-time.After(3000 * time.Millisecond)
				game.levelThreeEnemyList[i].projectileHold = false
			}()
			game.projectileAndWallCollision = false

			tempFireball := game.fireball

			if game.levelThreeEnemyList[i].direction == "up" {
				tempFireball.xLoc = game.levelThreeEnemyList[i].xLoc + 20
				tempFireball.yLoc = game.levelThreeEnemyList[i].yLoc - 18
				tempFireball.dx = 0
				tempFireball.dy = -3
				game.levelThreeEnemyList[i].enemyProjectileList = append(game.levelThreeEnemyList[i].enemyProjectileList, tempFireball)
			} else if game.levelThreeEnemyList[i].direction == "down" {
				tempFireball.xLoc = game.levelThreeEnemyList[i].xLoc + 20
				tempFireball.yLoc = game.levelThreeEnemyList[i].yLoc + 55
				tempFireball.dx = 0
				tempFireball.dy = 3
				game.levelThreeEnemyList[i].enemyProjectileList = append(game.levelThreeEnemyList[i].enemyProjectileList, tempFireball)
			} else if game.levelThreeEnemyList[i].direction == "left" {
				tempFireball.xLoc = game.levelThreeEnemyList[i].xLoc - 15
				tempFireball.yLoc = game.levelThreeEnemyList[i].yLoc + 18
				tempFireball.dx = -3
				tempFireball.dy = 0
				game.levelThreeEnemyList[i].enemyProjectileList = append(game.levelThreeEnemyList[i].enemyProjectileList, tempFireball)
			} else if game.levelThreeEnemyList[i].direction == "right" {
				tempFireball.xLoc = game.levelThreeEnemyList[i].xLoc + 55
				tempFireball.yLoc = game.levelThreeEnemyList[i].yLoc + 18
				tempFireball.dx = 3
				tempFireball.dy = 0
				game.levelThreeEnemyList[i].enemyProjectileList = append(game.levelThreeEnemyList[i].enemyProjectileList, tempFireball)
			} else {
				tempFireball.xLoc = game.levelThreeEnemyList[i].xLoc + 20
				tempFireball.yLoc = game.levelThreeEnemyList[i].yLoc - 18
				tempFireball.dx = 0
				tempFireball.dy = -3
				game.levelThreeEnemyList[i].enemyProjectileList = append(game.levelThreeEnemyList[i].enemyProjectileList, tempFireball)
			}
		}
		return game.levelThreeEnemyList[i].enemyProjectileList
	} else {
		return game.levelOneEnemyList
	}
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

func (game *Game) spawnLevel1Enemies() {
	if game.spawnedLevel1Enemies == false {
		personEnemy1 := game.personEnemy
		personEnemy2 := game.personEnemy
		monsterEnemy1 := game.monsterEnemy
		monsterEnemy2 := game.monsterEnemy
		personEnemy1.direction = "down"
		personEnemy2.direction = "left"
		monsterEnemy1.direction = "right"
		monsterEnemy2.direction = "left"
		personEnemy1.health = 1
		personEnemy2.health = 1
		monsterEnemy1.health = 2
		monsterEnemy2.health = 2
		personEnemy1.xLoc = 90
		personEnemy1.yLoc = 40
		personEnemy2.xLoc = 425
		personEnemy2.yLoc = 285
		monsterEnemy1.xLoc = 300
		monsterEnemy1.yLoc = 85
		monsterEnemy2.xLoc = 650
		monsterEnemy2.yLoc = 600

		game.levelOneEnemyList = append(game.levelOneEnemyList, personEnemy1)
		game.levelOneEnemyList = append(game.levelOneEnemyList, personEnemy2)
		game.levelOneEnemyList = append(game.levelOneEnemyList, monsterEnemy1)
		game.levelOneEnemyList = append(game.levelOneEnemyList, monsterEnemy2)
	}
	game.spawnedLevel1Enemies = true
}

func (game *Game) spawnLevel2Enemies() {
	if game.spawnedLevel2Enemies == false {
		game.collectedGold = false
		rand.Seed(int64(time.Now().Second()))
		game.coinSprite.xLoc = rand.Intn(ScreenWidth - 72)
		game.coinSprite.yLoc = rand.Intn(ScreenHeight - 72)
		personEnemy1 := game.personEnemy
		personEnemy2 := game.personEnemy
		monsterEnemy1 := game.monsterEnemy
		monsterEnemy2 := game.monsterEnemy
		personEnemy1.direction = "left"
		personEnemy2.direction = "up"
		monsterEnemy1.direction = "up"
		monsterEnemy2.direction = "left"
		personEnemy1.health = 1
		personEnemy2.health = 1
		monsterEnemy1.health = 2
		monsterEnemy2.health = 2
		personEnemy1.xLoc = 365
		personEnemy1.yLoc = 585
		personEnemy2.xLoc = 665
		personEnemy2.yLoc = 550
		monsterEnemy1.xLoc = 80
		monsterEnemy1.yLoc = 600
		monsterEnemy2.xLoc = 650
		monsterEnemy2.yLoc = 100

		game.levelTwoEnemyList = append(game.levelTwoEnemyList, personEnemy1)
		game.levelTwoEnemyList = append(game.levelTwoEnemyList, personEnemy2)
		game.levelTwoEnemyList = append(game.levelTwoEnemyList, monsterEnemy1)
		game.levelTwoEnemyList = append(game.levelTwoEnemyList, monsterEnemy2)
	}
	game.spawnedLevel2Enemies = true
}

func (game *Game) spawnLevel3Enemies() {
	if game.spawnedLevel3Enemies == false {
		if game.extraLifeAwarded == false {
			game.collectedGold = false
			rand.Seed(int64(time.Now().Second()))
			game.coinSprite.xLoc = rand.Intn(ScreenWidth - 72)
			game.coinSprite.yLoc = rand.Intn(ScreenHeight - 72)
		}
		personEnemy1 := game.personEnemy
		personEnemy2 := game.personEnemy
		monsterEnemy1 := game.monsterEnemy
		monsterEnemy2 := game.monsterEnemy
		personEnemy1.direction = "right"
		personEnemy2.direction = "left"
		monsterEnemy1.direction = "up"
		monsterEnemy2.direction = "right"
		personEnemy1.health = 1
		personEnemy2.health = 1
		monsterEnemy1.health = 2
		monsterEnemy2.health = 2
		personEnemy1.xLoc = 100
		personEnemy1.yLoc = 100
		personEnemy2.xLoc = 665
		personEnemy2.yLoc = 585
		monsterEnemy1.xLoc = 85
		monsterEnemy1.yLoc = 585
		monsterEnemy2.xLoc = 350
		monsterEnemy2.yLoc = 325

		game.levelThreeEnemyList = append(game.levelThreeEnemyList, personEnemy1)
		game.levelThreeEnemyList = append(game.levelThreeEnemyList, personEnemy2)
		game.levelThreeEnemyList = append(game.levelThreeEnemyList, monsterEnemy1)
		game.levelThreeEnemyList = append(game.levelThreeEnemyList, monsterEnemy2)
	}
	game.spawnedLevel3Enemies = true
}

func (game *Game) movementLevel1Enemies() {
	personEnemyMovementSpeed := 1
	if len(game.levelOneEnemyList) == 4 {
		for i := 0; i < len(game.levelOneEnemyList); i++ {
			//personEnemy1 moves up and down along left side
			if math.Abs(float64(game.levelOneEnemyList[i].xLoc-game.playerSprite.xLoc)) < 150 &&
				math.Abs(float64(game.levelOneEnemyList[i].yLoc-game.playerSprite.yLoc)) < 150 {
				//enemy is to the left and above player
				game.levelOneEnemyList[i].inPlayerProximity = true
				if game.levelOneEnemyList[i].xLoc <= game.playerSprite.xLoc &&
					game.levelOneEnemyList[i].yLoc <= game.playerSprite.yLoc {
					game.levelOneEnemyList[i].dx = 1
					game.levelOneEnemyList[i].dy = 1
					if math.Abs(float64(game.levelOneEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelOneEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelOneEnemyList[i].direction = "right"
					} else {
						game.levelOneEnemyList[i].direction = "down"
					}
					game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
					game.levelOneEnemyList[i].yLoc += game.levelOneEnemyList[i].dy
					game.levelOneEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)

				} else if game.levelOneEnemyList[i].xLoc >= game.playerSprite.xLoc &&
					game.levelOneEnemyList[i].yLoc <= game.playerSprite.yLoc {
					//enemy to the right and above player
					game.levelOneEnemyList[i].dx = -1
					game.levelOneEnemyList[i].dy = 1
					if math.Abs(float64(game.levelOneEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelOneEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelOneEnemyList[i].direction = "left"
					} else {
						game.levelOneEnemyList[i].direction = "down"
					}
					game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
					game.levelOneEnemyList[i].yLoc += game.levelOneEnemyList[i].dy
					game.levelOneEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				} else if game.levelOneEnemyList[i].xLoc <= game.playerSprite.xLoc &&
					game.levelOneEnemyList[i].yLoc >= game.playerSprite.yLoc {
					//enemy to the left and below player
					game.levelOneEnemyList[i].dx = 1
					game.levelOneEnemyList[i].dy = -1
					if math.Abs(float64(game.levelOneEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelOneEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelOneEnemyList[i].direction = "right"
					} else {
						game.levelOneEnemyList[i].direction = "up"
					}
					game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
					game.levelOneEnemyList[i].yLoc += game.levelOneEnemyList[i].dy
					game.levelOneEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				} else if game.levelOneEnemyList[i].xLoc >= game.playerSprite.xLoc &&
					game.levelOneEnemyList[i].yLoc >= game.playerSprite.yLoc {
					//enemy location to the right and below
					game.levelOneEnemyList[i].dx = -1
					game.levelOneEnemyList[i].dy = -1
					if math.Abs(float64(game.levelOneEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelOneEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelOneEnemyList[i].direction = "left"
					} else {
						game.levelOneEnemyList[i].direction = "up"
					}
					game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
					game.levelOneEnemyList[i].yLoc += game.levelOneEnemyList[i].dy
					game.levelOneEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				}
			} else if math.Abs(float64(game.levelOneEnemyList[i].xLoc-game.playerSprite.xLoc)) >= 150 ||
				math.Abs(float64(game.levelOneEnemyList[i].yLoc-game.playerSprite.yLoc)) >= 150 &&
					game.levelOneEnemyList[i].inPlayerProximity == false {

				if i == 0 {
					if game.levelOneEnemyList[i].direction == "down" && game.levelOneEnemyList[i].yLoc < 500 {
						game.levelOneEnemyList[i].dy = personEnemyMovementSpeed
						game.levelOneEnemyList[i].yLoc += game.levelOneEnemyList[i].dy
					} else if game.levelOneEnemyList[i].direction == "down" && game.levelOneEnemyList[i].yLoc >= 500 {
						game.levelOneEnemyList[i].direction = "up"
						game.levelOneEnemyList[i].dy = 0
					} else if game.levelOneEnemyList[i].direction == "up" &&
						game.levelOneEnemyList[i].yLoc > 40 {
						game.levelOneEnemyList[i].dy = -personEnemyMovementSpeed
						game.levelOneEnemyList[i].yLoc += game.levelOneEnemyList[i].dy
					} else if game.levelOneEnemyList[i].direction == "up" &&
						game.levelOneEnemyList[i].yLoc <= 40 {
						game.levelOneEnemyList[i].direction = "down"
						game.levelOneEnemyList[i].dy = 0
					}
				} else if i == 1 {
					// personEnemy2 moves around in a square
					if game.levelOneEnemyList[i].direction == "left" && game.levelOneEnemyList[i].yLoc <= 285 &&
						game.levelOneEnemyList[i].xLoc <= 425 && game.levelOneEnemyList[i].xLoc > 285 {
						game.levelOneEnemyList[i].dx = -personEnemyMovementSpeed
						game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
					} else if game.levelOneEnemyList[i].direction == "left" && game.levelOneEnemyList[i].yLoc <= 285 &&
						game.levelOneEnemyList[i].xLoc <= 285 {
						game.levelOneEnemyList[i].direction = "down"
						game.levelOneEnemyList[i].dx = 0
					} else if game.levelOneEnemyList[i].direction == "down" &&
						game.levelOneEnemyList[i].yLoc < 425 {
						game.levelOneEnemyList[i].dy = personEnemyMovementSpeed
						game.levelOneEnemyList[i].yLoc += game.levelOneEnemyList[i].dy
					} else if game.levelOneEnemyList[i].direction == "down" && game.levelOneEnemyList[i].yLoc >= 425 {
						game.levelOneEnemyList[i].direction = "right"
						game.levelOneEnemyList[i].dy = 0
					} else if game.levelOneEnemyList[i].direction == "right" && game.levelOneEnemyList[i].yLoc >= 425 &&
						game.levelOneEnemyList[i].xLoc < 425 {
						game.levelOneEnemyList[i].dx = personEnemyMovementSpeed
						game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
					} else if game.levelOneEnemyList[i].direction == "right" && game.levelOneEnemyList[i].yLoc >= 425 &&
						game.levelOneEnemyList[i].xLoc >= 425 {
						game.levelOneEnemyList[i].direction = "up"
						game.levelOneEnemyList[i].dx = 0
					} else if game.levelOneEnemyList[i].direction == "up" && game.levelOneEnemyList[i].xLoc >= 425 &&
						game.levelOneEnemyList[i].yLoc > 285 {
						game.levelOneEnemyList[i].dy = -personEnemyMovementSpeed
						game.levelOneEnemyList[i].yLoc += game.levelOneEnemyList[i].dy
					} else if game.levelOneEnemyList[i].direction == "up" && game.levelOneEnemyList[i].yLoc <= 285 {
						game.levelOneEnemyList[i].direction = "left"
						game.levelOneEnemyList[i].dy = 0
					}
				} else if i == 2 {
					//monsterEnemy1 moves back and forth left and right at the top and chases if in certain proximity
					if game.levelOneEnemyList[i].direction == "right" && game.levelOneEnemyList[i].xLoc < 600 {
						game.levelOneEnemyList[i].dx = personEnemyMovementSpeed
						game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
					} else if game.levelOneEnemyList[i].direction == "right" && game.levelOneEnemyList[i].xLoc >= 600 {
						game.levelOneEnemyList[i].direction = "left"
						game.levelOneEnemyList[i].dx = 0
					} else if game.levelOneEnemyList[i].direction == "left" &&
						game.levelOneEnemyList[i].xLoc > 300 {
						game.levelOneEnemyList[i].dx = -personEnemyMovementSpeed
						game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
					} else if game.levelOneEnemyList[i].direction == "left" &&
						game.levelOneEnemyList[i].xLoc <= 300 {
						game.levelOneEnemyList[i].direction = "right"
						game.levelOneEnemyList[i].dx = 0
					}
				} else if i == 3 {
					//monsterEnemy2 moves back and forth left and right at the bottom and chases if in certain proximity
					if game.levelOneEnemyList[i].direction == "left" && game.levelOneEnemyList[i].xLoc > 100 {
						game.levelOneEnemyList[i].dx = -personEnemyMovementSpeed
						game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
					} else if game.levelOneEnemyList[i].direction == "left" && game.levelOneEnemyList[i].xLoc <= 100 {
						game.levelOneEnemyList[i].direction = "right"
						game.levelOneEnemyList[i].dx = 0
					} else if game.levelOneEnemyList[i].direction == "right" &&
						game.levelOneEnemyList[i].xLoc < 700 {
						game.levelOneEnemyList[i].dx = personEnemyMovementSpeed
						game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
					} else if game.levelOneEnemyList[i].direction == "right" &&
						game.levelOneEnemyList[i].xLoc >= 700 {
						game.levelOneEnemyList[i].direction = "left"
						game.levelOneEnemyList[i].dx = 0
					}
				}
			} else {
				//enemy is to the left and above player
				game.levelOneEnemyList[i].inPlayerProximity = true
				if game.levelOneEnemyList[i].xLoc <= game.playerSprite.xLoc &&
					game.levelOneEnemyList[i].yLoc <= game.playerSprite.yLoc {
					game.levelOneEnemyList[i].dx = 1
					game.levelOneEnemyList[i].dy = 1
					if math.Abs(float64(game.levelOneEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelOneEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelOneEnemyList[i].direction = "right"
					} else {
						game.levelOneEnemyList[i].direction = "down"
					}
					game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
					game.levelOneEnemyList[i].yLoc += game.levelOneEnemyList[i].dy
					game.levelOneEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				} else if game.levelOneEnemyList[i].xLoc >= game.playerSprite.xLoc &&
					game.levelOneEnemyList[i].yLoc <= game.playerSprite.yLoc {
					game.levelOneEnemyList[i].dx = -1
					game.levelOneEnemyList[i].dy = 1
					if math.Abs(float64(game.levelOneEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelOneEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelOneEnemyList[i].direction = "left"
					} else {
						game.levelOneEnemyList[i].direction = "down"
					}
					game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
					game.levelOneEnemyList[i].yLoc += game.levelOneEnemyList[i].dy
					game.levelOneEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				} else if game.levelOneEnemyList[i].xLoc <= game.playerSprite.xLoc &&
					game.levelOneEnemyList[i].yLoc >= game.playerSprite.yLoc {
					game.levelOneEnemyList[i].dx = 1
					game.levelOneEnemyList[i].dy = -1
					if math.Abs(float64(game.levelOneEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelOneEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelOneEnemyList[i].direction = "right"
					} else {
						game.levelOneEnemyList[i].direction = "up"
					}
					game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
					game.levelOneEnemyList[i].yLoc += game.levelOneEnemyList[i].dy
					game.levelOneEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				} else if game.levelOneEnemyList[i].xLoc >= game.playerSprite.xLoc &&
					game.levelOneEnemyList[i].yLoc >= game.playerSprite.yLoc {
					game.levelOneEnemyList[i].dx = -1
					game.levelOneEnemyList[i].dy = -1
					if math.Abs(float64(game.levelOneEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelOneEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelOneEnemyList[i].direction = "left"
					} else {
						game.levelOneEnemyList[i].direction = "up"
					}
					game.levelOneEnemyList[i].xLoc += game.levelOneEnemyList[i].dx
					game.levelOneEnemyList[i].yLoc += game.levelOneEnemyList[i].dy
					game.levelOneEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				}
			}
		}
	}
}

func (game *Game) movementLevel2Enemies() {
	personEnemyMovementSpeed := 1
	if len(game.levelTwoEnemyList) == 4 {
		for i := 0; i < len(game.levelTwoEnemyList); i++ {
			//personEnemy1 moves up and down along left side
			if math.Abs(float64(game.levelTwoEnemyList[i].xLoc-game.playerSprite.xLoc)) < 150 &&
				math.Abs(float64(game.levelTwoEnemyList[i].yLoc-game.playerSprite.yLoc)) < 150 {
				//enemy is to the left and above player
				game.levelTwoEnemyList[i].inPlayerProximity = true
				if game.levelTwoEnemyList[i].xLoc <= game.playerSprite.xLoc &&
					game.levelTwoEnemyList[i].yLoc <= game.playerSprite.yLoc {
					game.levelTwoEnemyList[i].dx = 1
					game.levelTwoEnemyList[i].dy = 1
					if math.Abs(float64(game.levelTwoEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelTwoEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelTwoEnemyList[i].direction = "right"
					} else {
						game.levelTwoEnemyList[i].direction = "down"
					}
					game.levelTwoEnemyList[i].xLoc += game.levelTwoEnemyList[i].dx
					game.levelTwoEnemyList[i].yLoc += game.levelTwoEnemyList[i].dy
					game.levelTwoEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)

				} else if game.levelTwoEnemyList[i].xLoc >= game.playerSprite.xLoc &&
					game.levelTwoEnemyList[i].yLoc <= game.playerSprite.yLoc {
					//enemy to the right and above player
					game.levelTwoEnemyList[i].dx = -1
					game.levelTwoEnemyList[i].dy = 1
					if math.Abs(float64(game.levelTwoEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelTwoEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelTwoEnemyList[i].direction = "left"
					} else {
						game.levelTwoEnemyList[i].direction = "down"
					}
					game.levelTwoEnemyList[i].xLoc += game.levelTwoEnemyList[i].dx
					game.levelTwoEnemyList[i].yLoc += game.levelTwoEnemyList[i].dy
					game.levelTwoEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				} else if game.levelTwoEnemyList[i].xLoc <= game.playerSprite.xLoc &&
					game.levelTwoEnemyList[i].yLoc >= game.playerSprite.yLoc {
					//enemy to the left and below player
					game.levelTwoEnemyList[i].dx = 1
					game.levelTwoEnemyList[i].dy = -1
					if math.Abs(float64(game.levelTwoEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelTwoEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelTwoEnemyList[i].direction = "right"
					} else {
						game.levelTwoEnemyList[i].direction = "up"
					}
					game.levelTwoEnemyList[i].xLoc += game.levelTwoEnemyList[i].dx
					game.levelTwoEnemyList[i].yLoc += game.levelTwoEnemyList[i].dy
					game.levelTwoEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				} else if game.levelTwoEnemyList[i].xLoc >= game.playerSprite.xLoc &&
					game.levelTwoEnemyList[i].yLoc >= game.playerSprite.yLoc {
					//enemy location to the right and below
					game.levelTwoEnemyList[i].dx = -1
					game.levelTwoEnemyList[i].dy = -1
					if math.Abs(float64(game.levelTwoEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelTwoEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelTwoEnemyList[i].direction = "left"
					} else {
						game.levelTwoEnemyList[i].direction = "up"
					}
					game.levelTwoEnemyList[i].xLoc += game.levelTwoEnemyList[i].dx
					game.levelTwoEnemyList[i].yLoc += game.levelTwoEnemyList[i].dy
					game.levelTwoEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				}
			} else if math.Abs(float64(game.levelTwoEnemyList[i].xLoc-game.playerSprite.xLoc)) >= 150 ||
				math.Abs(float64(game.levelTwoEnemyList[i].yLoc-game.playerSprite.yLoc)) >= 150 &&
					game.levelTwoEnemyList[i].inPlayerProximity == false {
				if i == 0 {
					if game.levelTwoEnemyList[i].direction == "left" && game.levelTwoEnemyList[i].xLoc > 150 {
						game.levelTwoEnemyList[i].dx = -personEnemyMovementSpeed
						game.levelTwoEnemyList[i].xLoc += game.levelTwoEnemyList[i].dx
					} else if game.levelTwoEnemyList[i].direction == "left" && game.levelTwoEnemyList[i].xLoc <= 150 {
						game.levelTwoEnemyList[i].direction = "right"
						game.levelTwoEnemyList[i].dx = 0
					} else if game.levelTwoEnemyList[i].direction == "right" &&
						game.levelTwoEnemyList[i].xLoc < 365 {
						game.levelTwoEnemyList[i].dx = personEnemyMovementSpeed
						game.levelTwoEnemyList[i].xLoc += game.levelTwoEnemyList[i].dx
					} else if game.levelTwoEnemyList[i].direction == "right" &&
						game.levelTwoEnemyList[i].xLoc >= 365 {
						game.levelTwoEnemyList[i].direction = "left"
						game.levelTwoEnemyList[i].dy = 0
					}
				} else if i == 1 {
					// personEnemy2 moves up and down on right side
					if game.levelTwoEnemyList[i].direction == "up" && game.levelTwoEnemyList[i].yLoc > 150 {
						game.levelTwoEnemyList[i].dy = -personEnemyMovementSpeed
						game.levelTwoEnemyList[i].yLoc += game.levelTwoEnemyList[i].dy
					} else if game.levelTwoEnemyList[i].direction == "up" && game.levelTwoEnemyList[i].yLoc <= 150 {
						game.levelTwoEnemyList[i].direction = "down"
						game.levelTwoEnemyList[i].dy = 0
					} else if game.levelTwoEnemyList[i].direction == "down" &&
						game.levelTwoEnemyList[i].yLoc < 550 {
						game.levelTwoEnemyList[i].dy = personEnemyMovementSpeed
						game.levelTwoEnemyList[i].yLoc += game.levelTwoEnemyList[i].dy
					} else if game.levelTwoEnemyList[i].direction == "down" &&
						game.levelTwoEnemyList[i].yLoc >= 550 {
						game.levelTwoEnemyList[i].direction = "up"
						game.levelTwoEnemyList[i].dy = 0
					}

				} else if i == 2 {
					//monsterEnemy1 moves up and down
					if game.levelTwoEnemyList[i].direction == "up" && game.levelTwoEnemyList[i].yLoc > 150 {
						game.levelTwoEnemyList[i].dy = -personEnemyMovementSpeed
						game.levelTwoEnemyList[i].yLoc += game.levelTwoEnemyList[i].dy
					} else if game.levelTwoEnemyList[i].direction == "up" && game.levelTwoEnemyList[i].yLoc <= 150 {
						game.levelTwoEnemyList[i].direction = "down"
						game.levelTwoEnemyList[i].dy = 0
					} else if game.levelTwoEnemyList[i].direction == "down" &&
						game.levelTwoEnemyList[i].yLoc < 400 {
						game.levelTwoEnemyList[i].dy = personEnemyMovementSpeed
						game.levelTwoEnemyList[i].yLoc += game.levelTwoEnemyList[i].dy
					} else if game.levelTwoEnemyList[i].direction == "down" &&
						game.levelTwoEnemyList[i].yLoc >= 400 {
						game.levelTwoEnemyList[i].direction = "up"
						game.levelTwoEnemyList[i].dy = 0
					}
				} else if i == 3 {
					//monsterEnemy2 moves back and forth left and right at the top and chases if in certain proximity
					if game.levelTwoEnemyList[i].direction == "left" && game.levelTwoEnemyList[i].xLoc > 400 {
						game.levelTwoEnemyList[i].dx = -personEnemyMovementSpeed
						game.levelTwoEnemyList[i].xLoc += game.levelTwoEnemyList[i].dx
					} else if game.levelTwoEnemyList[i].direction == "left" && game.levelTwoEnemyList[i].xLoc <= 400 {
						game.levelTwoEnemyList[i].direction = "right"
						game.levelTwoEnemyList[i].dx = 0
					} else if game.levelTwoEnemyList[i].direction == "right" &&
						game.levelTwoEnemyList[i].xLoc < 650 {
						game.levelTwoEnemyList[i].dx = personEnemyMovementSpeed
						game.levelTwoEnemyList[i].xLoc += game.levelTwoEnemyList[i].dx
					} else if game.levelTwoEnemyList[i].direction == "right" &&
						game.levelTwoEnemyList[i].xLoc >= 650 {
						game.levelTwoEnemyList[i].direction = "left"
						game.levelTwoEnemyList[i].dy = 0
					}
				}
			} else {
				//enemy is to the left and above player
				game.levelTwoEnemyList[i].inPlayerProximity = true
				if game.levelTwoEnemyList[i].xLoc <= game.playerSprite.xLoc &&
					game.levelTwoEnemyList[i].yLoc <= game.playerSprite.yLoc {
					game.levelTwoEnemyList[i].dx = 1
					game.levelTwoEnemyList[i].dy = 1
					if math.Abs(float64(game.levelTwoEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelTwoEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelTwoEnemyList[i].direction = "right"
					} else {
						game.levelTwoEnemyList[i].direction = "down"
					}
					game.levelTwoEnemyList[i].xLoc += game.levelTwoEnemyList[i].dx
					game.levelTwoEnemyList[i].yLoc += game.levelTwoEnemyList[i].dy
					game.levelTwoEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				} else if game.levelTwoEnemyList[i].xLoc >= game.playerSprite.xLoc &&
					game.levelTwoEnemyList[i].yLoc <= game.playerSprite.yLoc {
					game.levelTwoEnemyList[i].dx = -1
					game.levelTwoEnemyList[i].dy = 1
					if math.Abs(float64(game.levelTwoEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelTwoEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelTwoEnemyList[i].direction = "left"
					} else {
						game.levelTwoEnemyList[i].direction = "down"
					}
					game.levelTwoEnemyList[i].xLoc += game.levelTwoEnemyList[i].dx
					game.levelTwoEnemyList[i].yLoc += game.levelTwoEnemyList[i].dy
					game.levelTwoEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				} else if game.levelTwoEnemyList[i].xLoc <= game.playerSprite.xLoc &&
					game.levelTwoEnemyList[i].yLoc >= game.playerSprite.yLoc {
					game.levelTwoEnemyList[i].dx = 1
					game.levelTwoEnemyList[i].dy = -1
					if math.Abs(float64(game.levelTwoEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelTwoEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelTwoEnemyList[i].direction = "right"
					} else {
						game.levelTwoEnemyList[i].direction = "up"
					}
					game.levelTwoEnemyList[i].xLoc += game.levelTwoEnemyList[i].dx
					game.levelTwoEnemyList[i].yLoc += game.levelTwoEnemyList[i].dy
					game.levelTwoEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				} else if game.levelTwoEnemyList[i].xLoc >= game.playerSprite.xLoc &&
					game.levelTwoEnemyList[i].yLoc >= game.playerSprite.yLoc {
					game.levelTwoEnemyList[i].dx = -1
					game.levelTwoEnemyList[i].dy = -1
					if math.Abs(float64(game.levelTwoEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelTwoEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelTwoEnemyList[i].direction = "left"
					} else {
						game.levelTwoEnemyList[i].direction = "up"
					}
					game.levelTwoEnemyList[i].xLoc += game.levelTwoEnemyList[i].dx
					game.levelTwoEnemyList[i].yLoc += game.levelTwoEnemyList[i].dy
					game.levelTwoEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				}
			}
		}
	}
}

func (game *Game) movementLevel3Enemies() {
	personEnemyMovementSpeed := 1
	if len(game.levelThreeEnemyList) == 4 {
		for i := 0; i < len(game.levelThreeEnemyList); i++ {
			//personEnemy1 moves up and down along left side
			if math.Abs(float64(game.levelThreeEnemyList[i].xLoc-game.playerSprite.xLoc)) < 150 &&
				math.Abs(float64(game.levelThreeEnemyList[i].yLoc-game.playerSprite.yLoc)) < 150 {
				//enemy is to the left and above player
				game.levelThreeEnemyList[i].inPlayerProximity = true
				if game.levelThreeEnemyList[i].xLoc <= game.playerSprite.xLoc &&
					game.levelThreeEnemyList[i].yLoc <= game.playerSprite.yLoc {
					game.levelThreeEnemyList[i].dx = 1
					game.levelThreeEnemyList[i].dy = 1
					if math.Abs(float64(game.levelThreeEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelThreeEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelThreeEnemyList[i].direction = "right"
					} else {
						game.levelThreeEnemyList[i].direction = "down"
					}
					game.levelThreeEnemyList[i].xLoc += game.levelThreeEnemyList[i].dx
					game.levelThreeEnemyList[i].yLoc += game.levelThreeEnemyList[i].dy
					game.levelThreeEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)

				} else if game.levelThreeEnemyList[i].xLoc >= game.playerSprite.xLoc &&
					game.levelThreeEnemyList[i].yLoc <= game.playerSprite.yLoc {
					//enemy to the right and above player
					game.levelThreeEnemyList[i].dx = -1
					game.levelThreeEnemyList[i].dy = 1
					if math.Abs(float64(game.levelThreeEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelThreeEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelThreeEnemyList[i].direction = "left"
					} else {
						game.levelThreeEnemyList[i].direction = "down"
					}
					game.levelThreeEnemyList[i].xLoc += game.levelThreeEnemyList[i].dx
					game.levelThreeEnemyList[i].yLoc += game.levelThreeEnemyList[i].dy
					game.levelThreeEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				} else if game.levelThreeEnemyList[i].xLoc <= game.playerSprite.xLoc &&
					game.levelThreeEnemyList[i].yLoc >= game.playerSprite.yLoc {
					//enemy to the left and below player
					game.levelThreeEnemyList[i].dx = 1
					game.levelThreeEnemyList[i].dy = -1
					if math.Abs(float64(game.levelThreeEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelThreeEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelThreeEnemyList[i].direction = "right"
					} else {
						game.levelThreeEnemyList[i].direction = "up"
					}
					game.levelThreeEnemyList[i].xLoc += game.levelThreeEnemyList[i].dx
					game.levelThreeEnemyList[i].yLoc += game.levelThreeEnemyList[i].dy
					game.levelThreeEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				} else if game.levelThreeEnemyList[i].xLoc >= game.playerSprite.xLoc &&
					game.levelThreeEnemyList[i].yLoc >= game.playerSprite.yLoc {
					//enemy location to the right and below
					game.levelThreeEnemyList[i].dx = -1
					game.levelThreeEnemyList[i].dy = -1
					if math.Abs(float64(game.levelThreeEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelThreeEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelThreeEnemyList[i].direction = "left"
					} else {
						game.levelThreeEnemyList[i].direction = "up"
					}
					game.levelThreeEnemyList[i].xLoc += game.levelThreeEnemyList[i].dx
					game.levelThreeEnemyList[i].yLoc += game.levelThreeEnemyList[i].dy
					game.levelThreeEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				}
			} else if math.Abs(float64(game.levelThreeEnemyList[i].xLoc-game.playerSprite.xLoc)) >= 150 ||
				math.Abs(float64(game.levelThreeEnemyList[i].yLoc-game.playerSprite.yLoc)) >= 150 &&
					game.levelThreeEnemyList[i].inPlayerProximity == false {
				if i == 0 {
					if game.levelThreeEnemyList[i].direction == "right" && game.levelThreeEnemyList[i].xLoc < 600 {
						game.levelThreeEnemyList[i].dx = personEnemyMovementSpeed
						game.levelThreeEnemyList[i].xLoc += game.levelThreeEnemyList[i].dx
					} else if game.levelThreeEnemyList[i].direction == "left" && game.levelThreeEnemyList[i].xLoc >= 600 {
						game.levelThreeEnemyList[i].direction = "left"
						game.levelThreeEnemyList[i].dx = 0
					} else if game.levelThreeEnemyList[i].direction == "left" &&
						game.levelThreeEnemyList[i].xLoc > 100 {
						game.levelThreeEnemyList[i].dx = -personEnemyMovementSpeed
						game.levelThreeEnemyList[i].xLoc += game.levelThreeEnemyList[i].dx
					} else if game.levelThreeEnemyList[i].direction == "left" &&
						game.levelThreeEnemyList[i].xLoc <= 100 {
						game.levelThreeEnemyList[i].direction = "left"
						game.levelThreeEnemyList[i].dy = 0
					}
				} else if i == 1 {
					// personEnemy2 moves up and down on right side
					if game.levelThreeEnemyList[i].direction == "left" && game.levelThreeEnemyList[i].xLoc > 200 {
						game.levelThreeEnemyList[i].dx = -personEnemyMovementSpeed
						game.levelThreeEnemyList[i].xLoc += game.levelThreeEnemyList[i].dx
					} else if game.levelThreeEnemyList[i].direction == "left" && game.levelThreeEnemyList[i].xLoc <= 200 {
						game.levelThreeEnemyList[i].direction = "right"
						game.levelThreeEnemyList[i].dx = 0
					} else if game.levelThreeEnemyList[i].direction == "right" &&
						game.levelThreeEnemyList[i].xLoc < 665 {
						game.levelThreeEnemyList[i].dx = personEnemyMovementSpeed
						game.levelThreeEnemyList[i].xLoc += game.levelThreeEnemyList[i].dx
					} else if game.levelThreeEnemyList[i].direction == "right" &&
						game.levelThreeEnemyList[i].xLoc >= 665 {
						game.levelThreeEnemyList[i].direction = "left"
						game.levelThreeEnemyList[i].dy = 0
					}
				} else if i == 2 {
					//monsterEnemy1 moves up and down
					if game.levelThreeEnemyList[i].direction == "up" && game.levelThreeEnemyList[i].yLoc > 150 {
						game.levelThreeEnemyList[i].dy = -personEnemyMovementSpeed
						game.levelThreeEnemyList[i].yLoc += game.levelThreeEnemyList[i].dy
					} else if game.levelThreeEnemyList[i].direction == "up" && game.levelThreeEnemyList[i].yLoc <= 150 {
						game.levelThreeEnemyList[i].direction = "down"
						game.levelThreeEnemyList[i].dy = 0
					} else if game.levelThreeEnemyList[i].direction == "down" &&
						game.levelThreeEnemyList[i].yLoc < 585 {
						game.levelThreeEnemyList[i].dy = personEnemyMovementSpeed
						game.levelThreeEnemyList[i].yLoc += game.levelThreeEnemyList[i].dy
					} else if game.levelThreeEnemyList[i].direction == "down" &&
						game.levelThreeEnemyList[i].yLoc >= 585 {
						game.levelThreeEnemyList[i].direction = "up"
						game.levelThreeEnemyList[i].dy = 0
					}
				} else if i == 3 {
					//monsterEnemy2 moves back and forth left and right at the top and chases if in certain proximity
					if game.levelThreeEnemyList[i].direction == "right" && game.levelThreeEnemyList[i].xLoc < 600 {
						game.levelThreeEnemyList[i].dx = personEnemyMovementSpeed
						game.levelThreeEnemyList[i].xLoc += game.levelThreeEnemyList[i].dx
					} else if game.levelThreeEnemyList[i].direction == "right" && game.levelThreeEnemyList[i].xLoc >= 600 {
						game.levelThreeEnemyList[i].direction = "left"
						game.levelThreeEnemyList[i].dx = 0
					} else if game.levelThreeEnemyList[i].direction == "left" &&
						game.levelThreeEnemyList[i].xLoc > 350 {
						game.levelThreeEnemyList[i].dx = personEnemyMovementSpeed
						game.levelThreeEnemyList[i].xLoc -= game.levelThreeEnemyList[i].dx
					} else if game.levelThreeEnemyList[i].direction == "left" &&
						game.levelThreeEnemyList[i].xLoc <= 350 {
						game.levelThreeEnemyList[i].direction = "right"
						game.levelThreeEnemyList[i].dy = 0
					}
				}
			} else {
				//enemy is to the left and above player
				game.levelThreeEnemyList[i].inPlayerProximity = true
				if game.levelThreeEnemyList[i].xLoc <= game.playerSprite.xLoc &&
					game.levelThreeEnemyList[i].yLoc <= game.playerSprite.yLoc {
					game.levelThreeEnemyList[i].dx = 1
					game.levelThreeEnemyList[i].dy = 1
					if math.Abs(float64(game.levelThreeEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelThreeEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelThreeEnemyList[i].direction = "right"
					} else {
						game.levelThreeEnemyList[i].direction = "down"
					}
					game.levelThreeEnemyList[i].xLoc += game.levelThreeEnemyList[i].dx
					game.levelThreeEnemyList[i].yLoc += game.levelThreeEnemyList[i].dy
					game.levelThreeEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				} else if game.levelThreeEnemyList[i].xLoc >= game.playerSprite.xLoc &&
					game.levelThreeEnemyList[i].yLoc <= game.playerSprite.yLoc {
					game.levelThreeEnemyList[i].dx = -1
					game.levelThreeEnemyList[i].dy = 1
					if math.Abs(float64(game.levelThreeEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelThreeEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelThreeEnemyList[i].direction = "left"
					} else {
						game.levelThreeEnemyList[i].direction = "down"
					}
					game.levelThreeEnemyList[i].xLoc += game.levelThreeEnemyList[i].dx
					game.levelThreeEnemyList[i].yLoc += game.levelThreeEnemyList[i].dy
					game.levelThreeEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				} else if game.levelThreeEnemyList[i].xLoc <= game.playerSprite.xLoc &&
					game.levelThreeEnemyList[i].yLoc >= game.playerSprite.yLoc {
					game.levelThreeEnemyList[i].dx = 1
					game.levelThreeEnemyList[i].dy = -1
					if math.Abs(float64(game.levelThreeEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelThreeEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelThreeEnemyList[i].direction = "right"
					} else {
						game.levelThreeEnemyList[i].direction = "up"
					}
					game.levelThreeEnemyList[i].xLoc += game.levelThreeEnemyList[i].dx
					game.levelThreeEnemyList[i].yLoc += game.levelThreeEnemyList[i].dy
					game.levelThreeEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				} else if game.levelThreeEnemyList[i].xLoc >= game.playerSprite.xLoc &&
					game.levelThreeEnemyList[i].yLoc >= game.playerSprite.yLoc {
					game.levelThreeEnemyList[i].dx = -1
					game.levelThreeEnemyList[i].dy = -1
					if math.Abs(float64(game.levelThreeEnemyList[i].xLoc-game.playerSprite.xLoc)) >
						math.Abs(float64(game.levelThreeEnemyList[i].yLoc-game.playerSprite.yLoc)) {
						game.levelThreeEnemyList[i].direction = "left"
					} else {
						game.levelThreeEnemyList[i].direction = "up"
					}
					game.levelThreeEnemyList[i].xLoc += game.levelThreeEnemyList[i].dx
					game.levelThreeEnemyList[i].yLoc += game.levelThreeEnemyList[i].dy
					game.levelThreeEnemyList[i].enemyProjectileList =
						game.enemyShootFireball(i)
				}
			}
		}
	}
}

func (game *Game) manageLevel1CollisionDetection() {
	if game.collectedGold == false {
		game.collectedGold = game.gotGold(game.playerSprite, game.coinSprite)
	}

	//player collision with wall check
	if game.playerAndWallCollision == false {
		game.playerAndWallCollision = wallCollisionCheckFirstLevel(game.playerSprite, 61)
	} else {
		if game.playerSprite.xLoc < ScreenWidth/2 {
			game.playerSprite.yLoc = 450
			game.playerSprite.xLoc = 650 //player width
		} else if game.playerSprite.xLoc > ScreenWidth/2 {
			game.playerSprite.yLoc = ScreenHeight / 2
			game.playerSprite.xLoc = 74 //player width
		}
		game.playerAndWallCollision = false
		game.deathCounter += 1
		g.playerDeathAudioPlayer.Rewind()
		g.playerDeathAudioPlayer.Play()
	}

	//enemy collision with wall check
	if len(game.levelOneEnemyList) > 0 {
		for i := 0; i < len(game.levelOneEnemyList); i++ {
			if game.levelOneEnemyList[i].collision == false {
				if game.levelOneEnemyList[i].direction == "left" {
					spriteWidth, _ := game.levelOneEnemyList[i].leftPict.Size()
					game.levelOneEnemyList[i].collision = wallCollisionCheckFirstLevel(game.levelOneEnemyList[i], spriteWidth)
					if game.levelOneEnemyList[i].collision == true && spriteWidth == 50 {
						g.monsterEnemyDeathAudioPlayer.Rewind()
						g.monsterEnemyDeathAudioPlayer.Play()
					} else if game.levelOneEnemyList[i].collision == true && spriteWidth != 50 {
						g.humanEnemyDeathAudioPlayer.Rewind()
						g.humanEnemyDeathAudioPlayer.Play()
					}
				} else if game.levelOneEnemyList[i].direction == "right" {
					spriteWidth, _ := game.levelOneEnemyList[i].rightPict.Size()
					game.levelOneEnemyList[i].collision = wallCollisionCheckFirstLevel(game.levelOneEnemyList[i], spriteWidth)
					if game.levelOneEnemyList[i].collision == true && spriteWidth == 50 {
						g.monsterEnemyDeathAudioPlayer.Rewind()
						g.monsterEnemyDeathAudioPlayer.Play()
					} else if game.levelOneEnemyList[i].collision == true && spriteWidth != 50 {
						g.humanEnemyDeathAudioPlayer.Rewind()
						g.humanEnemyDeathAudioPlayer.Play()
					}
				} else if game.levelOneEnemyList[i].direction == "up" {
					spriteWidth, _ := game.levelOneEnemyList[i].upPict.Size()
					game.levelOneEnemyList[i].collision = wallCollisionCheckFirstLevel(game.levelOneEnemyList[i], spriteWidth)
					if game.levelOneEnemyList[i].collision == true && spriteWidth == 50 {
						g.monsterEnemyDeathAudioPlayer.Rewind()
						g.monsterEnemyDeathAudioPlayer.Play()
					} else if game.levelOneEnemyList[i].collision == true && spriteWidth != 50 {
						g.humanEnemyDeathAudioPlayer.Rewind()
						g.humanEnemyDeathAudioPlayer.Play()
					}
				} else if game.levelOneEnemyList[i].direction == "down" {
					spriteWidth, _ := game.levelOneEnemyList[i].downPict.Size()
					game.levelOneEnemyList[i].collision = wallCollisionCheckFirstLevel(game.levelOneEnemyList[i], spriteWidth)
					if game.levelOneEnemyList[i].collision == true && spriteWidth == 50 {
						g.monsterEnemyDeathAudioPlayer.Rewind()
						g.monsterEnemyDeathAudioPlayer.Play()
					} else if game.levelOneEnemyList[i].collision == true && spriteWidth != 50 {
						g.humanEnemyDeathAudioPlayer.Rewind()
						g.humanEnemyDeathAudioPlayer.Play()
					}
				}
			} else {
				if game.levelOneEnemyList[i].health == 2 && game.levelOneEnemyList[i].collision == true {
					game.score += 300
					game.levelOneEnemyList[i].health = 0
				}
				if game.levelOneEnemyList[i].health == 1 && game.levelOneEnemyList[i].collision == true {
					game.score += 200
					game.levelOneEnemyList[i].health = 0
				}
				game.levelOneEnemyList[i].dx = 0
				game.levelOneEnemyList[i].dy = 0
			}
		}
	}

	//player projectile collides with wall check
	if len(game.projectileList) > 0 {
		for i := 0; i < len(game.projectileList); i++ {
			if game.projectileList[i].collision == false {
				game.projectileList[i].xLoc += game.projectileList[i].dx
				game.projectileList[i].yLoc += game.projectileList[i].dy
				game.projectileList[i].collision = wallCollisionCheckFirstLevel(game.projectileList[i], 20)
			}
		}
	}

	//enemy projectile collides with wall check
	if len(game.levelOneEnemyList) > 0 {
		for i := 0; i < len(game.levelOneEnemyList); i++ {
			if len(game.levelOneEnemyList[i].enemyProjectileList) > 0 {
				for j := 0; j < len(game.levelOneEnemyList[i].enemyProjectileList); j++ {
					if game.levelOneEnemyList[i].enemyProjectileList[j].collision == false {
						game.levelOneEnemyList[i].enemyProjectileList[j].xLoc += game.levelOneEnemyList[i].enemyProjectileList[j].dx
						game.levelOneEnemyList[i].enemyProjectileList[j].yLoc += game.levelOneEnemyList[i].enemyProjectileList[j].dy
						game.levelOneEnemyList[i].enemyProjectileList[j].collision =
							wallCollisionCheckFirstLevel(game.levelOneEnemyList[i].enemyProjectileList[j], 20)
					}
				}
			}
		}
	}

	//player collides with enemy check
	if len(game.levelOneEnemyList) > 0 {
		for i := 0; i < len(game.levelOneEnemyList); i++ {
			if game.levelOneEnemyList[i].collision == false {
				enemyWidth, _ := game.levelOneEnemyList[i].leftPict.Size()
				playerWidth, _ := game.playerSprite.upPict.Size()
				death := playerCollisionWithEnemy(game.levelOneEnemyList[i], game.playerSprite, enemyWidth, playerWidth)
				if death == 1 {
					g.enemyAndPlayerCollisionAudioPlayer.Rewind()
					g.enemyAndPlayerCollisionAudioPlayer.Play()
					game.playerSprite.xLoc, game.playerSprite.yLoc = 190, ScreenHeight*0.72
					game.deathCounter += death
				}
			}
		}
	}

	//enemy projectile collides with player check
	if len(game.levelOneEnemyList) > 0 {
		for i := 0; i < len(game.levelOneEnemyList); i++ {
			if len(game.levelOneEnemyList[i].enemyProjectileList) > 0 {
				for j := 0; j < len(game.levelOneEnemyList[i].enemyProjectileList); j++ {
					if game.levelOneEnemyList[i].enemyProjectileList[j].collision == false {
						death := 0
						game.levelOneEnemyList[i].enemyProjectileList[j].collision, death =
							projectileCollisionWithPlayer(game.playerSprite,
								game.levelOneEnemyList[i].enemyProjectileList[j], 61, 20)
						if death == 1 {
							g.playerDeathAudioPlayer.Rewind()
							g.playerDeathAudioPlayer.Play()
							game.playerSprite.xLoc, game.playerSprite.yLoc = 190, ScreenHeight*0.72
							game.deathCounter += death
						}

					}
				}
			}
		}
	}

	//player projectile collides with enemy check
	if len(game.projectileList) > 0 && len(game.levelOneEnemyList) > 0 {
		for i := 0; i < len(game.projectileList); i++ {
			for j := 0; j < len(game.levelOneEnemyList); j++ {
				enemyWidth, _ := game.levelOneEnemyList[j].upPict.Size()
				if game.levelOneEnemyList[j].collision == false && game.projectileList[i].collision == false {
					additionalScore := 0
					game.levelOneEnemyList[j].collision, game.projectileList[i].collision, game.levelOneEnemyList[j].health, additionalScore =
						projectileCollisionWithEnemy(game.levelOneEnemyList[j], game.projectileList[i], enemyWidth, 20)
					game.score += additionalScore
				}
			}
		}
	}
}

func (game *Game) manageLevel2CollisionDetection() {
	if game.collectedGold == false {
		game.collectedGold = game.gotGold(game.playerSprite, game.coinSprite)
	}

	//player collision with wall check
	if game.playerAndWallCollision == false {
		game.playerAndWallCollision = wallCollisionCheckSecondLevel(game.playerSprite, 61)
	} else {
		game.playerSprite.xLoc, game.playerSprite.yLoc = 100, 100
		game.playerAndWallCollision = false
		g.playerDeathAudioPlayer.Rewind()
		g.playerDeathAudioPlayer.Play()
		game.deathCounter += 1
	}

	//enemy collision with wall check
	if len(game.levelTwoEnemyList) > 0 {
		for i := 0; i < len(game.levelTwoEnemyList); i++ {
			if game.levelTwoEnemyList[i].collision == false {
				if game.levelTwoEnemyList[i].direction == "left" {
					spriteWidth, _ := game.levelTwoEnemyList[i].leftPict.Size()
					game.levelTwoEnemyList[i].collision = wallCollisionCheckSecondLevel(game.levelTwoEnemyList[i], spriteWidth)
					if game.levelTwoEnemyList[i].collision == true && spriteWidth == 50 {
						g.monsterEnemyDeathAudioPlayer.Rewind()
						g.monsterEnemyDeathAudioPlayer.Play()
					} else if game.levelTwoEnemyList[i].collision == true && spriteWidth != 50 {
						g.humanEnemyDeathAudioPlayer.Rewind()
						g.humanEnemyDeathAudioPlayer.Play()
					}
				} else if game.levelTwoEnemyList[i].direction == "right" {
					spriteWidth, _ := game.levelTwoEnemyList[i].rightPict.Size()
					game.levelTwoEnemyList[i].collision = wallCollisionCheckSecondLevel(game.levelTwoEnemyList[i], spriteWidth)
					if game.levelTwoEnemyList[i].collision == true && spriteWidth == 50 {
						g.monsterEnemyDeathAudioPlayer.Rewind()
						g.monsterEnemyDeathAudioPlayer.Play()
					} else if game.levelTwoEnemyList[i].collision == true && spriteWidth != 50 {
						g.humanEnemyDeathAudioPlayer.Rewind()
						g.humanEnemyDeathAudioPlayer.Play()
					}
				} else if game.levelTwoEnemyList[i].direction == "up" {
					spriteWidth, _ := game.levelTwoEnemyList[i].upPict.Size()
					game.levelTwoEnemyList[i].collision = wallCollisionCheckSecondLevel(game.levelTwoEnemyList[i], spriteWidth)
					if game.levelTwoEnemyList[i].collision == true && spriteWidth == 50 {
						g.monsterEnemyDeathAudioPlayer.Rewind()
						g.monsterEnemyDeathAudioPlayer.Play()
					} else if game.levelTwoEnemyList[i].collision == true && spriteWidth != 50 {
						g.humanEnemyDeathAudioPlayer.Rewind()
						g.humanEnemyDeathAudioPlayer.Play()
					}
				} else if game.levelTwoEnemyList[i].direction == "down" {
					spriteWidth, _ := game.levelTwoEnemyList[i].downPict.Size()
					game.levelTwoEnemyList[i].collision = wallCollisionCheckSecondLevel(game.levelTwoEnemyList[i], spriteWidth)
					if game.levelTwoEnemyList[i].collision == true && spriteWidth == 50 {
						g.monsterEnemyDeathAudioPlayer.Rewind()
						g.monsterEnemyDeathAudioPlayer.Play()
					} else if game.levelTwoEnemyList[i].collision == true && spriteWidth != 50 {
						g.humanEnemyDeathAudioPlayer.Rewind()
						g.humanEnemyDeathAudioPlayer.Play()
					}
				}
			} else {
				if game.levelTwoEnemyList[i].health == 2 && game.levelTwoEnemyList[i].collision == true {
					game.score += 300
					game.levelTwoEnemyList[i].health = 0
				}
				if game.levelTwoEnemyList[i].health == 1 && game.levelTwoEnemyList[i].collision == true {
					game.score += 200
					game.levelTwoEnemyList[i].health = 0
				}
				game.levelTwoEnemyList[i].dx = 0
				game.levelTwoEnemyList[i].dy = 0
			}
		}
	}

	//player projectile collides with wall check
	if len(game.projectileList) > 0 {
		for i := 0; i < len(game.projectileList); i++ {
			if game.projectileList[i].collision == false {
				game.projectileList[i].xLoc += game.projectileList[i].dx
				game.projectileList[i].yLoc += game.projectileList[i].dy
				game.projectileList[i].collision = wallCollisionCheckSecondLevel(game.projectileList[i], 20)
			}
		}
	}

	//enemy projectile collides with wall check
	if len(game.levelTwoEnemyList) > 0 {
		for i := 0; i < len(game.levelTwoEnemyList); i++ {
			if len(game.levelTwoEnemyList[i].enemyProjectileList) > 0 {
				for j := 0; j < len(game.levelTwoEnemyList[i].enemyProjectileList); j++ {
					if game.levelTwoEnemyList[i].enemyProjectileList[j].collision == false {
						game.levelTwoEnemyList[i].enemyProjectileList[j].xLoc += game.levelTwoEnemyList[i].enemyProjectileList[j].dx
						game.levelTwoEnemyList[i].enemyProjectileList[j].yLoc += game.levelTwoEnemyList[i].enemyProjectileList[j].dy
						game.levelTwoEnemyList[i].enemyProjectileList[j].collision =
							wallCollisionCheckSecondLevel(game.levelTwoEnemyList[i].enemyProjectileList[j], 20)
					}
				}
			}
		}
	}

	//player collides with enemy check
	if len(game.levelTwoEnemyList) > 0 {
		for i := 0; i < len(game.levelTwoEnemyList); i++ {
			if game.levelTwoEnemyList[i].collision == false {
				enemyWidth, _ := game.levelTwoEnemyList[i].leftPict.Size()
				playerWidth, _ := game.playerSprite.upPict.Size()
				death := playerCollisionWithEnemy(game.levelTwoEnemyList[i], game.playerSprite, enemyWidth, playerWidth)
				if death == 1 {
					g.enemyAndPlayerCollisionAudioPlayer.Rewind()
					g.enemyAndPlayerCollisionAudioPlayer.Play()
					game.playerSprite.xLoc, game.playerSprite.yLoc = 100, 100
					game.deathCounter += death
				}
			}
		}
	}

	//enemy projectile collides with player check
	if len(game.levelTwoEnemyList) > 0 {
		for i := 0; i < len(game.levelTwoEnemyList); i++ {
			if len(game.levelTwoEnemyList[i].enemyProjectileList) > 0 {
				for j := 0; j < len(game.levelTwoEnemyList[i].enemyProjectileList); j++ {
					if game.levelTwoEnemyList[i].enemyProjectileList[j].collision == false {
						death := 0
						game.levelTwoEnemyList[i].enemyProjectileList[j].collision, death =
							projectileCollisionWithPlayer(game.playerSprite,
								game.levelTwoEnemyList[i].enemyProjectileList[j], 61, 20)
						if death == 1 {
							g.playerDeathAudioPlayer.Rewind()
							g.playerDeathAudioPlayer.Play()
							game.playerSprite.xLoc, game.playerSprite.yLoc = 100, 100
							game.deathCounter += death
						}

					}
				}
			}
		}
	}

	//player projectile collides with enemy check
	if len(game.projectileList) > 0 && len(game.levelTwoEnemyList) > 0 {
		for i := 0; i < len(game.projectileList); i++ {
			for j := 0; j < len(game.levelTwoEnemyList); j++ {
				enemyWidth, _ := game.levelTwoEnemyList[j].upPict.Size()
				if game.levelTwoEnemyList[j].collision == false && game.projectileList[i].collision == false {
					additionalScore := 0
					game.levelTwoEnemyList[j].collision, game.projectileList[i].collision, game.levelTwoEnemyList[j].health, additionalScore =
						projectileCollisionWithEnemy(game.levelTwoEnemyList[j], game.projectileList[i], enemyWidth, 20)
					game.score += additionalScore
				}
			}
		}
	}
}

func (game *Game) manageLevel3CollisionDetection() {
	if game.collectedGold == false {
		game.collectedGold = game.gotGold(game.playerSprite, game.coinSprite)
	}

	//player collision with wall check
	if game.playerAndWallCollision == false {
		game.playerAndWallCollision = wallCollisionCheckThirdLevel(game.playerSprite, 61)
	} else {
		game.playerSprite.xLoc, game.playerSprite.yLoc = 600, 100
		game.playerAndWallCollision = false
		g.playerDeathAudioPlayer.Rewind()
		g.playerDeathAudioPlayer.Play()
		game.deathCounter += 1
	}

	//enemy collision with wall check
	if len(game.levelThreeEnemyList) > 0 {
		for i := 0; i < len(game.levelThreeEnemyList); i++ {
			if game.levelThreeEnemyList[i].collision == false {
				if game.levelThreeEnemyList[i].direction == "left" {
					spriteWidth, _ := game.levelThreeEnemyList[i].leftPict.Size()
					game.levelThreeEnemyList[i].collision = wallCollisionCheckThirdLevel(game.levelThreeEnemyList[i], spriteWidth)
					if game.levelThreeEnemyList[i].collision == true && spriteWidth == 50 {
						g.monsterEnemyDeathAudioPlayer.Rewind()
						g.monsterEnemyDeathAudioPlayer.Play()
					} else if game.levelThreeEnemyList[i].collision == true && spriteWidth != 50 {
						g.humanEnemyDeathAudioPlayer.Rewind()
						g.humanEnemyDeathAudioPlayer.Play()
					}
				} else if game.levelThreeEnemyList[i].direction == "right" {
					spriteWidth, _ := game.levelThreeEnemyList[i].rightPict.Size()
					game.levelThreeEnemyList[i].collision = wallCollisionCheckThirdLevel(game.levelThreeEnemyList[i], spriteWidth)
					if game.levelThreeEnemyList[i].collision == true && spriteWidth == 50 {
						g.monsterEnemyDeathAudioPlayer.Rewind()
						g.monsterEnemyDeathAudioPlayer.Play()
					} else if game.levelThreeEnemyList[i].collision == true && spriteWidth != 50 {
						g.humanEnemyDeathAudioPlayer.Rewind()
						g.humanEnemyDeathAudioPlayer.Play()
					}
				} else if game.levelThreeEnemyList[i].direction == "up" {
					spriteWidth, _ := game.levelThreeEnemyList[i].upPict.Size()
					game.levelThreeEnemyList[i].collision = wallCollisionCheckThirdLevel(game.levelThreeEnemyList[i], spriteWidth)
					if game.levelThreeEnemyList[i].collision == true && spriteWidth == 50 {
						g.monsterEnemyDeathAudioPlayer.Rewind()
						g.monsterEnemyDeathAudioPlayer.Play()
					} else if game.levelThreeEnemyList[i].collision == true && spriteWidth != 50 {
						g.humanEnemyDeathAudioPlayer.Rewind()
						g.humanEnemyDeathAudioPlayer.Play()
					}
				} else if game.levelThreeEnemyList[i].direction == "down" {
					spriteWidth, _ := game.levelThreeEnemyList[i].downPict.Size()
					game.levelThreeEnemyList[i].collision = wallCollisionCheckThirdLevel(game.levelThreeEnemyList[i], spriteWidth)
					if game.levelThreeEnemyList[i].collision == true && spriteWidth == 50 {
						g.monsterEnemyDeathAudioPlayer.Rewind()
						g.monsterEnemyDeathAudioPlayer.Play()
					} else if game.levelThreeEnemyList[i].collision == true && spriteWidth != 50 {
						g.humanEnemyDeathAudioPlayer.Rewind()
						g.humanEnemyDeathAudioPlayer.Play()
					}
				}
			} else {
				if game.levelThreeEnemyList[i].health == 2 && game.levelThreeEnemyList[i].collision == true {
					game.score += 300
					game.levelThreeEnemyList[i].health = 0
				}
				if game.levelThreeEnemyList[i].health == 1 && game.levelThreeEnemyList[i].collision == true {
					game.score += 200
					game.levelThreeEnemyList[i].health = 0
				}
				game.levelThreeEnemyList[i].dx = 0
				game.levelThreeEnemyList[i].dy = 0
			}
		}
	}

	//player projectile collides with wall check
	if len(game.projectileList) > 0 {
		for i := 0; i < len(game.projectileList); i++ {
			if game.projectileList[i].collision == false {
				game.projectileList[i].xLoc += game.projectileList[i].dx
				game.projectileList[i].yLoc += game.projectileList[i].dy
				game.projectileList[i].collision = wallCollisionCheckThirdLevel(game.projectileList[i], 20)
			}
		}
	}

	//enemy projectile collides with wall check
	if len(game.levelThreeEnemyList) > 0 {
		for i := 0; i < len(game.levelThreeEnemyList); i++ {
			if len(game.levelThreeEnemyList[i].enemyProjectileList) > 0 {
				for j := 0; j < len(game.levelThreeEnemyList[i].enemyProjectileList); j++ {
					if game.levelThreeEnemyList[i].enemyProjectileList[j].collision == false {
						game.levelThreeEnemyList[i].enemyProjectileList[j].xLoc += game.levelThreeEnemyList[i].enemyProjectileList[j].dx
						game.levelThreeEnemyList[i].enemyProjectileList[j].yLoc += game.levelThreeEnemyList[i].enemyProjectileList[j].dy
						game.levelThreeEnemyList[i].enemyProjectileList[j].collision =
							wallCollisionCheckThirdLevel(game.levelThreeEnemyList[i].enemyProjectileList[j], 20)
					}
				}
			}
		}
	}

	//player collides with enemy check
	if len(game.levelThreeEnemyList) > 0 {
		for i := 0; i < len(game.levelThreeEnemyList); i++ {
			if game.levelThreeEnemyList[i].collision == false {
				enemyWidth, _ := game.levelThreeEnemyList[i].leftPict.Size()
				playerWidth, _ := game.playerSprite.upPict.Size()
				death := playerCollisionWithEnemy(game.levelThreeEnemyList[i], game.playerSprite, enemyWidth, playerWidth)
				if death == 1 {
					g.enemyAndPlayerCollisionAudioPlayer.Rewind()
					g.enemyAndPlayerCollisionAudioPlayer.Play()
					game.playerSprite.xLoc, game.playerSprite.yLoc = 600, 100
					game.deathCounter += death
				}
			}
		}
	}

	//enemy projectile collides with player check
	if len(game.levelThreeEnemyList) > 0 {
		for i := 0; i < len(game.levelThreeEnemyList); i++ {
			if len(game.levelThreeEnemyList[i].enemyProjectileList) > 0 {
				for j := 0; j < len(game.levelThreeEnemyList[i].enemyProjectileList); j++ {
					if game.levelThreeEnemyList[i].enemyProjectileList[j].collision == false {
						death := 0
						game.levelThreeEnemyList[i].enemyProjectileList[j].collision, death =
							projectileCollisionWithPlayer(game.playerSprite,
								game.levelThreeEnemyList[i].enemyProjectileList[j], 61, 20)
						if death == 1 {
							game.playerSprite.xLoc, game.playerSprite.yLoc = 600, 100
							g.playerDeathAudioPlayer.Rewind()
							g.playerDeathAudioPlayer.Play()
							game.deathCounter += death
						}

					}
				}
			}
		}
	}

	//player projectile collides with enemy check
	if len(game.projectileList) > 0 && len(game.levelThreeEnemyList) > 0 {
		for i := 0; i < len(game.projectileList); i++ {
			for j := 0; j < len(game.levelThreeEnemyList); j++ {
				enemyWidth, _ := game.levelThreeEnemyList[j].upPict.Size()
				if game.levelThreeEnemyList[j].collision == false && game.projectileList[i].collision == false {
					additionalScore := 0
					game.levelThreeEnemyList[j].collision, game.projectileList[i].collision, game.levelThreeEnemyList[j].health, additionalScore =
						projectileCollisionWithEnemy(game.levelThreeEnemyList[j], game.projectileList[i], enemyWidth, 20)
					game.score += additionalScore
				}
			}
		}
	}
}

func (game *Game) checkLevel() {
	if game.gameOver == false {
		if game.score < 1000 {
			game.levelOneIsActive = true
			game.levelTwoIsActive = false
			game.levelThreeIsActive = false
		} else if game.score >= 1000 && game.score < 2000 && game.levelOneIsActive == true {
			game.levelOneIsActive = false
			game.levelTwoIsActive = true
			game.levelThreeIsActive = false
			game.playerSprite.xLoc, game.playerSprite.yLoc = 100, 100
		} else if game.score >= 1000 && game.score < 2000 && game.levelOneIsActive == false {
			game.levelOneIsActive = false
			game.levelTwoIsActive = true
			game.levelThreeIsActive = false
		} else if game.score >= 2000 && game.levelTwoIsActive == true {
			game.levelOneIsActive = false
			game.levelTwoIsActive = false
			game.levelThreeIsActive = true
			game.playerSprite.xLoc, game.playerSprite.yLoc = 600, 100
		} else if game.score >= 2000 && game.score < 3000 && game.levelThreeIsActive == true {
			game.levelOneIsActive = false
			game.levelTwoIsActive = false
			game.levelThreeIsActive = true
		} else if game.score >= 3000 && game.gameWon == false {
			game.levelOneIsActive = false
			game.levelTwoIsActive = false
			game.levelThreeIsActive = false
			game.gameWon = true
		} else {
			game.gameWon = true
		}
	} else {
		game.levelOneIsActive = true
		game.levelTwoIsActive = false
		game.levelThreeIsActive = false
	}
}

func (game *Game) Update() error {
	game.checkLevel()

	if game.deathCounter >= 3 && game.gameWon == false {
		if game.playedLoseSound == false {
			game.playedLoseSound = true
			g.loseAudioPlayer.Rewind()
			g.loseAudioPlayer.Play()
		}
		game.gameOver = true
	} else {
		game.gameOver = false
	}
	if game.score >= 3000 {
		if game.playedWinSound == false {
			game.playedWinSound = true
			g.winAudioPlayer.Rewind()
			g.winAudioPlayer.Play()
		}
		game.gameWon = true
		game.gameOver = false
	}

	if game.startGame == false {
		game.getUserName()
	} else if game.startGame == true && game.levelOneIsActive == true && game.gameOver == false {
		game.spawnLevel1Enemies()
		game.movementLevel1Enemies()
		game.changeTankDirection()
		game.changeTankTopperDirection()
		game.playerShootFireball()
		game.manageTankTopperOffset()
		game.manageLevel1CollisionDetection()
	} else if game.startGame == true && game.levelTwoIsActive == true && game.gameOver == false {
		game.spawnLevel2Enemies()
		game.movementLevel2Enemies()
		game.changeTankDirection()
		game.changeTankTopperDirection()
		game.playerShootFireball()
		game.manageTankTopperOffset()
		game.manageLevel2CollisionDetection()
	} else if game.startGame == true && game.levelThreeIsActive == true && game.gameOver == false && game.gameWon == false {
		game.spawnLevel3Enemies()
		game.movementLevel3Enemies()
		game.changeTankDirection()
		game.changeTankTopperDirection()
		game.playerShootFireball()
		game.manageTankTopperOffset()
		game.manageLevel3CollisionDetection()
	} else if game.startGame == true && game.gameOver == true && game.dbEntryComplete == false {
		myDatabase := OpenDataBase("./LeaderBoard.db")
		create_tables(myDatabase)
		game.addGameEntry(myDatabase)
		game.dbEntryComplete = true
		game.allScores = true
		myDatabase.Close()
	} else if game.startGame == true && game.gameWon == true && game.dbEntryComplete == false {
		myDatabase := OpenDataBase("./LeaderBoard.db")
		create_tables(myDatabase)
		game.addGameEntry(myDatabase)
		game.dbEntryComplete = true
		game.allScores = true
		myDatabase.Close()
	} else if game.startGame == true && game.gameWon == true && game.dbEntryComplete == true && game.processedDB == false {
		game.processDBtoMaps()
		game.processedDB = true
	} else if game.startGame == true && game.gameOver == true && game.dbEntryComplete == true && game.processedDB == false {
		game.processDBtoMaps()
		game.processedDB = true
	} else if game.startGame == true && game.gameWon == true && game.dbEntryComplete == true && game.processedDB == true {
		game.getLeaderBoardFormat()
	} else if game.startGame == true && game.gameOver == true && game.dbEntryComplete == true && game.processedDB == true {
		game.getLeaderBoardFormat()
	} else {
		game.spawnLevel1Enemies()
		game.changeTankDirection()
		game.changeTankTopperDirection()
		game.playerShootFireball()
		game.manageTankTopperOffset()
		game.manageLevel1CollisionDetection()
	}
	return nil
}

func (game Game) Draw(screen *ebiten.Image) {
	if game.startGame == false {
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.titleScreenBackground.xLoc), float64(game.titleScreenBackground.yLoc))
		screen.DrawImage(game.titleScreenBackground.upPict, &game.drawOps)

		game.drawOps.GeoM.Reset()
		text.Draw(screen, "Enter Username: ", mplusNormalFont, ScreenWidth*0.20, ScreenHeight*0.25, colornames.White)

		if len(game.userNameList) > 0 {
			game.iterateAndStoreUserName()
			game.drawOps.GeoM.Reset()
			text.Draw(screen, "Enter Username: "+game.userName, mplusNormalFont, ScreenWidth*0.20, ScreenHeight*0.25, colornames.White)
			game.drawOps.GeoM.Reset()
			text.Draw(screen, "Press ENTER to start Berserk/Tank game.", mplusNormalFont, ScreenWidth*0.20, ScreenHeight*0.45, color.Black)
		}
	}
	if game.startGame == true && game.gameOver == false && game.gameWon == false {

		if game.levelOneIsActive {
			game.drawOps.GeoM.Reset()
			game.drawOps.GeoM.Translate(float64(game.firstMap.xLoc), float64(game.firstMap.yLoc))
			screen.DrawImage(game.firstMap.upPict, &game.drawOps)
			game.drawOps.GeoM.Reset()
			text.Draw(screen, "Score: "+strconv.Itoa(game.score), mplusNormalFont, ScreenWidth*0.77, ScreenHeight*0.08, colornames.White)

			if len(game.levelOneEnemyList) > 0 {
				for i := 0; i < len(game.levelOneEnemyList); i++ {
					if game.levelOneEnemyList[i].collision == false {
						game.drawOps.GeoM.Reset()
						game.drawOps.GeoM.Translate(float64(game.levelOneEnemyList[i].xLoc), float64(game.levelOneEnemyList[i].yLoc))
						if game.levelOneEnemyList[i].direction == "left" {
							screen.DrawImage(game.levelOneEnemyList[i].leftPict, &game.drawOps)
						} else if game.levelOneEnemyList[i].direction == "right" {
							screen.DrawImage(game.levelOneEnemyList[i].rightPict, &game.drawOps)
						} else if game.levelOneEnemyList[i].direction == "up" {
							screen.DrawImage(game.levelOneEnemyList[i].upPict, &game.drawOps)
						} else if game.levelOneEnemyList[i].direction == "down" {
							screen.DrawImage(game.levelOneEnemyList[i].downPict, &game.drawOps)
						} else {
							screen.DrawImage(game.levelOneEnemyList[i].upPict, &game.drawOps)
						}
					}
				}
			}

			if len(game.levelOneEnemyList) > 0 {
				for i := 0; i < len(game.levelOneEnemyList); i++ {
					if len(game.levelOneEnemyList[i].enemyProjectileList) > 0 {
						for j := 0; j < len(game.levelOneEnemyList[i].enemyProjectileList); j++ {
							if game.levelOneEnemyList[i].enemyProjectileList[j].collision == false {
								game.drawOps.GeoM.Reset()
								game.drawOps.GeoM.Translate(float64(game.levelOneEnemyList[i].enemyProjectileList[j].xLoc),
									float64(game.levelOneEnemyList[i].enemyProjectileList[j].yLoc))
								screen.DrawImage(game.levelOneEnemyList[i].enemyProjectileList[j].upPict, &game.drawOps)
							}
						}
					}
				}
			}

			if len(game.projectileList) > 0 {
				for i := 0; i < len(game.projectileList); i++ {
					if game.projectileList[i].collision == false {
						game.drawOps.GeoM.Reset()
						game.drawOps.GeoM.Translate(float64(game.projectileList[i].xLoc), float64(game.projectileList[i].yLoc))
						screen.DrawImage(game.projectileList[i].upPict, &game.drawOps)
					}
				}
			}
		} else if game.levelTwoIsActive {
			game.drawOps.GeoM.Reset()
			game.drawOps.GeoM.Translate(float64(game.secondMap.xLoc), float64(game.secondMap.yLoc))
			screen.DrawImage(game.secondMap.upPict, &game.drawOps)

			game.drawOps.GeoM.Reset()
			text.Draw(screen, "Score: "+strconv.Itoa(game.score), mplusNormalFont, ScreenWidth*0.77, ScreenHeight*0.08, colornames.White)

			if len(game.levelTwoEnemyList) > 0 {
				for i := 0; i < len(game.levelTwoEnemyList); i++ {
					if game.levelTwoEnemyList[i].collision == false {
						game.drawOps.GeoM.Reset()
						game.drawOps.GeoM.Translate(float64(game.levelTwoEnemyList[i].xLoc), float64(game.levelTwoEnemyList[i].yLoc))
						if game.levelTwoEnemyList[i].direction == "left" {
							screen.DrawImage(game.levelTwoEnemyList[i].leftPict, &game.drawOps)
						} else if game.levelTwoEnemyList[i].direction == "right" {
							screen.DrawImage(game.levelTwoEnemyList[i].rightPict, &game.drawOps)
						} else if game.levelTwoEnemyList[i].direction == "up" {
							screen.DrawImage(game.levelTwoEnemyList[i].upPict, &game.drawOps)
						} else if game.levelTwoEnemyList[i].direction == "down" {
							screen.DrawImage(game.levelTwoEnemyList[i].downPict, &game.drawOps)
						} else {
							screen.DrawImage(game.levelTwoEnemyList[i].upPict, &game.drawOps)
						}
					}
				}
			}

			if len(game.levelTwoEnemyList) > 0 {
				for i := 0; i < len(game.levelTwoEnemyList); i++ {
					if len(game.levelTwoEnemyList[i].enemyProjectileList) > 0 {
						for j := 0; j < len(game.levelTwoEnemyList[i].enemyProjectileList); j++ {
							if game.levelTwoEnemyList[i].enemyProjectileList[j].collision == false {
								game.drawOps.GeoM.Reset()
								game.drawOps.GeoM.Translate(float64(game.levelTwoEnemyList[i].enemyProjectileList[j].xLoc),
									float64(game.levelTwoEnemyList[i].enemyProjectileList[j].yLoc))
								screen.DrawImage(game.levelTwoEnemyList[i].enemyProjectileList[j].upPict, &game.drawOps)
							}
						}
					}
				}
			}

			if len(game.projectileList) > 0 {
				for i := 0; i < len(game.projectileList); i++ {
					if game.projectileList[i].collision == false {
						game.drawOps.GeoM.Reset()
						game.drawOps.GeoM.Translate(float64(game.projectileList[i].xLoc), float64(game.projectileList[i].yLoc))
						screen.DrawImage(game.projectileList[i].upPict, &game.drawOps)
					}
				}
			}
		} else if game.levelThreeIsActive {
			game.drawOps.GeoM.Reset()
			game.drawOps.GeoM.Translate(float64(game.thirdMap.xLoc), float64(game.thirdMap.yLoc))
			screen.DrawImage(game.thirdMap.upPict, &game.drawOps)

			game.drawOps.GeoM.Reset()
			text.Draw(screen, "Score: "+strconv.Itoa(game.score), mplusNormalFont, ScreenWidth*0.77, ScreenHeight*0.08, colornames.White)

			if len(game.levelThreeEnemyList) > 0 {
				for i := 0; i < len(game.levelThreeEnemyList); i++ {
					if game.levelThreeEnemyList[i].collision == false {
						game.drawOps.GeoM.Reset()
						game.drawOps.GeoM.Translate(float64(game.levelThreeEnemyList[i].xLoc), float64(game.levelThreeEnemyList[i].yLoc))
						if game.levelThreeEnemyList[i].direction == "left" {
							screen.DrawImage(game.levelThreeEnemyList[i].leftPict, &game.drawOps)
						} else if game.levelThreeEnemyList[i].direction == "right" {
							screen.DrawImage(game.levelThreeEnemyList[i].rightPict, &game.drawOps)
						} else if game.levelThreeEnemyList[i].direction == "up" {
							screen.DrawImage(game.levelThreeEnemyList[i].upPict, &game.drawOps)
						} else if game.levelThreeEnemyList[i].direction == "down" {
							screen.DrawImage(game.levelThreeEnemyList[i].downPict, &game.drawOps)
						} else {
							screen.DrawImage(game.levelThreeEnemyList[i].upPict, &game.drawOps)
						}
					}
				}
			}

			if len(game.levelThreeEnemyList) > 0 {
				for i := 0; i < len(game.levelThreeEnemyList); i++ {
					if len(game.levelThreeEnemyList[i].enemyProjectileList) > 0 {
						for j := 0; j < len(game.levelThreeEnemyList[i].enemyProjectileList); j++ {
							if game.levelThreeEnemyList[i].enemyProjectileList[j].collision == false {
								game.drawOps.GeoM.Reset()
								game.drawOps.GeoM.Translate(float64(game.levelThreeEnemyList[i].enemyProjectileList[j].xLoc),
									float64(game.levelThreeEnemyList[i].enemyProjectileList[j].yLoc))
								screen.DrawImage(game.levelThreeEnemyList[i].enemyProjectileList[j].upPict, &game.drawOps)
							}
						}
					}
				}
			}
		}

		if len(game.projectileList) > 0 {
			for i := 0; i < len(game.projectileList); i++ {
				if game.projectileList[i].collision == false {
					game.drawOps.GeoM.Reset()
					game.drawOps.GeoM.Translate(float64(game.projectileList[i].xLoc), float64(game.projectileList[i].yLoc))
					screen.DrawImage(game.projectileList[i].upPict, &game.drawOps)
				}
			}
		}
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
	} else if game.startGame == true && game.gameOver == true && game.gameWon == false && game.processedDB == true {
		tempHeight := 150
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.loserScreen.xLoc), float64(game.loserScreen.yLoc))
		screen.DrawImage(game.loserScreen.upPict, &game.drawOps)
		game.drawOps.GeoM.Reset()
		text.Draw(screen, "LEADERBOARD", mplusNormalFont, ScreenWidth*0.40, ScreenHeight*0.08, colornames.White)
		if game.allScores == true && game.playerScores == false {
			game.drawOps.GeoM.Reset()
			text.Draw(screen, "Press SPACE to switch to your top 5 scores.", mplusNormalFont, ScreenWidth*0.17, ScreenHeight*0.90, colornames.White)
			if len(userNameMap) > 0 {
				for i := 0; i < len(userNameMap) && i < 5; i++ {
					if (game.currentPlayerAndScoreLeaderboard == false) && (userNameMap[i][0] ==
						game.userName) && (scoreMap[i][0] == game.score) {
						game.drawOps.GeoM.Reset()
						text.Draw(screen, strconv.Itoa(i+1)+". "+userNameMap[i][0]+": "+strconv.Itoa(scoreMap[i][0]), mplusNormalFont, ScreenWidth*0.15, tempHeight, colornames.Red)
						tempHeight += 100
						game.currentPlayerAndScoreLeaderboard = true
					} else {
						game.drawOps.GeoM.Reset()
						text.Draw(screen, strconv.Itoa(i+1)+". "+userNameMap[i][0]+": "+strconv.Itoa(scoreMap[i][0]), mplusNormalFont, ScreenWidth*0.15, tempHeight, colornames.White)
						tempHeight += 100
					}
				}
			}
		} else if game.playerScores == true && game.allScores == false {
			game.drawOps.GeoM.Reset()
			text.Draw(screen, "Press SPACE to switch to all players top 5 scores.", mplusNormalFont, ScreenWidth*0.15, ScreenHeight*0.90, colornames.White)
			if len(currentPlayerMap) > 0 {
				for i := 0; i < len(currentPlayerMap) && i < 5; i++ {
					if (game.currentPlayerAndScoreLeaderboard == false) && (currentPlayerMap[i][0] ==
						game.userName) && (currentPlayerScoreMap[i][0] == game.score) {
						game.drawOps.GeoM.Reset()
						text.Draw(screen, strconv.Itoa(i+1)+". "+currentPlayerMap[i][0]+": "+strconv.Itoa(currentPlayerScoreMap[i][0]), mplusNormalFont, ScreenWidth*0.15, tempHeight, colornames.Red)
						tempHeight += 100
						game.currentPlayerAndScoreLeaderboard = true

					} else {
						game.drawOps.GeoM.Reset()
						text.Draw(screen, strconv.Itoa(i+1)+". "+currentPlayerMap[i][0]+": "+strconv.Itoa(currentPlayerScoreMap[i][0]), mplusNormalFont, ScreenWidth*0.15, tempHeight, colornames.White)
						tempHeight += 100
					}
				}
			}
		}
	} else if game.startGame == true && game.gameOver == false && game.gameWon == true && game.processedDB == true {
		tempHeight := 150
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.winnerScreen.xLoc), float64(game.winnerScreen.yLoc))
		screen.DrawImage(game.winnerScreen.upPict, &game.drawOps)
		game.drawOps.GeoM.Reset()
		text.Draw(screen, "LEADERBOARD", mplusNormalFont, ScreenWidth*0.40, ScreenHeight*0.08, colornames.White)
		if game.allScores == true && game.playerScores == false {
			game.drawOps.GeoM.Reset()
			text.Draw(screen, "Press SPACE to switch to your top 5 scores.", mplusNormalFont, ScreenWidth*0.17, ScreenHeight*0.90, colornames.White)
			if len(userNameMap) > 0 {
				for i := 0; i < len(userNameMap) && i < 5; i++ {
					if (game.currentPlayerAndScoreLeaderboard == false) && (userNameMap[i][0] ==
						game.userName) && (scoreMap[i][0] == game.score) {
						game.drawOps.GeoM.Reset()
						text.Draw(screen, strconv.Itoa(i+1)+". "+userNameMap[i][0]+": "+strconv.Itoa(scoreMap[i][0]), mplusNormalFont, ScreenWidth*0.15, tempHeight, colornames.Red)
						tempHeight += 100
						game.currentPlayerAndScoreLeaderboard = true
					} else {
						game.drawOps.GeoM.Reset()
						text.Draw(screen, strconv.Itoa(i+1)+". "+userNameMap[i][0]+": "+strconv.Itoa(scoreMap[i][0]), mplusNormalFont, ScreenWidth*0.15, tempHeight, colornames.White)
						tempHeight += 100
					}
				}
			}
		} else if game.playerScores == true && game.allScores == false {
			game.drawOps.GeoM.Reset()
			text.Draw(screen, "Press SPACE to switch to all players top 5 scores.", mplusNormalFont, ScreenWidth*0.15, ScreenHeight*0.90, colornames.White)
			if len(currentPlayerMap) > 0 {
				for i := 0; i < len(currentPlayerMap) && i < 5; i++ {
					if (game.currentPlayerAndScoreLeaderboard == false) && (currentPlayerMap[i][0] ==
						game.userName) && (currentPlayerScoreMap[i][0] == game.score) {
						game.drawOps.GeoM.Reset()
						text.Draw(screen, strconv.Itoa(i+1)+". "+currentPlayerMap[i][0]+": "+strconv.Itoa(currentPlayerScoreMap[i][0]), mplusNormalFont, ScreenWidth*0.15, tempHeight, colornames.Red)
						tempHeight += 100
						game.currentPlayerAndScoreLeaderboard = true
					} else {
						game.drawOps.GeoM.Reset()
						text.Draw(screen, strconv.Itoa(i+1)+". "+currentPlayerMap[i][0]+": "+strconv.Itoa(currentPlayerScoreMap[i][0]), mplusNormalFont, ScreenWidth*0.15, tempHeight, colornames.White)
						tempHeight += 100
					}
				}
			}
		}
	}
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func OpenDataBase(dbfile string) *sql.DB {
	database, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatal(err)
	}
	return database
}

func create_tables(database *sql.DB) {
	createStatement1 := "CREATE TABLE IF NOT EXISTS players(    " +
		"user_name TEXT NOT NULL," +
		"score INTEGER DEFAULT 0);"
	database.Exec(createStatement1)
}

func (game Game) addGameEntry(database *sql.DB) {
	insertStatement := "INSERT INTO PLAYERS (user_name, score) VALUES (?,?);"
	preppedStatement, err := database.Prepare(insertStatement)
	if err != nil {
		log.Fatal(err)
	}
	preppedStatement.Exec(game.userName, game.score)
}

func (game Game) processDBtoMaps() {
	db, err := sql.Open("sqlite3", "./LeaderBoard.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM players ORDER BY score DESC")
	if err != nil {
		panic(err)
	}

	var temp_user_name string
	var temp_score int
	row_number := 0
	current_player_row_number := 0

	for rows.Next() {
		err = rows.Scan(&temp_user_name, &temp_score)
		userNameMap[row_number] = append(userNameMap[row_number], temp_user_name)
		scoreMap[row_number] = append(scoreMap[row_number], temp_score)
		if temp_user_name == game.userName {
			currentPlayerMap[current_player_row_number] = append(currentPlayerMap[current_player_row_number], temp_user_name)
			currentPlayerScoreMap[current_player_row_number] = append(currentPlayerScoreMap[current_player_row_number], temp_score)
			current_player_row_number += 1
		}
		row_number += 1
	}
	rows.Close()
	db.Close()
	game.processedDB = true

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Berserk/Tank Game by Trevor Wysong")
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
	titleScreenBackground, _, err := ebitenutil.NewImageFromFile("art assets/background.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.titleScreenBackground.upPict = titleScreenBackground

	winnerScreen, _, err := ebitenutil.NewImageFromFile("art assets/winner.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.winnerScreen.upPict = winnerScreen

	loserScreen, _, err := ebitenutil.NewImageFromFile("art assets/loser.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.loserScreen.upPict = loserScreen

	firstMap, _, err := ebitenutil.NewImageFromFile("art assets/Level1Correct.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.firstMap.upPict = firstMap

	secondMap, _, err := ebitenutil.NewImageFromFile("art assets/Level2.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.secondMap.upPict = secondMap

	thirdMap, _, err := ebitenutil.NewImageFromFile("art assets/Level3.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.thirdMap.upPict = thirdMap

	upPlayer, _, err := ebitenutil.NewImageFromFile("art assets/tankFilledTopSquare.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	downPlayer, _, err := ebitenutil.NewImageFromFile("art assets/tankFilledTopSquareDown.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	leftPlayer, _, err := ebitenutil.NewImageFromFile("art assets/tankFilledTopSquareLeft.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	rightPlayer, _, err := ebitenutil.NewImageFromFile("art assets/tankFilledTopSquareRight.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.playerSprite.upPict = upPlayer
	game.playerSprite.downPict = downPlayer
	game.playerSprite.leftPict = leftPlayer
	game.playerSprite.rightPict = rightPlayer

	tankTopperUp, _, err := ebitenutil.NewImageFromFile("art assets/tankTopper.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	tankTopperDown, _, err := ebitenutil.NewImageFromFile("art assets/tankTopperDown.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	tankTopperLeft, _, err := ebitenutil.NewImageFromFile("art assets/tankTopperLeft.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	tankTopperRight, _, err := ebitenutil.NewImageFromFile("art assets/tankTopperRight.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.tankTopper.upPict = tankTopperUp
	game.tankTopper.downPict = tankTopperDown
	game.tankTopper.leftPict = tankTopperLeft
	game.tankTopper.rightPict = tankTopperRight

	fireball, _, err := ebitenutil.NewImageFromFile("art assets/fireball.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.fireball.upPict = fireball

	coins, _, err := ebitenutil.NewImageFromFile("art assets/gold-coins-large.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.coinSprite.upPict = coins

	personEnemyUp, _, err := ebitenutil.NewImageFromFile("art assets/personEnemyUp.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	personEnemyDown, _, err := ebitenutil.NewImageFromFile("art assets/personEnemyDown.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	personEnemyLeft, _, err := ebitenutil.NewImageFromFile("art assets/personEnemyLeft.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	personEnemyRight, _, err := ebitenutil.NewImageFromFile("art assets/personEnemyRight.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.personEnemy.upPict = personEnemyUp
	game.personEnemy.downPict = personEnemyDown
	game.personEnemy.leftPict = personEnemyLeft
	game.personEnemy.rightPict = personEnemyRight

	monsterEnemyUp, _, err := ebitenutil.NewImageFromFile("art assets/monsterEnemyUp.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	monsterEnemyDown, _, err := ebitenutil.NewImageFromFile("art assets/monsterEnemyDown.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	monsterEnemyLeft, _, err := ebitenutil.NewImageFromFile("art assets/monsterEnemyLeft.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	monsterEnemyRight, _, err := ebitenutil.NewImageFromFile("art assets/monsterEnemyRight.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.monsterEnemy.upPict = monsterEnemyUp
	game.monsterEnemy.downPict = monsterEnemyDown
	game.monsterEnemy.leftPict = monsterEnemyLeft
	game.monsterEnemy.rightPict = monsterEnemyRight

	heart, _, err := ebitenutil.NewImageFromFile("art assets/heartScaled.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.heartSprite1.upPict = heart
	game.heartSprite2.upPict = heart
	game.heartSprite3.upPict = heart
}
