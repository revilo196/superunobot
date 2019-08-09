package superunobot

import (
	"fmt"
	"testing"
)

func TestTestServer(t *testing.T)  {
	go serverGame("http://localhost:3000", "P1" )
	serverGame("http://localhost:3000",  "P2")
}

func TestTestMain(t *testing.T) {



	pl:= make([]Player, 4)
	pl[0] = new(BaselineBot)
	pl[0].reset()
	pl[1] = new(BaselineBot)
	pl[1].reset()


	pl[3] = new(BaselineBot)
	pl[3].reset()








	for j:= 0; j< 10 ; j++ {
		score:= make([]int, 4)
		plx := new(ExpBot)
		plx.x = 2.5
		plx.y = 0.5 * float64(j)
		plx.nn = 1.0
		plx.ww = 2.0

		pl[2] = plx
		pl[2].reset()

		for i := 0; i < 10000; i++ {
			Game(pl)

			if pl[0].final() {
				score[0] += pl[1].points() + pl[2].points() + pl[3].points()
			} else if pl[1].final() {
				score[1] += pl[0].points() + pl[2].points() + pl[3].points()
			} else if pl[2].final() {
				score[2] += pl[0].points() + pl[1].points() + pl[3].points()
			} else if pl[3].final() {
				score[3] += pl[0].points() + pl[2].points() + pl[1].points()
			}



			pl[0].reset()
			pl[1].reset()
			pl[2].reset()
			pl[3].reset()
		}
		fmt.Print(0.05 * float64(j))
		fmt.Print(" ")
		fmt.Println(score[2])
	}
}
