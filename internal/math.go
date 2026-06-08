package internal

import "math/cmplx"

func StepTowards(start, destination, step float32) float32 {
	if start == destination {
		return destination
	}

	if start < destination {
		if start+step > destination {
			return destination
		} else {
			return start + step
		}
	} else { // start > destination
		if start-step < destination {
			return destination
		} else {
			return start - step
		}
	}
}

func C64Abs(n complex64) float32 {
	// ughhh go please
	return float32(cmplx.Abs(complex128(n)))
}
