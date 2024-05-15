package mq

import (
	"errors"
	"strings"

	log "xxx-server/application/logger"
	"xxx-server/domain/entity"
	"xxx-server/domain/repository"
	"xxx-server/infrastructure/config"

	"github.com/google/uuid"
	json "github.com/json-iterator/go"
	"github.com/wgdzlh/mqlib"
	"go.uber.org/zap"
)

const (
	MQ_LOG_TAG = "MqApi:"

	topicWorkflow = "workflow-service"
	tagSubmitTask = "service-task"

	TopicWorkflowResp = "rsim"
	TagWorkflowResp   = config.APP

	TopicTileServiceStatus = "tileServiceStatus"

	TopicInlayTask = "c_task"

	SepInTag    = "@"
	TagWildcard = "*"

	// defaultTimeout = 10 // in seconds
	// defaultConcurrent = 1
)

var (
	ErrMsgQueTimeout  = errors.New("msg queue timeout")
	ErrMissingHandler = errors.New("missing handler")
	ErrMqNotInitiated = errors.New("mq not initiated")

	AllTopicTags = []string{TagWildcard}
)

type SubHandler struct {
	Topic   string
	Tags    []string
	TagPre  string
	Matcher func(tag string, keys []string) bool
	F       func(tag string, keys []string, body []byte) error
}

type mqApiRepo struct {
	app    string
	client mqlib.PubSubClient
	uid    string // 本次运行态的独特id，用于数据隔离
	logTag string
	// handlers map[string]SubHandler // 订阅的每个topic对应的处理函数
	// msgQue   chan *mqlib.Message
	// timeout  time.Duration         // 排队超时
	// keySuffix string                // 发出/接收的消息带的key后缀
	// concurrent int
}

func NewMqService(handlers ...SubHandler) repository.MqApi {
	cfg := config.C.Mq
	if cfg.Disable {
		return nil
	}
	r := &mqApiRepo{
		app: config.APP,
		// msgQue:   make(chan *mqlib.Message, cfg.QueSize),
		// timeout:  time.Second * defaultTimeout,
		// handlers: map[string]SubHandler{},
		uid:    uuid.NewString(),
		logTag: MQ_LOG_TAG,
	}
	if cfg.App != "" {
		r.app = cfg.App
	}
	// r.keySuffix = SepInTag + r.app
	// if cfg.QueTimeout > 0 {
	// 	r.timeout = time.Second * time.Duration(cfg.QueTimeout)
	// }
	c, err := mqlib.NewPubSubClient(cfg.NameServer, r.app, r.subTopics(handlers)...)
	if err != nil {
		log.Fatal(r.logTag+"init mq client failed", zap.Error(err))
	}
	if config.C.Server.DevMode {
		c.SetDeDup(false)
	}
	r.client = c
	// go r.msgLoop()
	config.AddCallback(r.updateDeDupOpt)
	log.Info(r.logTag+"start mq api succeed", zap.String("app", r.app), zap.String("uid", r.uid))
	return r
}

func (r *mqApiRepo) subTopics(handlers []SubHandler) (topics []mqlib.Topic) {
	for _, h := range handlers {
		if len(h.Tags) == 0 {
			h.Tags = []string{r.uid} // 如果未指定，每个topic订阅当前运行态的unique tag，避免获取到全量消息
		}
		if h.Tags[0] == TagWildcard {
			h.Tags = nil
		}
		// r.handlers[h.Topic] = h
		topics = append(topics, mqlib.Topic{
			Name:     h.Topic,
			Tags:     h.Tags,
			Callback: r.getSubCallback(h),
		})
		// log.Info(r.logTag+"new sub", zap.String("topic", h.Topic), zap.Any("tags", h.Tags))
	}
	return
}

func (r *mqApiRepo) SendWorkflowReq(req *entity.WorkflowReq) (err error) {
	// h, ok := r.handlers[TopicWorkflowResp]
	// if !ok {
	// 	err = ErrMissingHandler
	// 	return
	// }
	mqReq := entity.WorkflowMqReq{
		WorkflowReq: req,
		Ext: entity.MQInfo{
			Topic:       TopicWorkflowResp,
			Tag:         TagWorkflowResp,
			CustomMsgId: req.TaskId,
			Keys:        []string{req.TaskId + SepInTag + req.TaskType, r.app},
		},
	}
	body, err := json.Marshal(mqReq)
	if err != nil {
		return
	}
	return r.SendMsg(&mqlib.Message{
		Topic: topicWorkflow,
		Tag:   tagSubmitTask,
		Body:  body,
	})
}

func (r *mqApiRepo) SendMsg(msg *mqlib.Message) error {
	if r.client == nil {
		return ErrMqNotInitiated
	}
	return r.client.SendMessage(msg)
}

// func (r *mqApiRepo) msgLoop() {
// 	var (
// 		h  SubHandler
// 		ok bool
// 	)
// 	for msg := range r.msgQue {
// 		if h, ok = r.handlers[msg.Topic]; ok {
// 			if h.Pre != "" && !strings.HasPrefix(msg.Tag, h.Pre) {
// 				continue
// 			}
// 			if h.Matcher != nil && !h.Matcher(msg.Tag, msg.Keys) {
// 				continue
// 			}
// 			log.Info(r.logTag+"got message", zap.Any("msg", msg))
// 			h.F(msg.Tag, msg.Keys, msg.Body)
// 			log.Info(r.logTag+"finish message", zap.Any("msg", msg))
// 		} else {
// 			log.Error(r.logTag+"unknown msg topic", zap.String("topic", msg.Topic))
// 		}
// 	}
// }

func (r *mqApiRepo) getSubCallback(h SubHandler) mqlib.SubCallback {
	return func(msg *mqlib.Message) (err error) {
		if h.TagPre != "" && !strings.HasPrefix(msg.Tag, h.TagPre) {
			return
		}
		if h.Matcher != nil && !h.Matcher(msg.Tag, msg.Keys) {
			return
		}
		// log.Info(r.logTag+"got message", zap.Any("msg", msg))
		err = h.F(msg.Tag, msg.Keys, msg.Body)
		log.Info(r.logTag+"finish message", zap.Any("msg", msg), zap.Error(err))
		return
	}
}

// func (r *mqApiRepo) ProcessMsg(msg *mqlib.Message) (err error) {
// 	select {
// 	case r.msgQue <- msg:
// 	case <-time.After(r.timeout):
// 		err = ErrMsgQueTimeout
// 	}
// 	return
// }

func (r *mqApiRepo) updateDeDupOpt(sc *config.SelfConfig) {
	if sc.Server.DevMode != config.C.Server.DevMode {
		if sc.Server.DevMode {
			r.client.SetDeDup(false)
		} else {
			r.client.SetDeDup(true)
		}
	}
}
