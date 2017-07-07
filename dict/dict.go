package dict

import "context"

type Dictionary interface {
	Search(ctx context.Context, word string) (*Word, error)
}

type Pronunciation struct {
	US        string `json:"us"`
	US_MP3URL string `json:"us_mp3url"`
	UK        string `json:"uk"`
	UK_MP3URL string `json:"uk_mp3url"`
}

type Def struct {
	PartOfSpeech string `json:"pos"`
	Def          string `json:"def"`
}

type SampleSentence struct {
	English string `json:"en"`
	Chinese string `json:"ch"`
	MP3URL  string `json:"mp3url"`
}

type Word struct {
	Word            string           `json:"word"`
	Pronunciation   Pronunciation    `json:"pronunciation"`
	Defs            []Def            `json:"defs"`
	SampleSentences []SampleSentence `json:"sample_sentences"`
}
