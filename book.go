package main

import (
	"encoding/json"
	"fmt"
	"github.com/pyrousnet/slash_commands/MatterMost"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Books struct {
	NumFound      int  `json:"numFound"`
	Start         int  `json:"start"`
	NumFoundExact bool `json:"numFoundExact"`
	Docs          []struct {
		Key                   string   `json:"key"`
		Type                  string   `json:"type"`
		Seed                  []string `json:"seed"`
		Title                 string   `json:"title"`
		TitleSuggest          string   `json:"title_suggest"`
		TitleSort             string   `json:"title_sort"`
		EditionCount          int      `json:"edition_count"`
		EditionKey            []string `json:"edition_key"`
		PublishDate           []string `json:"publish_date"`
		PublishYear           []int    `json:"publish_year"`
		FirstPublishYear      int      `json:"first_publish_year"`
		NumberOfPagesMedian   int      `json:"number_of_pages_median"`
		Lccn                  []string `json:"lccn"`
		Oclc                  []string `json:"oclc"`
		Lcc                   []string `json:"lcc"`
		Isbn                  []string `json:"isbn"`
		LastModifiedI         int      `json:"last_modified_i"`
		EbookCountI           int      `json:"ebook_count_i"`
		EbookAccess           string   `json:"ebook_access"`
		HasFulltext           bool     `json:"has_fulltext"`
		PublicScanB           bool     `json:"public_scan_b"`
		RatingsAverage        float64  `json:"ratings_average"`
		RatingsSortable       float64  `json:"ratings_sortable"`
		RatingsCount          int      `json:"ratings_count"`
		RatingsCount1         int      `json:"ratings_count_1"`
		RatingsCount2         int      `json:"ratings_count_2"`
		RatingsCount3         int      `json:"ratings_count_3"`
		RatingsCount4         int      `json:"ratings_count_4"`
		RatingsCount5         int      `json:"ratings_count_5"`
		ReadinglogCount       int      `json:"readinglog_count"`
		WantToReadCount       int      `json:"want_to_read_count"`
		CurrentlyReadingCount int      `json:"currently_reading_count"`
		AlreadyReadCount      int      `json:"already_read_count"`
		CoverEditionKey       string   `json:"cover_edition_key"`
		CoverI                int      `json:"cover_i"`
		Publisher             []string `json:"publisher"`
		Language              []string `json:"language"`
		AuthorKey             []string `json:"author_key"`
		AuthorName            []string `json:"author_name"`
		Subject               []string `json:"subject"`
		IDAmazon              []string `json:"id_amazon"`
		PublisherFacet        []string `json:"publisher_facet"`
		SubjectFacet          []string `json:"subject_facet"`
		Version               int64    `json:"_version_"`
		LccSort               string   `json:"lcc_sort"`
		AuthorFacet           []string `json:"author_facet"`
		SubjectKey            []string `json:"subject_key"`
	} `json:"docs"`
	NumFound0 int         `json:"num_found"`
	Q         string      `json:"q"`
	Offset    interface{} `json:"offset"`
}

func getBookInfo(w http.ResponseWriter, r *http.Request) {
	var b Books
	rawText := r.URL.Query().Get("text")
	text := strings.Replace(rawText, " ", "+", -1)

	fmt.Printf("\033[32mIncomming book request for:\033[0m\033[36m " + rawText + "\033[0m\n")

	resp, err := http.Get("https://openlibrary.org/search.json?q=" + text)
	if err != nil {
		fmt.Printf("\033[31mError: " + err.Error() + "\033[0m\n")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("\033[31mError: " + err.Error() + "\033[0m\n")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = json.Unmarshal([]byte(body), &b)
	if err != nil {
		fmt.Printf("\033[31mError: " + err.Error() + "\033[0m\n")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	rs := MatterMost.Response{
		ResponseType: "in_channel",
		Text:         "",
	}

	if b.NumFound <= 0 {
		rs.Text = "Your book 404'd"
	} else {
		formattedResult := ""
		for i := 0; i < b.NumFound; i++ {
			book := b.Docs[i]
			formattedResult += "| Data | Value |\n"
			formattedResult += "| :------ | :-------|"
			formattedResult += "\n| Book " + strconv.Itoa(i+1) + " Title | " + book.Title + " |"
			formattedResult += "\n| Book " + strconv.Itoa(i+1) + " Author | " + book.AuthorName[0] + " |"
			formattedResult += "\n| Book " + strconv.Itoa(i+1) + " ISBN | " + book.Isbn[0] + " |"
			formattedResult += "\n\n"
			if i > 1 { // Only print 3 books 0, 1, and 2 (breaks when i == 2)
				break
			}
		}

		rs.Text = formattedResult
	}

	toMM, err := json.Marshal(rs)
	if err != nil {
		fmt.Printf("\033[31mError: " + err.Error() + "\033[0m\n")
		fmt.Printf("Error: %s", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(toMM))
}
