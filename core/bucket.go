package core

import (
	"github.com/jmorgan1321/SpaceRep/utils"
)

// Bucket represents the frequency this card should be reviewed.
type Bucket int

func (b Bucket) String() string {
	switch b {
	case Daily:
		return "Daily"
	case Weekly:
		return "Weekly"
	case Monthly:
		return "Monthly"
	case Yearly:
		return "Yearly"
	default:
		panic("Unknown bucket")
	}
}

func (b Bucket) GetMaxCount() int {
	return utils.Pow(2, int(bucket_count-b)-1)
}

func (b Bucket) NextBucket() Bucket {
	switch b {
	case Daily:
		return Weekly
	case Weekly:
		return Monthly
	case Monthly:
		return Yearly
	case Yearly:
		return Yearly
	}
	return Daily
}
func (b Bucket) PrevBucket() Bucket {
	switch b {
	case Daily:
		return Daily
	case Weekly:
		return Daily
	case Monthly:
		return Weekly
	case Yearly:
		return Monthly
	}
	return Daily
}

const (
	Daily Bucket = iota
	Weekly
	Monthly
	Yearly
	bucket_count
)

const (
	BucketCount = int(bucket_count)
)
