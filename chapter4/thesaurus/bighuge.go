package thesaurus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BigHuge struct {
	APIKey string
}

type synonyms struct {
	Noun *words `json:"noun"`
	Verb *words `json:"verb"`
}

type words struct {
	Syn []string `json:"syn"`
}

func NewBigHuge(key string) *BigHuge {
	return &BigHuge{APIKey: key}
}

func (b *BigHuge) Synonyms(term string) ([]string, error) {
	var syns []string
	url := fmt.Sprintf("https://words.bighugelabs.com/api/2/%s/word/json", b.APIKey)
	res, err := http.Get(url)
	if err != nil {
		return syns, err
	}
	defer res.Body.Close()

	var data synonyms
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return syns, err
	}

	if data.Noun != nil {
		syns = append(syns, data.Noun.Syn...)
	}

	if data.Verb != nil {
		syns = append(syns, data.Verb.Syn...)
	}

	return syns, nil
}
