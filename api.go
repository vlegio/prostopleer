package prostopleer

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const TOKEN_URL = "http://api.pleer.com/token.php"
const API_URL = "http://api.pleer.com/index.php"

type Api struct {
	User        string    `json:"-"`
	Password    string    `json:"-"`
	AccessToken string    `json:"access_token"`
	ExpiresIn   int64     `json:"expires_in"`
	ExpirieTime time.Time `json:"-"`
}

//Create api var with AccessToken
func (api *Api) newApi() (err error) {
	client := &http.Client{}
	data := url.Values{
		"grant_type": []string{"client_credentials"},
	}
	curTime := time.Now()
	req, err := http.NewRequest("POST", TOKEN_URL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(api.User, api.Password)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	jsonResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonResponse, &api)
	if err != nil {
		return err
	}
	api.ExpirieTime = curTime.Add(time.Duration(api.ExpiresIn) * time.Minute)
	return err
}

func (a *Api) updateToken() (err error) {
	if a.ExpirieTime.UnixNano() < time.Now().UnixNano() {
		var err error
		err = a.newApi()
		if err != nil {
			return err
		}
	}
	return err
}

type searchResults struct {
	Success bool `json:"success"`
	Count  json.Number `json:"count"`
	Tracks map[string]Track `json:"tracks"`
}

type searchResultsEmpty struct {
	Success bool `json:"success"`
	Count  json.Number `json:"count"`
	Tracks []Track `json:"tracks"`
}

//Search track by query
//
//quality = all|bad|good|best
//
//page = current page
//
//perPage = result per page
func (a *Api) SearchTrack(query, quality string, page, perPage int) (tracks []Track, count int64, err error) {
	searchResult := searchResults{}
	data := map[string][]string{
		"query":          []string{query},
		"quality":        []string{quality},
		"page":           []string{strconv.Itoa(page)},
		"result_on_page": []string{strconv.Itoa(perPage)},
		"method":         []string{"tracks_search"},
	}
	body, err := a.sendPost(data, API_URL)
	if err != nil {
		return tracks, count, err
	}

	err = json.Unmarshal(body, &searchResult)
	if err != nil {		
		if _, ok := err.(*json.UnmarshalTypeError); ok {		
			searchResultEmpty := searchResultsEmpty{}
			err = json.Unmarshal(body, &searchResultEmpty)
			if err == nil {
				return tracks, 0, nil	
			}
		}

		return tracks, count, err
	}

	count, err = searchResult.Count.Int64()
	for _, track := range searchResult.Tracks {
		track.api = a
		tracks = append(tracks, track)
	}
	return tracks, count, err
}

//Get Top List
//
//list_type = 1- 1 week, 2 -  1 month, 3 - 3 months, 4 - 6 months, 5 - 1 year
//
//lang = Type of TopList (ru, en)
func (a *Api) GetTopList(list_type, page int, lang string) (tracks []Track, count int64, err error) {
	searchResult := searchResults{}
	data := map[string][]string{
		"list_type": []string{strconv.Itoa(list_type)},
		"page":      []string{strconv.Itoa(page)},
		"lang":      []string{lang},
		"method":    []string{"get_top_list"},
	}
	body, err := a.sendPost(data, API_URL)
	if err != nil {
		return tracks, count, err
	}
	err = json.Unmarshal(body, &searchResult)
	if err != nil {
		return tracks, count, err
	}
	count, err = searchResult.Count.Int64()
	for _, track := range searchResult.Tracks {
		track.api = a
		tracks = append(tracks, track)
	}
	return tracks, count, err
}

//Get suggest for user input
//
//part - user input
func (a *Api) Autocomplete(part string) (suggest []string, err error) {
	data := map[string][]string{
		"part":   []string{part},
		"method": []string{"get_suggest"},
	}
	body, err := a.sendPost(data, API_URL)
	if err != nil {
		return suggest, err
	}
	var answer struct {
		Success bool     `json:"success"`
		Suggest []string `json:"suggest"`
	}
	err = json.Unmarshal(body, &answer)
	if err != nil {
		return suggest, err
	}
	suggest = answer.Suggest
	return suggest, err
}

func (a *Api) sendPost(values map[string][]string, uri string) (body []byte, err error) {
	var req *http.Request
	var resp *http.Response
	err = a.updateToken()
	if err != nil {
		return body, err
	}
	values["access_token"] = []string{a.AccessToken}
	client := &http.Client{}
	data := url.Values(values)
	req, err = http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return body, err
	}
	resp, err = client.Do(req)
	if err != nil {
		return body, err
	}
	body, err = ioutil.ReadAll(resp.Body)
	return body, err
}
