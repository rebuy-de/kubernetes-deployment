package api

import (
	"regexp"
	"strings"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/kubeutil"
	"github.com/rebuy-de/kubernetes-deployment/pkg/templates"
)

func (app *App) Render(fetched *FetchResult) ([]runtime.Object, error) {
	defaultVariables := templates.Variables{
		"gitBranchName": fetched.Branch.Name,
		"gitCommitID":   fetched.Branch.SHA,
	}

	fetched.Service.Variables.Defaults(defaultVariables)

	log.WithFields(log.Fields{
		"Values": fetched.Service.Variables,
	}).Debug("collected template values")

	return app.decode(fetched.Templates, fetched.Service.Variables)
}

func (app *App) decode(files []gh.File, vars templates.Variables) ([]runtime.Object, error) {
	var objects []runtime.Object

	for _, file := range files {
		log.WithFields(log.Fields{
			"Location": file.Location,
		}).Debug("decoding file")

		switch {
		case strings.HasSuffix(file.Name(), ".yml"):
			fallthrough

		case strings.HasSuffix(file.Name(), ".yaml"):
			objs, err := app.decodeYAML(file, vars)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			objects = append(objects, objs...)

		case strings.HasSuffix(file.Name(), ".jsonnet"):
			objs, err := app.decodeJsonnet(file, vars, files)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			objects = append(objects, objs...)

		default:
			log.WithFields(log.Fields{
				"Name": file.Name(),
			}).Debug("Ignoring file with wrong extension.")
		}
	}

	return objects, nil
}

func (app *App) decodeYAML(file gh.File, vars templates.Variables) ([]runtime.Object, error) {
	rendered, err := templates.Render(file.Content, vars)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var objects []runtime.Object
	splitted := regexp.MustCompile("[\n\r]---").Split(rendered, -1)

	for _, part := range splitted {
		if strings.TrimSpace(part) == "" {
			log.WithFields(log.Fields{
				"Name": file.Name(),
			}).Debug("Ignoring empty file.")
			continue
		}

		obj, err := kubeutil.Decode([]byte(part))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		obj, err = app.Interceptors.PostManifestRender(obj)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		objects = append(objects, obj)
	}

	return objects, nil
}

func (app *App) decodeJsonnet(file gh.File, vars templates.Variables, all []gh.File) ([]runtime.Object, error) {
	var objects []runtime.Object

	importer := gh.NewJsonnetImporter(app.Clients.GitHub)

	vm := jsonnet.MakeVM()
	vm.Importer(importer)
	for k, v := range vars {
		vm.ExtVar(k, v)
	}

	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "resolveGitSHA",
		Params: ast.Identifiers{"location"},
		Func: func(x []interface{}) (interface{}, error) {
			location, err := gh.NewLocation(x[0].(string))
			if err != nil {
				return nil, err
			}

			branch, err := app.Clients.GitHub.GetBranch(location)
			if err != nil {
				return nil, err
			}

			return branch.SHA, nil
		},
	})

	docs, err := vm.EvaluateSnippetStream(file.Location.String(), file.Content)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, doc := range docs {
		obj, err := kubeutil.Decode([]byte(doc))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		obj, err = app.Interceptors.PostManifestRender(obj)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		objects = append(objects, obj)
	}

	return objects, nil
}
