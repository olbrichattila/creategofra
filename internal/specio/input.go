package specio

import (
	"fmt"
	"log"
	"strings"

	"github.com/eiannone/keyboard"
)

const (
	hideCursor  = "\033[?25l"
	showCursor  = "\033[?25h"
	cursorBegin = "\033[45m"
	cursorEnd   = "\033[0m"
)

func Input(prompt, defaultTxt string) string {
	result := defaultTxt
	curPos := len(defaultTxt)
	maxLen := curPos

	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	for {
		displayText(prompt, result, curPos, maxLen)
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if key == 3 {
			fmt.Print(showCursor)
			panic("Application manually terminated")
			//	os.Exit(0) // This brakes the console
		}

		if key == 0 {
			result = recordKeyPress(result, curPos, char)
			curPos++
			if curPos > maxLen {
				maxLen = curPos
			}
			continue
		}

		if key == keyboard.KeyArrowLeft && curPos > 0 {
			curPos--
			continue
		}

		if key == keyboard.KeyArrowRight && curPos < len(result) {
			curPos++
			if curPos > maxLen {
				maxLen = curPos
			}
			continue
		}

		if (key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2) && curPos > 0 {
			result = removeCharAtCursor(result, curPos)
			curPos--
			continue
		}

		if key == keyboard.KeyDelete && curPos < len(result) {
			result = removeCharAfterCursor(result, curPos)
			continue
		}

		if key == keyboard.KeyHome {
			curPos = 0
			continue
		}

		if key == keyboard.KeyEnd {
			curPos = len(result)
			continue
		}

		if key == keyboard.KeyEsc || key == keyboard.KeyEnter || key == keyboard.KeyTab {
			fmt.Print(showCursor)
			break
		}
	}

	return result
}

func displayText(prompt, s string, p, m int) {
	fmt.Print("\r" + prompt + strings.Repeat(" ", m+2))
	result := hideCursor + "\r" + prompt
	if len(s) == p {
		fmt.Print(result + s + cursorBegin + " " + cursorEnd)
		return
	}
	for i, c := range s {
		if i == p {
			result = result + cursorBegin + string(c) + cursorEnd
			continue
		}
		result = result + string(c)
	}

	fmt.Print(result)
}

func recordKeyPress(s string, p int, c rune) string {
	if len(s) == p {
		return s + string(c)
	}

	result := ""
	for i, l := range s {
		if i == p {
			result = result + string(c)
		}
		result = result + string(l)
	}

	return result
}

func removeCharAtCursor(s string, p int) string {
	result := ""
	for i, l := range s {
		if i != p-1 {
			result = result + string(l)
		}
	}

	return result
}

func removeCharAfterCursor(s string, p int) string {
	result := ""
	for i, l := range s {
		if i != p {
			result = result + string(l)
		}
	}

	return result
}
