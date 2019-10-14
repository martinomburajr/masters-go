package utils

func PermutationsWithRepetitions(values []int) [][]int {
	size := len(values)
	pn := make([]int, size)
	p := make([]int, size)
	response := make([][]int, 0)
	for {
		// generate permutaton
		p = make([]int, size)
		for i, x := range pn {
			p[i] = values[x]
		}
		// show progress
		//fmt.Println(p)
		response = append(response, p)
		// pass to deciding function
		//if decide(p) {
		//	return // terminate early
		//}
		// increment permutation number
		for i := 0; ; {
			pn[i]++
			if pn[i] < size {
				break
			}
			pn[i] = 0
			i++

			if i == size {
				return response // all permutations generated
			}
		}
	}
	return response
}
