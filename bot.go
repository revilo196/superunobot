package superunobot

import "math"

/*
type Player interface {
	giveCardToPlayer(c Card)
	update(c Action)
	final() bool
	points() int
	turn() Action
	reset()
}*/

//Parameter


type BaselineBot struct {
	karten []Card
	lastAction Action
}

func filterColAndSym(c []Card, color Color, symbol Symbol, black bool) []Card {
	newC := make([]Card,0)

	for _,k := range c {
		if k.color == color || k.symbol == symbol || (k.color == SCHWARZ && black){
			newC = append(newC, k)
		}
	}
	return  newC
}

func filterCol(c []Card, color Color, black bool) []Card {
	newC := make([]Card,0)

	for _,k := range c {
		if k.color == color  || (k.color == SCHWARZ && black){
			newC = append(newC, k)
		}
	}
	return  newC
}

func (b *BaselineBot) giveCardToPlayer(c Card)  {
 	b.karten = append(b.karten, c)
}

func (b *BaselineBot) update(a Action)  {
	b.lastAction = a
}

func (b *BaselineBot) final() bool {
	return len(b.karten) == 0
}

func (b *BaselineBot) points() int {
	sum := 0.0

	for _,k := range b.karten {
		sum = sum + math.Abs( float64(k.symbol) )
	}

	return int(sum)
}



func (b *BaselineBot) turn() Action  {

	var auswahl []Card

 	if b.lastAction.karte.symbol == WAHL || b.lastAction.karte.symbol == PLUS4 {
	 	auswahl = filterCol(b.karten, b.lastAction.wish, false)
 	} else {
 		auswahl = filterColAndSym(b.karten, b.lastAction.karte.color, b.lastAction.karte.symbol, true)
 	}


	max := -99
	idx := -1
	for i ,k := range auswahl {
		if k.symbol > max  {
			max = k.symbol
			idx = i
		}
	}

	if idx == -1 {
		return Action{
			karte: b.karten[0],
			block: true,
			wish:  b.karten[0].color,
		}
	} else {
		k := auswahl[idx]

		idxK := -1
		for i, ka := range b.karten {
			if k.symbol == ka.symbol && k.color == ka.color {
				idxK = i; break
			}
		}

		b.karten[idxK] = b.karten[len(b.karten)-1]
		b.karten = b.karten[:len(b.karten)-1]




		if k.color == SCHWARZ {

			sumROT := 0
			sumGRE := 0
			sumYEL := 0
			sumBLU := 0

			for _, ka := range b.karten {
				if ka.color == RED {
					sumROT++
				} else if ka.color == GREEN {
					sumGRE++
				} else if ka.color == GELB {
					sumYEL++
				} else if ka.color == BLUE {
					sumBLU++
				}
			}

			if sumROT >= sumGRE && sumROT >= sumYEL &&  sumROT >= sumBLU{
				return Action{
					karte: k,
					block: false,
					wish:  RED,}
			} else if sumGRE >= sumROT && sumGRE >= sumYEL &&  sumGRE >= sumBLU {
					return Action{
						karte: k,
						block: false,
						wish:  GREEN,}
			} else if sumYEL >= sumGRE && sumYEL >= sumROT &&  sumYEL >= sumBLU {
						return Action{
							karte: k,
							block: false,
							wish:  GELB,}
			} else if sumBLU >= sumGRE && sumBLU >= sumYEL &&  sumBLU >= sumROT {
							return Action{
								karte: k,
								block: false,
								wish:  BLUE,}
			}
		}

		return Action{
			karte: k,
			block: false,
			wish:  k.color,}

		}

	}



func (b *BaselineBot) reset()  {
	b.karten = make([]Card,0)
	b.lastAction = Action{
		karte: Card{0,0},
		block: false,
		wish:  0,
	}
}