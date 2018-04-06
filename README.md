# topstories
[![Build Status](https://travis-ci.org/yanndr/topstories.svg?branch=master)](https://travis-ci.org/yanndr/topstories) 
[![Go Report Card](https://goreportcard.com/badge/github.com/yanndr/topstories)](https://goreportcard.com/report/github.com/yanndr/topstories)

## About
topstories dispays a given number of top stories from a news aggregator. The current verion implements only Hakernews.

## Install
```
go get github.com/yanndr/topstories
go install
```

## Usage
```
topstories 
Flags:
 -c int
        max concurency allowed (default 20)
  -csv
        Save the result to a csv file.
  -n int
        number of stories to display (default 20)
  -o string
        output file name (default "outupt.csv")
```