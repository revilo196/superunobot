package superunobot

type (
	Stack struct {
		top *node
		length int
	}
	node struct {
		value Card
		prev *node
	}
)
// Create a new stack
func NewCardStack() *Stack {
	return &Stack{nil,0}
}
// Return the number of items in the stack
func (this *Stack) Len() int {
	return this.length
}
// View the top item on the stack
func (this *Stack) Peek() Card {
	if this.length == 0 {
		return Card{
			color:  -1,
			symbol: -1,
		}
	}
	return this.top.value
}
// Pop the top item of the stack and return it
func (this *Stack) Pop() Card {
	if this.length == 0 {
		return Card{
			color:  -1,
			symbol: -1,
		}
	}

	n := this.top
	this.top = n.prev
	this.length--
	return n.value
}
// Push a value onto the top of the stack
func (this *Stack) Push(value Card) {
	n := &node{value,this.top}
	this.top = n
	this.length++
}
