package png

// Paeth Predictor algorithm from the RFC 2083
// https://www.rfc-editor.org/rfc/rfc2083#section-6
// 6.6 [Page 36]
//
// a = left, b = above, c = upper left
func paethPredictor(a, b, c int) int {
	p := a + b - c   // initial estimate
	pa := abs(p - a) // distance to a, b, c
	pb := abs(p - b)
	pc := abs(p - c)
	// return nearest of a, b, c
	// breaking ties in order a, b, c.
	if pa <= pb && pa <= pc {
		return a
	} else if pb <= pc {
		return b
	} else {
		return c
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
