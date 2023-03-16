// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/lastingasset/wallet-service/go-iden3-auth/state (interfaces: StateGetter)

// Package mock_state is a generated GoMock package.
package mock_state

import (
	big "math/big"
	reflect "reflect"

	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	gomock "github.com/golang/mock/gomock"
	state "github.com/lastingasset/wallet-service/go-iden3-auth/state"
)

// MockStateGetter is a mock of StateGetter interface.
type MockStateGetter struct {
	ctrl     *gomock.Controller
	recorder *MockStateGetterMockRecorder
}

// MockStateGetterMockRecorder is the mock recorder for MockStateGetter.
type MockStateGetterMockRecorder struct {
	mock *MockStateGetter
}

// NewMockStateGetter creates a new mock instance.
func NewMockStateGetter(ctrl *gomock.Controller) *MockStateGetter {
	mock := &MockStateGetter{ctrl: ctrl}
	mock.recorder = &MockStateGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStateGetter) EXPECT() *MockStateGetterMockRecorder {
	return m.recorder
}

// GetStateInfoByState mocks base method.
func (m *MockStateGetter) GetStateInfoByState(arg0 *bind.CallOpts, arg1 *big.Int) (state.StateV2StateInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateInfoByState", arg0, arg1)
	ret0, _ := ret[0].(state.StateV2StateInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStateInfoByState indicates an expected call of GetStateInfoByState.
func (mr *MockStateGetterMockRecorder) GetStateInfoByState(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateInfoByState", reflect.TypeOf((*MockStateGetter)(nil).GetStateInfoByState), arg0, arg1)
}
