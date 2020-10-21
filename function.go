// Package p contains an HTTP Cloud Function.
package p

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"google.golang.org/api/sheets/v4"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

var sheetID = os.Getenv("SHEET_ID")

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

	charLimit, err := strconv.Atoi(os.Getenv("CHAR_LIMIT"))
	if err != nil {
		log.Fatalf("char limit: %v", err)
	}

	if len(insult) > charLimit {
		// text to speech api ain't cheap
		insult = insult[:charLimit]
	}

	type Response struct {
		Message  string `json:"message"`
		Subtitle string `json:"subtitle"`
	}

	j, err := json.Marshal(Response{Message: insult, Subtitle: insult})
	if err != nil {
		panic(err)
	}

	// insultButInMp3, err := createAudio(ctx, insult)
	// if err != nil {
	// log.Fatalf("audio: %v", err)
	// }

	w.Header().Add("Content-Type", "application/json")
	if _, err = w.Write(j); err != nil {
		log.Fatalf("writing: %v", err)
	}
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
	i, err := strconv.Atoi(strings.Trim(string(b), " \n\t"))
	if err != nil {
		return 0, fmt.Errorf("tf did i get: val=%v err=%v", string(b), err)
	}
	return i, nil
}

// copied from text to speech api docs
func createAudio(ctx context.Context, text string) ([]byte, error) {
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_FEMALE,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp.GetAudioContent(), nil
}
