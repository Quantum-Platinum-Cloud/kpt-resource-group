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

package handler

import (
	"time"

	"github.com/go-logr/glogr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/event"

	"kpt.dev/resourcegroup/apis/kpt.dev/v1alpha1"
	"kpt.dev/resourcegroup/controllers/resourcemap"
)

type fakeMapping struct{}

func (m fakeMapping) Get(gknn v1alpha1.ObjMetadata) []types.NamespacedName {
	if gknn.Kind == "MyKind" {
		return []types.NamespacedName{
			{
				Name:      "my-name",
				Namespace: "my-namespace",
			},
			{
				Name:      "name2",
				Namespace: "namespace2",
			},
		}
	}
	return []types.NamespacedName{
		{
			Name:      "name1",
			Namespace: "namespace1",
		},
		{
			Name:      "name2",
			Namespace: "namespace2",
		},
	}

}

func (m fakeMapping) GetResources(_ schema.GroupKind) []v1alpha1.ObjMetadata {
	return nil
}

func (m fakeMapping) SetStatus(_ v1alpha1.ObjMetadata, _ *resourcemap.CachedStatus) {
	return
}

var _ = Describe("Util tests", func() {
	Describe("EventHandler", func() {
		It("push events from one event handler", func() {
			ch := make(chan event.GenericEvent)
			h := EnqueueEventToChannel{
				Mapping: fakeMapping{},
				Channel: ch,
				Log:     glogr.New(),
			}
			u := &unstructured.Unstructured{}

			// Push an event to channel
			go func() { h.OnAdd(u) }()

			// Consume an event from the channel
			e := <-ch
			Expect(e.Object.GetName()).Should(Equal("name1"))
			Expect(e.Object.GetNamespace()).Should(Equal("namespace1"))

			// Consume another event from the channel
			e = <-ch
			Expect(e.Object.GetName()).Should(Equal("name2"))
			Expect(e.Object.GetNamespace()).Should(Equal("namespace2"))
		})

		It("push events from multiple event handler", func() {
			ch := make(chan event.GenericEvent)
			h1 := EnqueueEventToChannel{
				Mapping: fakeMapping{},
				Channel: ch,
				Log:     glogr.New(),
			}

			h2 := EnqueueEventToChannel{
				Mapping: fakeMapping{},
				Channel: ch,
				Log:     glogr.New(),
				GVK:     schema.GroupVersionKind{Kind: "MyKind"},
			}

			u1 := &unstructured.Unstructured{}
			u2 := &unstructured.Unstructured{}
			u2.SetGroupVersionKind(schema.GroupVersionKind{Kind: "MyKind"})
			// Push an event to channel
			go func() { h1.OnAdd(u1) }()
			go func() { h2.OnDelete(u2) }()
			time.Sleep(time.Second)

			// Consume four events from the channel
			// These events should be from h1 and h2.
			// There should be
			//   1 event for "name1"
			//   1 event for "my-name"
			//   2 events for "name2"
			names := map[string]int{
				"name1":   0,
				"name2":   0,
				"my-name": 0,
			}
			for i := 1; i <= 4; i++ {
				e := <-ch
				name := e.Object.GetName()
				names[name] += 1
			}

			Expect(names["name1"]).Should(Equal(1))
			Expect(names["name2"]).Should(Equal(2))
			Expect(names["my-name"]).Should(Equal(1))
		})
	})
})