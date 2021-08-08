package mailservices

// import (
// 	"context"

// 	"dating/internal/app/api/types"

// 	// "github.com/stretchr/testify/mock"
// )

// type RepositoryMock struct {
// 	mock.Mock
// }

// func (mock *RepositoryMock) FindByEmail(ctx context.Context, email string) (*types.EmailVerification, error) {
// 	args := mock.Called()
// 	result := args.Get(0)
// 	if result == nil {
// 		return nil, args.Error(1)
// 	}
// 	return result.(*types.EmailVerification), args.Error(1)
// }
// func (mock *RepositoryMock) Insert(ctx context.Context, emailVerification types.EmailVerification) error {

// }
