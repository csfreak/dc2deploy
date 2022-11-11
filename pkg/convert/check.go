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
	"github.com/csfreak/dc2deploy/pkg/writer"
	ocappsv1 "github.com/openshift/api/apps/v1"
)

type Warning struct {
	Name        string
	Path        string
	Description string
}

var (
	OwnerReferenceWarning = &Warning{
		Name:        "OwnerReference",
		Path:        "metadata.ownerReferences",
		Description: "DeploymentConfigs with OwnerReferences set are managed by another process.",
	}
	UnsupportedFeatureTestWarning = &Warning{
		Name:        "UnsupportedFeature - Test",
		Path:        "spec.test",
		Description: "The test feature is not supported on Deployments.",
	}
	UnsupportedFeatureCustomWarning = &Warning{
		Name:        "UnsupportedFeature - Custom Strategy",
		Path:        "spec.strategy.type",
		Description: "The custom deployment strategy is not supported on Deployments.",
	}
	UnsupportedFeatureHooksWarning = &Warning{
		Name:        "UnsupportedFeature - Lifecycle Hooks",
		Path:        "spec.strategy.*Params.['pre','mid','post']",
		Description: "The LifeCycleHooks are not supported on Deployments.",
	}
	UnsupportedFeatureRollingIntervalSecondsWarning = &Warning{
		Name:        "UnsupportedFeature - Rolling IntervalSeconds ",
		Path:        "spec.strategy.RollingParams.IntervalSeconds",
		Description: "The IntervalSeconds setting is not supported on Deployments.",
	}
	UnsupportedFeatureRollingUpdatePeriodSecondsWarning = &Warning{
		Name:        "UnsupportedFeature - Rolling UpdatePeriodSeconds ",
		Path:        "spec.strategy.RollingParams.UpdatePeriodSeconds",
		Description: "The UpdatePeriodSeconds setting is not supported on Deployments.",
	}
	ChangedLabelWarning = &Warning{
		Name:        "Selector Labels Changed",
		Path:        "spec.selector",
		Description: "The Selector label 'deploymentconfig' will be changed to 'deployment'.",
	}
)

func CheckFeatures(orig *ocappsv1.DeploymentConfig) []*Warning {
	var result []*Warning

	if len(orig.OwnerReferences) != 0 {
		result = append(result, OwnerReferenceWarning)
	}

	if orig.Spec.Test {
		result = append(result, UnsupportedFeatureTestWarning)
	}

	switch {
	case orig.Spec.Strategy.Type == ocappsv1.DeploymentStrategyTypeCustom:
		result = append(result, UnsupportedFeatureCustomWarning)
	case orig.Spec.Strategy.RollingParams != nil:
		if orig.Spec.Strategy.RollingParams.Pre != nil || orig.Spec.Strategy.RollingParams.Post != nil {
			result = append(result, UnsupportedFeatureHooksWarning)
		}

		if orig.Spec.Strategy.RollingParams.IntervalSeconds != nil {
			result = append(result, UnsupportedFeatureRollingIntervalSecondsWarning)
		}

		if orig.Spec.Strategy.RollingParams.UpdatePeriodSeconds != nil {
			result = append(result, UnsupportedFeatureRollingUpdatePeriodSecondsWarning)
		}

	case orig.Spec.Strategy.RecreateParams != nil:
		if orig.Spec.Strategy.RecreateParams.Pre != nil || orig.Spec.Strategy.RecreateParams.Mid != nil || orig.Spec.Strategy.RecreateParams.Post != nil {
			result = append(result, UnsupportedFeatureHooksWarning)
		}
	}

	if _, ok := orig.Spec.Selector[DeploymentConfigPodLabel]; ok {
		result = append(result, ChangedLabelWarning)
	}

	return result
}

func (w *Warning) Print(level uint8) {
	writer.WriteErr(level, "Conversion Warning: %s\n%s: path - %s", w.Name, w.Description, w.Path)
}
