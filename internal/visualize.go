package internal

import (
	"fmt"
	"math"
)

func VisualizeFrequencyDomain(frequencyDomain []complex64, height int) {
	strengths := []float32{}
	var graphScale float32 = 1

	for _, freq := range frequencyDomain {
		strength := C64Abs(freq)
		if strength > graphScale {
			graphScale = strength
		}
		strengths = append(strengths, strength)
	}

	// make scale incress in steps
	// y = 2 ^ (ceil(log2(x)))
	graphScale = float32(math.Pow(2, math.Ceil(math.Log2(float64(graphScale)))))

	rowHeight := 1 / float32(height)

	for y := range height {
		rowTopStrength := 1 - float32(y)*rowHeight

		fmt.Printf("%5.1f ", rowTopStrength*graphScale)

		for _, strength := range strengths {
			strength /= graphScale

			if strength < rowTopStrength-rowHeight {
				print(" ")
			} else if strength > rowTopStrength {
				print("█")
			} else if strength > rowTopStrength-(rowHeight/8*1) {
				print("▇")
			} else if strength > rowTopStrength-(rowHeight/4) {
				print("▆")
			} else if strength > rowTopStrength-(rowHeight/8*3) {
				print("▅")
			} else if strength > rowTopStrength-(rowHeight/2) {
				print("▄")
			} else if strength > rowTopStrength-(rowHeight/8*5) {
				print("▃")
			} else if strength > rowTopStrength-(rowHeight/4*3) {
				print("▂")
			} else { // strength > rowTopStrength-(rowHeight/8*7)
				print("▁")
			}
		}
		println("")
	}
}
