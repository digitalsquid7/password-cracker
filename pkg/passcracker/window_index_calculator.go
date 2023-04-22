package passcracker

type windowIndexCalculator struct {
	totalSize  int
	windowSize int
}

func (c *windowIndexCalculator) calculate(i int) (left int, right int) {
	left = i * c.windowSize
	if i == c.totalSize/c.windowSize-1 {
		right = c.totalSize
	} else {
		right = left + c.windowSize
	}
	return left, right
}
