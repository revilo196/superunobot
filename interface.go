package superunobot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)



type Kart struct {
	Color string `json:"color"`
	Value string `json:"value"`
}


type oPlayer struct {
	PlayerName string  `json:"player_name"`
	CardCount int `json:"card_count"`
}

type gamestate struct {
	MyTurn bool 			`json:"my_turn"`
	OtherPlayers []oPlayer `json:"other_players"`
	Hand []Kart					`json:"hand"`
	DiscardedCard Kart			`json:"discarded_card"`
}

func convertReverse(card Card, wish Color) * Kart  {
	valuemapping := make(map[int]string)
	colormapping := make(map[int]string)


	valuemapping[0]  ="ZERO"
	valuemapping[1] = "ONE"
	valuemapping[2] = "TWO"
	valuemapping[3] = "THREE"
	valuemapping[4] = "FOUR"
	valuemapping[5] = "FIVE"
	valuemapping[6] = "SIX"
	valuemapping[7] = "SEVEN"
	valuemapping[8] = "EIGHT"
	valuemapping[9] = "NINE"
	valuemapping[PLUS2] = "DRAW_TWO"
	valuemapping[WECHSEL] = "REVERSE"
	valuemapping[STOPP] = "SKIP"
	valuemapping[WAHL] = "WILD"
	valuemapping[PLUS4] = "WILD_DRAW_FOUR"


	colormapping[RED] = "RED"
	colormapping[BLUE] = "BLUE"
	colormapping[GREEN] = "GREEN"
	colormapping[GELB] = "YELLOW"

	v , ok := valuemapping[card.symbol]
	if !ok {
		v  = "WILD"
	}

	c , ok := colormapping[card.color]
	if !ok {
		c = colormapping[wish]
	}

	k := new(Kart)
		k.Color = c
		k.Value = v
	return k

}

func convert(kart Kart) Card  {

	valuemapping := make(map[string]int)
	colormapping := make(map[string]int)

	valuemapping["ZERO"] = ZERO
	valuemapping["ONE"] = 1
	valuemapping["TWO"] = 2
	valuemapping["THREE"] = 3
	valuemapping["FOUR"] = 4
	valuemapping["FIVE"] = 5
	valuemapping["SIX"] = 6
	valuemapping["SEVEN"] = 7
	valuemapping["EIGHT"] = 8
	valuemapping["NINE"] = NEUN
	valuemapping["DRAW_TWO"] = PLUS2
	valuemapping["REVERSE"] = WECHSEL
	valuemapping["SKIP"] = STOPP
	valuemapping["WILD"] = WAHL
	valuemapping["WILD_DRAW_FOUR"] = PLUS4


	colormapping["RED"] = RED
	colormapping["BLUE"] = BLUE
	colormapping["GREEN"] = GREEN
	colormapping["YELLOW"] = GELB

	v , ok := valuemapping[kart.Value]
	if !ok {
		v  = -1
	}

	c , ok := colormapping[kart.Color]
	if !ok {
		c = -1
	}

	return Card{
		color:  c,
		symbol: v,
	}

}



type serverConnection struct {
		ipad net.IPAddr

}
/*
resp, err := http.Get("http://example.com/")
...
resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
...
resp, err := http.PostForm("http://example.com/form",
url.Values{"key": {"Value"}, "id": {"123"}})
*/
func serverGame(addresse string , name string)  {
	//bot := new(BaselineBot)

	kvPairs := make(map[string]string)
	kvPairs["name"] = name

	// Make this JSON
	postJson, err := json.Marshal(kvPairs)
	if err != nil { panic(err) }
	fmt.Printf("Sending JSON string '%s'\n", string(postJson))

	// http.POST expects an io.Reader, which a byte buffer does
	postContent := bytes.NewBuffer(postJson)
	fmt.Println(postContent)

	resp, err := http.Post(addresse + "/join","application/json", postContent)

	if err != nil{
		log.Fatal(err)
	}



	fmt.Printf("Status: %s\n", resp.Status)
	buf, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(buf))


	idandplayername := make(map[string]string)
	_ = json.Unmarshal(buf, &idandplayername)
	fmt.Println(idandplayername)


	id := idandplayername["player_id"]

	for ; ;  {

		time.Sleep(time.Millisecond * 100)
		resp1, _ := http.Get(addresse + "/games" + "?id=" + id)
		fmt.Println(addresse + "/games" + "?id=" + id)

		fmt.Printf("Status: %s\n", resp1.Status)
		buf, _ := ioutil.ReadAll(resp1.Body)


		//zugdata := make(map[string]string)

		game := gamestate{}





		err := json.Unmarshal(buf, &game)
		fmt.Println(err)
		//fmt.Println(string(buf))
		fmt.Println(game)

		if game.MyTurn {

			cards := make([]Card, len(game.Hand))

			for i,k := range game.Hand {
				cards[i] = convert(k)
			}


			a := Action{
				karte: convert(game.DiscardedCard),
				block: false,
				wish:  convert(game.DiscardedCard).color,
			}

			bot := 	BaselineBot{
					karten:     cards,
					lastAction: a,
				}



			anext := bot.turn()

			kvPairs := make(map[string]*Kart)

			if anext.block {
				kvPairs["play_card"] = nil

			} else {
				kvPairs["play_card"] = convertReverse(anext.karte, anext.wish)

			}

			// Make this JSON
			postJson, err := json.Marshal(kvPairs)
			if err != nil { panic(err) }
			fmt.Printf("Sending JSON string '%s'\n", string(postJson))

			// http.POST expects an io.Reader, which a byte buffer does
			postContent := bytes.NewBuffer(postJson)
			fmt.Println(postContent)

			resp, err := http.Post(addresse + "/games" + "?id=" + id,"application/json", postContent)
			fmt.Println(resp)


			if err != nil{
				log.Fatal(err)
			}


		}




	}

}