package pipeline

import (
	"context"

	broker "github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/logger"
	internalconfig "github.com/layer5io/meshsync/internal/config"
	"github.com/layer5io/meshsync/pkg/model"
	"github.com/myntra/pipeline"

	kubeerror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type LocalResource struct {
	pipeline.StepContext
	log           logger.Handler
	dynamicClient dynamic.Interface
	brokerClient  broker.Handler
	config        internalconfig.PipelineConfig
}

func NewLocalResource(log logger.Handler, dclient dynamic.Interface, bclient broker.Handler, config internalconfig.PipelineConfig) *LocalResource {
	return &LocalResource{
		log:           log,
		dynamicClient: dclient,
		brokerClient:  bclient,
		config:        config,
	}
}

// Exec - step interface
func (c *LocalResource) Exec(request *pipeline.Request) *pipeline.Result {
	result, err := c.dynamicClient.Resource(schema.GroupVersionResource{
		Group:    c.config.Group,
		Version:  c.config.Version,
		Resource: c.config.Resource,
	}).Namespace(c.config.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.log.Error(ErrDynamicClient(c.config.Resource, err))
		if !kubeerror.IsNotFound(err) {
			return &pipeline.Result{
				Error: ErrList(c.config.Resource, err),
			}
		}
		return &pipeline.Result{
			Error: nil,
		}
	}
	c.log.Info("discovering: ", c.config.Resource)

	for _, item := range result.Items {
		err = c.brokerClient.Publish(c.config.PublishSubject, &broker.Message{
			ObjectType: broker.Single,
			EventType:  broker.Add,
			Object:     model.ParseList(item),
		})
		if err != nil {
			c.log.Error(err)
			return &pipeline.Result{
				Error: ErrPublish(c.config.Resource, err),
			}
		}
	}

	return &pipeline.Result{
		Error: nil,
	}
}

// Cancel - step interface
func (c *LocalResource) Cancel() error {
	c.Status("cancel step")
	return nil
}
