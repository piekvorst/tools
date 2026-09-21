package main

import (
	"fmt"
	"math/rand/v2"
)

const (
	consonants        = "bcdfghjkmnpqrstvwxz" // 'l' (lowercase 'L') excluded
	capitalConsonants = "BCDFGHJKLMNPQRSTVWXZ"
	vowels            = "aeiouy"
	capitalVowels     = "AEUY"     // 'I' (capital 'i') and 'O' (capital 'o') excluded
	digits            = "23456789" // '0' (zero) and '1' (one) excluded
)

func pattern() []string {
	var consonantSets, vowelSets []string
	for i := 0; i < 12; i++ {
		consonantSets = append(consonantSets, consonants)
	}
	for i := 0; i < 6; i++ {
		vowelSets = append(vowelSets, vowels)
	}

	var orderedSets []*string
	var i, j int
	for k := 0; k < 6; k++ {
		orderedSets = append(orderedSets, &consonantSets[i])
		i++

		orderedSets = append(orderedSets, &vowelSets[j])
		j++

		orderedSets = append(orderedSets, &consonantSets[i])
		i++
	}

	capital := orderedSets[rand.IntN(len(orderedSets))]
	switch *capital {
	case consonants:
		*capital = capitalConsonants
	case vowels:
		*capital = capitalVowels
	}

	if *capital == capitalConsonants {
		consonantSets[0], *capital = *capital, consonantSets[0]
		consonantSets[1+rand.IntN(len(consonantSets)-1)] = digits
		consonantSets[0], *capital = *capital, consonantSets[0]
	} else {
		consonantSets[rand.IntN(len(consonantSets))] = digits
	}

	var out []string
	for _, s := range orderedSets[:6] {
		out = append(out, *s)
	}
	out = append(out, "-")
	for _, s := range orderedSets[6:12] {
		out = append(out, *s)
	}
	out = append(out, "-")
	for _, s := range orderedSets[12:] {
		out = append(out, *s)
	}
	return out
}

func main() {
	for _, s := range pattern() {
		fmt.Printf("%c", s[rand.IntN(len(s))])
	}

	fmt.Printf("\n")
}
