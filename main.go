// solve a card game quetstion
// if a player has the chance to check the top 8 cards from the card stack
// and distribute the card to himself his teammate, how could the player and his teammate get the more cards and higher card quality
// means list the highest card quantity and print several choices the palyer could made from the choices.
// considering the case that it is possible to have repeated cards (number).
package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"
)

func main() {
	run()
}

func run() {
	cards := cardGenerator(8)
	fmt.Printf("All cards: %v\n------------------\n", cards)
	if len(cards) < 2 {
		fmt.Printf("Only %v cards available, can not divide to 2 players", len(cards))
		return
	}
	a := lessCardPlayer(cards)
	distros := []distro{}
	for _, i := range a {
		d, ok := newDistro(i, cards)
		if ok {
			distros = append(distros, d)
		}
	}
	switch len(distros) {
	case 0:
		fmt.Println("No Match Found!")
		return
	case 1:
		fmt.Println(distros[0])
	default:
		sort.SliceStable(distros, func(i, j int) bool {
			return distros[i].MaxCardNum() < distros[j].MaxCardNum()
		})
		// printout top 3 distros, not exactly the most card quantity,
		for i, j := len(distros)-1, 0; i > 0; i-- {
			if j < 3 {
				fmt.Println(distros[i])
			}
			j++
		}
	}

}

func lessCardPlayer(s []int) [][]int {
	n := len(s) / 2
	cards, err := nCardMix(s, n)
	if err != nil {
		log.Fatal(err)
	}
	return cards
}

func sum(s []int) int {
	var sum int
	for _, v := range s {
		sum += v
	}
	return sum
}

type distro struct {
	// p1 contains the card player 1 take from the available cards
	p1 []int
	//p2 contains all the max card number (max length of []int) cases the cards' sum == the cards sum of p1
	p2 [][]int
}

func (d distro) MaxCardNum() int {
	max := 0
	idx := 0
	for i, v := range d.p2 {
		if max < len(v) {
			max = len(v)
			idx = i
		}
	}
	return len(d.p1) + len(d.p2[idx])
}

func (d distro) String() string {
	s1 := fmt.Sprintf("P1: %v\n", d.p1)
	s2 := "P2: \n"
	for _, v := range d.p2 {
		s2 += fmt.Sprintf("--> %v\n", v)
	}
	return s1 + s2
}

func newDistro(s []int, src []int) (distro, bool) {
	d := distro{
		p1: s,
	}
	sumOfP1 := sum(s)
	r := exclude(s, src)
	allP2, err := nCardMix(r, len(r))
	if err != nil {
		log.Fatal(err)
	}
	filterP2 := [][]int{}
	for _, items := range allP2 {
		if sum(items) == sumOfP1 {
			filterP2 = append(filterP2, items)
		}
	}
	// if there is a match for p1, set bool to true, else set it to false
	if len(filterP2) == 0 {
		return d, false
	}
	d.p2 = filterP2
	return d, true
}

func exclude(s1 []int, src []int) []int {
	s2 := []int{}
	for _, v := range src {
		var inside bool
		for _, v1 := range s1 {
			if v == v1 {
				inside = true
			}
		}
		if !inside {
			s2 = append(s2, v)
		}
	}
	return s2
}

func appendSlices(target [][]int, srcs ...[][]int) [][]int {
	for _, src := range srcs {
		for _, i := range src {
			target = append(target, i)
		}
	}
	return target
}

func nCardMix(s []int, n int) ([][]int, error) {
	p := [][]int{}
	for i := 1; i <= n; i++ {
		p = append(p, cardMix(s, i)...)
	}
	return p, nil
}

// this is the core of the problem: how to generate all cases if the cards to draw from is given and number of cards to draw is given
func cardMix(s []int, pick int) [][]int {
	p := [][]int{}
	// base case one: len(s) == pick, return [][]int{s}
	if len(s) == pick {
		return [][]int{s}
	}
	// base case two: pick 1 value from slice,
	if pick == 1 {
		r := [][]int{}
		for _, item := range s {
			r = append(r, []int{item})
		}
		return r
	}
	// left the first value from the slice out and get all cases
	leaveOne := cardMix(s[1:], pick)
	// start consider the first value from the slice and get all cases with first value included
	leaveOnePickOneLess := cardMix(s[1:], pick-1)
	addBack := [][]int{}
	for _, item := range leaveOnePickOneLess {
		newItem := append(item, s[0])
		addBack = append(addBack, newItem)
	}
	// add the above 2 cases together, get the final result
	p = appendSlices(p, leaveOne, addBack)
	return p
}

func cardGenerator(num int) []int {
	seed := rand.NewSource(time.Now().Unix())
	r := rand.New(seed)
	cards := make([]int, num)
	for i := 0; i < num; i++ {
		cards[i] = r.Intn(13) + 1
	}
	return cards
}
