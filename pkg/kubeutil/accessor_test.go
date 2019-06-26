package kubeutil

import (
	"testing"

	apps_v1 "k8s.io/api/apps/v1"
	batch_v1beta1 "k8s.io/api/batch/v1beta1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPodTemplateSpecAccessor(t *testing.T) {
	t.Run("Deployment", func(t *testing.T) {
		deployment := &apps_v1.Deployment{
			ObjectMeta: meta.ObjectMeta{Name: "deployment"},
			Spec:       apps_v1.DeploymentSpec{},
		}

		template := PodTemplateSpecAccessor(deployment)
		if template == nil {
			t.Fatal("did not find template")
		}

		template.Name = "blub"
		if deployment.Spec.Template.Name != "blub" {
			t.Fatal("resulted template does not use same reference")
		}
	})

	t.Run("CronJob", func(t *testing.T) {
		cronjob := &batch_v1beta1.CronJob{
			ObjectMeta: meta.ObjectMeta{Name: "cronjob"},
		}

		template := PodTemplateSpecAccessor(cronjob)
		if template == nil {
			t.Fatal("did not find template")
		}

		template.Name = "blub"
		if cronjob.Spec.JobTemplate.Spec.Template.Name != "blub" {
			t.Fatal("resulted template does not use same reference")
		}
	})
}
