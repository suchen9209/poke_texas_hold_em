package models

type Card struct {
	Value int
	Color string
}

const (
	HEART   = "heart"
	DIAMOND = "diamond"
	SPADE   = "spade"
	CLUB    = "club"
)

const (
	POKER_NUMBER_2  = 2
	POKER_NUMBER_3  = 3
	POKER_NUMBER_4  = 4
	POKER_NUMBER_5  = 5
	POKER_NUMBER_6  = 6
	POKER_NUMBER_7  = 7
	POKER_NUMBER_8  = 8
	POKER_NUMBER_9  = 9
	POKER_NUMBER_10 = 10
	POKER_NUMBER_J  = 11
	POKER_NUMBER_Q  = 12
	POKER_NUMBER_K  = 13
	POKER_NUMBER_A  = 14
)

var cardMap map[int]Card

func InitCardMap() {
	cardMap = make(map[int]Card, 52)
	suitsArr := [4]string{HEART, DIAMOND, SPADE, CLUB}
	index := 1
	for i := POKER_NUMBER_2; i <= POKER_NUMBER_A; i++ {
		for _, v := range suitsArr {
			poker := Card{Value: i, Color: v}
			cardMap[index] = poker
			index++
		}
	}
}

func GetOneCard() *Card {
	if len(cardMap) == 0 {
		panic("no cards error")
	}
	var card = new(Card)
	for key, v := range cardMap {
		card = &v
		delete(cardMap, key)
		break
	}
	return card

}
