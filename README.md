# codepipeline-discord-notificator

This notificator receives notifications from CodePipeline and notify them to Discord via Webhook API.

## Build & Deploy

1. Run `make build` in the project directory.
1. Upload `function.zip` created to AWS Lambda.
    - Ensure that enviroment variables `DISCORD_WEBHOOK_ID`, `DISCORD_WEBHOOK_TOKEN`, `TZ`(optional) are set.
1. Create a SNS Topic and create a relation between CodePipeline and the SNS Topic.
    - Recommended events:
      - codepipeline-pipeline-pipeline-execution-failed
      - codepipeline-pipeline-pipeline-execution-canceled
      - codepipeline-pipeline-pipeline-execution-started
      - codepipeline-pipeline-pipeline-execution-resumed
      - codepipeline-pipeline-pipeline-execution-succeeded
      - codepipeline-pipeline-pipeline-execution-superseded
