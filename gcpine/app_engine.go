package gcpine

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/gcp-kit/gcpen"
	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

// AppEngineProps - props for App Engine.
type AppEngineProps interface {
	SetTasksClient(client *cloudtasks.Client)
	SetSecret(secret string)
	SetService(service string)
	ParentEvent(ctx context.Context, body []byte) error
	ChildEvent(ctx context.Context, body []byte) error
	Props
}

type appEngineProps struct {
	queuePath   string
	relativeURI string
	service     string
	pine        *GCPine
	client      *cloudtasks.Client
	secret      string
}

// NewAppEngineProps - constructor
func NewAppEngineProps(queuePath, relativeURI string) AppEngineProps {
	gcpen.Reload()
	return &appEngineProps{
		queuePath:   queuePath,
		relativeURI: relativeURI,
		service:     gcpen.ServiceName,
	}
}

// SetTasksClient - setter
func (ae *appEngineProps) SetTasksClient(client *cloudtasks.Client) {
	ae.client = client
}

// SetSecret - setter
func (ae *appEngineProps) SetSecret(secret string) {
	ae.secret = secret
}

// SetService - setter
func (ae *appEngineProps) SetService(service string) {
	ae.secret = service
}

// SetGCPine - setter
func (ae *appEngineProps) SetGCPine(pine *GCPine) {
	ae.pine = pine
}

func (ae *appEngineProps) createTask(ctx context.Context, data []byte) error {
	req := &tasks.CreateTaskRequest{
		Parent: ae.queuePath,
		Task: &tasks.Task{
			MessageType: &tasks.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &tasks.AppEngineHttpRequest{
					Body:        data,
					HttpMethod:  tasks.HttpMethod_POST,
					RelativeUri: ae.relativeURI,
				},
			},
		},
	}

	if len(ae.service) > 0 {
		gaeReq := req.Task.GetAppEngineHttpRequest()
		if gaeReq.AppEngineRouting == nil {
			gaeReq.AppEngineRouting = new(tasks.AppEngineRouting)
		}
		gaeReq.AppEngineRouting.Service = ae.service
	}

	if _, err := ae.client.CreateTask(ctx, req); err != nil {
		return fmt.Errorf("failed to create tasks: %w", err)
	}

	return nil
}

// ReceiveWebHook - receive webhooks of LINE on App Engine.
func (ae *appEngineProps) ReceiveWebHook(r *http.Request, w http.ResponseWriter) error {
	defer r.Body.Close()

	// guard
	if ae.secret == "" {
		return fmt.Errorf("secret is required")
	}
	if ae.client == nil {
		return fmt.Errorf("cloud tasks client is required")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "NG", http.StatusInternalServerError)
		return fmt.Errorf("failed to read all of the body: %w", err)
	}

	if !ValidateSignature(ae.secret, r.Header.Get("X-Line-Signature"), body) {
		http.Error(w, "NG", http.StatusBadRequest)
		return fmt.Errorf("failed to signature verification")
	}

	if err = ae.createTask(r.Context(), body); err != nil {
		http.Error(w, "NG", http.StatusInternalServerError)
		return fmt.Errorf("failed to creating a task: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, http.StatusText(http.StatusOK))
	return nil
}

// ParentEvent - receive parent events on Cloud Tasks.
func (ae *appEngineProps) ParentEvent(ctx context.Context, body []byte) error {
	// guard
	if ae.client == nil {
		return fmt.Errorf("cloud tasks client is required")
	}

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
func (ae *appEngineProps) ChildEvent(ctx context.Context, body []byte) error {
	// guard
	if ae.pine == nil {
		return fmt.Errorf("GCPine is required")
	}

	event := new(linebot.Event)
	if err := event.UnmarshalJSON(body); err != nil {
		return err
	}

	if err := ae.pine.Execute(ctx, event); err != nil {
		if len(ae.pine.ErrMessages) > 0 {
			if er := ae.pine.SendReplyMessage(event.ReplyToken, ae.pine.ErrMessages); er != nil {
				return fmt.Errorf("failed to send error messages: %w", err)
			}
		}
		return fmt.Errorf("failed to function execution: %w", err)
	}

	return nil
}
