package commands

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mindwingx/gocally"
	"github.com/mindwingx/graph-generator/constants"
	"github.com/spf13/cobra"
	"log"
	"sync"
	"time"
)

const (
	goroutineWorkers = 50
	goroutineCap     = 200
)

var messageGeneratorCmd = &cobra.Command{
	Use:   "msg:gen",
	Short: "generate messages",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		var total int

		//time.Sleep(time.Second * 2)

		totalChan := make(chan int)
		st := time.Now()

		fmt.Println("[started]")

		go func() {
			wg.Wait()
			close(totalChan)
		}()

		for i := 0; i < goroutineWorkers; i++ {
			wg.Add(1)
			go msgGen(&wg, totalChan, goroutineCap)
		}

		for i := range totalChan {
			total += i
		}

		end := time.Since(st)
		res := fmt.Sprintf(
			"\n%d messages were generated and sent within %s seconds.",
			total, end,
		)

		fmt.Println(res)
	},
}

// HELPER METHODS

func msgGen(wg *sync.WaitGroup, tch chan int, wgCap int) {
	defer wg.Done()
	rq := gocally.SetRequest().WithUrl(constants.AggregatorUrl)

	// prepare the nested goroutine
	var mutex sync.Mutex
	var innerWg sync.WaitGroup

	for i := 0; i < wgCap; i++ {
		innerWg.Add(1)

		go func() {
			defer func() {
				if rec := recover(); rec != nil {
					// Handle the panic and log the recovered value.
					res := fmt.Sprintf("[cmd-generator][recovered] %s", rec)
					fmt.Println(res)
					//todo: add to the failure channel to retry
				}
			}()

			mutex.Lock()
			res, err := rq.SetBody(map[string]interface{}{
				"message": gofakeit.HackerPhrase(),
			}).Post().Payload()

			innerWg.Done()
			mutex.Unlock()

			if err != nil {
				log.Println("[cmd-generator][api-response] error:", err.Error())
			} else if res["status_code"].(int) != 201 {
				log.Printf(
					"[cmd-generator][api] response: %d - %s", res["status_code"].(int),
					res["payload"].(map[string]interface{})["data"].(string),
				)
			} else {
				// the totalChan value, indicate the acknowledgment
				//per request with in the related goroutine
				tch <- 1
			}
		}()
	}

	innerWg.Wait()
}
