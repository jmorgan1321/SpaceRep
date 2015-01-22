package core

import "math/rand"

type Distribution [BucketCount]float32
type Deck [BucketCount][]*Card

func (d *Deck) GetCards(seed int64, dist Distribution) <-chan *Card {
	r := rand.New(rand.NewSource(42))
	ch := make(chan *Card)

	go func() {
		for {
			// get random distribution selection between [0...1]
			x := r.Float32()

			bucket := Daily
			switch {
			case x < dist[Yearly]:
				bucket = Yearly
			case x < dist[Monthly]:
				bucket = Monthly
			case x < dist[Weekly]:
				bucket = Weekly
			}

			if len(d[bucket]) == 0 {
				continue
			}

			// pick random card from bucket
			i := r.Intn(len(d[bucket]))
			ch <- d[bucket][i]
		}
	}()

	return ch
}

// TODO: spec this
func (d *Deck) Count() int {
	i := 0
	for _, b := range d {
		for _ = range b {
			i++
		}
	}
	return i
}
