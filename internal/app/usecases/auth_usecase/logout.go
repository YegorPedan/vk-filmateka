package authUseCase

import "context"

func (a *authUseCase) Logout(ctx context.Context, refreshToken string) error {
	err := a.TokenService.DeleteByValue(ctx, refreshToken)
	if err != nil {
		return err
	}
	return nil
}
