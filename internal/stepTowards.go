package internal

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
