package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
)

type CodePipelineEventDetail struct {
	Pipeline         string `json:"pipeline"`
	ExecutionId      string `json:"execution-id"`
	ExecutionTrigger struct {
		TriggerType   string `json:"trigger-type"`
		TriggerDetail string `json:"trigger-detail"`
	} `json:"execution-trigger"`
	State   string  `json:"state"`
	Version float64 `json:"version"`
}

type CodePipelineEvent struct {
	Account              string                  `json:"account"`
	DetailType           string                  `json:"detailType"`
	Region               string                  `json:"region"`
	Source               string                  `json:"source"`
	Time                 time.Time               `json:"time"`
	NotificationRuleArn  string                  `json:"notificationRuleArn"`
	Detail               CodePipelineEventDetail `json:"detail"`
	Resources            []string                `json:"resources"`
	AdditionalAttributes map[string]interface{}  `json:"additionalAttributes"`
}

func createEmbed(event CodePipelineEvent) discord.Embed {
	builder := discord.NewEmbedBuilder().
		SetTitlef("%s: %s", event.Detail.Pipeline, event.Detail.State).
		AddField("ExecutionId", event.Detail.ExecutionId, false).
		AddField("URL", fmt.Sprintf("https://%s.console.aws.amazon.com/codesuite/codepipeline/pipelines/%s/executions/%s/timeline?region=%s", event.Region, event.Detail.Pipeline, event.Detail.ExecutionId, event.Region), false).
		AddField("Time", event.Time.Local().String(), false)

	switch event.Detail.State {
	case "STARTED":
		return builder.SetColor(7506394).
			AddField("ExecutionTrigger", event.Detail.ExecutionTrigger.TriggerType+":"+event.Detail.ExecutionTrigger.TriggerDetail, false).Build()
	case "SUCCEEDED":
		return builder.SetColor(3066993).Build()
	case "FAILED":
		return builder.SetColor(15158332).Build()
	case "CANCELED":
		return builder.SetColor(10070709).Build()
	default:
		return builder.SetColor(7506394).Build()
	}
}

func HandleRequest(ctx context.Context, event events.SNSEvent) (string, error) {
	webhookId := os.Getenv("DISCORD_WEBHOOK_ID")
	webhookToken := os.Getenv("DISCORD_WEBHOOK_TOKEN")

	if webhookId == "" || webhookToken == "" {
		return "error", fmt.Errorf("environment variables is not set")
	}

	client := webhook.New(snowflake.MustParse(webhookId), webhookToken)
	defer client.Close(ctx)

	for _, record := range event.Records {
		var pipelineEvent CodePipelineEvent
		if err := json.Unmarshal([]byte(record.SNS.Message), &pipelineEvent); err != nil {
			return "error", fmt.Errorf("failed to unmarshal event: %s, record: %+v", err, record)
		}

		if _, err := client.CreateEmbeds([]discord.Embed{createEmbed(pipelineEvent)}, rest.WithDelay(2*time.Second)); err != nil {
			return "error", fmt.Errorf("failed to send message to Discord: %s", err)
		}
	}
	return "ok", nil
}

func main() {
	lambda.Start(HandleRequest)
}
