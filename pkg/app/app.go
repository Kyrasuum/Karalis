package app

import (
	"godev/pkg/stage"
)

var (
	CurApp App
)

type App interface {
	GetWidth() int32
	GetHeight() int32
	Running() bool
	Start()
	Exit()
	SetStage(nextStage stage.Stage)
	GetStage() stage.Stage
}
