package cmd

import "fmt"

type Command func(app *App) error

var (
	CommandMapping = map[string]Command{
		"fetch":  FetchServicesCommand,
		"render": RenderTemplatesCommand,
		"deploy": DeployServicesCommand,
	}
	CommandAliases = map[string][]string{
		"all": []string{"fetch", "render", "deploy"},
	}
	CommandOrder = []string{
		"fetch", "render", "deploy",
	}
)

func GetCommands(names ...string) ([]Command, error) {
	unaliased := []string{}
	for _, name := range names {
		replace, ok := CommandAliases[name]
		if ok {
			unaliased = append(unaliased, replace...)
		} else {
			unaliased = append(unaliased, name)
		}
	}

	commands := make(map[string]bool)
	for _, name := range CommandOrder {
		commands[name] = false
	}

	for _, name := range unaliased {
		_, ok := CommandMapping[name]
		if !ok {
			return nil, fmt.Errorf("Unknown command %s", name)
		}

		commands[name] = true
	}

	result := []Command{}
	for _, name := range CommandOrder {
		if commands[name] {
			result = append(result, CommandMapping[name])
		}
	}

	return result, nil
}
