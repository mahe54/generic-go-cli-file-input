package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type Thing struct {
	Name        string
	Description string
	Value       string
}

func createThing(thing Thing, progress chan<- int, cancel <-chan struct{}) {
	// Simulate long-running task, updates progress channel every 500-1500ms by 10%
	for i := 1; i <= 10; i++ {
		select {
		case progress <- i * 10:
			sleepTime := time.Duration(rand.Intn(1000)+500) * time.Millisecond
			time.Sleep(sleepTime)
		case <-cancel:
			return
		}
	}
}

func readThingsFromFile(fileName string) ([]Thing, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	thingRecords, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var things []Thing
	for _, record := range thingRecords {
		thing := Thing{
			Name:        record[0],
			Description: record[1],
			Value:       record[2],
		}
		things = append(things, thing)
	}

	return things, nil
}

func launchThings(things []Thing, progressChannels []chan int, cancelChannels []chan struct{}, semaphore chan struct{}, p *mpb.Progress) {
	totalThings := len(things)
	bars := make([]*mpb.Bar, totalThings)

	for i, thing := range things {
		currentThingNumber := i + 1
		progress := make(chan int)
		progressChannels[i] = progress
		cancel := make(chan struct{})
		cancelChannels[i] = cancel

		barm := p.New(100,
			// BarFillerBuilder with custom style
			mpb.BarStyle().Lbound("[\u001b[32;1m█").Filler("█").Tip("█\u001b[36;1m").Padding("\u001b[0m░").Rbound("\u001b[0m]╟"),
			mpb.PrependDecorators(
				// display our name with one space on the right
				decor.Name(fmt.Sprintf("Creating %d/%d | %s", currentThingNumber, totalThings, thing.Name)),
				// replace ETA decorator with "done" message, OnComplete event
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.Percentage(decor.WC{W: 5}), "Done",
				),
			),
		)
		// bar := p.AddBar(100,
		// 	mpb.PrependDecorators(
		// 		decor.Name(fmt.Sprintf("Creating %d/%d | %s", currentThingNumber, totalThings, thing.Name)),
		// 	),
		// 	mpb.AppendDecorators(
		// 		// decor.EwmaETA(decor.ET_STYLE_MMSS, 0, decor.WCSyncWidth),
		// 		decor.OnComplete(
		// 			// ETA decorator with ewma age of 60
		// 			decor.Percentage(decor.WC{W: 5}), "Done",
		// 		),
		// 		// decor.Percentage(decor.WC{W: 5}),
		// 	),
		// )
		bars[i] = barm

		semaphore <- struct{}{}
		go func(thing Thing, progress chan<- int, cancel <-chan struct{}) {
			defer func() { <-semaphore }()
			createThing(thing, progress, cancel)
			close(progress)
		}(thing, progress, cancel)

		go func(bar *mpb.Bar, progress <-chan int) {
			for p := range progress {
				bar.SetCurrent(int64(p))
			}
			bar.SetTotal(100, true)
		}(barm, progress)
	}

	p.Wait()
}

func main() {
	things, err := readThingsFromFile("input.csv")
	if err != nil {
		log.Fatal(err)
	}

	p := mpb.New(mpb.WithWaitGroup(&sync.WaitGroup{}), mpb.WithWidth(25))
	progressChannels := make([]chan int, len(things))
	cancelChannels := make([]chan struct{}, len(things))

	semaphore := make(chan struct{}, 3)
	launchThings(things, progressChannels, cancelChannels, semaphore, p)
}
