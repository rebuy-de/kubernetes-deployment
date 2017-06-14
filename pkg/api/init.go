package api

// This file configures the deserializer of the Kubernetes go-client. It
// basically registeres all API namespace to the file decoders.
//
// See also: https://stackoverflow.com/q/44306554/393157

import (
	_ "k8s.io/client-go/pkg/api/install"
	_ "k8s.io/client-go/pkg/apis/apps/install"
	_ "k8s.io/client-go/pkg/apis/authentication/install"
	_ "k8s.io/client-go/pkg/apis/autoscaling/install"
	_ "k8s.io/client-go/pkg/apis/batch/install"
	_ "k8s.io/client-go/pkg/apis/extensions/install"
	_ "k8s.io/client-go/pkg/apis/policy/install"
	_ "k8s.io/client-go/pkg/apis/storage/install"
)
