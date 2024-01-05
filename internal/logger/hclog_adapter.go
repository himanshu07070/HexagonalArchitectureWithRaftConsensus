package logger

import (
	"hexagonal-architecture/internal"

	"github.com/hashicorp/go-hclog"
)

type HclogAdapter struct {
	logger hclog.Logger
}

func NewHclogAdapter() internal.LoggerPort {
	hcLogger := hclog.New(&hclog.LoggerOptions{
		Name:  "hexagonal-architecture",
		Level: hclog.Trace,
	})

	return &HclogAdapter{logger: hcLogger}
}
func (h *HclogAdapter) Trace(message string, args ...interface{}) {
	h.logger.Trace(message, args...)
}
func (h *HclogAdapter) Info(message string, args ...interface{}) {
	h.logger.Info(message, args...)
}
func (h *HclogAdapter) Debug(message string, args ...interface{}) {
	h.logger.Debug(message, args...)
}
func (h *HclogAdapter) Error(message string, args ...interface{}) {
	h.logger.Error(message, args...)
}
