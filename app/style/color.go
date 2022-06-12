package style

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
)

func GetColorWeb(str string) tcell.Color {
	ReadHexString := func(str string) (int32, error) {
		if str == "" {
			return 0, nil
		}
		n, err := strconv.ParseUint(str[1:], 16, 64)
		if err != nil {
			return 0, err
		}
		return int32(n), nil
	}

	if str == "default" {
		return tcell.ColorDefault
	}
	i, err := ReadHexString(str)
	if err != nil {
		return tcell.ColorDefault
	}
	return tcell.NewHexColor(i)
}
