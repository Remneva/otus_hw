package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/cheggaaa/pb"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	var endCh = make(chan struct{})
	var startCh = make(chan struct{})

	bar := pb.New(847).SetUnits(pb.U_BYTES)
	bar.ShowPercent = true
	bar.ShowCounters = false
	bar.SetRefreshRate(time.Millisecond)

	go func() {
		startCh <- struct{}{}
		bar.Start()
		for {
			select {
			case <-endCh:
				bar.Finish()
			default:
				bar.Increment()
				time.Sleep(time.Millisecond)
			}
		}
	}()

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fromPath := dir + "/" + from

	<-startCh
	err = Copy(fromPath, to, offset, limit)
	if err != nil {
		fmt.Println(err)
	}
	endCh <- struct{}{}

	close(endCh)
	close(startCh)
}
