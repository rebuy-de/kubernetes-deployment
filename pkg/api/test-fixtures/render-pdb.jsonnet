[
    {
        apiVersion: "policy/v1beta1",
        kind: "PodDisruptionBudget",
        metadata: {
            name: "my-app",
            namespace: "default",
            labels: {
                app: "my-app",
                team: "me",
                test: std.extVar("testString"),
            },
        },
        spec: {
            maxUnavailable: 1,
            selector: {
                matchLabels: {
                    app: "my-app",
                },
            },
        },
    },
]
