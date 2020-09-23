// Package p contains an HTTP Cloud Function.
package p

import (
	"context"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"google.golang.org/api/sheets/v4"
)

const sheetID = "18J1dfIk2ckKd8885XvytVONG1cYu0Bjo_NP69ZmB6co"

// InsultMe insults me when I need those slap back to reality.
func InsultMe(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	srv, err := sheets.NewService(ctx)
	if err != nil {
		log.Fatalf("Sheet client wtf: %v", err)
	}

	sizeLocation := "F1"
	sizeStr, err := getCell(sizeLocation, srv)
	if err != nil {
		log.Fatalf("can't get size:\n%v", err)
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		log.Fatalf("can't convert size: %v", err)
	}

	limit, err := strconv.Atoi(os.Getenv("PLEASE_NO_MORE"))
	if err != nil {
		log.Fatalf("cannot even limit: %v", err)
	}

	if size > limit {
		size = limit
	}

	i, err := randomNum(size)
	if err != nil {
		log.Fatal(err)
	}

	insult, err := getCell(fmt.Sprintf("A%d", i), srv)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, html.EscapeString(insult))
}

func getCell(cell string, srv *sheets.Service) (string, error) {
	stuff, err := srv.Spreadsheets.Values.Get(sheetID, cell).Do()
	if err != nil {
		return "", err
	}
	if len(stuff.Values) == 0 {
		return "", fmt.Errorf("empty: cell=%v stuff=%v", cell, stuff)
	}
	if len(stuff.Values[0]) == 0 {
		return "", fmt.Errorf("empty again: cell=%v stuff=%v", cell, stuff)
	}
	return stuff.Values[0][0].(string), nil
}

// need this because random number can't be random in cloud functions
func randomNum(max int) (int, error) {
	random, err := http.Get(
		fmt.Sprintf(
			"https://www.random.org/integers/?num=1&min=1&max=%d&col=1&base=10&format=plain&rnd=new",
			max,
		),
	)
	if err != nil {
		return 0, fmt.Errorf("can't even get random number, cringe: %v", err)
	}
	b, err := ioutil.ReadAll(random.Body)
	if err != nil {
		return 0, fmt.Errorf("dafuq: %v", err)
	}
	i, err := strconv.Atoi(string(b))
	if err != nil {
		return 0, fmt.Errorf("tf did i get: val=%v err=%v", string(b), err)
	}
	return i, nil
}
