# go-wordle

An application for solving Wordle puzzles. Currently, the project analyses wordlists taken from the New York Times'
Javascript, and presents the results. The actual player is still in progress.

## Playing the Game

To play the game, run `go run ./cmd/main.go SPARE`, or with any other five-letter word that is a valid solution
to a Wordle.

```
go run main.go play SPARE
Reading valid guesses from file...
Found 12947 words
Reading valid solutions from file...
Found 2309 words

There are currently 2309 possible solutions
The next guess should be "RAISE"
Worst-case scenario for guess is the feedback "-----". Carry-over ratio for possible solutions list would be 7.23%
Guess number 1: "RAISE"
Feedback from guess was: "YY-YG"

There are currently 5 possible solutions
The remaining possible solutions are: [SCARE, SHARE, SNARE, SPARE, STARE]
The next guess should be "CHANT"
Guess number 2: "NICHT"
Feedback from guess was: "--G--"

There are currently 1 possible solutions
The remaining possible solutions are: [SPARE]
The next guess should be "SPARE"
Guess number 3: "SPARE"
Feedback from guess was: "GGGGG"

Won the Wordle in 3 turns
Finished!
```

## Wordlists

`./data/wordlist-valid-guesses.csv` and `./data/wordlist-valid-solutions.csv` are the input files. They are taken from
the New York Times website's Javascript.
