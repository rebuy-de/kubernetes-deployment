package main

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
				testGetNameOfFn(FetchServicesCommand),
				testGetNameOfFn(RenderTemplatesCommand),
				testGetNameOfFn(DeployServicesCommand),
			},
		},
		{
			[]string{"fetch", "render", "deploy"},
			[]string{
				testGetNameOfFn(FetchServicesCommand),
				testGetNameOfFn(RenderTemplatesCommand),
				testGetNameOfFn(DeployServicesCommand),
			},
		},
		{
			[]string{"all", "fetch", "render", "deploy"},
			[]string{
				testGetNameOfFn(FetchServicesCommand),
				testGetNameOfFn(RenderTemplatesCommand),
				testGetNameOfFn(DeployServicesCommand),
			},
		},
		{
			[]string{"deploy", "render", "deploy", "fetch"},
			[]string{
				testGetNameOfFn(FetchServicesCommand),
				testGetNameOfFn(RenderTemplatesCommand),
				testGetNameOfFn(DeployServicesCommand),
			},
		},
		{
			[]string{"fetch"},
			[]string{
				testGetNameOfFn(FetchServicesCommand),
			},
		},
	}

	for i, tc := range cases {
		commands, err := GetCommands(tc.names...)
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
