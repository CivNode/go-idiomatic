package preferrangeintbad

func Sum(xs []int) int {
	sum := 0
	for i := 0; i < len(xs); i++ { // want `prefer .for i := range xs.`
		sum += xs[i]
	}
	return sum
}
