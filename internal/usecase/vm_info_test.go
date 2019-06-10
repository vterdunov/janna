package usecase_test

import (
	"testing"

	"github.com/vterdunov/janna/internal/usecase"
)

func TestUsecase_VMInfo(t *testing.T) {
	// TODO: switch to table test
	vmWareRepositoryMock := new(usecase.MockVMWareRepository)
	vmWareRepositoryMock.On("VMInfo", "ddd").Return(usecase.VMInfoResponse{}, nil)
	u := usecase.NewUsecase(nil, vmWareRepositoryMock)
	_, _ = u.VMInfo("ddd")

	vmWareRepositoryMock.AssertExpectations(t)
}
