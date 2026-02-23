package common

import (
	"errors"
	"fmt"
	"github.com/saichler/l8types/go/ifs"
)

type StatusTransitionConfig[T any] struct {
	StatusGetter  func(*T) int32
	StatusSetter  func(*T, int32)
	FilterBuilder func(*T) *T
	ServiceName   string
	ServiceArea   byte
	InitialStatus int32
	Transitions   map[int32][]int32
	StatusNames   map[int32]string
}

func (cfg *StatusTransitionConfig[T]) BuildValidator() ActionValidateFunc[T] {
	return func(entity *T, action ifs.Action, vnic ifs.IVNic) error {
		if action == ifs.POST {
			if cfg.InitialStatus > 0 {
				cfg.StatusSetter(entity, cfg.InitialStatus)
			}
			return nil
		}
		if action != ifs.PUT && action != ifs.PATCH {
			return nil
		}

		newStatus := cfg.StatusGetter(entity)
		filter := cfg.FilterBuilder(entity)
		old, err := GetEntity[T](cfg.ServiceName, cfg.ServiceArea, filter, vnic)
		if err != nil {
			return fmt.Errorf("failed to fetch existing entity for status validation: %w", err)
		}
		if old == nil {
			return errors.New("entity not found for status transition validation")
		}

		oldStatus := cfg.StatusGetter(old)
		if oldStatus == newStatus {
			return nil
		}

		allowed, ok := cfg.Transitions[oldStatus]
		if !ok {
			return fmt.Errorf("status %s is terminal and cannot be changed",
				cfg.statusName(oldStatus))
		}
		for _, a := range allowed {
			if a == newStatus {
				return nil
			}
		}
		return fmt.Errorf("invalid status transition from %s to %s",
			cfg.statusName(oldStatus), cfg.statusName(newStatus))
	}
}

func (cfg *StatusTransitionConfig[T]) statusName(val int32) string {
	if n, ok := cfg.StatusNames[val]; ok {
		return n
	}
	return fmt.Sprintf("UNKNOWN(%d)", val)
}
