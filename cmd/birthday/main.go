package main

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jsamunderu/salty/pkg/saltfile"
)

var text = `
[
    ["Doe", "John", "1982/10/08"],
    ["Wayne", "Bruce", "1965/01/30"],
    ["Gaga", "Lady", "1986/03/28"],
    ["Curry", "Mark", "1988/02/29"]
]
`

func main() {
	r := bufio.NewReader(strings.NewReader(text))

	lines, err := saltfile.Document(r)
	if err != nil {
		log.Fatal(err.Error())
	}

	thisYear, thisMonth, today := time.Now().Date()

	for _, line := range lines {
		line = strings.ReplaceAll(line, "\"", "")
		fmt.Println(line)
		name, surname, birthday, err := saltfile.GetEntryValues(line, ',')
		if err != nil {
			log.Println(err.Error())
		} else {
			fmt.Printf("name[%s] surname[%s] birthday[%s]\n", name, surname, birthday)
			year, month, date, err := saltfile.GetEntryValues(birthday, '/')
			if err != nil {
				log.Println(err.Error())
			} else {
				fmt.Printf("\tyear[%s] month[%s] date[%s]\n", year, month, date)
				dt, err := strconv.Atoi(date)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				mnth, err := strconv.Atoi(month)
				if err != nil {
					log.Println(err.Error())
					continue
				}

				if saltfile.IsBirthday(today, int(thisMonth), thisYear, dt, mnth) {

					fmt.Println("Happy birthday")
				}
			}
		}
	}
}
