// Package p contains an HTTP Cloud Function.
package p

import (
	"context"
	"fmt"
	"html"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"google.golang.org/api/sheets/v4"
)

// InsultMe insults me when I need those slap back to reality.
func InsultMe(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	srv, err := sheets.NewService(ctx)
	if err != nil {
		log.Fatalf("Sheet client wtf: %v", err)
	}
	sheetID := "18J1dfIk2ckKd8885XvytVONG1cYu0Bjo_NP69ZmB6co"
	s, err := srv.Spreadsheets.Get(sheetID).Fields("sheets.properties").Do()
	if err != nil {
		log.Fatalf("Sheet data cannot pls: %v", err)
	}
	size := s.Sheets[0].Properties.GridProperties.RowCount

	limit, err := strconv.Atoi(os.Getenv("PLEASE_NO_MORE"))
	if err != nil {
		log.Fatalf("cannot even limit: %v", err)
	}
	if l := int64(limit); size > l {
		size = l
	}

	cell := fmt.Sprintf("A%d", rand.New(rand.NewSource(size)).Intn(int(size)))
	insult, err := srv.Spreadsheets.Values.Get(sheetID, cell).Do()
	if err != nil {
		log.Fatalf("cannot insult reeeee: %v", err)
	}
	if len(insult.Values) == 0 {
		log.Fatalf("it's wrong wtf: %s %v", cell, insult)
	}
	if len(insult.Values[0]) == 0 {
		log.Fatalf("wrong again: %s %v", cell, insult)
	}
	fmt.Fprintf(w, html.EscapeString(insult.Values[0][0].(string)))
}
