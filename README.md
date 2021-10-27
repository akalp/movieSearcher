# OMDb Movie Searcher

  

It is a terminal application to search movies by term using OMDb API. 

It needs some arguments to run.
```
-api-key string
	OMDB api key
-search string
	Movie title to search for.
-size int
	Number of movies to search (default 5)
```

Usage:
```
go run .\omdbSearcher.go --search=<searchTerm> --api-key=<yourAPIKey>
```

Example Scenario:
```
go run .\omdbSearcher.go --search=batman --api-key=APIKEY
Here are the results for "batman"

[1] Batman Begins (2005)
[2] Batman v Superman: Dawn of Justice (2016)        
[3] Batman (1989)
[4] Batman Returns (1992)
[5] Batman Forever (1995)

Type the number for movie detail [1-5] [q for quit]: 2

Batman v Superman: Dawn of Justice (2016)
  Fearing that the actions of Superman are left unchecked, Batman takes on the Man of Steel, while the world wrestles with what kind of a hero it really needs.
Genre       : Action, Adventure, Sci-Fi
IMDB Rating : 6.4
MetaScore   : 44
Director    : Zack Snyder
Writers     : Chris Terrio, David S. Goyer, Bob Kane
Actors      : Ben Affleck, Henry Cavill, Amy Adams
Awards      : 14 wins & 33 nominations
BoxOffice   : $330,360,194
IMDB Page   : https://www.imdb.com/title/tt2975590/

Bye!
```