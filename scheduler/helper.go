package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nafisalfiani/ketson-go-lib/appcontext"
	"github.com/nafisalfiani/ketson-go-lib/auth"
	"github.com/nafisalfiani/ketson-go-lib/codes"
	"github.com/nafisalfiani/ketson-go-lib/errors"
)

const (
	SystemID   string = "system"
	SystemName string = "system"
)

const (
	schedulerUserAgent string = "Cron Scheduler : %s"

	schedulerAssignError string = "Assigning Scheduler %s error: %s"

	schedulerRunning       string = "Running scheduler %s is running"
	schedulerDoneError     string = "Running scheduler %s error: %v"
	schedulerDoneSuccess   string = "Running scheduler %s success"
	schedulerTimeExecution string = "Scheduler %s done in %v"

	schedulerTimeTypeExact    string = "daily"
	schedulerTimeTypeInterval string = "interval"
)

type handlerFunc func(ctx context.Context) error

func (s *scheduler) AssignTask(conf TaskConf, task handlerFunc) {
	if conf.Enabled {
		var err error
		ctx := context.Background()
		schedulerFunc := s.taskWrapper(conf, task)

		switch conf.TimeType {
		case schedulerTimeTypeInterval:
			_, err = s.cron.Every(conf.Interval).Tag(conf.Name).Do(schedulerFunc)
		case schedulerTimeTypeExact:
			_, err = s.cron.Every(1).Day().Tag(conf.Name).At(conf.ScheduledTime).Do(schedulerFunc)
		default:
			err = errors.NewWithCode(codes.CodeInternalServerError, "Unknown Scheduler Task Time Type")
		}

		if err != nil {
			s.log.Fatal(ctx, fmt.Sprintf(schedulerAssignError, conf.Name, err.Error()))
		}

	}
}

func (s *scheduler) taskWrapper(conf TaskConf, task handlerFunc) func() {
	return func() {
		ctx := s.createContext(conf)
		s.log.Info(ctx, fmt.Sprintf(schedulerRunning, conf.Name))
		if err := task(ctx); err != nil {
			s.log.Error(ctx, fmt.Sprintf(schedulerDoneError, conf.Name, err))
		} else {
			s.log.Info(ctx, fmt.Sprintf(schedulerDoneSuccess, conf.Name))
		}

		startTime := appcontext.GetRequestStartTime(ctx)
		s.log.Info(ctx, fmt.Sprintf(schedulerTimeExecution, conf.Name, time.Since(startTime)))
	}
}

func (s *scheduler) createContext(conf TaskConf) context.Context {
	ctx := context.Background()
	ctx = appcontext.SetUserAgent(ctx, fmt.Sprintf(schedulerUserAgent, conf.Name))
	ctx = appcontext.SetRequestId(ctx, uuid.New().String())
	ctx = appcontext.SetRequestStartTime(ctx, time.Now())

	schedulerUser := auth.User{
		Id:   SystemID,
		Name: SystemName,
	}
	ctx = s.auth.SetUserAuthInfo(ctx, schedulerUser, nil)

	return ctx
}
