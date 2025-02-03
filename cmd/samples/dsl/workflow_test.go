package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/testsuite"
	"go.uber.org/cadence/workflow"
)

func TestActivitySequenceParallelStatements(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}

	tests := []struct {
		name     string
		fields   Statement
		bindings map[string]string
		wantErr  bool
	}{
		{
			name: "Test Activity Invocation",
			fields: Statement{
				Activity: &ActivityInvocation{
					Name:      "sampleActivity",
					Arguments: []string{"var1", "var2"},
					Result:    "resultVar",
				},
			},
			bindings: map[string]string{
				"var1": "value1",
				"var2": "value2",
			},
			wantErr: false,
		},
		{
			name: "Test Sequence Execution",
			fields: Statement{
				Sequence: &Sequence{
					Elements: []*Statement{
						{
							Activity: &ActivityInvocation{
								Name:      "sampleActivity",
								Arguments: []string{"var1"},
								Result:    "resultVar1",
							},
						},
						{
							Activity: &ActivityInvocation{
								Name:      "sampleActivity",
								Arguments: []string{"var2"},
								Result:    "resultVar2",
							},
						},
					},
				},
			},
			bindings: map[string]string{
				"var1": "value1",
				"var2": "value2",
			},
			wantErr: false,
		},
		{
			name: "Test Parallel Execution",
			fields: Statement{
				Parallel: &Parallel{
					Branches: []*Statement{
						{
							Activity: &ActivityInvocation{
								Name:      "sampleActivity",
								Arguments: []string{"var1"},
								Result:    "resultVar1",
							},
						},
						{
							Activity: &ActivityInvocation{
								Name:      "sampleActivity",
								Arguments: []string{"var2"},
								Result:    "resultVar2",
							},
						},
					},
				},
			},
			bindings: map[string]string{
				"var1": "value1",
				"var2": "value2",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := testSuite.NewTestWorkflowEnvironment()
			env.RegisterActivityWithOptions(sampleActivity, activity.RegisterOptions{
				Name: "sampleActivity",
			})
			env.ExecuteWorkflow(func(ctx workflow.Context) error {
				ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
					ScheduleToStartTimeout: time.Minute,
					StartToCloseTimeout:    time.Minute,
				})
				return tt.fields.execute(ctx, tt.bindings)
			})

			require.True(t, env.IsWorkflowCompleted())
			if tt.wantErr {
				require.Error(t, env.GetWorkflowError())
			} else {
				require.NoError(t, env.GetWorkflowError())
			}
		})
	}
}

func TestSequenceFlow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}

	tests := []struct {
		name     string
		fields   Sequence
		bindings map[string]string
		wantErr  bool
	}{
		{
			name: "Test Sequence Execution with Single Activity",
			fields: Sequence{
				Elements: []*Statement{
					{
						Activity: &ActivityInvocation{
							Name:      "sampleActivity",
							Arguments: []string{"var1"},
							Result:    "resultVar1",
						},
					},
				},
			},
			bindings: map[string]string{
				"var1": "value1",
			},
			wantErr: false,
		},
		{
			name: "Test Sequence Execution with Multiple Activities",
			fields: Sequence{
				Elements: []*Statement{
					{
						Activity: &ActivityInvocation{
							Name:      "sampleActivity",
							Arguments: []string{"var1"},
							Result:    "resultVar1",
						},
					},
					{
						Activity: &ActivityInvocation{
							Name:      "sampleActivity",
							Arguments: []string{"var2"},
							Result:    "resultVar2",
						},
					},
				},
			},
			bindings: map[string]string{
				"var1": "value1",
				"var2": "value2",
			},
			wantErr: false,
		},
		{
			name: "Test Sequence Execution with Error",
			fields: Sequence{
				Elements: []*Statement{
					{
						Activity: &ActivityInvocation{
							Name:      "sampleActivity",
							Arguments: []string{"var1"},
							Result:    "resultVar1",
						},
					},
					{
						Activity: &ActivityInvocation{
							Name:      "nonExistentActivity",
							Arguments: []string{"var2"},
							Result:    "resultVar2",
						},
					},
				},
			},
			bindings: map[string]string{
				"var1": "value1",
				"var2": "value2",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := testSuite.NewTestWorkflowEnvironment()
			env.RegisterActivityWithOptions(sampleActivity, activity.RegisterOptions{
				Name: "sampleActivity",
			})

			env.ExecuteWorkflow(func(ctx workflow.Context) error {
				ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
					ScheduleToStartTimeout: time.Minute,
					StartToCloseTimeout:    time.Minute,
				})
				return tt.fields.execute(ctx, tt.bindings)
			})

			require.True(t, env.IsWorkflowCompleted())
			if tt.wantErr {
				require.Error(t, env.GetWorkflowError())
			} else {
				require.NoError(t, env.GetWorkflowError())
			}
		})
	}
}

func TestParallelFlow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}

	tests := []struct {
		name     string
		fields   Parallel
		bindings map[string]string
		wantErr  bool
	}{
		{
			name: "Test Parallel Execution with Single Activity",
			fields: Parallel{
				Branches: []*Statement{
					{
						Activity: &ActivityInvocation{
							Name:      "sampleActivity",
							Arguments: []string{"var1"},
							Result:    "resultVar1",
						},
					},
				},
			},
			bindings: map[string]string{
				"var1": "value1",
			},
			wantErr: false,
		},
		{
			name: "Test Parallel Execution with Multiple Activities",
			fields: Parallel{
				Branches: []*Statement{
					{
						Activity: &ActivityInvocation{
							Name:      "sampleActivity",
							Arguments: []string{"var1"},
							Result:    "resultVar1",
						},
					},
					{
						Activity: &ActivityInvocation{
							Name:      "sampleActivity",
							Arguments: []string{"var2"},
							Result:    "resultVar2",
						},
					},
				},
			},
			bindings: map[string]string{
				"var1": "value1",
				"var2": "value2",
			},
			wantErr: false,
		},
		{
			name: "Test Parallel Execution with Error",
			fields: Parallel{
				Branches: []*Statement{
					{
						Activity: &ActivityInvocation{
							Name:      "sampleActivity",
							Arguments: []string{"var1"},
							Result:    "resultVar1",
						},
					},
					{
						Activity: &ActivityInvocation{
							Name:      "nonExistentActivity",
							Arguments: []string{"var2"},
							Result:    "resultVar2",
						},
					},
				},
			},
			bindings: map[string]string{
				"var1": "value1",
				"var2": "value2",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := testSuite.NewTestWorkflowEnvironment()
			env.RegisterActivityWithOptions(sampleActivity, activity.RegisterOptions{
				Name: "sampleActivity",
			})

			env.ExecuteWorkflow(func(ctx workflow.Context) error {
				ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
					ScheduleToStartTimeout: time.Minute,
					StartToCloseTimeout:    time.Minute,
				})
				return tt.fields.execute(ctx, tt.bindings)
			})

			require.True(t, env.IsWorkflowCompleted())
			if tt.wantErr {
				require.Error(t, env.GetWorkflowError())
			} else {
				require.NoError(t, env.GetWorkflowError())
			}
		})
	}
}

func TestActivityInvocationFlow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}

	tests := []struct {
		name     string
		fields   ActivityInvocation
		bindings map[string]string
		wantErr  bool
	}{
		{
			name: "Test Activity Invocation Success",
			fields: ActivityInvocation{
				Name:      "sampleActivity",
				Arguments: []string{"var1"},
				Result:    "resultVar",
			},
			bindings: map[string]string{
				"var1": "value1",
			},
			wantErr: false,
		},
		{
			name: "Test Activity Invocation with Error",
			fields: ActivityInvocation{
				Name:      "nonExistentActivity",
				Arguments: []string{"var1"},
				Result:    "resultVar",
			},
			bindings: map[string]string{
				"var1": "value1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := testSuite.NewTestWorkflowEnvironment()
			env.RegisterActivityWithOptions(sampleActivity, activity.RegisterOptions{
				Name: "sampleActivity",
			})

			env.ExecuteWorkflow(func(ctx workflow.Context) error {
				ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
					ScheduleToStartTimeout: time.Minute,
					StartToCloseTimeout:    time.Minute,
				})
				return tt.fields.execute(ctx, tt.bindings)
			})

			require.True(t, env.IsWorkflowCompleted())
			if tt.wantErr {
				require.Error(t, env.GetWorkflowError())
			} else {
				require.NoError(t, env.GetWorkflowError())
			}
		})
	}
}

func Test_SimpleDSLWorkflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Define a sample DSL workflow
	dslWorkflow := Workflow{
		Variables: map[string]string{
			"var1": "value1",
			"var2": "value2",
		},
		Root: Statement{
			Activity: &ActivityInvocation{
				Name:      "sampleActivity",
				Arguments: []string{"var1", "var2"},
				Result:    "resultVar",
			},
		},
	}

	// Register a sample activity
	env.RegisterActivityWithOptions(sampleActivity, activity.RegisterOptions{
		Name: "sampleActivity",
	})

	env.ExecuteWorkflow(simpleDSLWorkflow, dslWorkflow)

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}

func sampleActivity(input []string) (string, error) {
	name := "sampleActivity"
	fmt.Printf("Run %s with input %v \n", name, input)
	return "Result_" + name, nil
}
