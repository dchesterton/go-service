package service

import (
	"math/rand"
)

type serviceList []*Service

func (l serviceList) Shuffle() {
	n := len(l)

	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		l[i], l[j] = l[j], l[i]
	}
}
