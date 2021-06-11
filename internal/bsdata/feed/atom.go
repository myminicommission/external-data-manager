package feed

import (
	"encoding/xml"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	Title   string   `xml:"title"`
	Link    []struct {
		Text string `xml:",chardata"`
		Rel  string `xml:"rel,attr"`
		Type string `xml:"type,attr"`
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Author struct {
		Text string `xml:",chardata"`
		Name string `xml:"name"`
		URI  string `xml:"uri"`
	} `xml:"author"`
	Subtitle struct {
		Text string `xml:",chardata"`
		Type string `xml:"type,attr"`
	} `xml:"subtitle"`
	ID      string `xml:"id"`
	Updated string `xml:"updated"`
	Entry   []struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
			Href string `xml:"href,attr"`
		} `xml:"link"`
		ID        string `xml:"id"`
		Updated   string `xml:"updated"`
		Published string `xml:"published"`
		Summary   struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
		} `xml:"summary"`
	} `xml:"entry"`
}

const (
	atomAllURL = "https://battlescribedata.appspot.com/repos/feeds/all.atom"
)

// GetAll fetches the All ATOM feed from BSData
func GetAll() (Feed, error) {
	var feed Feed
	// TODO: get the atom feed
	res, err := http.Get(atomAllURL)
	if err != nil {
		logrus.Error("http.Get", err)
		return feed, err
	}

	defer res.Body.Close()

	dec := xml.NewDecoder(res.Body)
	err = dec.Decode(&feed)
	if err != nil {
		logrus.Error("dec.Decode", err)
		return feed, err
	}

	return feed, nil
}
