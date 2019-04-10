// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crd

import (
	"log"
	"reflect"
	"testing"

	"github.com/ghodss/yaml"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type multiVersionTestcase struct {
	inputPackage string
	types        []string

	listFilesFn listFilesFn
	listDirsFn  listDirsFn
	// map of path to file content.
	inputFiles map[string][]byte

	expectedCrdSpecs map[schema.GroupKind][]byte
	expectedErr      error
}

func TestMultiVerGenerate(t *testing.T) {
	testcases := []multiVersionTestcase{
		{
			inputPackage: "github.com/myorg/myapi",
			types:        []string{"Toy"},
			listDirsFn: func(pkgPath string) (strings []string, e error) {
				return []string{"v1", "v1alpha1"}, nil
			},
			listFilesFn: func(pkgPath string) (s string, strings []string, e error) {
				return pkgPath, []string{"types.go"}, nil
			},
			inputFiles: map[string][]byte{
				"github.com/myorg/myapi/v1/types.go": []byte(`
package v1

// +groupName=foo.bar.com
// +versionName=v1

// +kubebuilder:resource:path=toys,shortName=to;ty
// +kubebuilder:singular=toy

// Toy is a toy struct
type Toy struct {
	// +kubebuilder:validation:Maximum=90
	// +kubebuilder:validation:Minimum=1

	// Replicas is a number
	Replicas int32 ` + "`" + `json:"replicas"` + "`" + `
}
`),
				"github.com/myorg/myapi/v1alpha1/types.go": []byte(`
package v1alpha1

// +groupName=fun.myk8s.io
// +versionName=v1alpha1

// +kubebuilder:resource

// Toy is a toy struct
type Toy struct {
	// +kubebuilder:validation:MaxLength=15
	// +kubebuilder:validation:MinLength=1

	// Name is a string
	Name string ` + "`" + `json:"name,omitempty"` + "`" + `

	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Minimum=1

	// Replicas is a number
	Replicas int32 ` + "`" + `json:"replicas"` + "`" + `
}
`),
			},
			expectedCrdSpecs: map[schema.GroupKind][]byte{
				schema.GroupKind{Group: "foo.bar.com"}: []byte(`Toy:
 description: Toy is a toy struct
 properties:
   replicas:
     description: Replicas is a number
     type: integer
 required:
 - replicas
 type: object`),
			},
		},
	}
	for _, tc := range testcases {
		fs, err := prepareTestFs(tc.inputFiles)
		if err != nil {
			t.Errorf("unable to prepare the in-memory fs for testing: %v", err)
			continue
		}

		op := &MultiVersionOptions{
			InputPackage: tc.inputPackage,
			Types:        tc.types,
			listDirsFn:   tc.listDirsFn,
			listFilesFn:  tc.listFilesFn,
			fs:           fs,
		}

		crdSpecs := op.parse()

		log.Println(len(crdSpecs))
		for gk, spec := range crdSpecs {
			log.Println("#########")
			log.Println(gk)
			yamlSpec, _ := yaml.Marshal(spec)
			log.Printf("yamlSpec: %s", yamlSpec)
			log.Println("#########")
		}

		if len(tc.expectedCrdSpecs) > 0 {
			expectedSpecsByKind := map[schema.GroupKind]*v1beta1.CustomResourceDefinitionSpec{}
			for gk := range tc.expectedCrdSpecs {
				var spec v1beta1.CustomResourceDefinitionSpec
				err = yaml.Unmarshal(tc.expectedCrdSpecs[gk], &spec)
				if err != nil {
					t.Errorf("unable to unmarshal the expected crd spec: %v", err)
					continue
				}
				expectedSpecsByKind[gk] = &spec
			}

			if !reflect.DeepEqual(crdSpecs, expectedSpecsByKind) {
				t.Errorf("expected: %s, but got: %s", expectedSpecsByKind, crdSpecs)
				continue
			}
		}
	}
}
