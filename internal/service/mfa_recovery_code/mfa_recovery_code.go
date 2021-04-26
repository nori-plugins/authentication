package mfa_recovery_code

import (
	"context"
	"time"

	"github.com/nori-plugins/authentication/internal/domain/service"

	"github.com/nori-plugins/authentication/internal/domain/entity"
	errors2 "github.com/nori-plugins/authentication/internal/domain/errors"
)

//@todo как передать сюда всю сессию? скорее, нужно извлечь пользовательский userID из контекста
//@todo что мы будем хранить в контексте?
func (srv *MfaRecoveryCodeService) GetMfaRecoveryCodes(ctx context.Context, data service.GetMfaRecoveryCodes) ([]*entity.MfaRecoveryCode, error) {
	//@todo будет ли использоваться паттерн?
	//@todo нужна ли максимальная длина, или указать всё в паттерне?
	//@todo указать ограничение на максимальную длину, связанную с базой данных?
	if err := data.Validate(); err != nil {
		return nil, err
	}
	var mfaRecoveryCodes []*entity.MfaRecoveryCode
	mfa_recovery_codes, err := srv.mfaRecoveryCodeHelper.Generate()
	if err != nil {
		return nil, err
	}
	for _, v := range mfa_recovery_codes {
		mfaRecoveryCodes = append(mfaRecoveryCodes, &entity.MfaRecoveryCode{
			ID:        0,
			UserID:    data.UserID,
			Code:      v,
			CreatedAt: time.Now(),
		})
	}

	if err := srv.transactor.Transact(ctx, func(tx context.Context) error {
		if err = srv.mfaRecoveryCodeRepository.DeleteMfaRecoveryCodes(ctx, data.UserID); err != nil {
			return err
		}
		if err = srv.mfaRecoveryCodeRepository.Create(ctx, mfaRecoveryCodes); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return mfaRecoveryCodes, nil
}

func (srv *MfaRecoveryCodeService) Apply(ctx context.Context, data service.ApplyData) error {
	//@todo проверить кэш и отп
	if err := data.Validate(); err != nil {
		return err
	}
	mfaRecoveryCode, err := srv.mfaRecoveryCodeRepository.FindByUserID(ctx, data.UserID, data.Code)
	if err != nil {
		return err
	}
	if mfaRecoveryCode == nil {
		return errors2.MfaRecoveryCodeNotFound
	}
	if err := srv.mfaRecoveryCodeRepository.DeleteMfaRecoveryCode(ctx, data.UserID, data.Code); err != nil {
		return err
	}
	return nil
}
