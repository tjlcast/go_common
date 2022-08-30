package rand_utils

import (
	"fmt"
	"testing"
)

func TestWeightChooser(t *testing.T) {
	choice1 := NewChoice("TJL", 1000)
	choice2 := NewChoice("WJ", 1000)
	choice3 := NewChoice("bupt", 500)

	chooser := NewChooser(choice1, choice2, choice3)

	for i:=0; i<100; i++ {
		pick := chooser.Pick()
		fmt.Println("First: ", pick)

		//pick, score := chooser.PickItemAndScore()
		//fmt.Println("Second: ", pick, " - ", score)
	}
}
