package cmd

import (
	"reflect"
	"runtime"
	"testing"
)

func testGetNameOfFn(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func TestGetCommands(t *testing.T) {
	var cases = []struct {
		names  []string
		expect []string
	}{
		{
			[]string{"all"},
			[]string{
				testGetNameOfFn(FetchServicesGoal),
				testGetNameOfFn(RenderTemplatesGoal),
				testGetNameOfFn(DeployServicesGoal),
			},
		},
		{
			[]string{"fetch", "render", "deploy"},
			[]string{
				testGetNameOfFn(FetchServicesGoal),
				testGetNameOfFn(RenderTemplatesGoal),
				testGetNameOfFn(DeployServicesGoal),
			},
		},
		{
			[]string{"all", "fetch", "render", "deploy"},
			[]string{
				testGetNameOfFn(FetchServicesGoal),
				testGetNameOfFn(RenderTemplatesGoal),
				testGetNameOfFn(DeployServicesGoal),
			},
		},
		{
			[]string{"deploy", "render", "deploy", "fetch"},
			[]string{
				testGetNameOfFn(FetchServicesGoal),
				testGetNameOfFn(RenderTemplatesGoal),
				testGetNameOfFn(DeployServicesGoal),
			},
		},
		{
			[]string{"fetch"},
			[]string{
				testGetNameOfFn(FetchServicesGoal),
			},
		},
	}

	for i, tc := range cases {
		commands, err := GetGoals(tc.names...)
		if err != nil {
			t.Errorf("Test case %d failed.", i)
			t.Errorf("  in:    %v", tc.names)
			t.Errorf("  error: %v", err)
		}

		commandNames := make([]string, len(commands))
		for i, command := range commands {
			commandNames[i] = testGetNameOfFn(command)
		}

		if !reflect.DeepEqual(commandNames, tc.expect) {
			t.Errorf("Test case %d failed.", i)
			t.Errorf("  in:       %v", tc.names)
			t.Errorf("  out:      %v", commandNames)
			t.Errorf("  expected: %v", tc.expect)
		}
	}
}
