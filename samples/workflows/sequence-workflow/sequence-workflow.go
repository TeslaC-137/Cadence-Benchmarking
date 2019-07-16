/**
 * This code implements a workflow to be used in benchmarking for sequence activities
 */
package main

import (
	"context"
	"errors"
	"math/rand"
	"time"
	"fmt"
	"BenchmarkCode/samples/common"

	//"github.com/pborman/uuid"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	//"go.uber.org/cadence/client"
	//"go.uber.org/cadence/worker"
	"go.uber.org/zap"
)

/*
TODOS:
	1. Add activityType parameter (regular and local) to SubWorkflow

*/

/*
QUESTIONS:
	1. Proper documentation for example for child workflow options.
*/

// WorkflowName is the task list.
const WorkflowName = "Sequence_Workflow"

//Added to resolve error produced by .Get : Not enough argument for get.
var dummy string

// Registering activity and workflow
func init() {
	workflow.Register(SubWorkflow)
	workflow.Register(MainWorkflow)
	activity.Register(SleepActivity)
}




/**
  * This is the workflow we will be using for benchmarking purpose
  * This workflow executes SubWorkflow specified number of times
  * Parameters:
		ctx: a context
		numWorkflow: number of workflow to be executed
		numActivities: number of activities the workflow should execute
		sleepTime: number of seconds each activity sleeps for
  * Returns:
  		error: A string describing error if workflow fails, nil if workflow completes
*/
//MainWorkflow is the main workflow.
func MainWorkflow(ctx workflow.Context, params *common.WorkflowParamsSequence) error {
	fmt.Println("MainWorkflow Entered")
	fmt.Println(params.NumWorkflow, params.NumActivity, params.SleepTime)
	logger := workflow.GetLogger(ctx)
	fmt.Println(params.NumWorkflow+1)
	fmt.Println("MainWorkflow Started")
	fmt.Println(params.NumWorkflow)

	// Execute SubWorkflow for numWorkflow time
	for i := 0; i < params.NumWorkflow; i++ {
		fmt.Printf("This is %dth iteration of for loop inside MainWorkflow", i)
		//execution := workflow.GetInfo(ctx).WorkflowExecution

		//childID := fmt.Sprintf("SubWorkflow:%v", execution.RunID)
		cwfo := workflow.ChildWorkflowOptions{

			//WorkflowID:                   childID,
			ExecutionStartToCloseTimeout: time.Duration(2 * params.SleepTime * params.NumActivity) * time.Second,
		}
		ctx = workflow.WithChildOptions(ctx, cwfo)
		err := workflow.ExecuteChildWorkflow(ctx, SubWorkflow, i, params.NumActivity, params.SleepTime, params.FailRate).Get(ctx, &dummy)
		if err != nil {
			logger.Error("Workflow failed")
			return err
		}
	}
	fmt.Println("MainWorkflow Ended")
	return nil
}

/**
SubWorkflow :
  * This is the workflow whose specified number of instances run inside the main workflow
  * Parameters:
		ctx: a context
		workflowID: an integer that uniquely identifies each SubWorkflow running inside of MainWorkflow
		numActivities: number of activities the workflow should execute
		sleepTime: number of seconds each activity sleeps for
  * Returns:
  		error: A string describing error if workflow fails, nil if workflow completes
*/
func SubWorkflow(ctx workflow.Context, workflowID int, numActivities int, sleepTime int, failRate float64) error {
	logger := workflow.GetLogger(ctx)
	fmt.Println("SubWorkflow Started")
	ao := workflow.ActivityOptions{
		TaskList: 				WorkflowName,
		ScheduleToStartTimeout: time.Duration(2*sleepTime) * time.Second,
		StartToCloseTimeout:    time.Duration(2*sleepTime) * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Execute activity for numActivities time
	for i := 0; i < numActivities; i++ {
		future := workflow.ExecuteActivity(ctx, SleepActivity, i, sleepTime, failRate)
		err := future.Get(ctx, &dummy)
		if err != nil {
			logger.Error("Activity failed", zap.Error(err))    //Fix this to display the id of failed activity
		}
	}
	fmt.Println("SubWorkflow Ended Successfully")
	return nil
}

/**
SleepActivity :
 * This is the activity used in workflow.
 * This activity does nothing beside sleeping for specified number of seconds
 * Parameters:
		ctx: Context
		sleepTime: The number of seconds the activity sleeps for
		activityID: an integer that uniquely identifies each activity running inside a particular SubWorkflow
 		failureRate: A number between 0 and 1 indicating the rate at which activities fail
*/

/*
	Questions:
	1. Is context required as parameter????
 */
func SleepActivity(ctx context.Context, activityID int, sleepTime int, failureRate float64) error {
	//logger := activity.GetLogger(ctx)
	fmt.Print("SleepActivity Started")
	if rand.Float64() < failureRate {
		return errors.New("The activity failed")
	}
	time.Sleep(time.Duration(sleepTime) * time.Second)
	fmt.Print("SleepActivity Suceeded")
	return nil
}
