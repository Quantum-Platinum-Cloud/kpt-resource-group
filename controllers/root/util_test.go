// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package root

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"kpt.dev/resourcegroup/apis/kpt.dev/v1alpha1"
)

var _ = Describe("Util tests", func() {
	Describe("GroupKind to GroupVersionKind", func() {

		It("should get groupkind's from spec", func() {
			resources := []v1alpha1.ObjMetadata{
				{
					GroupKind: v1alpha1.GroupKind{
						Group: "",
						Kind:  "ConfigMap",
					},
					Name:      "nm1",
					Namespace: "ns1",
				},
				{
					GroupKind: v1alpha1.GroupKind{
						Group: "apps",
						Kind:  "Deployment",
					},
					Name:      "nm2",
					Namespace: "ns2",
				},
				{
					GroupKind: v1alpha1.GroupKind{
						Group: "groupname",
						Kind:  "KindName",
					},
					Name:      "nm3",
					Namespace: "ns3",
				},
			}
			spec := v1alpha1.ResourceGroupSpec{
				Resources: resources,
			}

			expected := map[schema.GroupKind]struct{}{
				schema.GroupKind{Group: "", Kind: "ConfigMap"}:         {},
				schema.GroupKind{Group: "apps", Kind: "Deployment"}:    {},
				schema.GroupKind{Group: "groupname", Kind: "KindName"}: {},
			}

			gkSet := getGroupKinds(spec)
			Expect(gkSet).Should(Equal(expected))
		})
	})
})
