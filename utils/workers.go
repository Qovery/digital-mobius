package utils

import (
	log "github.com/sirupsen/logrus"
	"sync"
)

type Job func(elems ...interface{}) interface{}

func worker(wg *sync.WaitGroup, jobParams *chan []interface{}, job Job, results *chan interface{}) {
	for param := range *jobParams {
		*results <- job(param...)
	}
	wg.Done()
}

func createWorkerPool(workerCount int, jobParams *chan []interface{}, job Job, results *chan interface{}) {
	var wg sync.WaitGroup
	for i:= 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(&wg, jobParams, job, results)
	}
	wg.Wait()
	close(*results)
}

func allocate(arrayParamIndex int, jobParams *chan []interface{} , job Job, params ...interface{}) {
	var array = params[arrayParamIndex].([]NodeInfos)

		for _, elem := range array {
			var newParams []interface{}

			for i := 0 ; i < len(params) ; i++ {
				if i == arrayParamIndex {
					newParams = append(newParams, elem)
				} else {
					newParams =append(newParams, params[i])
				}
			}

			*jobParams <- newParams
		}

		close(*jobParams)
}

func result(done chan bool, results chan interface{}) {
	for result := range results {
		err, errOk := result.(error)
		if errOk {
			log.Error(err.Error())
			return
		}
			log.Debug(result)
	}
	done <- true
}

/**
CAREFUL : the arrayParamIndex should have the same value that the index of the element used in params
 */
func RunWorkers (workerCount int, endMessage string, arrayParamIndex int, job Job, params ...interface{}) {
		jobParams := make(chan []interface{})
		results := make(chan interface{})
		go allocate(arrayParamIndex, &jobParams, job, params...)
		done := make(chan bool)
		go result(done, results)
		createWorkerPool(workerCount, &jobParams, job, &results)
		<- done

		log.Debug(endMessage)
}

