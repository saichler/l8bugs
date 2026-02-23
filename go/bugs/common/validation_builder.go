package common

import (
	"github.com/saichler/l8types/go/ifs"
)

type VB[T any] struct {
	typeName         string
	setID            SetIDFunc[T]
	validators       []func(*T, ifs.IVNic) error
	actionValidators []ActionValidateFunc[T]
	afterActions     []ActionValidateFunc[T]
}

func NewValidation[T any](typeName string, setID SetIDFunc[T]) *VB[T] {
	return &VB[T]{typeName: typeName, setID: setID}
}

func (b *VB[T]) Require(getter func(*T) string, name string) *VB[T] {
	b.validators = append(b.validators, func(e *T, _ ifs.IVNic) error {
		return ValidateRequired(getter(e), name)
	})
	return b
}

func (b *VB[T]) RequireInt64(getter func(*T) int64, name string) *VB[T] {
	b.validators = append(b.validators, func(e *T, _ ifs.IVNic) error {
		return ValidateRequiredInt64(getter(e), name)
	})
	return b
}

func (b *VB[T]) Enum(getter func(*T) int32, nameMap map[int32]string, name string) *VB[T] {
	b.validators = append(b.validators, func(e *T, _ ifs.IVNic) error {
		return ValidateEnum(getter(e), nameMap, name)
	})
	return b
}

func (b *VB[T]) DateNotZero(getter func(*T) int64, name string) *VB[T] {
	b.validators = append(b.validators, func(e *T, _ ifs.IVNic) error {
		return ValidateDateNotZero(getter(e), name)
	})
	return b
}

func (b *VB[T]) DateAfter(getter1, getter2 func(*T) int64, name1, name2 string) *VB[T] {
	b.validators = append(b.validators, func(e *T, _ ifs.IVNic) error {
		d1, d2 := getter1(e), getter2(e)
		if d1 == 0 || d2 == 0 {
			return nil
		}
		return ValidateDateAfter(d1, d2, name1, name2)
	})
	return b
}

func (b *VB[T]) Custom(fn func(*T, ifs.IVNic) error) *VB[T] {
	b.validators = append(b.validators, fn)
	return b
}

func (b *VB[T]) StatusTransition(cfg *StatusTransitionConfig[T]) *VB[T] {
	b.actionValidators = append(b.actionValidators, cfg.BuildValidator())
	return b
}

func (b *VB[T]) After(fn ActionValidateFunc[T]) *VB[T] {
	b.afterActions = append(b.afterActions, fn)
	return b
}

func (b *VB[T]) Build() ifs.IServiceCallback {
	validate := func(item *T, vnic ifs.IVNic) error {
		for _, v := range b.validators {
			if err := v(item, vnic); err != nil {
				return err
			}
		}
		return nil
	}
	if len(b.afterActions) > 0 {
		return NewServiceCallbackWithAfter(b.typeName, b.setID, validate,
			b.actionValidators, b.afterActions)
	}
	return NewServiceCallback(b.typeName, b.setID, validate, b.actionValidators...)
}
