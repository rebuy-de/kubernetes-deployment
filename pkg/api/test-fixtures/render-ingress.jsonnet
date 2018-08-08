{
  "apiVersion": "v1",
  "kind": "Service",

  "metadata": {
    "name": "zipkin",
    "namespace": "default",
    "labels": {
      "app": "jaeger",
      "jaeger-infra": "zipkin-service",
    },
  },
  "spec": {
    "type": "ClusterIP"
    "selector": {
      "app": "jaeger-collector",
    },
    "ports": [
      {
        "name": "jaeger-collector-zipkin",
        "port": 9411,
        "protocol": "TCP",
        "targetPort": 9411,
      }
    ],
  },
}
