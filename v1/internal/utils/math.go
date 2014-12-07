package utils

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Clamp(val, min, max int) int {
	return Min(Max(val, min), max)
}

func Pow(base, exp int) int {
	if exp == 0 {
		return 1
	}
	ret := base
	for i := exp - 1; i > 0; i-- {
		ret *= base
	}
	return ret
}
