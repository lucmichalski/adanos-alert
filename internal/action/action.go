package action

import (
	"encoding/json"
	"sync"

	"github.com/mylxsw/adanos-alert/internal/queue"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
)

type Action interface {
	Validate(meta string) error
	Handle(rule repository.Rule, trigger repository.Trigger, grp repository.MessageGroup) error
}

type Manager interface {
	Resolve(f interface{}) error
	MustResolve(f interface{})
	Dispatch(action string) Action
	Run(action string) Action
	Register(name string, action Action)
}

type actionManager struct {
	cc      *container.Container
	lock    sync.RWMutex
	actions map[string]Action
}

func NewManager(cc *container.Container) Manager {
	return &actionManager{cc: cc, actions: make(map[string]Action)}
}

func (manager *actionManager) Resolve(f interface{}) error {
	return manager.cc.ResolveWithError(f)
}

func (manager *actionManager) MustResolve(f interface{}) {
	manager.cc.MustResolve(f)
}

// Dispatch dispatch a action to queue
func (manager *actionManager) Dispatch(action string) Action {
	return &QueueAction{
		action:  action,
		manager: manager,
	}
}

// Run execute a action
func (manager *actionManager) Run(action string) Action {
	manager.lock.RLock()
	defer manager.lock.RUnlock()

	return manager.actions[action]
}

// Register register a new action
func (manager *actionManager) Register(name string, action Action) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	manager.actions[name] = action
}

type QueueAction struct {
	action  string
	manager Manager
}

func (q *QueueAction) Validate(meta string) error {
	return nil
}

type Payload struct {
	Action  string                  `json:"action"`
	Rule    repository.Rule         `json:"rule"`
	Trigger repository.Trigger      `json:"trigger"`
	Group   repository.MessageGroup `json:"group"`
}

func (payload *Payload) Encode() []byte {
	data, _ := json.Marshal(payload)
	return data
}

func (payload *Payload) Decode(data []byte) error {
	return json.Unmarshal(data, payload)
}

func (q *QueueAction) Handle(rule repository.Rule, trigger repository.Trigger, grp repository.MessageGroup) error {
	return q.manager.Resolve(func(queueManager queue.Manager) error {
		payload := Payload{
			Action:  q.action,
			Trigger: trigger,
			Group:   grp,
			Rule:    rule,
		}

		id, err := queueManager.Enqueue(repository.QueueJob{
			Name:    "action",
			Payload: string(payload.Encode()),
		})
		if err != nil {
			return err
		}

		log.WithFields(log.Fields{
			"action": q.action,
			"id":     id,
		}).Debug("enqueue a action to queue")

		return nil
	})
}
