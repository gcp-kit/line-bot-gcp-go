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

// ReceiveWebHook - receive webhooks of LINE on Cloud Functions.
// CloudFunctions(Trigger: HTTP)
func ReceiveWebHook(r *http.Request, w http.ResponseWriter, secret string, topic *pubsub.Topic) error {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "NG", http.StatusInternalServerError)
		return fmt.Errorf("failed to read all of the body: %w", err)
	}

	if !ValidateSignature(secret, r.Header.Get("X-Line-Signature"), body) {
		http.Error(w, "NG", http.StatusBadRequest)
		return fmt.Errorf("failed to signature verification")
	}

	ctx := r.Context()
	if err = publishMessage(ctx, topic, body); err != nil {
		http.Error(w, "NG", http.StatusInternalServerError)
		return fmt.Errorf("could not publish message: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
	return nil
}

// ParentEvent - receive parent events on Cloud Functions.
// CloudFunctions(Trigger: Pub/Sub)
func ParentEvent(ctx context.Context, message *pubsub.Message, topic *pubsub.Topic) error {
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

			if err = publishMessage(ctx, topic, data); err != nil {
				return
			}
		}(event)
	}

	wg.Wait()

	return nil
}

// ChildEvent - receive child event on Cloud Functions.
// CloudFunctions(Trigger: Pub/Sub)
func ChildEvent(ctx context.Context, message *pubsub.Message, pine *GCPine) error {
	event := new(linebot.Event)
	if err := event.UnmarshalJSON(message.Data); err != nil {
		return fmt.Errorf("faild to json unmarshal: %w", err)
	}

	if err := pine.Execute(ctx, event); err != nil {
		if len(pine.ErrMessages) > 0 {
			if err = pine.SendReplyMessage(event.ReplyToken, pine.ErrMessages); err != nil {
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
