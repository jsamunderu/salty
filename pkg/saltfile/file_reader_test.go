package saltfile

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"testing"
)

func Test_NextEof(t *testing.T) {
	var text = ``
	r := bufio.NewReader(strings.NewReader(text))
	_, err := next(r)
	if err != io.EOF {
		t.Errorf("next: expected: %v got %v", 'v', err.Error())
	}
}

func Test_NextEofSpace(t *testing.T) {
	var text = `
	
	
	                                                                   
	
	`
	r := bufio.NewReader(strings.NewReader(text))
	_, err := next(r)
	if err != io.EOF {
		t.Errorf("next: expected: %v got %v", 'v', err.Error())
	}
}

func Test_Next(t *testing.T) {
	var text = `v`
	r := bufio.NewReader(strings.NewReader(text))
	c, err := next(r)
	if err != nil {
		t.Errorf("next: expected: %v got %v", 'v', err.Error())
	}
	if c != 'v' {
		t.Errorf("next: expected: %v got %v", 'v', string(c))
	}
	t.Logf("next: Read %v", string(c))
}

func Test_NextSpace(t *testing.T) {
	var text = `
	
	v
	
	`
	r := bufio.NewReader(strings.NewReader(text))
	c, err := next(r)
	if err != nil {
		t.Errorf("next: expected: %v got %v", 'v', err.Error())
	}
	if c != 'v' {
		t.Errorf("next: expected: %v got %v", 'v', string(c))
	}
	t.Logf("next: Read %v", string(c))
}

func Test_entryEmpty(t *testing.T) {
	var text = ``
	r := bufio.NewReader(strings.NewReader(text))
	_, err := entry(r)
	if err != io.EOF {
		t.Errorf("entry: expected: %v got %v", io.EOF.Error(), err.Error())
		return
	}
	t.Logf("entry: expected: %v", err.Error())
}

func Test_entryNonEmpty(t *testing.T) {
	var text = `abcdefghijklmnopqrstuv xy z`
	r := bufio.NewReader(strings.NewReader(text))
	_, err := entry(r)
	if err != io.EOF {
		t.Errorf("entry: expected: %v got %v", io.EOF.Error(), err.Error())
		return
	}
	t.Logf("entry: expected: %v", err.Error())
}

func Test_entryEmptyTeminated(t *testing.T) {
	var text = `]`
	r := bufio.NewReader(strings.NewReader(text))
	line, err := entry(r)
	if err != nil {
		t.Errorf("entry: expected: %v got %v", io.EOF.Error(), err.Error())
		return
	}
	t.Logf("entry: expected: %v", line)
}

func Test_entryNonEmptyTeminated(t *testing.T) {
	var text = `abcdefghijklmnopqrstuv xy z]`
	r := bufio.NewReader(strings.NewReader(text))
	line, err := entry(r)
	if err != nil {
		t.Errorf("entry: expected: %v got %v", io.EOF.Error(), err.Error())
		return
	}
	t.Logf("entry: expected: %v", string(line))
}

func Test_entryReentryTerm(t *testing.T) {
	var text = `[`
	r := bufio.NewReader(strings.NewReader(text))
	_, err := entry(r)
	if err == nil {
		t.Errorf("entry: expected: %v got nil", errors.New("Malformed file").Error())
		return
	}
	t.Logf("entry: expected: %v", err.Error())
}

func Test_entryGetFirst(t *testing.T) {
	var text = `xxxx],[yyyy]`
	r := bufio.NewReader(strings.NewReader(text))
	line, err := entry(r)
	if err != nil {
		t.Errorf("entry: expected: %v got %v", io.EOF.Error(), err.Error())
		return
	}
	if string(line) != "xxxx" {
		t.Errorf("entry: expected: %v got %v", "xxxx", string(line))
	}
	t.Logf("entry: expected: %v", string(line))
}

func Test_entryEmptyTerm(t *testing.T) {
	var text = `],[yyyy]`
	r := bufio.NewReader(strings.NewReader(text))
	line, err := entry(r)
	if err != nil {
		t.Errorf("entry: expected: %v got %v", io.EOF.Error(), err.Error())
		return
	}
	if len(line) != 0 {
		t.Errorf("entry: expected: %v got %v", "xxxx", string(line))
	}
	t.Logf("entry: expected: %v", string(line))
}

func Test_entriesEmpty(t *testing.T) {
	var text = ``
	r := bufio.NewReader(strings.NewReader(text))
	_, err := entries(r)
	if err == nil {
		t.Errorf("entries: expected: %v got nil", errors.New("Malformed file").Error())
		return
	}

	t.Logf("entries: expected: %v", err.Error())
}

func Test_entriesTerm(t *testing.T) {
	var text = `]`
	r := bufio.NewReader(strings.NewReader(text))
	_, err := entries(r)
	if err == nil {
		t.Errorf("entries: expected: %v got nil", io.EOF)
		return
	}

	t.Logf("entries: expected: %v", err.Error())
}

func Test_entriesEmptyTerms(t *testing.T) {
	var text = `[]`
	r := bufio.NewReader(strings.NewReader(text))
	_, err := entries(r)
	if err != io.EOF {
		t.Errorf("entries: expected: %v got %v", io.EOF, err.Error())
		return
	}
}

func Test_entriesEmptyTermsMulti(t *testing.T) {
	var text = `[],[]`
	r := bufio.NewReader(strings.NewReader(text))
	_, err := entries(r)
	if err != io.EOF {
		t.Errorf("entries: expected: %v got %v", io.EOF, err.Error())
		return
	}
}

func Test_entriesEmptyTermsMalformed(t *testing.T) {
	var text = `[],`
	r := bufio.NewReader(strings.NewReader(text))
	_, err := entries(r)
	if err == io.EOF {
		t.Errorf("entries: expected: %v got %v", errors.New("Malformed file").Error(), err)
		return
	}
}

func Test_documentEmpty(t *testing.T) {
	var text = ``
	r := bufio.NewReader(strings.NewReader(text))
	_, err := Document(r)
	if err != nil {
		t.Errorf("document: expected: nil got %v", err.Error())
		return
	}
}

func Test_documentEmptyTerms(t *testing.T) {
	var text = `[]`
	r := bufio.NewReader(strings.NewReader(text))
	_, err := Document(r)
	if err != nil {
		t.Errorf("document: expected: nil got %v", err.Error())
		return
	}
}

func Test_documentTermsEmpty(t *testing.T) {
	var text = `[[]]`
	r := bufio.NewReader(strings.NewReader(text))
	lines, err := Document(r)
	if err != nil {
		t.Errorf("document: expected: nil got %v", err.Error())
		return
	}
	if len(lines) != 1 {
		t.Errorf("document: expected: 1 got %v", len(lines))
	}
}

func Test_documentMalformedOuter(t *testing.T) {
	var text = `p[[]]`
	r := bufio.NewReader(strings.NewReader(text))
	_, err := Document(r)
	if err == nil {
		t.Errorf("document: expected: %v got nil", errors.New("Malformed file").Error())
		return
	}
}
func Test_documentMalformedInner(t *testing.T) {
	var text = `[p[]]`
	r := bufio.NewReader(strings.NewReader(text))
	_, err := Document(r)
	if err == nil {
		t.Errorf("document: expected: %v got nil", errors.New("Malformed file").Error())
		return
	}
}

func TestGetEntryValuesEntry(t *testing.T) {
	var text = `Gaga,Lady,1986/03/28`
	name, surname, birthday, err := GetEntryValues(text, ',')
	if err != nil {
		t.Errorf("GetEntryValues expected: nil got: %v", err)
	}
	if name != "Gaga" && surname != "Lady" && birthday != "1986/03/28" {
		t.Errorf("name[%v] surname[%v] birthname[%v]", name, surname, birthday)
	}
}

func TestGetEntryValuesBirthday(t *testing.T) {
	var text = `1986/03/28`
	year, month, day, err := GetEntryValues(text, '/')
	if err != nil {
		t.Errorf("GetEntryValues expected: nil got: %v", err)
	}
	if year != "1986" && month != "03" && day != "28" {
		t.Errorf("year[%v] month[%v] day[%v]", year, month, day)
	}
}

func TestIsBirthday(t *testing.T) {
	thisMonth, today, thisYear, day, month := 2, 28, 2005, 29, 2
	if IsBirthday(today, thisMonth, thisYear, day, month) {
		t.Errorf("Expected: false got: true")
	}
}

func TestIsBirthdayLeapYear(t *testing.T) {
	thisMonth, today, thisYear, day, month := 2, 28, 2004, 29, 2
	if IsBirthday(today, thisMonth, thisYear, day, month) {
		t.Errorf("Expected: true got: false")
	}
}
