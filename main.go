package main
// TO DO RECHECKER TOUS LES USAGES DE POINTEURS


import (
	// "unsafe"
	"os"
	"log"
	// "reflect"
	"fmt"
    "strconv"
    "encoding/json"
	// "encoding/binary"
	"time"
	// "image"
	// "image/draw"
	// "image/color"
	// "math"
	// "math/rand"

	"github.com/gofiber/fiber/v2"
    "github.com/gofiber/template/html"
	// "gocv.io/x/gocv"
	// "github.com/gordonklaus/portaudio"
	// "github.com/xthexder/go-jack"
	// "github.com/gofiber/fiber/v2"
	// "go.mongodb.org/mongo-driver/mongo"
)

const defaultMicId = 1

func check(message string, err error) {
	if err != nil {
		fmt.Printf("check() err : " + message)
        log.Fatal(err)
	}
}

func main() {
	fmt.Println("App is Starting...")

	gameEngine := NewGameEngine()
	// fmt.Println("gameEngine : ", gameEngine)



    // Initialize standard Go html template engine
    engine := html.New("./views", ".html")
	// app := fiber.New()
    app := fiber.New(fiber.Config{
        Views: engine,
    })

	app.Get("/stats", func(c *fiber.Ctx) error {
		statsJSON, err := json.Marshal(gameEngine.makeStats())
		if err != nil {
			return err
		}
		return c.SendString(string(statsJSON))
	})

	app.Get("/states", func(c *fiber.Ctx) error {
		statesJSON, err := json.Marshal(gameEngine.States)
		if err != nil {
			return err
		}
		return c.SendString(string(statesJSON))
	})

	app.Get("/history", func(c *fiber.Ctx) error {
		historyJSON, err := json.Marshal(gameEngine.EndGameReqs)
		if err != nil {
			return err
		}
		return c.SendString(string(historyJSON))
	})

	app.Get("/", func(c *fiber.Ctx) error {
		startGameRes := gameEngine.startGame()
		// fmt.Println("startGameRes : ", startGameRes)

		DoILive, err := json.Marshal(startGameRes.DoILive)
		if err != nil {
			return err
		}
		WhoStartsFirst, err := json.Marshal(startGameRes.WhoStartsFirst)
		if err != nil {
			return err
		}
		StatesJSON, err := json.Marshal(startGameRes.States)
		if err != nil {
			return err
		}
		StatsJSON, err := json.Marshal(startGameRes.Stats)
		if err != nil {
			return err
		}

        // Render index template
        return c.Render("index", fiber.Map{
            "DoILive": string(DoILive),
            "WhoStartsNext": string(WhoStartsFirst),
			"States": string(StatesJSON),
			"Stats": string(StatsJSON),
        })
	})
	app.Post("/", func(c *fiber.Ctx) error {
		// c.Accepts("application/json") // "application/json"
		payload := EndGameReq{}
		if err := json.Unmarshal(c.Body(), &payload); err != nil {
			fmt.Println("err : ", err)
			return c.SendString(err.Error())
		}
		// fmt.Println("payload after : ", payload)
		// fmt.Println("payload.Winner : ", payload.Winner)
		// fmt.Println("payload.History : ", payload.History)

		gameEngine.endGame(payload.Winner, payload.History)
		// fmt.Println("endGameRes : ", endGameRes)
		return c.SendString("endGameRes")
	})
	gameEngine.saveGameEngine()
	go func() {
		for range time.NewTicker(time.Hour).C {
			gameEngine.saveGameEngine()
		}
	}()
	app.Listen(":3000")
	fmt.Println("Server listening on port 3000")

	// charset := []string{"0", "1", "2"}
	// res := []string{}
	// goDeeper(func (value string) {
	// 	fmt.Println("value : ", value)
	// 	res = append(res, value)
	// }, charset, 9 - 1, "")
	// // }, charset, 9 - 1, "")
	// fmt.Println("len(res) : ", len(res))
	// fmt.Println("res : ", res)


	fmt.Println("This app is Exiting...")
}

// STRUCTS
func NewGameEngine() GameEngine {
	whoStartsNext := "human"
	var states []State
	charset := []string{"0", "1", "2"}
	goDeeper(func (value string) {
		state := NewState(value)
		states = append(states, state)
	}, charset, 9 - 1, "")
	
	fmt.Println("states length : ", len(states))

	return GameEngine{
		whoStartsNext,
		states,
		[]EndGameReq{},
	}
}
type GameEngine struct {
	// Possible values "human" | "robot"
	WhoStartsNext string
	States []State
	EndGameReqs []EndGameReq
}
func (this *GameEngine) makeStats() map[string]string {
	stats := make(map[string]string)

	stats["numberOfGamesPlayed"] = strconv.Itoa(len(this.EndGameReqs))
	humanWinsCount := 0
	robotWinsCount := 0
	drawCount := 0
	
	for i := range this.EndGameReqs {
		if this.EndGameReqs[i].Winner == "human"  {
			humanWinsCount += 1
		} else if this.EndGameReqs[i].Winner == "robot" {
			robotWinsCount += 1
		} else if this.EndGameReqs[i].Winner == "draw" {
			drawCount += 1
		}
	}

	stats["humanWinsCount"] = strconv.Itoa(humanWinsCount)
	stats["robotWinsCount"] = strconv.Itoa(robotWinsCount)
	stats["drawCount"] = strconv.Itoa(drawCount)
	// stats["winRate"] = fmt.Sprintf("%f", float64(robotWinsCount) / float64(len(this.EndGameReqs)), 64)
	stats["winRate"] = fmt.Sprintf("%f", float64(robotWinsCount) / (float64(robotWinsCount) + float64(humanWinsCount)), 64)

	return stats
}
func (this *GameEngine) doILive() bool {
	return true
}

func (this *GameEngine) startGame() StartGameRes {
	stats := this.makeStats()

	humanWinsCount := 0
	robotWinsCount := 0
	drawCount := 0
	
	for i := range this.EndGameReqs {
		if this.EndGameReqs[i].Winner == "human"  {
			humanWinsCount += 1
		} else if this.EndGameReqs[i].Winner == "robot" {
			robotWinsCount += 1
		} else if this.EndGameReqs[i].Winner == "draw" {
			drawCount += 1
		}
	}

	return StartGameRes{
		this.doILive(),
		this.WhoStartsNext,
		this.States,
		stats,
	}
}
func (this *GameEngine) endGame(winner string, history []Turn) {
	robotTurns := []Turn{}
    for i := range history {
        if history[i].WhosTurn == "robot" {
            robotTurns = append(robotTurns, history[i])
        }
    }

	for i := range robotTurns {
		var currentStateIndex int
		for stateIndex, state := range this.States {
			if state.Id == robotTurns[i].CurrentState {
				currentStateIndex = stateIndex
			}
		}
		if winner == "robot" {
			this.States[currentStateIndex].reward(robotTurns[i].CurrentState, robotTurns[i].ChosenMove)
		} else if winner == "human" {
			this.States[currentStateIndex].punish(robotTurns[i].CurrentState, robotTurns[i].ChosenMove)
		} else if winner == "draw" {
			// ?
		}
	}

	this.EndGameReqs = append(this.EndGameReqs, EndGameReq{
		winner,
		history,
	})

	switch this.WhoStartsNext {
    case "robot":
		this.WhoStartsNext = "human"
    case "human":
		this.WhoStartsNext = "robot"
	}
}
func (this *GameEngine) saveGameEngine() {
	fmt.Println("saveGameEngine")
	// File
	now := time.Now()

    f, err := os.Create("./saves/" + strconv.Itoa(now.Year()) + "-" + strconv.Itoa(int(now.Month())) + "-" + strconv.Itoa(now.Day()) + "." + strconv.Itoa(now.Hour()) + "-" + strconv.Itoa(now.Minute()) + ".txt")
    check("saveGameEngine() os.Create : ", err)
    defer f.Close()


	gameEngineJSON, err := json.Marshal(this)
    check("saveGameEngine() json.Marshal : ", err)

    _, err2 := f.WriteString(string(gameEngineJSON))
    check("saveGameEngine() f.WriteString : ", err2)

	fmt.Println("saveGameEngine done")
	// TO DO db
}

func NewState(id string) State {
	chars := []rune(id)
	weights := Weights{}

	// fmt.Println("string(chars) : ", string(chars))
	if string(chars[0]) == "0" {
		weights.A = 5
	}
	if string(chars[1]) == "0" {
		weights.B = 5
	}
	if string(chars[2]) == "0" {
		weights.C = 5
	}
	if string(chars[3]) == "0" {
		weights.D = 5
	}
	if string(chars[4]) == "0" {
		weights.E = 5
	}
	if string(chars[5]) == "0" {
		weights.F = 5
	}
	if string(chars[6]) == "0" {
		weights.G = 5
	}
	if string(chars[7]) == "0" {
		weights.H = 5
	}
	if string(chars[8]) == "0" {
		weights.I = 5
	}
	// fmt.Println("weights : ", weights)

	return State{
		id,
		weights,
	}
}
type State struct {
	Id string
	Weights Weights
}
func (this *State) reward(stateId string, chosenMove string) {
	switch chosenMove {
	case "a":
		this.Weights.A += 1
	case "b":
		this.Weights.B += 1
	case "c":
		this.Weights.C += 1
	case "d":
		this.Weights.D += 1
	case "e":
		this.Weights.E += 1
	case "f":
		this.Weights.F += 1
	case "g":
		this.Weights.G += 1
	case "h":
		this.Weights.H += 1
	case "i":
		this.Weights.I += 1
	}
}
func (this *State) punish(stateId string, chosenMove string) {
	switch chosenMove {
	case "a":
		this.Weights.A -= 1
	case "b":
		this.Weights.B -= 1
	case "c":
		this.Weights.C -= 1
	case "d":
		this.Weights.D -= 1
	case "e":
		this.Weights.E -= 1
	case "f":
		this.Weights.F -= 1
	case "g":
		this.Weights.G -= 1
	case "h":
		this.Weights.H -= 1
	case "i":
		this.Weights.I -= 1
	}
}

// TYPES
type Weights struct {
	A int
	B int
	C int
	D int
	E int
	F int
	G int
	H int
	I int
}

type StartGameRes struct {
	DoILive bool
	WhoStartsFirst string
	States []State
	Stats map[string]string
}

type EndGameReq struct {
	Winner string
	History []Turn
}
type Turn struct {
	WhosTurn string
	CurrentState string
	ChosenMove string
}

// UTILS
func goDeeper(callback func(value string), charSet []string, k int, current string) {
	for i := 0; i < len(charSet); i++ {
		if k == 0 {
			// await requestSemaphore.acquire();
			res := current + charSet[i]

			charsMap := map[string]int{
				"0": 0,
				"1": 0,
				"2": 0,
			}
		
			for i := 0; i < len(res); i++ { 
				charsMap[string(res[i])] += 1; 
			} 
			condA := charsMap["1"] - charsMap["2"] < 2
			condB := charsMap["2"] - charsMap["1"] < 2

			if condA {
				if condB {
					callback(res); 
				}
			}
		} else {
			goDeeper(callback, charSet, k - 1, current + charSet[i])
		}   
	}
}
