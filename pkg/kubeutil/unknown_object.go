package kubeutil

import (
	"encoding/json"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type UnknownObject struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata"`

	Raw map[string]interface{}
}

func (o *UnknownObject) GetObjectKind() schema.ObjectKind {
	return nil
}

func (o *UnknownObject) DeepCopyObject() runtime.Object {
	return nil
}

func (o *UnknownObject) FromJson(raw []byte) error {
	err := json.Unmarshal(raw, o)
	if err != nil {
		return err
	}

	o.Raw = map[string]interface{}{}
	return json.Unmarshal(raw, &o.Raw)
}

func (o *UnknownObject) syncRaw() {
	o.Raw["metadata"] = o.ObjectMeta
	o.Raw["apiVersion"] = o.TypeMeta.APIVersion
	o.Raw["kind"] = o.TypeMeta.Kind
}

func (o *UnknownObject) MarshalJSON() ([]byte, error) {
	o.syncRaw()
	return json.Marshal(o.Raw)
}
