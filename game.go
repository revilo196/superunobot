package superunobot

import (
	"math/rand"
	"time"
)

type Color = int
type Symbol = int

const (
	RED 	Color = 1
	GREEN 	Color = 2
	BLUE 	Color = 3
	GELB 	Color = 4
	SCHWARZ Color = 0
)

const (
	ZERO 	Symbol = 0
	EINS 	Symbol = 1
	ZWEI 	Symbol = 2
	DREI 	Symbol = 3
	VIER    Symbol = 4
	FUMF	Symbol = 5
	SECHS	Symbol = 6
	SIEBEN	Symbol = 7
	ACHT	Symbol = 8
	NEUN	Symbol = 9
	WECHSEL Symbol = 19
	STOPP 	Symbol = 20
	PLUS2	Symbol = 21
	WAHL    Symbol = -50
	PLUS4   Symbol = -51
)

type Card struct {
	color Color
	symbol Symbol
}



type Action struct {
	karte Card
	block bool
	wish Color
}

type Player interface {
	giveCardToPlayer(c Card)
	update(c Action)
	final() bool
	points() int
 	turn() Action
	reset()
}

type SpielData struct {
	karten Stack
	ablage Stack
	spieler []Player
	zugnummer int
	currentPlayer int
}

func ErzeugeKarten() []Card  {
	cards := make([]Card,112)
	counter := 0


	for i := 1 ; i <= 4; i++ {
		c := Color(i)

		for  j:=0; j<=9 ; j++  {
			s := Symbol(j)
			cards[counter] = Card{
				color:  c,
				symbol: s,
			}
			counter++
		}

		for  j:=1; j<=9 ; j++  {
			s := Symbol(j)
			cards[counter] = Card{
				color:  c,
				symbol: s,
			}
			counter++
		}

		cards[counter] = Card{c, WECHSEL}; counter++
		cards[counter] = Card{c, WECHSEL}; counter++
		cards[counter] = Card{c, PLUS2}; counter++
		cards[counter] = Card{c, PLUS2}; counter++
		cards[counter] = Card{c, STOPP}; counter++
		cards[counter] = Card{c, STOPP}; counter++

	}

	cards[counter] = Card{SCHWARZ, WAHL}; counter++
	cards[counter] = Card{SCHWARZ, WAHL}; counter++
	cards[counter] = Card{SCHWARZ, WAHL}; counter++
	cards[counter] = Card{SCHWARZ, WAHL}; counter++
	cards[counter] = Card{SCHWARZ, WAHL}; counter++
	cards[counter] = Card{SCHWARZ, WAHL}; counter++
	cards[counter] = Card{SCHWARZ, WAHL}; counter++
	cards[counter] = Card{SCHWARZ, WAHL}; counter++

	cards[counter] = Card{SCHWARZ, PLUS4}; counter++
	cards[counter] = Card{SCHWARZ, PLUS4}; counter++
	cards[counter] = Card{SCHWARZ, PLUS4}; counter++
	cards[counter] = Card{SCHWARZ, PLUS4}; counter++

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	return cards
}


func Game(player []Player)  {

	kartenRAW := ErzeugeKarten()
	karten := NewCardStack()
	ablage := NewCardStack()
	spieler := 0
	counter := 0

	dir := 1

	//Stapel Erstellen
	for _,k := range kartenRAW {
		karten.Push(k)
	}

	//Karten austeilen
	for i:=0 ; i < 7 ; i++  {
		for i,_ := range  player {
			player[i].giveCardToPlayer(karten.Pop())
		}

	}

	for i := SCHWARZ; i==SCHWARZ ;  {
		ablage.Push(karten.Pop())
		i = ablage.Peek().color
	}

	lastAction := Action{
		karte: ablage.Peek(),
		block: false,
		wish:  ablage.Peek().color,
	}

	for ;  ;  {

		if karten.Len() <= 6  {
			kartenNEU := make([]Card,0)

			for ;karten.Len()>0; {
				ablage.Push(karten.Pop())
			}

			for ;ablage.Len()>0; {
				kartenNEU = append(kartenNEU,ablage.Pop())
			}

			rand.Shuffle(len(kartenNEU), func(i, j int) {
				kartenNEU[i], kartenNEU[j] = kartenNEU[j], kartenNEU[i]
			})

			for _,k := range kartenNEU {
				karten.Push(k)
			}
		}

	//	fmt.Print(LastAction)
	//	fmt.Println(spieler)

		for i,_ := range  player {
			player[i].update(lastAction)
		}

		if spieler < 0 {
			spieler = spieler + len(player)
		}

		a := player[spieler].turn()

		if player[spieler].final() {
			//fmt.Printf("Spieler %d hat gewonnen", spieler)
			break
		}
		//Spieler kann nicht
		if a.block {
			//Spieler nimmt karte
			player[spieler].giveCardToPlayer(karten.Pop())

			//2. Versuch
			a2 := player[spieler].turn()
			if a2.block {
				//nothing
			} else {
				ablage.Push(a2.karte)
				//fmt.Println(a2)

				if a.karte.symbol == WECHSEL {
					dir = dir * -1
					spieler = (spieler + dir) % len(player)
				} else if a.karte.symbol == STOPP {
					spieler = (spieler + dir + dir) % len(player)
				} else if a.karte.symbol == PLUS2 {
					give := (spieler + dir) %  len(player)
					if give < 0 {give += len(player)}
					player[give].giveCardToPlayer(karten.Pop())
					player[give].giveCardToPlayer(karten.Pop())
					spieler = (spieler + dir + dir) % len(player)
				} else if a.karte.symbol == WAHL {

					spieler = (spieler + dir) % len(player)
				} else if a.karte.symbol == PLUS4 {
					give := (spieler + dir) %  len(player)
					if give < 0 {give += len(player)}
					player[give].giveCardToPlayer(karten.Pop())
					player[give].giveCardToPlayer(karten.Pop())
					player[give].giveCardToPlayer(karten.Pop())
					player[give].giveCardToPlayer(karten.Pop())
					spieler = (spieler + dir + dir) % len(player)
				} else {
					spieler = (spieler + dir) % len(player)
				}

				lastAction = a2
			}
		}  else {
			ablage.Push(a.karte)
			//fmt.Println(a)

			if a.karte.symbol == WECHSEL {
				dir = dir * -1
				spieler = (spieler + dir) % len(player)
			} else if a.karte.symbol == STOPP {
				spieler = (spieler + dir + dir) % len(player)
			} else if a.karte.symbol == PLUS2 {
				give := (spieler + dir) %  len(player)
				if give < 0 {give += len(player)}
				player[give].giveCardToPlayer(karten.Pop())
				player[give].giveCardToPlayer(karten.Pop())
				spieler = (spieler + dir + dir) % len(player)
			} else if a.karte.symbol == WAHL {

				spieler = (spieler + dir) % len(player)
			} else if a.karte.symbol == PLUS4 {
				give := (spieler + dir) %  len(player)
				if give < 0 {give += len(player)}
				player[give].giveCardToPlayer(karten.Pop())
				player[give].giveCardToPlayer(karten.Pop())
				player[give].giveCardToPlayer(karten.Pop())
				player[give].giveCardToPlayer(karten.Pop())
				spieler = (spieler + dir + dir) % len(player)
			} else {
				spieler = (spieler + dir) % len(player)
			}

			lastAction = a
		}



		counter++
	}



}
