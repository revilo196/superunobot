package superunobot

import "math"

type ExpBot struct {
	Karten     []Card
	LastAction Action
	nn float64
	ww float64
	x float64
	y float64
}




func (b *ExpBot)nutzen(c Card, liegt Card,  myCards []Card) float64  {
	if c.color == SCHWARZ {
		return 0
	}

	sumOfColor := 0.0
	countOfColor := 0.0
	for _,k := range  myCards {
		if c.color == k.color {
			sumOfColor =+ float64(k.symbol)
			countOfColor += 1
		}


	}

	n :=    (b.x * sumOfColor) + (b.y * (10 - countOfColor))

	return n
}

func (b *ExpBot)nutzwert(c Card, liegt Card,  myCards []Card, wert int, nut float64) float64 {

	nw := (b.ww * float64(wert))  +  (b.nn *  nut)
	return nw
}


func (b *ExpBot) giveCardToPlayer(c Card)  {
	b.Karten = append(b.Karten, c)
}

func (b *ExpBot) update(a Action)  {
	b.LastAction = a
}

func (b *ExpBot) final() bool {
	return len(b.Karten) == 0
}

func (b *ExpBot) points() int {
	sum := 0.0

	for _,k := range b.Karten {
		sum = sum + math.Abs( float64(k.symbol) )
	}

	return int(sum)
}



func (b *ExpBot) turn() Action  {

	var auswahl []Card

	if b.LastAction.karte.symbol == WAHL || b.LastAction.karte.symbol == PLUS4 {
		auswahl = filterCol(b.Karten, b.LastAction.wish, false)
	} else {
		auswahl = filterColAndSym(b.Karten, b.LastAction.karte.color, b.LastAction.karte.symbol, true)
	}

	max := -999.00
	idx := -1
	for i ,k := range auswahl {

		nut := b.nutzen(k,b.LastAction.karte,b.Karten)
		wert:= k.symbol
		nw := b.nutzwert(k,b.LastAction.karte,b.Karten,wert,nut)

		if nw > max {
			max = nw
			idx = i
		}

	}

/*

*/
	if idx == -1 {
		return Action{
			karte: Card{-1,-1},
			block: true,
			wish:  -1,
		}
	} else {
		k := auswahl[idx]

		idxK := -1
		for i, ka := range b.Karten {
			if k.symbol == ka.symbol && k.color == ka.color {
				idxK = i; break
			}
		}

		b.Karten[idxK] = b.Karten[len(b.Karten)-1]
		b.Karten = b.Karten[:len(b.Karten)-1]




		if k.color == SCHWARZ {

			sumROT := 0
			sumGRE := 0
			sumYEL := 0
			sumBLU := 0

			for _, ka := range b.Karten {
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
			} else {
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



func (b *ExpBot) reset()  {
	b.Karten = make([]Card,0)
	b.LastAction = Action{
		karte: Card{0,0},
		block: false,
		wish:  0,
	}
}
