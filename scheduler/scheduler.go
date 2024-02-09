package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/nafisalfiani/ketson-go-lib/auth"
	"github.com/nafisalfiani/ketson-go-lib/log"
)

var (
	once = &sync.Once{}
)

type Interface interface {
	Run()
	AssignTask(conf TaskConf, task handlerFunc)
	TriggerScheduler(name string) error
}

type TaskConf struct {
	Name          string
	Enabled       bool
	TimeType      string
	Interval      time.Duration
	ScheduledTime string
}

type scheduler struct {
	cron *gocron.Scheduler
	log  log.Interface
	auth auth.Interface
}

func Init(log log.Interface, auth auth.Interface) Interface {
	s := &scheduler{}
	once.Do(func() {
		cron := gocron.NewScheduler(time.UTC)
		cron.TagsUnique()

		s = &scheduler{
			cron: cron,
			log:  log,
			auth: auth,
		}
	})

	return s
}

func (s *scheduler) Run() {
	s.cron.StartAsync()
	s.log.Info(context.Background(), "Scheduler is running")
}

func (s *scheduler) TriggerScheduler(name string) error {
	return s.cron.RunByTag(name)
}
