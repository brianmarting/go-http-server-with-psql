package goroutines

import (
	"github.com/rs/zerolog/log"
	"go-backend/model"
)

var WithdrawalQueue = make(chan model.WithdrawRequest, 100)

var WithdrawalWorkerQueue chan chan model.WithdrawRequest

func StartDispatcher(amount int) {
	WithdrawalWorkerQueue = make(chan chan model.WithdrawRequest, 5)

	for i := 0; i < amount; i++ {
		worker := NewWorker(i, WithdrawalWorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case withdrawal := <-WithdrawalQueue:
				log.Info().Msg("Incoming withdrawal req")
				// Start the withdrawal work
				go func() {
					// get idle worker from queue
					withdrawalWorker := <-WithdrawalWorkerQueue

					log.Info().Msg("Got idle worker from queue")

					// add withdrawal work to worker so it can process it
					withdrawalWorker <- withdrawal
				}()
			}
		}
	}()
}
