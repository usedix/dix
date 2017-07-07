package dict

import (
	"context"
	"net/http"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

var urlRegexp = regexp.MustCompile(`[a-zA-Z]+://[a-zA-Z0-9./]+`)
var pronunciationRegexp = regexp.MustCompile(`\[[^\]]*\]`)

type Bing struct {
	Client *http.Client
}

func (bing Bing) Search(ctx context.Context, rawword string) (*Word, error) {
	// build request
	req, err := http.NewRequest("GET",
		"http://cn.bing.com/dict/search?q="+url.QueryEscape(rawword), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	// get http client
	client := bing.Client
	if client == nil {
		client = http.DefaultClient
	}

	// do http request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	word := &Word{}

	// get word
	word.Word = doc.Find("#headword > h1 > strong").Text()

	// get pronunciation
	us, us_mp3url, uk, uk_mp3url := getPronunciation(doc)

	word.Pronunciation = Pronunciation{
		US:        us,
		US_MP3URL: us_mp3url,
		UK:        uk,
		UK_MP3URL: uk_mp3url,
	}

	// get defs
	doc.Find("body > div.contentPadding > div > div > div.lf_area > div.qdef > ul > li").
		Each(func(i int, s *goquery.Selection) {
			pos := s.Find(".pos").Text()
			// map pos from 网络 => Web
			if pos == "网络" {
				pos = "Web"
			}

			def := s.Find(".def").Text()

			word.Defs = append(word.Defs, Def{
				PartOfSpeech: pos,
				Def:          def,
			})
		})

		// get sentences
	doc.Find("#sentenceSeg > .se_li").Each(func(i int, s *goquery.Selection) {
		en := s.Find(".sen_en").Text()
		cn := s.Find(".sen_cn").Text()

		mp3url := s.Find(".mm_div .bigaud").AttrOr("onmousedown", "")
		mp3url = urlRegexp.FindString(mp3url)

		word.SampleSentences = append(word.SampleSentences, SampleSentence{
			English: en,
			Chinese: cn,
			MP3URL:  mp3url,
		})
	})
	return word, nil
}

func getPronunciation(doc *goquery.Document) (us, us_mp3url, uk, uk_mp3url string) {
	us = doc.Find("body > div.contentPadding > div > div > div.lf_area > div.qdef > div.hd_area > div.hd_tf_lh > div > div.hd_prUS").Text()
	us = pronunciationRegexp.FindString(us)
	if len(us) > 2 {
		us = us[1 : len(us)-1]
	} else {
		us = ""
	}
	// us_mp3url
	us_mp3url = doc.Find("body > div.contentPadding > div > div > div.lf_area > div.qdef > div.hd_area > div.hd_tf_lh > div > div:nth-child(2) > a").AttrOr("onclick", "")
	us_mp3url = urlRegexp.FindString(us_mp3url)
	// uk
	uk = doc.Find("body > div.contentPadding > div > div > div.lf_area > div.qdef > div.hd_area > div.hd_tf_lh > div > div.hd_pr").Text()
	uk = pronunciationRegexp.FindString(uk)
	if len(uk) > 2 {
		uk = uk[1 : len(uk)-1]
	} else {
		uk = ""
	}
	// uk_mp3url
	uk_mp3url = doc.Find("body > div.contentPadding > div > div > div.lf_area > div.qdef > div.hd_area > div.hd_tf_lh > div > div:nth-child(4) > a").AttrOr("onclick", "")
	uk_mp3url = urlRegexp.FindString(uk_mp3url)

	return us, us_mp3url, uk, uk_mp3url
}
