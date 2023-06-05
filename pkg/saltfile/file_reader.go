package saltfile

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode"
)

//
// This is the implementation of the below gramma
// Its an recursive gramma, but the implementaion
// is an iterative version of a basic recursive descent
// parsing algorithm, because if the file is large,
// recursion will use up a lot of stack space and
// might result in a stack overflow
//
// DOCUMENT = '[' ENTRIES ']'
// ENTRIES = '[' ENTRY ']' | '[' ENTRY ']', ENTRIES
// ENTRY =  "NAME", "SURNAME", "DATE"
// DATE = YYYY/DD/MM
// NAME = [A-Za-z]+
// SURNAME = [A-Za-z]+
//
// Another solution is to read the whole file into
// memory and split it by the '[' & ']' characters
// to extract the entries, but if the file is very
// big this will use up a lot of memory and result
// in paging which is slow; if if the files are small
// this is the best solution and might be done quickly
// even in a scripting language like python
// The optimal solution is process each entry as
// you read it, similar to what is done
// in streaming xml instead of DOM
//
// One could also read line by line
// but that requires that the file is
// well structured such that entries are at end of
// line boundaries: i chose not to make such an
// assumption.
//
// I chose the below solution because it can
// be adapted to streaming processing where memory
// space is an issue becuase of the size of the file.

// read upto the next none space character
func next(r *bufio.Reader) (rune, error) {
	for {
		c, _, err := r.ReadRune()
		if err == io.EOF {
			return 0, err
		}
		if !unicode.IsSpace(c) {
			return c, nil
		}
	}
}

// golang does not have a peek rune
// the simpliest way to do this is to read and unread a rune
func peekRune(r *bufio.Reader) (rune, error) {
	c, err := next(r)
	if err != nil {
		return 0, err
	}
	err = r.UnreadRune()
	if err != nil {
		return 0, err
	}
	return c, nil
}

func entry(r *bufio.Reader) (string, error) {
	line := make([]rune, 0)
	for {
		c, err := next(r)
		if err != nil {
			return "", err
		}
		if c == '[' { // unexpected start terminal when inside an entry
			return "", errors.New("Unexpected '['")
		}
		if c == ']' {
			return string(line), nil
		}
		line = append(line, c)
	}
}

func entries(r *bufio.Reader) ([]string, error) {
	lines := make([]string, 0)
	for {
		c, err := next(r)
		if err != nil {
			return nil, errors.New("Malformed file")
		}
		if c != '[' {
			return nil, errors.New("Expecting '['")
		}

		line, err := entry(r)
		if err != nil {
			return nil, err
		}

		lines = append(lines, string(line))

		c, err = peekRune(r)
		if err != nil {
			return nil, err
		}
		if c == ',' {
			_, err = next(r) // Eat up the the coma separator
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return lines, nil
}

func Document(r *bufio.Reader) ([]string, error) {
	c, err := next(r)
	if err == io.EOF {
		return nil, nil
	}
	if c != '[' {
		return nil, errors.New("Malformed file")
	}
	c, err = peekRune(r)
	if err != nil {
		return nil, errors.New("Malformed file")
	}
	if c == ']' { // cater for the case of []: file with no entries
		_, err = next(r) // Eat up the the squar bracket
		if err != nil {
			return nil, errors.New("File error")
		}
		_, err := next(r)
		if err == io.EOF {
			return nil, nil
		}
		return nil, errors.New("Malformed file")
	}

	lines, err := entries(r)

	c, err = next(r)
	if err != nil {
		return nil, errors.New("Malformed file")
	}
	if c == ']' {
		_, err := next(r)
		if err == io.EOF {
			return lines, nil
		}
	}

	return nil, errors.New("Malformed file")
}

func GetEntryValues(entry string, c rune) (string, string, string, error) {
	values := strings.Split(entry, string(c))
	if len(values) != 3 {
		return "", "", "", errors.New("Malformed entry")
	}
	return values[0], values[1], values[2], nil
}

func IsBirthday(day, month, year, birthDay, birthMonth int) bool {
	return month == birthMonth && day == birthDay || month == 2 && month == birthMonth && birthDay == 29 && day == 28 && year%4 > 0
}
