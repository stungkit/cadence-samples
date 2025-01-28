package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/encoded"
	"go.uber.org/cadence/testsuite"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}

	env := testSuite.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(delayStartWorkflow)
	env.RegisterActivity(delayStartActivity)

	var activityMessage string
	env.SetOnActivityCompletedListener(func(activityInfo *activity.Info, result encoded.Value, err error) {
		result.Get(&activityMessage)
	})

	delayStart := 30 * time.Second
	env.ExecuteWorkflow(delayStartWorkflow, delayStart)

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	require.Equal(t, "Activity started after "+delayStart.String(), activityMessage)
}
