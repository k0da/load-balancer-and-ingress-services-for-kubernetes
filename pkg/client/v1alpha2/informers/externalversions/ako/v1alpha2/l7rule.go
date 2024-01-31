/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha2

import (
	"context"
	time "time"

	akov1alpha2 "github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/apis/ako/v1alpha2"
	versioned "github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/client/v1alpha2/clientset/versioned"
	internalinterfaces "github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/client/v1alpha2/informers/externalversions/internalinterfaces"
	v1alpha2 "github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/client/v1alpha2/listers/ako/v1alpha2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// L7RuleInformer provides access to a shared informer and lister for
// L7Rules.
type L7RuleInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha2.L7RuleLister
}

type l7RuleInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewL7RuleInformer constructs a new informer for L7Rule type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewL7RuleInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredL7RuleInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredL7RuleInformer constructs a new informer for L7Rule type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredL7RuleInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AkoV1alpha2().L7Rules(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AkoV1alpha2().L7Rules(namespace).Watch(context.TODO(), options)
			},
		},
		&akov1alpha2.L7Rule{},
		resyncPeriod,
		indexers,
	)
}

func (f *l7RuleInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredL7RuleInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *l7RuleInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&akov1alpha2.L7Rule{}, f.defaultInformer)
}

func (f *l7RuleInformer) Lister() v1alpha2.L7RuleLister {
	return v1alpha2.NewL7RuleLister(f.Informer().GetIndexer())
}
