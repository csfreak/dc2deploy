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
	"fmt"

	ocappsv1 "github.com/openshift/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/resource"
)

func LoadDC(path string) (*ocappsv1.DeploymentConfig, error) {
	scheme := runtime.NewScheme()
	gvocappsv1 := schema.GroupVersion{Group: "apps.openshift.io", Version: "v1"}

	scheme.AddKnownTypes(gvocappsv1, &ocappsv1.DeploymentConfig{})

	b := resource.NewLocalBuilder().
		FilenameParam(false, &resource.FilenameOptions{Filenames: []string{path}}).
		WithScheme(scheme, gvocappsv1).
		Flatten().
		Do()

	r, err := b.Infos()
	if err != nil {
		return nil, fmt.Errorf("unable to build dc from file: %w", err)
	}

	return r[0].Object.(*ocappsv1.DeploymentConfig), nil
}
