local metadata(role, lang) = {
    name: "my-app",
    namespace: "default",
    labels: {
        app: "my-app",
        team: "me",
        role: role,
        lang: lang,
        test: std.extVar("testString"),
    },
};

local spec(lang) = {
    backoffLimit: 4,
    template: {
        spec: {
            restartPolicy: "OnFailure",
            containers: [
                {
                    name: "process",
                    image: "my-app",
                    args: [
                        lang,
                    ],
                    resources: {
                        limits: {
                            cpu: "2048m",
                            memory: "2048Mi",
                        },
                        requests: {
                            cpu: "1024m",
                            memory: "1024Mi",
                        },
                    },
                },
            ],
        },
    },
};

local langs = ["de", "fr", "nl"];

[
    {
        apiVersion: "batch/v1",
        kind: "Job",
        metadata: metadata("job", lang),
        spec: spec(lang),
    }
    for lang in langs
] + [
    {
        apiVersion: "batch/v2alpha1",
        kind: "CronJob",
        metadata: metadata("job", "de"),
        spec: {
            schedule: "0 */6 * * *",
            failedJobsHistoryLimit: 5,
            successfulJobsHistoryLimit: 1,

            jobTemplate: {
                metadata: metadata("job", "de"),
                spec: spec(lang),
            },
        },
    }
    for lang in langs
]
