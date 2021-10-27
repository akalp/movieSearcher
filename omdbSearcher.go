package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type config struct {
	apiKey     string
	searchTerm string
	searchSize int
}

type SearchResponse struct {
	Search []struct {
		Title  string `json:"Title"`
		Year   string `json:"Year"`
		ImdbID string `json:"imdbID"`
		Type   string `json:"Type"`
		Poster string `json:"Poster"`
	} `json:"Search"`
	TotalResults string `json:"totalResults"`
	Response     string `json:"Response"`
}

type MovieResponse struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Rated    string `json:"Rated"`
	Released string `json:"Released"`
	Runtime  string `json:"Runtime"`
	Genre    string `json:"Genre"`
	Director string `json:"Director"`
	Writer   string `json:"Writer"`
	Actors   string `json:"Actors"`
	Plot     string `json:"Plot"`
	Language string `json:"Language"`
	Country  string `json:"Country"`
	Awards   string `json:"Awards"`
	Poster   string `json:"Poster"`
	Ratings  []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings"`
	Metascore  string `json:"Metascore"`
	ImdbRating string `json:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes"`
	ImdbID     string `json:"imdbID"`
	Type       string `json:"Type"`
	Dvd        string `json:"DVD"`
	BoxOffice  string `json:"BoxOffice"`
	Production string `json:"Production"`
	Website    string `json:"Website"`
	Response   string `json:"Response"`
}

var (
	omdbSearchURL string = "http://www.omdbapi.com/?apikey=%s&s=%s"
	omdbGetURL    string = "http://www.omdbapi.com/?apikey=%s&i=%s"
)

var moviePrint string = `
%s (%s)
  %s
Genre       : %s
IMDB Rating : %s
MetaScore   : %s
Director    : %s
Writers     : %s
Actors      : %s
Awards      : %s
BoxOffice   : %s
IMDB Page   : https://www.imdb.com/title/%s/
`

func main() {
	if err := omdbProcess(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func omdbProcess() error {
	conf, err := parseConfig(os.Args[1:])

	if err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	searchResponse, err := searchMovies(conf)
	if err != nil {
		return err
	}

	choice, err := intMenu(searchResponse, conf)
	if err != nil {
		return err
	}
	if choice > 0 {
		MovieResponse, err := getMovie(conf, searchResponse.Search[choice-1].ImdbID)

		if err != nil {
			return err
		}

		printMovie(MovieResponse)
	}

	fmt.Printf("Bye!\n\n")

	return nil
}

func parseConfig(args []string) (*config, error) {
	var (
		flagAPIKey     = flag.String("api-key", "", "OMDB api key")
		flagSearchTerm = flag.String("search", "", "Movie title to search for.")
		flagSearchSize = flag.Int("size", 5, "Number of movies to search")
	)

	if err := flag.CommandLine.Parse(args); err != nil {
		return nil, err
	}

	if flag.NFlag() == 0 {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return nil, flag.ErrHelp
	}

	if len(*flagAPIKey) == 0 {
		fmt.Fprintln(os.Stderr, "api-key is mandatory!")
		return nil, flag.ErrHelp
	}

	if len(*flagSearchTerm) == 0 {
		fmt.Fprintln(os.Stderr, "search field is mandatory!")
		return nil, flag.ErrHelp
	}

	conf := &config{
		apiKey:     *flagAPIKey,
		searchTerm: *flagSearchTerm,
		searchSize: *flagSearchSize,
	}

	return conf, nil
}

func searchMovies(conf *config) (*SearchResponse, error) {
	resp, err := http.Get(fmt.Sprintf(omdbSearchURL, conf.apiKey, conf.searchTerm))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pResp SearchResponse
	if err := json.Unmarshal(body, &pResp); err != nil {
		return nil, err
	}

	return &pResp, nil
}

func intMenu(searchResponse *SearchResponse, conf *config) (int, error) {
	fmt.Printf("Here are the results for \"%s\"\n\n", conf.searchTerm)

	for i, rec := range searchResponse.Search {
		fmt.Printf("[%d] %s (%s)\n", i+1, rec.Title, rec.Year)

		if i+1 == conf.searchSize {
			break
		}
	}

	reader := bufio.NewReader(os.Stdin)
	var resp int

	for {
		fmt.Printf("\nType the number for movie detail [1-%d] [q for quit]: ", conf.searchSize)

		input, err := reader.ReadString('\n')

		if err != nil {
			return -1, err
		}

		input = strings.TrimSpace(input)

		if input == "q" {
			return -1, nil
		}

		choice, err := strconv.ParseInt(input, 10, 0)

		if err != nil {
			return -1, err
		}

		if int(choice) > conf.searchSize || choice < 1 {
			fmt.Println("Value entered is not valid, try again!")
		} else {
			resp = int(choice)
			break
		}
	}

	return resp, nil
}

func getMovie(conf *config, id string) (*MovieResponse, error) {
	resp, err := http.Get(fmt.Sprintf(omdbGetURL, conf.apiKey, id))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pResp MovieResponse
	if err := json.Unmarshal(body, &pResp); err != nil {
		return nil, err
	}

	return &pResp, nil
}

func printMovie(movie *MovieResponse) {
	printText := fmt.Sprintf(moviePrint,
		movie.Title, movie.Year,
		movie.Plot,
		movie.Genre,
		movie.ImdbRating,
		movie.Metascore,
		movie.Director,
		movie.Writer,
		movie.Actors,
		movie.Awards,
		movie.BoxOffice,
		movie.ImdbID)
	fmt.Println(printText)
}
