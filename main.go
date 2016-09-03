package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func timeOffset(from time.Time, offset time.Duration) string {
	return from.UTC().Add(-offset).Format("2006-01-02T15:04:05")
}

func timeOffsetNow(offset time.Duration) string {
	return timeOffset(time.Now(), offset)
}

func main() {

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")

	var service string
	flag.StringVar(&service, "service", "wfs.geonet.org.nz", "earthquake query service")

	var minmag float64
	flag.Float64Var(&minmag, "minmag", 3.0, "minimum magnitude to process, use 0.0 for no limit")

	var maxmag float64
	flag.Float64Var(&maxmag, "maxmag", 0.0, "maximum magnitude to process, use 0.0 for no limit")

	var since time.Duration
	flag.DurationVar(&since, "since", 30*time.Minute, "modified event search window since this time offset, use 0 for no offset")

	var ago time.Duration
	flag.DurationVar(&ago, "ago", 0, "modified event search window at least this time offset ago, use 0 for no offset")

	var eventType string
	flag.StringVar(&eventType, "type", "earthquake", "event type query parameter")

	var evaluationStatus string
	flag.StringVar(&evaluationStatus, "status", "confirmed", "event status query parameter")

	var evaluationMode string
	flag.StringVar(&evaluationMode, "mode", "manual", "event mode query parameter")

	var limit int
	flag.IntVar(&limit, "limit", 0, "maximum number of records to process before filters, use 0 for no limit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Provide a list of recent earthquakes suitable for shakemap processing\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()

	q := NewQuery(service, limit)

	q.AddFilter("eventtype+=+" + eventType)
	q.AddFilter("evaluationstatus+=+" + evaluationStatus)
	q.AddFilter("evaluationmode+=+" + evaluationMode)

	if minmag > 0.0 {
		q.AddFilter("magnitude+>=+" + strconv.FormatFloat(minmag, 'f', -1, 64))
	}
	if maxmag > 0.0 {
		q.AddFilter("magnitude+<=+" + strconv.FormatFloat(maxmag, 'f', -1, 64))
	}

	if since > 0 {
		q.AddFilter("modificationtime+>=+" + timeOffsetNow(since))
	}
	if ago > 0 {
		q.AddFilter("modificationtime+<=+" + timeOffsetNow(ago))
	}

	if verbose {
		log.Println(q.URL().String())
	}

	search, err := q.Search()
	if err != nil {
		log.Fatal(err)
	}

	var ids []string

	for _, feature := range search.Features {
		if feature.Properties.PublicID != nil {
			if verbose {
				log.Printf("id: %s", *feature.Properties.PublicID)
			}
			ids = append(ids, *feature.Properties.PublicID)
		}
	}

	if len(ids) > 0 {
		fmt.Fprintf(os.Stdout, "%s\n", strings.Join(ids, " "))
	}
}
