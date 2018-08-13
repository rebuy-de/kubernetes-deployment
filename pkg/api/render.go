package api

import (
	"regexp"
	"strings"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

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
	decode := scheme.Codecs.UniversalDeserializer().Decode

	for _, part := range splitted {
		if strings.TrimSpace(part) == "" {
			log.WithFields(log.Fields{
				"Name": file.Name(),
			}).Debug("Ignoring empty file.")
			continue
		}

		obj, _, err := decode([]byte(part), nil, nil)

		// Fallback to UnknownObject if the API/Kind is not registered, so the
		// interceptors still work. In case the Kind actually does not exist,
		// kubectl will fail later anyway.
		if runtime.IsNotRegisteredError(err) {
			unknown := new(kubeutil.UnknownObject)
			err = unknown.FromYAML([]byte(part))
			obj = unknown
		}

		if err != nil {
			log.Warnf("%#v", err)
			return nil, errors.Wrapf(err, "unable to decode file '%s'", file.Name())
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
	decode := scheme.Codecs.UniversalDeserializer().Decode

	importer := new(jsonnet.MemoryImporter)
	importer.Data = make(map[string]string)
	for _, f := range all {
		importer.Data[f.Name()] = f.Content
	}

	vm := jsonnet.MakeVM()
	vm.Importer(importer)
	for k, v := range vars {
		vm.ExtVar(k, v)
	}

	docs, err := vm.EvaluateSnippetStream(file.Name(), file.Content)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, doc := range docs {
		obj, _, err := decode([]byte(doc), nil, nil)

		// Fallback to UnknownObject if the API/Kind is not registered, so the
		// interceptors still work. In case the Kind actually does not exist,
		// kubectl will fail later anyway.
		if runtime.IsNotRegisteredError(err) {
			unknown := new(kubeutil.UnknownObject)
			err = unknown.FromJSON([]byte(doc))
			obj = unknown
		}

		if err != nil {
			return nil, errors.Wrapf(err, "unable to decode file '%s'", file.Name())
		}

		obj, err = app.Interceptors.PostManifestRender(obj)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		objects = append(objects, obj)
	}

	return objects, nil
}
