package core

func StandardDistribution(cards [BucketCount][]*Card, sessionLength float32) [BucketCount]float32 {
	ret := [BucketCount]float32{}

	d := [BucketCount]int{
		len(cards[Daily]),
		len(cards[Weekly]),
		len(cards[Monthly]),
		len(cards[Yearly]),
	}

	bucket_multiplier := [BucketCount]float32{8, 4, 2, 1}

	// scale each bucket by duration * bucket_multiplier, keeping track of
	// total.
	var total float32 = 0.0
	for i := 0; i < int(BucketCount); i++ {
		ret[i] = float32(d[i]) * bucket_multiplier[i] * sessionLength
		// keep running total of all buckets
		total += ret[i]
	}

	// normalize and shift each bucket
	bucket_shift := [BucketCount]float32{-.02, -.01, .02, .02}
	for i := 0; i < int(BucketCount); i++ {
		ret[i] /= total
		ret[i] += bucket_shift[i]
	}

	return ret
}
