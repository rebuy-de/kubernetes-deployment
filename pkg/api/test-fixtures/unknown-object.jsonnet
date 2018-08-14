[
    {
        apiVersion: "lifecycle.rebuy.com/v1alpha1",
        kind: "PodRestarter",
        metadata: {
            labels: { app: "thumbs", team: "platform" },
            name: "thumbs",
        },
        spec: {
            cooldownPeriod: "1h",
            maxUnavailable: 1,
            minAvailable: 2,
            restartCriteria: {
                maxAge: "1h",
            },
            selector: {
                matchLabels: { app: "thumbs" },
            },
        },
    },
]
