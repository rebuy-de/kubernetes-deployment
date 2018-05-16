{
    silo(name, team, replicas, cpuRequests, memoryRequests)::
        local metadata = {
            name: name,
            labels: {
                app: name,
                team: team,
                role: "silo",
            },
        };
        [

            {
                apiVersion: "v1",
                kind: "Service",
                metadata: metadata,
                spec: {
                    selector: {
                        app: name,
                    },
                    ports: [{
                        port: 80,
                        targetPort: 8080,
                    }],
                },
            },

            {
                apiVersion: "extensions/v1beta1",
                kind: "Ingress",
                metadata: metadata,

                spec: {
                    rules: [{
                        host: std.format("%s.%s.rebuy.io", [name, std.extVar("clusterName")]),
                        http: {
                            paths: [{
                                path: "/",
                                backend: {
                                    serviceName: name,
                                    servicePort: 80,
                                },
                            }],
                        },
                    }],
                },
            },

            {
                apiVersion: "extensions/v1beta1",
                kind: "Deployment",
                metadata: metadata,

                spec: {
                    revisionHistoryLimit: 5,
                    replicas: replicas,

                    strategy: {
                        rollingUpdate: {
                            maxUnavailable: 0,
                        },
                    },

                    selector: {
                        matchLabels: {
                            app: name,
                        },
                    },

                    template: {
                        metadata: metadata,

                        spec: {
                            terminationGracePeriodSeconds: 120,
                            containers: [
                                {
                                    name: name,
                                    image: std.format("my-registry.loc/%s:%s", [name, std.extVar("gitCommitID")]),
                                    imagePullPolicy: "Always",

                                    env: [{
                                        name: "SILO_PROFILE",
                                        value: "kubernetes",
                                    }],

                                    ports: [{
                                        containerPort: 8080,
                                    }],

                                    resources: {
                                        limits: {
                                            cpu: cpuRequests,
                                            memory: memoryRequests,
                                        },
                                        requests: {
                                            cpu: cpuRequests,
                                            memory: memoryRequests,
                                        },
                                    },

                                    readinessProbe: {
                                        httpGet: {
                                            path: "/health",
                                            port: 8080,
                                        },
                                        initialDelaySeconds: 15,
                                        timeoutSeconds: 1,
                                    },

                                    livenessProbe: {
                                        httpGet: {
                                            path: "/health",
                                            port: 8080,
                                        },
                                        initialDelaySeconds: 120,
                                        timeoutSeconds: 1,
                                    },
                                },
                            ],
                        },
                    },
                },
            },
        ],
}
