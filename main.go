package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func isTTY() bool {
	fi, _ := os.Stdout.Stat()
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func red(s string) string {
	if !isTTY() {
		return s
	}
	return "\033[31m" + s + "\033[0m"
}

func green(s string) string {
	if !isTTY() {
		return s
	}
	return "\033[32m" + s + "\033[0m"
}

func generateTestChars(charset []Char, target []Char, rnd *rand.Rand) {
	lastIndex := -1
	for i := 0; i < len(target); i++ {
		charIndex := rand.Intn(len(charset))
		for len(charset) > 1 && charIndex == lastIndex {
			charIndex = rnd.Intn(len(charset))
		}
		target[i] = charset[charIndex]
		lastIndex = charIndex
	}
}

func run(charset []Char, minLen int, maxLen int, requireType bool, fastMode bool, acc bool, rnd *rand.Rand) {
	scanner := bufio.NewScanner(os.Stdin)
	cumCorrect := 0
	cumTotal := 0
	numRounds := 0
	for {
		numRounds++
		targetLen := minLen + rand.Intn(maxLen-minLen+1)
		testChars := make([]Char, targetLen)
		generateTestChars(charset, testChars, rnd)

		testStrBuf := bytes.NewBufferString(fmt.Sprintf("[%02d] ?> What is the romaji for: ", numRounds))
		for _, c := range testChars {
			testStrBuf.WriteRune(c.Char)
		}
		testStrBuf.WriteByte('?')
		testStr := testStrBuf.String()

		fmt.Println(testStr)

		fmt.Printf("[%02d] !> ", numRounds)
		os.Stdout.Sync()

		start := time.Now()
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "exit" || input == "quit" || input == ".q" {
			break
		} else if err := scanner.Err(); err != nil {
			panic(scanner.Err())
		}

		cumTotal += targetLen
		var inputRomajiList []string
		if fastMode {
			input = strings.TrimSpace(input)
			for _, c := range testChars {
				answerLen := len(c.Romaji)
				if len(input) >= answerLen {
					inputRomajiList = append(inputRomajiList, input[:answerLen])
					input = input[answerLen:]
				} else {
					inputRomajiList = append(inputRomajiList, input)
					break
				}
			}
		} else {
			inputRomajiList = strings.Split(input, " ")
		}
		for i, testChar := range testChars {
			inputRomaji := "?"
			if i < len(inputRomajiList) {
				inputRomaji = inputRomajiList[i]
			}
			result := red("incorrect")
			if charMatch(testChar, inputRomaji, requireType) {
				cumCorrect++
				result = green("correct")
			}
			fmt.Printf("%s (%s, %s) - %s was %s.\n", string(testChar.Char), strings.ToLower(testChar.Romaji), testChar.Type, inputRomaji, result)
		}

		dur := time.Now().Sub(start)
		if acc {
			fmt.Printf("%.2fs elapsed. (%.2f ch/s)\n", float64(dur.Nanoseconds())/1000000000., float64(targetLen)/float64(dur.Nanoseconds()/1000000000))
			fmt.Printf("Cumulative score: %d/%d (%.2f%%)\n", cumCorrect, cumTotal, float64(cumCorrect)/float64(cumTotal)*100)
		}
		fmt.Println()
	}
	fmt.Printf("Final score: %d/%d (%.2f%%)\n", cumCorrect, cumTotal, float64(cumCorrect)/float64(cumTotal)*100)
}

func mustAtoi(s string) int {
	res, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return res
}

func main() {
	charset := flag.String("charset", "50", "charset to use, a comma-separated list of characters.")
	sentenceLenDesc := flag.String("length", "5:10", "sentence length")
	requireType := flag.Bool("type", false, "require :h/:k notation when inputting romaji for testing hiragana/katakana classification")
	fastMode := flag.Bool("fast", false, "fast mode, do not require split by space, only use this when your error rate is low and want better speed")
	seed := flag.Int64("seed", time.Now().Unix(), "seed for random number generator")
	acc := flag.Bool("acc", false, "enable accuracy mode, show statistics after each round")

	rnd := rand.New(rand.NewSource(*seed))

	flag.Parse()

	log.Printf("charset: %s, sentence length: %s, require type: %v, fast mode: %v, seed: %d\n", *charset, *sentenceLenDesc, *requireType, *fastMode, *seed)
	log.Printf("Type %s or %s or %s to exit.\n", green("exit"), green("quit"), green(".q"))

	chars := resolveCharsets(*charset)
	minLen, maxLen := 0, 0
	if colon := strings.Index(*sentenceLenDesc, ":"); colon != -1 {
		minLen, maxLen = mustAtoi((*sentenceLenDesc)[:colon]), mustAtoi((*sentenceLenDesc)[colon+1:])
	} else {
		minLen, maxLen = mustAtoi(*sentenceLenDesc), mustAtoi(*sentenceLenDesc)
	}
	run(chars, minLen, maxLen, *requireType, *fastMode, *acc, rnd)
}
