package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const PORT = ":5000"

var lobbies []*Lobby
var games []*Game

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func homePage(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)
	fmt.Fprint(w, "Homepage Endpoint Hit")
}

func PlayerFromJWT(jwt string) (*Player, error) {

	// Call localhost:5000/user and extract player
	var player *Player
	req, err := http.NewRequest("GET", "http://localhost:8080/user", nil)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header
	req.Header.Add("Authorization", "Bearer "+jwt)

	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("body of user from jwt is", string(body))

	if err := json.Unmarshal(body, &player); err != nil {
		fmt.Println("Error when unmarshalling jwt", err)
		return nil, err
	}
	return player, nil
}

func GetLobby(gameCode string, lobbies *[]*Lobby) *Lobby {
	for _, lobby := range *lobbies {
		if lobby.GameCode == gameCode {
			return lobby
		}
	}
	return nil
}
func GetGame(gameCode string, games *[]*Game) *Game {
	for _, game := range *games {
		if game.GameCode == gameCode {
			return game
		}
	}
	return nil
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)
	//gameCode=12345&intent=lobby&jwt="{\"id\":\"123\",\"name\":\"test\"}"
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	//Hanle upgrade error
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error when upgrading", err)
		return
	}

	//Handle getting game code
	code := r.URL.Query().Get("gameCode")
	if len(code) != 5 {
		str, _ := NewErrorEventPayload("Game code len is not 5", true).ToJSON()
		fmt.Println(str)
		return
	}

	//Handle getting player from JWT
	player, err := PlayerFromJWT(r.URL.Query().Get("jwt"))
	if err != nil {
		str, _ := NewErrorEventPayload(err.Error(), true).ToJSON()
		fmt.Println(str)
		return
	}

	// Handle creating player with a ws connection
	player, err = NewPlayer(player.Id, player.Name, con)
	if err != nil {
		str, _ := NewErrorEventPayload(err.Error(), true).ToJSON()
		fmt.Println(str)
		return
	}

	//Player has ws now

	intent := r.URL.Query().Get("intent")
	var errorPayloadStr string

	//No other known intent
	if intent != "lobby" && intent != "game" {
		errorPayloadStr, _ = NewErrorEventPayload("Invalid intent "+intent, true).ToJSON()
	}

	//Handle lobby join or rejoin
	if intent == "lobby" {
		lobby := GetLobby(code, &lobbies)

		if lobby == nil {
			errorPayloadStr, _ = NewErrorEventPayload("Lobby not found", true).ToJSON()
		} else if !lobby.HasRoom() {
			errorPayloadStr, _ = NewErrorEventPayload("Not enough room in lobby"+lobby.GameCode, true).ToJSON()
		} else if err := lobby.Join(player); err != nil {
			errorPayloadStr, _ = NewErrorEventPayload(err.Error(), true).ToJSON()
		}
	}

	//Handle game rejoin (not join, that happens when lobby upgrades to game)
	if intent == "game" {
		game := GetGame(code, &games)
		if game == nil {
			errorPayloadStr, _ = NewErrorEventPayload("Game not found", true).ToJSON()
		} else if !game.CanRejoin(code, player) {
			errorPayloadStr, _ = NewErrorEventPayload("Player can not rejoin game "+player.Name+game.GameCode, true).ToJSON()
		} else if err := game.Rejoin(player); err != nil {
			errorPayloadStr, _ = NewErrorEventPayload(err.Error(), true).ToJSON()
		}
	}

	//If any error so far, print and send to player
	if errorPayloadStr != "" {
		fmt.Println(errorPayloadStr)
		if err := player.SendMessage(errorPayloadStr); err != nil {
			fmt.Println("Error when sending error message lol", err)
		}

		if err := player.Disconnect(); err != nil {
			fmt.Println("Error when disconnecting player", err)
		}
		return
	}

	fmt.Println("Client Successfully Connected...", player.Name)

}

type GameCode struct {
	GameCode string `json:"gameCode"`
}

func fiveDigitCode() string {
	randomNumber := rand.Intn(90000) + 10000
	return strconv.Itoa(randomNumber)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func createLobbiesRoute(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	gameCode := GameCode{fiveDigitCode()}
	lobby, err := NewLobby(gameCode.GameCode, &lobbies, &games)
	if err != nil {
		fmt.Println("Error when creating lobby", err)
		return
	}
	lobbies = append(lobbies, lobby)
	fmt.Println("Added lobby. Now, there are ", len(lobbies), "lobbies")
	json.NewEncoder(w).Encode(gameCode)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/lobbies", createLobbiesRoute)
	http.HandleFunc("/ws", wsEndpoint)
}

func StartServing() {
	fmt.Println("Hello World")
	setupRoutes()
	fmt.Println("server running on port ", PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(PORT, nil))

}

type PawnEvent struct {
	PawnNumber int `json:"pawnNumber"`
}

type ErrorEvent struct {
	Message string `json:"message"`
	IsFatal bool   `json:"isFatal"`
}

func NewErrorEventPayload(message string, isFatal bool) *EventPayload {
	str, _ := json.Marshal(ErrorEvent{message, isFatal}) //error free assume
	return &EventPayload{"error", string(str)}
}

type EventPayload struct {
	Event   string `json:"event"`
	Content string `json:"content"`
}

func (e *EventPayload) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		fmt.Println("Error when marshalling event payload", err)
		return "", err
	}
	return string(jsonBytes), nil
}

type Player struct {
	Id   string          `json:"id"`
	Name string          `json:"name"`
	Con  *websocket.Conn `json:"-"`

	OnDisconnect  func() error    `json:"-"`
	OnPawn        func(int) error `json:"-"`
	OnDice        func() error    `json:"-"`
	OnStartGame   func() error    `json:"-"`
	keepListening bool            `json:"-"`
}

func (p *Player) SendMessage(message string) error {
	err := p.Con.WriteMessage(websocket.TextMessage, []byte(message))
	return err
}

func (p *Player) Listen() error {

	con := p.Con

	//Call onDisconnect
	con.SetCloseHandler(func(code int, text string) error {
		fmt.Println("Connection ended bc of ", text, "code: ", code)
		p.keepListening = false
		p.OnDisconnect()
		return nil
	})

	//Listen for events
	for p.keepListening {
		_, bytes, err := con.ReadMessage()
		//If no connection, dont be surprised by ReadMessage error
		if err != nil {
			if !p.keepListening {
				return nil
			} else {
				fmt.Println("Error when reading message", err)
				return err
			}
		}
		str := string(bytes)
		var eventPayload EventPayload
		//unmarshall json from str
		if err := json.Unmarshal([]byte(str), &eventPayload); err != nil {
			fmt.Println("Could not decode event json in conn for loop ", err)
			return err
		}

		if eventPayload.Event == "pawn" {

			//If pawn event, call onPawn
			var pawnEvent PawnEvent
			if err := json.Unmarshal([]byte(eventPayload.Content), &pawnEvent); err != nil {
				fmt.Println("Could not decode pawn json in conn for loop ", err)
				return err
			}
			p.OnPawn(pawnEvent.PawnNumber)
		} else if eventPayload.Event == "dice" {

			//If dice event, call onDice
			err := p.OnDice()
			if err != nil {
				fmt.Println("Error when calling OnDice", err)
				return err
			}
		} else if eventPayload.Event == "startGame" {
			err := p.OnStartGame()
			if err != nil {
				fmt.Println("Error when calling OnStartGame", err)
				return err
			}

		} else {
			//If unknown event, print error
			fmt.Println("Error Unknown event: ", eventPayload.Event)
		}

		fmt.Println("Received Message: " + str)

	}

	fmt.Println("Listen() ended")
	return nil

}

func (p *Player) Disconnect() error {
	p.keepListening = false
	if p.Con == nil {
		return nil
	}
	return p.Con.Close()
}

func DefaultOnDisconnect() error {
	fmt.Println("onDisconnect() called")
	return nil
}

func DefaultOnPawn(pawnNumber int) error {
	fmt.Println("onPawn() called", pawnNumber)
	return nil
}

func DefaultOnDice() error {
	fmt.Println("onDice() called")
	return nil
}

func DefaultOnStartGame() error {
	fmt.Println("onStartGame() called")
	return nil
}

func NewPlayer(id string, name string, con *websocket.Conn) (*Player, error) {

	player := &Player{id, name, con, DefaultOnDisconnect, DefaultOnPawn, DefaultOnDice, DefaultOnStartGame, true}

	go func() {
		err := player.Listen()
		if err != nil {
			fmt.Println("Error when calling Listen", err)
		}
	}()

	return player, nil
}

func NewLobby(gameCode string, lobbies *[]*Lobby, games *[]*Game) (*Lobby, error) {

	onSelfRemove := func() error {
		for i, lobby := range *lobbies {
			if lobby.GameCode == gameCode {
				*lobbies = append((*lobbies)[:i], (*lobbies)[i+1:]...)

				fmt.Println("There are", len(*lobbies), "lobbies left", "after removing lobby", gameCode)
				return nil
			}
		}

		return fmt.Errorf("Lobby not found in lobbies list to self remove %s", gameCode)
	}

	onStartGame := func() error {
		fmt.Println("onStartGame() called on lobby")

		for _, lobby := range *lobbies {
			if lobby.GameCode == gameCode {
				fmt.Println("Lobby starting game", gameCode)
				game, err := NewGame(gameCode, lobby.Players)
				if err != nil {
					return err
				}
				*games = append(*games, game)
				return onSelfRemove()
			}
		}
		return fmt.Errorf("Lobby not found in lobbies list to start game %s", gameCode)

	}

	lobby := &Lobby{gameCode, make([]*Player, 0), onStartGame, onSelfRemove, nil, 5 * time.Second}
	lobby.StartSelfRemoveTimer()
	fmt.Println("Lobby created with gameCode", gameCode)
	return lobby, nil
}

type Lobby struct {
	GameCode        string        `json:"gameCode"`
	Players         []*Player     `json:"players"`
	OnStartGame     func() error  `json:"-"`
	OnSelfRemove    func() error  `json:"-"`
	SelfRemoveTimer *time.Timer   `json:"-"`
	TimerDuration   time.Duration `json:"-"`
}

func (l *Lobby) HasRoom() bool {
	return len(l.Players) < 4
}

func (l *Lobby) Join(player *Player) error {

	//If full, return error
	if !l.HasRoom() {
		return fmt.Errorf("Lobby is full %s", l.GameCode)
	}

	//If already in lobby, replace
	hasReplaced := false
	for i, otherPlayer := range l.Players {
		if otherPlayer.Id == player.Id {
			l.Players[i] = player
			hasReplaced = true
			break
		}
	}

	//Append if not replaced
	if !hasReplaced {
		l.Players = append(l.Players, player)
	}

	fmt.Println("Player added to lobby, now has length: ", len(l.Players))

	//Add listeners
	l.ListenToStartRequest()
	player.OnDisconnect = func() error {
		return l.Leave(player)

	}

	//Let player know of successful join
	str, _ := (&EventPayload{"ok", ""}).ToJSON()
	player.SendMessage(str)

	//Let players know of new neighbor
	l.UpdateLobbyPlayersState()

	//Cancel self remove timer
	if l.SelfRemoveTimer != nil {
		fmt.Println("Cancelling self remove timer")
		l.SelfRemoveTimer.Stop()
		l.SelfRemoveTimer = nil
	}

	return nil
}

func (l *Lobby) HasPlayer(p *Player) bool {
	for _, player := range l.Players {
		if player == p {
			return true
		}
	}
	return false
}

func (l *Lobby) Leave(p *Player) error {
	if !l.HasPlayer(p) {
		fmt.Println("Player not in lobby", p.Name)
		return fmt.Errorf("Player not in lobby")
	}

	for i, player := range l.Players {
		if player == p {
			fmt.Println("Player removed from lobby", p.Name)
			l.Players = append(l.Players[:i], l.Players[i+1:]...)
			l.ListenToStartRequest()
			break
		}
	}

	l.StartSelfRemoveTimer()
	l.UpdateLobbyPlayersState()
	return nil
}

func (l *Lobby) StartSelfRemoveTimer() {

	//Remove lobby from lobbies list if no player in lobby for eg) 1 minute
	l.SelfRemoveTimer = time.AfterFunc(l.TimerDuration, func() {
		if len(l.Players) == 0 {
			err := l.OnSelfRemove()
			if err != nil {
				fmt.Println("Error when calling OnRemoveSelf", err)
			}
		}
	})
}

func (l *Lobby) ListenToStartRequest() error {
	for i, player := range l.Players {
		//First player is the 'owner' or 'host', so he has access to the start button
		fmt.Println("Setting OnStartGame for player", player.Name)
		index := i //Closure?
		player.OnStartGame = func() error {
			if index != 0 {
				str, _ := NewErrorEventPayload("Not owner to start game", false).ToJSON()
				player.SendMessage(str)
				return fmt.Errorf("not owner to start game")
			}

			//TODO
			// if len(l.Players) == 1 {
			// 	str, _ := NewErrorEventPayload("Not enough players to start game", false).ToJSON()
			// 	player.SendMessage(str)
			// 	return fmt.Errorf("not enough players to start game")
			// }

			fmt.Println("Starting game by", player.Name)
			l.OnStartGame()
			return nil

		}
	}

	return nil
}

func (l *Lobby) UpdateLobbyPlayersState() error {

	str, _ := json.Marshal(l)

	eventPayload := EventPayload{"lobbyPlayers", string(str)}
	eventPayloadString, _ := eventPayload.ToJSON()

	for _, player := range l.Players {
		err := player.SendMessage(eventPayloadString)
		if err != nil {
			fmt.Println("Error when sending lobbyPlayersEvent", err)
			return err
		}
	}
	return nil
}

// <game engine>
var HOME_POSITIONS = [16]int{32, 33, 47, 48, 41, 42, 56, 57, 176, 177, 191, 192, 167, 168, 182, 183}
var WIN_POSITIONS = [4]int{111, 97, 113, 127}
var START_POSITIONS = [4]int{91, 23, 133, 201}
var START_TO_END_POSITIONS = generateStartToEndPositions()

var loopPositions = []int{90, 91, 92, 93, 94, 95, 81, 66, 51, 36, 21, 6, 7, 8, 23, 38, 53, 68, 83, 99, 100, 101, 102, 103, 104, 119, 134, 133, 132, 131, 130, 129, 143, 158, 173, 188, 203, 218, 217, 216, 201, 186, 171, 156, 141, 125, 124, 123, 122, 121, 120, 105}
var startToColored = [4][2]int{{91, 105}, {23, 7}, {133, 119}, {201, 217}}
var coloredPaths = [4][]int{
	{106, 107, 108, 109, 110, 111},
	{22, 37, 52, 67, 82, 97},
	{118, 117, 116, 115, 114, 113},
	{202, 187, 172, 157, 142, 127},
}

func Contains(element int, array []int) bool {
	for _, e := range array {
		if e == element {
			return true
		}
	}
	return false
}

func generateStartToEndPositions() [16][]int {
	startToEndPositions := [16][]int{}

	for playerNumber := 0; playerNumber < 4; playerNumber++ {
		for i := 0; i < 4; i++ {
			pawnNumber := playerNumber*4 + i
			homePath := []int{HOME_POSITIONS[pawnNumber]}
			toColoredPath := reslice(startToColored[playerNumber][0], startToColored[playerNumber][1])
			coloredPath := coloredPaths[playerNumber]
			fullPath := append(append(homePath, toColoredPath...), coloredPath...)
			startToEndPositions[pawnNumber] = fullPath
		}

	}

	return startToEndPositions

}
func findIndex(slice []int, value int) int {
	for i, item := range slice {
		if item == value {
			return i
		}
	}
	return -1
}

func lastElement(slice []int) int {
	return slice[len(slice)-1]
}

func reslice(startPos int, endPos int) []int {
	i1 := findIndex(loopPositions, startPos)
	i2 := findIndex(loopPositions, endPos)

	if i1 == -1 || i2 == -1 {
		fmt.Println("Invalid start or end position", startPos, endPos)
		return nil
	}

	resliced := []int{}
	for i := i1; i < 1000; i++ {
		rotatedIndex := i % len(loopPositions)
		resliced = append(resliced, loopPositions[rotatedIndex])
		if rotatedIndex == i2 {
			break
		}
	}

	return resliced
}

func ErrorIfInvalidPosition(pawnNumber int, position int) error {
	startIndex := findIndex(START_TO_END_POSITIONS[pawnNumber], position)
	if startIndex == -1 {
		return fmt.Errorf("pawn %d can never be positioned in position %d", pawnNumber, position)
	}
	return nil
}

func GetMovementPath(pawnNumber int, position int, steps int) ([]int, error) {

	//If it cant find an index, it wont find a path
	if err := ErrorIfInvalidPosition(pawnNumber, position); err != nil {
		return nil, err
	}

	startIndex := findIndex(START_TO_END_POSITIONS[pawnNumber], position)

	path := []int{}
	for i := 0; i <= steps; i++ {
		curIndex := startIndex + i
		if curIndex >= len(START_TO_END_POSITIONS[pawnNumber]) {
			return nil, fmt.Errorf("pawn %d can not move %d steps from position %d", pawnNumber, steps, position)
		}
		path = append(path, START_TO_END_POSITIONS[pawnNumber][curIndex])
	}
	return path, nil
}

func GetBackToHomePath(pawnNumber, position int) ([]int, error) {

	err := ErrorIfInvalidPosition(pawnNumber, position)
	if err != nil {
		return nil, err
	}
	startIndex := findIndex(START_TO_END_POSITIONS[pawnNumber], position)

	path := []int{}

	for i := startIndex; i >= 0; i-- {
		path = append(path, START_TO_END_POSITIONS[pawnNumber][i])
	}
	return path, nil
}

func GetPlayerNumber(pawnNumber int) int {
	return pawnNumber / 4
}

func ErrorIfInvalidMove(pawnNumber int, position int, dice int) error {
	if err := ErrorIfInvalidPosition(pawnNumber, position); err != nil {
		return err
	}
	//If pawn is home, cant move unless dice is 6
	if position == HOME_POSITIONS[pawnNumber] && dice != 6 {
		return fmt.Errorf("invalid move, pawn %d is in home and dice %d is not 6", pawnNumber, dice)
	}

	//If pawn is at last position, dont move it
	if position == WIN_POSITIONS[GetPlayerNumber(pawnNumber)] {
		return fmt.Errorf("invalid move, pawn is at last position")
	}

	//If pawn is going to over run the last position, dont move it
	posIndex := findIndex(START_TO_END_POSITIONS[pawnNumber], position)
	if posIndex+dice >= len(START_TO_END_POSITIONS[pawnNumber]) {
		return fmt.Errorf("pawn %d can not move %d steps from position %d", pawnNumber, dice, position)
	}
	return nil
}

// </engine path>

func NewGame(gameCode string, players []*Player) (*Game, error) {
	newPlayers := [4]*Player{}
	onlineStatuses := [4]bool{}

	for i, player := range players {
		newPlayers[i] = player
		onlineStatuses[i] = true
	}

	pawnPositions := HOME_POSITIONS
	turn := 0
	dice := 5
	hasDiced := false
	movablePawns := []int{}

	winner := -1
	pawnAnimationPaths := [16][]int{}
	offlineTimeout := 10
	playerRemoveTimers := [4]*time.Timer{nil, nil, nil, nil}

	game := &Game{gameCode, newPlayers, pawnPositions, turn, dice, hasDiced, movablePawns, winner, onlineStatuses, pawnAnimationPaths, offlineTimeout, playerRemoveTimers}
	game.ReAddListeners()
	game.RemoveNullPawnPositions()
	game.UpdateState()
	return game, nil

}

type Game struct {
	//nil or -1 or false if non existent
	GameCode      string     `json:"gameCode"`
	Players       [4]*Player `json:"players"`
	PawnPositions [16]int    `json:"pawnPositions"`
	Turn          int        `json:"turn"`         //Current Player
	Dice          int        `json:"dice"`         //Value of dice (can be before or after rolling dice)
	HasDiced      bool       `json:"hasDiced"`     //Either ready to dice or ready to move pawn
	MovablePawns  []int      `json:"movablePawns"` //To help frontend know which pawns can be chosen

	Winner             int            `json:"winner"`
	OnlineStatuses     [4]bool        `json:"onlineStatuses"`     //To help frontend show online status
	PawnAnimationPaths [16][]int      `json:"pawnAnimationPaths"` //To help frontend animate pawn jumps. null if no animation
	OfflineTimeout     int            `json:"offlineTimeout"`     // Remove from game if offline for this long
	PlayerRemoveTimers [4]*time.Timer `json:"-"`                  //To stop player from being removed if he reconnects
}

func (g *Game) ReAddListeners() {
	g.ListenToPawn()
	g.ListenToDice()
	g.ListenToOffline()
}

func (g *Game) ListenToPawn() {
	for i, p := range g.Players {

		playerNumber := i
		player := p

		if player == nil {
			continue
		}

		player.OnPawn = func(pawnNumber int) error {
			//If not your turn, return error
			if g.Turn != playerNumber {
				str, _ := NewErrorEventPayload("Not your turn", false).ToJSON()
				player.SendMessage(str)
				fmt.Printf("not your turn %s", player.Name)
				return nil
			}

			//If not your pawn, return error
			if GetPlayerNumber(pawnNumber) != playerNumber {
				str, _ := NewErrorEventPayload("Not your pawn", false).ToJSON()
				player.SendMessage(str)
				fmt.Printf("%d is not your pawn %s ", pawnNumber, player.Name)
				return nil
			}

			//If not diced, return error
			if !g.HasDiced {
				str, _ := NewErrorEventPayload("You have not diced yet", false).ToJSON()
				player.SendMessage(str)
				fmt.Printf("you have not diced yet %s", g.Players[g.Turn].Name)
				return nil
			}

			//Error if pawn can not be moved
			if err := ErrorIfInvalidMove(pawnNumber, g.PawnPositions[pawnNumber], g.Dice); err != nil {
				str, _ := NewErrorEventPayload(err.Error(), false).ToJSON()
				player.SendMessage(str)
				fmt.Printf("invalid move %s", err.Error())
				return nil
			}

			//This is a valid movable pawn
			if err := g.HandleMovePawn(pawnNumber); err != nil {
				str, _ := NewErrorEventPayload(err.Error(), false).ToJSON()
				player.SendMessage(str)
				fmt.Printf("error when handling move pawn %s", err.Error())
				return err
			}

			g.UpdateState()
			return nil
		}
	}
}

func (g *Game) ListenToDice() {
	for i, p := range g.Players {

		playerNumber := i
		player := p

		if player == nil {
			continue
		}
		player.OnDice = func() error {
			if g.Turn != playerNumber {
				str, _ := NewErrorEventPayload("Not your turn to dice", false).ToJSON()
				player.SendMessage(str)
				fmt.Println(str)
				return nil
			}

			if g.HasDiced {
				str, _ := NewErrorEventPayload("You have already diced", false).ToJSON()
				player.SendMessage(str)
				fmt.Println(str)
				return nil
			}

			//set dice and the fact that the player diced
			g.HasDiced = true

			g.Dice = rand.Intn(6) + 1 // (0 to 5 ) + 1

			//Let the player know the pawns that can be moved
			g.SetMovablePawns()

			//But if the player cant move any pawn, turn to the next player
			if len(g.MovablePawns) == 0 {
				g.NextTurn()
			}

			g.UpdateState()
			return nil
		}
	}
}

func (g *Game) ListenToOffline() {
	for i, p := range g.Players {

		playerNumber := i
		player := p

		if player == nil {
			continue
		}
		player.OnDisconnect = func() error {

			//announce offline status
			g.OnlineStatuses[playerNumber] = false
			g.UpdateState()

			//set up timer to self destruct
			duration := time.Duration(g.OfflineTimeout) * time.Second
			g.PlayerRemoveTimers[playerNumber] = time.AfterFunc(duration, func() {

				//remove player
				g.Players[playerNumber] = nil

				//remove pawns
				g.RemoveNullPawnPositions()

				//if was my turn, change turn
				if g.Turn == playerNumber {
					g.NextTurn()
				}

				//one player my be left, make him winner if so
				possiblyWinner := g.GetWinner()
				if possiblyWinner != -1 {
					g.Winner = possiblyWinner
				}
				g.UpdateState()
			})
			return nil
		}
	}
}

func (g *Game) RemoveNullPawnPositions() {
	for pawnNumber, _ := range g.PawnPositions {
		if g.Players[GetPlayerNumber(pawnNumber)] == nil {
			g.PawnPositions[pawnNumber] = -1
		}
	}
}

func (g *Game) SetMovablePawns() {
	playerNumber := g.Turn
	g.MovablePawns = []int{}
	for pawnNumber, position := range g.PawnPositions {
		isPlayers := GetPlayerNumber(pawnNumber) == playerNumber
		isMovable := ErrorIfInvalidMove(pawnNumber, position, g.Dice) == nil
		if isPlayers && isMovable {
			g.MovablePawns = append(g.MovablePawns, pawnNumber)
		}
	}
}

func (g *Game) HandleMovePawn(pawnNumber int) error {
	//pawn is valid to move at this point

	//Reset animation paths
	g.PawnAnimationPaths = [16][]int{}

	position := g.PawnPositions[pawnNumber]
	playerNumber := GetPlayerNumber(pawnNumber)

	//1 step to get out of home, dice steps for the rest
	steps := g.Dice
	if position == HOME_POSITIONS[pawnNumber] && g.Dice == 6 {
		steps = 1
	}

	//Update path and position
	path, err := GetMovementPath(pawnNumber, position, steps)
	if err != nil {
		return err
	}
	g.PawnAnimationPaths[pawnNumber] = path
	position = lastElement(path)
	g.PawnPositions[pawnNumber] = position

	//Displace other pawns if possible
	hasDisplaced := false
	for otherPawnNumber, otherPosition := range g.PawnPositions {
		samePosition := position == otherPosition
		differentPlayer := GetPlayerNumber(otherPawnNumber) != playerNumber

		if samePosition && differentPlayer {
			hasDisplaced = true
			path, err := GetBackToHomePath(otherPawnNumber, otherPosition)
			if err != nil {
				return err
			}
			g.PawnAnimationPaths[otherPawnNumber] = path
			g.PawnPositions[otherPawnNumber] = lastElement(path)
		}
	}

	possiblyWinner := g.GetWinner()
	hasWon := possiblyWinner != -1
	keepTurn := g.Dice == 6 || hasDisplaced

	if hasWon {
		//If won, set winner. Dont allow dicing
		g.Winner = possiblyWinner
	} else if keepTurn {
		//another chance to dice
		g.HasDiced = false
		g.MovablePawns = []int{}
	} else {
		//change turn
		g.NextTurn()
	}

	//If won, set winner
	if hasWon {
		g.Winner = possiblyWinner
	}

	fmt.Println("Animation paths are", g.PawnAnimationPaths)

	return g.UpdateState()
}

func (g *Game) CountPlayers() int {
	playerCount := 0
	for _, player := range g.Players {
		if player != nil {
			playerCount++
		}
	}
	return playerCount
}

func (g *Game) CanRejoin(gameCode string, player *Player) bool {
	if g.GameCode != gameCode {
		return false
	}
	for _, p := range g.Players {
		if p != nil && p.Id == player.Id {
			return true
		}
	}
	return false
}

func (g *Game) NextTurn() {
	//also works for 1 player

	for i := 0; i < 4; i++ {
		turn := (g.Turn + 1 + i) % 4
		fmt.Println("trying turn", turn, "with player", g.Players[turn])
		if g.Players[turn] != nil {
			fmt.Println("settled with player", g.Players[turn].Name)
			g.Turn = turn
			g.HasDiced = false
			g.MovablePawns = []int{}
			return
		}
	}

}

func (g *Game) Rejoin(player *Player) error {

	for playerNumber, p := range g.Players {

		if p != nil && p.Id == player.Id {
			//replace old player and set online status to true
			g.Players[playerNumber] = player
			g.OnlineStatuses[playerNumber] = true

			//cancel offline timer
			timer := g.PlayerRemoveTimers[playerNumber]
			if timer != nil {
				timer.Stop()
			}

			//Let player know of successful join
			str, _ := (&EventPayload{"ok", ""}).ToJSON()
			player.SendMessage(str)

			g.ReAddListeners()
			return nil
		}
	}

	return fmt.Errorf("Player %s not found in game %s when rejoining", player.Name, g.GameCode)
}

func (g *Game) GetWinner() int {
	//If one player left, he is the winner

	// TODO
	// if g.CountPlayers() == 1 {
	// 	for playerNumber, player := range g.Players {
	// 		if player != nil {
	// 			return playerNumber
	// 		}
	// 	}
	// }

	//4 pawns wins means the player wins
	wc := [4]int{} //to have a [0,0,0,0] with out doing [4]int{}
	winCounts := wc[:]
	for pawnNumber, position := range g.PawnPositions {
		playerNumber := GetPlayerNumber(pawnNumber)
		if position == WIN_POSITIONS[playerNumber] {
			winCounts[playerNumber]++
		}
	}

	return findIndex(winCounts, 4)
}

func (g *Game) UpdateState() error {
	//game state to json
	gameState, err := json.Marshal(g)
	if err != nil {
		fmt.Println("Error when marshalling game state", err)
		return err
	}

	//make game state event and toString it
	eventStr, _ := (&EventPayload{"gameState", string(gameState)}).ToJSON()

	//Send state to each player
	for _, player := range g.Players {
		if player == nil {
			continue
		}
		if err := player.SendMessage(eventStr); err != nil {
			fmt.Println("Error when sending gameStateEvent to player", player.Name, err)
		}
	}
	return nil
}
