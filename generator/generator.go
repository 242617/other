package generator

type charset uint64

const (
	// C_09 is Digits: [0-9]
	C_09 charset = 1 << iota
	// C_AZ is Upper case alphabet: [A-Z]
	C_AZ
	// C_az is Lower case alphabet: [a-z]
	C_az
)

func New(length int, chrst charset) chan string {
	ch := make(chan string)

	var charMap []byte

	src := []byte{}
	for i := 0; i < 123; i++ {
		src = append(src, byte(i))
	}

	if chrst&C_09 != 0 {
		charMap = append(charMap, src[48:58]...)
	}
	if chrst&C_az != 0 {
		charMap = append(charMap, src[97:123]...)
	}
	if chrst&C_AZ != 0 {
		charMap = append(charMap, src[65:91]...)
	}

	barr := make([]byte, length)
	end := charMap[len(charMap)-1]

	var fill func(n int)
	fill = func(n int) {

		for i := 0; i < len(charMap); i++ {
			barr[n] = charMap[i]
			if n == length-1 {
				ch <- string(barr)

				if filled(barr, end) {
					close(ch)
					return
				}

			} else {
				fill(n + 1)
			}
		}

	}

	go fill(0)

	return ch
}

func filled(barr []byte, end byte) bool {
	for _, b := range barr {
		if b != end {
			return false
		}
	}
	return true
}
