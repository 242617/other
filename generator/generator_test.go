package generator

import "testing"

var data = map[charset]map[int][]string{
	C_az: {
		1: {"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"},
		6: {"aaaaaa", "aaaaab"},
	},
	C_09: {
		1:  {"0", "1"},
		2:  {"00", "01"},
		10: {"0000000000", "0000000001"},
	},
	C_09 | C_AZ | C_az: {
		1: {
			"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
			"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
			"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		},
		2: {
			"00", "01", "02", "03", "04", "05", "06", "07", "08", "09",
			"0a", "0b", "0c", "0d", "0e", "0f", "0g", "0h", "0i", "0j", "0k", "0l", "0m", "0n", "0o", "0p", "0q", "0r", "0s", "0t", "0u", "0v", "0w", "0x", "0y", "0z",
			"0A", "0B", "0C", "0D", "0E", "0F", "0G", "0H", "0I", "0J", "0K", "0L", "0M", "0N", "0O", "0P", "0Q", "0R", "0S", "0T", "0U", "0V", "0W", "0X", "0Y", "0Z",
			"10", "11", "12", "13", "14", "15", "16", "17", "18", "19",
			"1a", "1b", "1c", "1d", "1e", "1f", "1g", "1h", "1i", "1j", "1k", "1l", "1m", "1n", "1o", "1p", "1q", "1r", "1s", "1t", "1u", "1v", "1w", "1x", "1y", "1z",
			"1A", "1B", "1C", "1D", "1E", "1F", "1G", "1H", "1I", "1J", "1K", "1L", "1M", "1N", "1O", "1P", "1Q", "1R", "1S", "1T", "1U", "1V", "1W", "1X", "1Y", "1Z",
			"20", "21",
		},
	},
}

func Test_generator(t *testing.T) {

	for charset := range data {
		for length, list := range data[charset] {

			gCh := New(length, charset)
			for _, ethalon := range list {

				result := <-gCh
				if result != ethalon {
					t.Fatalf("result %s, wanted %s", result, ethalon)
				}

			}
		}
	}

}
