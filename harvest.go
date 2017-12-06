package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

const defaultFile = "data.csv"
const layout = "15:04:05 2006-01-02"

var jst *time.Location
var twitter *anaconda.TwitterApi

func main() {
	fileName := flag.String("f", defaultFile, "write file name")
	flag.Parse()

	setup()

	f, err := os.OpenFile(*fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)

	stream := twitter.UserStream(nil)
	for {
		x := <-stream.C
		switch x := x.(type) {
		case anaconda.Tweet:
			w.Write(tweetToRecord(x))
		default:
		}
		w.Flush()
	}
}

func setup() {
	setupJST()
	setupTwitter()
}

func setupJST() {
	l, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	jst = l
}

func setupTwitter() {
	anaconda.SetConsumerKey(os.Getenv("CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("CONSUMER_SECRET"))
	twitter = anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
}

func tweetToRecord(tweet anaconda.Tweet) []string {
	var r *Record
	if tweet.RetweetedStatus == nil {
		r = newRecord(tweet)
	} else {
		r = newRecord(*tweet.RetweetedStatus)
		r.SharedBy = tweet.User.ScreenName
	}
	return r.String()
}

func newRecord(tweet anaconda.Tweet) *Record {
	at := createdAtJST(tweet.CreatedAt)
	url := "https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IdStr
	text := strings.Replace(tweet.FullText, "\n", " ", -1)

	r := &Record{
		CreatedAt:   at,
		ServiceName: "Twitter",
		SourceID:    tweet.IdStr,
		SourceURL:   url,
		UserID:      tweet.User.ScreenName,
		UserName:    tweet.User.Name,
		Text:        text,
		CSVWriteAt:  time.Now().String(),
	}

	return r
}

func createdAtJST(ut string) string {
	jt, err := time.Parse(time.RubyDate, ut)
	if err != nil {
		return ut
	}
	return jt.In(jst).String()
}
