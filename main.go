package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// definition of structure for the river sides
// containing variable 'm' and 'c'.
// The former is the number of missionaries and
// the latter is the number of cannibals.
type RiverSide struct {
	m, c int
}

// definition of the state
// which contains 2 river sides, left and right,
// and string variable 'boat' that describes
// which river side the boat is berthed
type State struct {
	left, right RiverSide
	boat        string // "left" or "right"
}

// The initial state:
// There are 3 missionaries, 3 cannibals, 1 boat
// in the left side of the river.
// There are none in the right side.
var initialState State = State{
	left:  RiverSide{3, 3},
	right: RiverSide{0, 0},
	boat:  "left",
}

// The goal state:
// There are 3 missionaries, 3 cannibals, 1 boat
// in the right side of the river.
// There are none in the left side.
var goalState State = State{
	left:  RiverSide{0, 0},
	right: RiverSide{3, 3},
	boat:  "right",
}

// Definition of structure for the operation
// which is applied to the state.
// The struct 'Operator' describes the persons
// boarding the boat.
// It contains two int variables, 'm' and 'c'.
// For example, Operator {1, 1} is meaning that
// 1 missionary and 1 cannibal board the
// boat shipping to opposite shore of the river.
type Operator struct {
	m, c int
}

// The variable 'Operators' is an array of Operator,
// and it contains the selectable options.
// Up to 2 people can ride on our boat, So the var
// 'Operators' is following:
var Operators = [5]Operator{
	{2, 0},
	{1, 0},
	{1, 1},
	{0, 1},
	{0, 2},
}

// Function 'valid' is the validator for the given state.
// It decide whether the given state is safe.
// The rule of this problem says, "On both river sides,
// if the number of cannibals is more than the number of
// missionaries, cannibals eat missionaries.
// So the function 'valid' returns 'true' when
// no missionaries is eaten.
// In addition, this 'valid' checks whether the variables for the
// numbers of people on both sides are not negative.
func valid(state State) bool {
	switch {
	case state.left.m < 0 || state.left.c < 0 || state.right.m < 0 || state.right.c < 0:
		return false
	case state.left.m > 0 && state.left.c > state.left.m:
		return false
	case state.right.m > 0 && state.right.c > state.right.m:
		return false
	default:
		return true
	}
}

// The 'stateTransition' is the state transition
// function of the finite automaton.
// So, it is given one state of current and
// one operator. If the given 'currentState' can
// accept the Operator 'op' and the 'currentState'
// can be properly changed into the next state,
// this function will returns 'nextState'.
// Otherwise, it should report error by 'ok' of false.
func stateTransition(currentState State, op Operator) (nextState State, ok bool) {

	var from, to *RiverSide

	if currentState.boat == "left" {
		from, to = &currentState.left, &currentState.right
		nextState.boat = "right"
		nextState.right = RiverSide{to.m + op.m, to.c + op.c}
		nextState.left = RiverSide{from.m - op.m, from.c - op.c}
	} else {
		from, to = &currentState.right, &currentState.left
		nextState.boat = "left"
		nextState.left = RiverSide{to.m + op.m, to.c + op.c}
		nextState.right = RiverSide{from.m - op.m, from.c - op.c}
	}

	ok = valid(nextState)

	return
}

var q []State

// prints (left side), (right side), (boat position)
func printState(currentState State) {
	fmt.Printf("State: {(%dM%dC), (%dM%dC), %s}\n", currentState.left.m, currentState.left.c, currentState.right.m, currentState.right.c, currentState.boat)
}

func insert(i int, s State) {
	var h1 []State
	var h2 []State

	for index, value := range q {
		if index < i {
			h1 = append(h1, value)
		} else {
			h2 = append(h2, value)
		}
	}
	q = append(h1, s)
	for i := 0; i < len(h2); i++ {
		q = append(q, h2[i])
	}

}

func printQueue() {
	for index := range q {
		fmt.Printf("index: %d ", index)
		printState(q[index])
	}
}

/*
	NON-DETERMINISTIC SEARCH IMPLEMENTATION
    --------------------------------------
	I remove the first element in the Queue
	then I find all reachable states from that element
	then I check if I've already generated that state before, (check if its in hashmap)
	if the state is "new", I insert that state at a random location in the Queue

	My state representation is as follows
		{(#M#C), (#M#C), left/right}

	Remark:
		A -> C -> B
		A -> D -> B (B is already discovered, cannot be added to queue)
		BUT
		A -> B
		A -> C -> B (B is a new state here because of the boat's position)

		e.g.
		{(2M2C), (1M1C), left} != {(2M2C), (1M1C), right}


*/
func nonDeterminism() bool {
	// 'history' is the map which contains the visited states
	history := map[State]bool{initialState: true}

	// random number generator
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	// loop until queue is empty or goal found
	for len(q) != 0 {
		fmt.Println("**********\nContents of Queue")
		printQueue()

		// dequeue
		s := q[0]
		q = q[1:]

		// check goal
		if s == goalState {
			fmt.Printf("Goal Found! - ")
			printState(s)
			return true
		}

		fmt.Printf("Expanding ")
		printState(s)
		length := len(q)

		// create new paths
		for _, op := range Operators {
			// create transitions, legal or not
			nextState, ok := stateTransition(s, op)
			if ok { // if it is legal transition
				if !history[nextState] {

					// add it to the history
					history[nextState] = true
					// non-determinism (rand-num)
					i := r1.Intn(len(q) + 1)

					fmt.Printf("\t-> Operation: (%d M, %d C) to %s\n", op.m, op.c, nextState.boat)
					fmt.Printf("\t     index: %d, ", i)
					printState(nextState)

					// insert in queue at i'th position
					insert(i, nextState)
				} else {
					fmt.Printf("\t-> Operation: (%d M, %d C) to %s\n", op.m, op.c, nextState.boat)
					fmt.Printf("\t     (!)repeated ")
					printState(nextState)
				}
			}
		}
		if len(q) == length {
			fmt.Println("\tdead-end")
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\n<continue: >")
		text, _ := reader.ReadString('\n')
		text += ""
	}
	return false
}

func main() {

	// queue
	q = append(q, initialState)

	if !nonDeterminism() {
		fmt.Println("No solutions exist!")
	}
}
