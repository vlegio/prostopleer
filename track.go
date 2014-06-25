package prostopleer

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Track struct {
	Id          string `json:"id"`
	Artist      string `json:"artist"`
	Name        string `json:"track"`
	Duration    string `json:"lenght"`
	Bitrate     string `json:"bitrate"`
	Size        string `json:"size"`
	Lyrics      string `json:"-"`
	ListenUrl   string `json:"-"`
	DownloadUrl string `json:"-"`
	MP3         []byte `json:"-"`
	api         *Api   `json:"-"`
}

//Update track info (I don't khow why)
func (t *Track) GetInfo() (err error) {
	data := map[string][]string{
		"track_id": []string{t.Id},
		"method":   []string{"tracks_get_info"},
	}
	body, err := t.api.sendPost(data, API_URL)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, t)
	if err != nil {
		return err
	}
	return err
}

//Get track lyrics
func (t *Track) GetLyrics() (err error) {
	data := map[string][]string{
		"track_id": []string{t.Id},
		"method":   []string{"tracks_get_lyrics"},
	}
	body, err := t.api.sendPost(data, API_URL)
	if err != nil {
		return err
	}
	var answer struct {
		Status bool   `json:"status"`
		Text   string `json:"text"`
	}
	err = json.Unmarshal(body, &answer)
	if err != nil {
		return err
	}
	t.Lyrics = answer.Text
	return err
}

//Get link
//
//linkType = listen|save
func (t *Track) GetLink(linkType string) (err error) {
	data := map[string][]string{
		"track_id": []string{t.Id},
		"method":   []string{"tracks_get_download_link"},
		"reason":   []string{linkType},
	}
	body, err := t.api.sendPost(data, API_URL)
	if err != nil {
		return err
	}
	var answer struct {
		Status bool   `json:"status"`
		Link   string `json:"url"`
	}
	err = json.Unmarshal(body, &answer)
	if err != nil {
		return err
	}
	switch linkType {
	case "save":
		t.DownloadUrl = answer.Link
	case "listen":
		t.ListenUrl = answer.Link
	}
	return err
}

//Save track to Track.MP3
func (t *Track) Download() (err error) {
	err = t.GetLink("save")
	if err != nil {
		return err
	}
	resp, err := http.Get(t.DownloadUrl)
	if err != nil {
		return err
	}
	t.MP3, err = ioutil.ReadAll(resp.Body)
	return err
}
