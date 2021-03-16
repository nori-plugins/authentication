package postgres

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nori-plugins/authentication/internal/domain/entity"
)

type MfaRecoveryCodeRepository struct {
	Db *gorm.DB
}

func (m MfaRecoveryCodeRepository) Create(ctx context.Context, e []entity.MfaRecoveryCode) error {
	var mfaRecoveryCodes []model

	for _, v := range e {
		mfaRecoveryCodes = append(mfaRecoveryCodes, NewModel(&v))
	}

	lastRecord := new(model)

	if err := m.Db.Create(mfaRecoveryCodes).Scan(&lastRecord).Error; err != nil {
		return err
	}
	lastRecord.Convert()

	return nil
}

func (m MfaRecoveryCodeRepository) Delete(ctx context.Context, e *entity.MfaRecoveryCode) error {
	panic("implement me")
}
