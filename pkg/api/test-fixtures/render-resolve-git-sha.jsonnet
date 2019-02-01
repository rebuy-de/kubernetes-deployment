[
    {
        apiVersion: "v1",
        kind: "Pod",
        metadata: {
            name: "myapp-pod",
            labels: { app: "myapp" },
        },
        spec: {
            containers: [{
                name: "myapp-container",
                image: std.format("busybox:%s", std.native("resolveGitSHA")("github.com/rebuy-de/test@master")),
            }],
        },
    },
]
