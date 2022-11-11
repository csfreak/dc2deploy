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

package k8s

import (
	"context"
	"fmt"

	"github.com/csfreak/dc2deploy/pkg/writer"
	ocappsv1 "github.com/openshift/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	dcresource = schema.GroupVersionResource{
		Group:    "apps.openshift.io",
		Version:  "v1",
		Resource: "deploymentconfigs",
	}
)

func LoadDC(name string, namespace string) (*ocappsv1.DeploymentConfig, error) {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}

	resp, err := Client.Resource(dcresource).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		writer.WriteOut(2, "unable to load deploymentconfig: name %s, namespace %s", name, namespace)
		return nil, fmt.Errorf("unable to load %s: %w", name, err)
	}

	unstructured := resp.UnstructuredContent()

	var dc ocappsv1.DeploymentConfig

	err = runtime.DefaultUnstructuredConverter.
		FromUnstructured(unstructured, &dc)
	if err != nil {
		return nil, fmt.Errorf("unable to parse deploymentconfig %s: %w", name, err)
	}

	return &dc, nil
}
