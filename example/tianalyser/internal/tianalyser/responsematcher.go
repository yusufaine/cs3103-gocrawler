package tianalyser

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func HasCountryRepresentationHeader(resp *http.Response) bool {
	if resp.StatusCode != 200 {
		return false
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return false
	}

	headerElement := doc.Find("#Country_Representation")
	// 0 indicates no match
	return headerElement.Length() > 0
}
