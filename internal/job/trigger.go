package job

import (
	"github.com/mylxsw/adanos-alert/internal/action"
	"github.com/mylxsw/adanos-alert/internal/repository"
	matcher "github.com/mylxsw/adanos-alert/internal/rule"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/go-toolkit/container"
	"go.mongodb.org/mongo-driver/bson"
)

type TriggerJob struct {
	app *container.Container
}

func NewTrigger(app *container.Container) *TriggerJob {
	return &TriggerJob{app: app}
}

func (a TriggerJob) Handle() {
	log.Debug("trigger actions...")

	a.app.MustResolve(a.processMessageGroups)
}

func (a TriggerJob) processMessageGroups(groupRepo repository.MessageGroupRepo, ruleRepo repository.RuleRepo) error {
	return groupRepo.Traverse(bson.M{"status": repository.MessageGroupStatusPending}, func(grp repository.MessageGroup) error {
		rule, err := ruleRepo.Get(grp.Rule.ID)
		if err != nil {
			log.Errorf("rule not exist: %s", err)
			return err
		}

		hasError := false
		maxFailedCount := 0
		triggers := make([]repository.Trigger, 0)
		for _, trigger := range rule.Triggers {
			// check whether the trigger has been executed
			for _, act := range grp.Actions {
				if act.ID == trigger.ID && act.Status == repository.TriggerStatusOK {
					continue
				}
			}

			tm, err := matcher.NewTriggerMatcher(trigger)
			if err != nil {
				log.Errorf("create matcher failed: %s", err)
				continue
			}

			matched, err := tm.Match(matcher.TriggerContext{Group: grp})
			if err != nil {
				log.Errorf("trigger matcher match failed: %s", err)
				continue
			}

			if matched {
				if err := action.Factory(trigger.Action).Handle(trigger); err != nil {
					trigger.Status = repository.TriggerStatusFailed
					trigger.FailedCount = trigger.FailedCount + 1
					trigger.FailedReason = err.Error()
					hasError = true
				} else {
					trigger.Status = repository.TriggerStatusOK
				}

				triggers = append(triggers, trigger)
				if trigger.FailedCount > maxFailedCount {
					maxFailedCount = trigger.FailedCount
				}
			}
		}

		if hasError {
			// if trigger failed count > 3, then set message group failed
			if maxFailedCount > 3 {
				grp.Status = repository.MessageGroupStatusFailed
			}
		} else {
			grp.Status = repository.MessageGroupStatusOK
		}

		grp.Actions = mergeActions(grp.Actions, triggers)
		return groupRepo.Update(grp.ID, grp)
	})
}

func mergeActions(actions []repository.Trigger, triggers []repository.Trigger) []repository.Trigger {
	newActions := make([]repository.Trigger, 0)
	for _, tr := range triggers {
		existed := false
		for i, act := range actions {
			if tr.ID == act.ID {
				actions[i] = tr
				existed = true
				break
			}
		}

		if existed {
			break
		}

		newActions = append(newActions, tr)
	}
	actions = append(actions, newActions...)
	return actions
}