package gcpine

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

// AppEngineProps - props for App Engine.
type AppEngineProps struct {
	QueuePath   string
	RelativeURI string
	Service     string

	client *cloudtasks.Client
	secret string
}

// NewAppEngineProps - constructor
func NewAppEngineProps(client *cloudtasks.Client, secret string) *AppEngineProps {
	return &AppEngineProps{
		client: client,
		secret: secret,
	}
}

// SetTQClient - setter
func (ae *AppEngineProps) SetTQClient(client *cloudtasks.Client) {
	ae.client = client
}

// SetSecret - setter
func (ae *AppEngineProps) SetSecret(secret string) {
	ae.secret = secret
}

func (ae *AppEngineProps) createTask(ctx context.Context, data []byte) error {
	req := &tasks.CreateTaskRequest{
		Parent: ae.QueuePath,
		Task: &tasks.Task{
			MessageType: &tasks.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &tasks.AppEngineHttpRequest{
					Body:        data,
					HttpMethod:  tasks.HttpMethod_POST,
					RelativeUri: ae.RelativeURI,
				},
			},
		},
	}

	if len(ae.Service) > 0 {
		gaeReq := req.Task.GetAppEngineHttpRequest()
		if gaeReq.AppEngineRouting == nil {
			gaeReq.AppEngineRouting = new(tasks.AppEngineRouting)
		}
		gaeReq.AppEngineRouting.Service = ae.Service
	}

	if _, err := ae.client.CreateTask(ctx, req); err != nil {
		return fmt.Errorf("failed to create tasks: %w", err)
	}

	return nil
}

// ReceiveWebHook - receive webhooks of LINE on App Engine.
func (ae *AppEngineProps) ReceiveWebHook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "NG", http.StatusInternalServerError)
		return
	}

	if !ValidateSignature(ae.secret, r.Header.Get("X-Line-Signature"), body) {
		http.Error(w, "NG", http.StatusBadRequest)
		return
	}

	if err = ae.createTask(r.Context(), body); err != nil {
		http.Error(w, "NG", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, http.StatusText(http.StatusOK))
}

// ParentEvent - receive parent events on Cloud Tasks.
func (ae *AppEngineProps) ParentEvent(ctx context.Context, body []byte) error {
	events, err := ParseEvents(body)
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

			if err = ae.createTask(ctx, data); err != nil {
				return
			}
		}(event)
	}

	wg.Wait()

	return nil
}

// ChildEvent - receive child event on Cloud Tasks.
func (ae *AppEngineProps) ChildEvent(ctx context.Context, pine *GCPine, body []byte) error {
	event := new(linebot.Event)
	if err := event.UnmarshalJSON(body); err != nil {
		return err
	}

	if err := pine.Execute(ctx, event); err != nil {
		if len(pine.ErrMessages) > 0 {
			if er := pine.SendReplyMessage(event.ReplyToken, pine.ErrMessages); er != nil {
				return fmt.Errorf("failed to send error messages: %w", err)
			}
		}
		return fmt.Errorf("failed to function execution: %w", err)
	}

	return nil
}
