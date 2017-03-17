package cmd

import "fmt"

type Goal func(app *App) error

var (
	GoalMapping = map[string]Goal{
		"fetch":  FetchServicesGoal,
		"render": RenderTemplatesGoal,
		"deploy": DeployServicesGoal,
	}
	GoalAliases = map[string][]string{
		"all": {"fetch", "render", "deploy"},
	}
	GoalOrder = []string{
		"fetch", "render", "deploy",
	}
)

func GetGoals(names ...string) ([]Goal, error) {
	unaliased := []string{}
	for _, name := range names {
		replace, ok := GoalAliases[name]
		if ok {
			unaliased = append(unaliased, replace...)
		} else {
			unaliased = append(unaliased, name)
		}
	}

	goals := make(map[string]bool)
	for _, name := range GoalOrder {
		goals[name] = false
	}

	for _, name := range unaliased {
		_, ok := GoalMapping[name]
		if !ok {
			return nil, fmt.Errorf("Unknown goal '%s'", name)
		}

		goals[name] = true
	}

	result := []Goal{}
	for _, name := range GoalOrder {
		if goals[name] {
			result = append(result, GoalMapping[name])
		}
	}

	return result, nil
}
