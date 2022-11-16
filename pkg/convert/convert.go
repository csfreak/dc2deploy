/*
Copyright © 2022 Jason Ross

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package convert

import (
	"encoding/json"
	"fmt"

	ocappsv1 "github.com/openshift/api/apps/v1"
	"github.com/openshift/library-go/pkg/image/trigger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yaml "sigs.k8s.io/yaml"
)

func ToDeploy(orig *ocappsv1.DeploymentConfig) (*appsv1.Deployment, error) {
	dc := orig.DeepCopy()
	deploy := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:         dc.Name,
			GenerateName: dc.GenerateName,
			Namespace:    dc.Namespace,
			Labels:       dc.Labels,
			Annotations:  cleanAnnotations(dc.Annotations),
		},
		Spec:   appsv1.DeploymentSpec{},
		Status: appsv1.DeploymentStatus{},
	}

	dc.Spec.Template.DeepCopyInto(&deploy.Spec.Template)
	deploy.Spec.Template.Annotations = cleanAnnotations(dc.Spec.Template.Annotations)
	deploy.Spec.Template.Labels = cleanLabels(dc.Spec.Template.Labels)
	deploy.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: cleanLabels(dc.Spec.Selector),
	}
	deploy.Spec.Paused = dc.Spec.Paused
	deploy.Spec.Replicas = &dc.Spec.Replicas
	deploy.Spec.RevisionHistoryLimit = dc.Spec.RevisionHistoryLimit
	deploy.Spec.MinReadySeconds = dc.Spec.MinReadySeconds

	if dc.Spec.Strategy.Type == ocappsv1.DeploymentStrategyTypeRolling {
		r := &appsv1.RollingUpdateDeployment{}

		if dc.Spec.Strategy.RollingParams != nil {
			r.MaxUnavailable = dc.Spec.Strategy.RollingParams.MaxUnavailable
			r.MaxSurge = dc.Spec.Strategy.RollingParams.MaxSurge

			if orig.Spec.Strategy.RollingParams.TimeoutSeconds != nil {
				timeout32 := int32(*orig.Spec.Strategy.RollingParams.TimeoutSeconds)
				deploy.Spec.ProgressDeadlineSeconds = &timeout32
			}
		}

		deploy.Spec.Strategy = appsv1.DeploymentStrategy{
			Type:          appsv1.RollingUpdateDeploymentStrategyType,
			RollingUpdate: r,
		}

	} else {
		deploy.Spec.Strategy = appsv1.DeploymentStrategy{
			Type: appsv1.RecreateDeploymentStrategyType,
		}

		if orig.Spec.Strategy.RecreateParams != nil && orig.Spec.Strategy.RecreateParams.TimeoutSeconds != nil {
			timeout32 := int32(*orig.Spec.Strategy.RecreateParams.TimeoutSeconds)
			deploy.Spec.ProgressDeadlineSeconds = &timeout32
		}
	}

	var triggers []trigger.ObjectFieldTrigger

	for _, dctrigger := range dc.Spec.Triggers {
		if dctrigger.Type == ocappsv1.DeploymentTriggerOnImageChange {
			for _, containername := range dctrigger.ImageChangeParams.ContainerNames {
				triggers = append(triggers, trigger.ObjectFieldTrigger{
					From: trigger.ObjectReference{
						Kind:       dctrigger.ImageChangeParams.From.Kind,
						Name:       dctrigger.ImageChangeParams.From.Name,
						Namespace:  dctrigger.ImageChangeParams.From.Namespace,
						APIVersion: dctrigger.ImageChangeParams.From.APIVersion,
					},
					FieldPath: fmt.Sprintf("spec.template.spec.containers[?(@.name==\"%s\")]", containername),
				})
			}
		}
	}

	if triggers != nil {
		triggersjson, err := json.Marshal(triggers)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal triggers: %w", err)
		}

		deploy.Annotations[trigger.TriggerAnnotationKey] = string(triggersjson)
	}

	return deploy, nil
}

func ToOuput(d *appsv1.Deployment, filetype string) ([]byte, error) {
	switch filetype {
	case "yaml":
		return yaml.Marshal(d)
	case "json":
		return json.Marshal(d)
	default:
		return nil, fmt.Errorf("unkown file type: %s (use json or yaml)", filetype)
	}
}

func cleanAnnotations(a map[string]string) map[string]string {
	var o = make(map[string]string)

	for key := range a {
		o[key] = a[key]
	}

	for i := range StripAnnotations {
		delete(o, StripAnnotations[i])
	}

	return o
}

func cleanLabels(l map[string]string) map[string]string {
	var o = make(map[string]string)

	for key := range l {
		o[key] = l[key]
	}

	for k, v := range ReplaceLabels {
		if _, ok := o[k]; ok {
			o[v] = o[k]
			delete(o, k)
		}
	}

	return o
}
