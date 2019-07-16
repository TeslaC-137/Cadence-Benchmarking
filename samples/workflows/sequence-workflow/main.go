package main

/*
 QUESTIONS:
	1. Default value for timeouts in activity options
 */



import (
	"BenchmarkCode/samples/common"
	//"context"
	"flag"
	"time"

	//"math/rand"
	//"errors"
	"fmt"

	"github.com/pborman/uuid"
	//"go.uber.org/cadence/activity"
	//"go.uber.org/cadence/workflow"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/worker"
	//"go.uber.org/zap"
)


func startWorkers(h *common.SampleHelper) {
	workerOptions := worker.Options{
		MetricsScope: h.Scope,
		Logger:       h.Logger,
	}
	h.StartWorkers(h.Config.DomainName, WorkflowName, workerOptions)

}


func startWorkflow(h *common.SampleHelper, params *common.WorkflowParamsSequence) {
	fmt.Println("startWorkflow function entered")
	fmt.Println(params.NumWorkflow, params.NumActivity, params.SleepTime)
	workflowOptions := client.StartWorkflowOptions{
		ID:                              "sequence_" + uuid.New(),
		TaskList:                        WorkflowName,
		ExecutionStartToCloseTimeout:    2 * time.Minute,
		DecisionTaskStartToCloseTimeout: 2 * time.Minute,
	}

	h.StartWorkflow(workflowOptions, MainWorkflow, params)
	fmt.Println(params.NumWorkflow, params.NumActivity, params.SleepTime)
	fmt.Println("startWorkflow function exited")

}

func main() {
	fmt.Println("Main Entered")
	var mode string
	fmt.Printf(mode)
	flag.StringVar(&mode, "m", "trigger", "Mode is worker or trigger.")
	/*
	Make arguments of MainWorkflow configurable via command line flags
	*/
	nActivity := flag.Int("nActivity", 10, "Number of activities each subworkflow will contain")
	nWorkflow := flag.Int("nWorkflow", 10, "Number of workflows to be executed")
	sleepTime := flag.Int("sleepTime", 10, "Number of seconds each activity will sleep for")
	//fRate := flag.Int("failRate", 0.0, "Number of activities each subworkflow will contain") // Because float cant be made a command line argument in this way

	fRate := 0.0
	flag.Parse()

	var h common.SampleHelper
	h.SetupServiceConfig()
	//fmt.Println(*nWorkflow, *nActivity, *sleepTime, fRate)

	//params := workflowParams{*nWorkflow, *nActivity, *sleepTime, fRate}
	//fmt.Println(params.numWorkflow, params.numActivity, params.sleepTime)

	switch mode {
	case "worker":
		startWorkers(&h)
		select {}

	case "trigger":
		params := common.WorkflowParamsSequence{*nWorkflow, *nActivity, *sleepTime, fRate}
		fmt.Println(params.NumWorkflow, params.NumActivity, params.SleepTime)
		startWorkflow(&h, &params)
		fmt.Println(params.NumWorkflow, params.NumActivity, params.SleepTime)
	}
	fmt.Println("Main Exited")

}
