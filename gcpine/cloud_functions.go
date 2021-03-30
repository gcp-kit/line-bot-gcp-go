package gcpine

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/line/line-bot-sdk-go/linebot"
)

// CloudFunctionsProps - props for Cloud Functions.
type CloudFunctionsProps interface {
	ParentEvent(ctx context.Context, message *pubsub.Message) error
	ChildEvent(ctx context.Context, message *pubsub.Message) error
	Props
}

type cloudFunctionsProps struct {
	pine        *GCPine
	parentTopic *pubsub.Topic
	childTopic  *pubsub.Topic
	secret      string
}

// NewCloudFunctionsProps - constructor
func NewCloudFunctionsProps(parent, child *pubsub.Topic) CloudFunctionsProps {
	return &cloudFunctionsProps{
		parentTopic: parent,
		childTopic:  child,
	}
}

// SetSecret - setter
func (cf *cloudFunctionsProps) SetSecret(secret string) {
	cf.secret = secret
}

// SetGCPine - setter
func (cf *cloudFunctionsProps) SetGCPine(pine *GCPine) {
	cf.pine = pine
}

// ReceiveWebHook - receive webhooks of LINE on Cloud Functions.
// CloudFunctions(Trigger: HTTP)
func (cf *cloudFunctionsProps) ReceiveWebHook(r *http.Request, w http.ResponseWriter) error {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "NG", http.StatusInternalServerError)
		return fmt.Errorf("failed to read all of the body: %w", err)
	}

	if !ValidateSignature(cf.secret, r.Header.Get("X-Line-Signature"), body) {
		http.Error(w, "NG", http.StatusBadRequest)
		return fmt.Errorf("failed to signature verification")
	}

	ctx := r.Context()
	if err = publishMessage(ctx, cf.parentTopic, body); err != nil {
		http.Error(w, "NG", http.StatusInternalServerError)
		return fmt.Errorf("could not publish message: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
	return nil
}

// ParentEvent - receive parent events on Cloud Functions.
// CloudFunctions(Trigger: Pub/Sub)
func (cf *cloudFunctionsProps) ParentEvent(ctx context.Context, message *pubsub.Message) error {
	events, err := ParseEvents(message.Data)
	if err != nil {
		return fmt.Errorf("could not parse the event: %w", err)
	}

	var wg sync.WaitGroup
	for _, event := range events {
		wg.Add(1)
		go func(ev *linebot.Event) {
			defer wg.Done()

			data, err := ev.MarshalJSON()
			if err != nil {
				return
			}

			if err = publishMessage(ctx, cf.childTopic, data); err != nil {
				return
			}
		}(event)
	}

	wg.Wait()

	return nil
}

// ChildEvent - receive child event on Cloud Functions.
// CloudFunctions(Trigger: Pub/Sub)
func (cf *cloudFunctionsProps) ChildEvent(ctx context.Context, message *pubsub.Message) error {
	event := new(linebot.Event)
	if err := event.UnmarshalJSON(message.Data); err != nil {
		return fmt.Errorf("faild to json unmarshal: %w", err)
	}

	if err := cf.pine.Execute(ctx, event); err != nil {
		if len(cf.pine.ErrMessages) > 0 {
			if err = cf.pine.SendReplyMessage(event.ReplyToken, cf.pine.ErrMessages); err != nil {
				return fmt.Errorf("failed to send error messages: %w", err)
			}
		}
		return fmt.Errorf("failed to function execution: %w", err)
	}

	return nil
}

func publishMessage(ctx context.Context, topic *pubsub.Topic, data []byte) error {
	msg := &pubsub.Message{Data: data}
	if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	return nil
}
