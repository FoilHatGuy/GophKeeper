//go:build unit

// Code generated by MockGen. DO NOT EDIT.
// Source: ./src/server/database/controller.go

package app

import (
	context "context"
	cfg "gophKeeper/src/server/cfg"
	database "gophKeeper/src/server/database"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockStorageController is a mock of StorageController interface.
type MockStorageController struct {
	ctrl     *gomock.Controller
	recorder *MockStorageControllerMockRecorder
}

// MockStorageControllerMockRecorder is the mock recorder for MockStorageController.
type MockStorageControllerMockRecorder struct {
	mock *MockStorageController
}

// NewMockStorageController creates a new mock instance.
func NewMockStorageController(ctrl *gomock.Controller) *MockStorageController {
	mock := &MockStorageController{ctrl: ctrl}
	mock.recorder = &MockStorageControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageController) EXPECT() *MockStorageControllerMockRecorder {
	return m.recorder
}

// AddCard mocks base method.
func (m *MockStorageController) AddCard(ctx context.Context, uid, dataID, metadata string, data []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCard", ctx, uid, dataID, metadata, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCard indicates an expected call of AddCard.
func (mr *MockStorageControllerMockRecorder) AddCard(ctx, uid, dataID, metadata, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCard", reflect.TypeOf((*MockStorageController)(nil).AddCard), ctx, uid, dataID, metadata, data)
}

// AddCredentials mocks base method.
func (m *MockStorageController) AddCredentials(ctx context.Context, uid, dataID, metadata string, data []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCredentials", ctx, uid, dataID, metadata, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCredentials indicates an expected call of AddCredentials.
func (mr *MockStorageControllerMockRecorder) AddCredentials(ctx, uid, dataID, metadata, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCredentials", reflect.TypeOf((*MockStorageController)(nil).AddCredentials), ctx, uid, dataID, metadata, data)
}

// AddFile mocks base method.
func (m *MockStorageController) AddFile(ctx context.Context, uid, dataID, metadata string, data []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFile", ctx, uid, dataID, metadata, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddFile indicates an expected call of AddFile.
func (mr *MockStorageControllerMockRecorder) AddFile(ctx, uid, dataID, metadata, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFile", reflect.TypeOf((*MockStorageController)(nil).AddFile), ctx, uid, dataID, metadata, data)
}

// AddSession mocks base method.
func (m *MockStorageController) AddSession(ctx context.Context, uid, sid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSession", ctx, uid, sid)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSession indicates an expected call of AddSession.
func (mr *MockStorageControllerMockRecorder) AddSession(ctx, uid, sid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSession", reflect.TypeOf((*MockStorageController)(nil).AddSession), ctx, uid, sid)
}

// AddText mocks base method.
func (m *MockStorageController) AddText(ctx context.Context, uid, dataID, metadata string, data []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddText", ctx, uid, dataID, metadata, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddText indicates an expected call of AddText.
func (mr *MockStorageControllerMockRecorder) AddText(ctx, uid, dataID, metadata, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddText", reflect.TypeOf((*MockStorageController)(nil).AddText), ctx, uid, dataID, metadata, data)
}

// AddUser mocks base method.
func (m *MockStorageController) AddUser(ctx context.Context, uid, login, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", ctx, uid, login, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUser indicates an expected call of AddUser.
func (mr *MockStorageControllerMockRecorder) AddUser(ctx, uid, login, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockStorageController)(nil).AddUser), ctx, uid, login, password)
}

// GetCard mocks base method.
func (m *MockStorageController) GetCard(ctx context.Context, uid, dataID string) (string, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCard", ctx, uid, dataID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCard indicates an expected call of GetCard.
func (mr *MockStorageControllerMockRecorder) GetCard(ctx, uid, dataID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCard", reflect.TypeOf((*MockStorageController)(nil).GetCard), ctx, uid, dataID)
}

// GetCardHead mocks base method.
func (m *MockStorageController) GetCardHead(ctx context.Context, uid string) (database.CategoryHead, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCardHead", ctx, uid)
	ret0, _ := ret[0].(database.CategoryHead)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCardHead indicates an expected call of GetCardHead.
func (mr *MockStorageControllerMockRecorder) GetCardHead(ctx, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCardHead", reflect.TypeOf((*MockStorageController)(nil).GetCardHead), ctx, uid)
}

// GetCredentials mocks base method.
func (m *MockStorageController) GetCredentials(ctx context.Context, uid, dataID string) (string, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCredentials", ctx, uid, dataID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCredentials indicates an expected call of GetCredentials.
func (mr *MockStorageControllerMockRecorder) GetCredentials(ctx, uid, dataID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCredentials", reflect.TypeOf((*MockStorageController)(nil).GetCredentials), ctx, uid, dataID)
}

// GetCredentialsHead mocks base method.
func (m *MockStorageController) GetCredentialsHead(ctx context.Context, uid string) (database.CategoryHead, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCredentialsHead", ctx, uid)
	ret0, _ := ret[0].(database.CategoryHead)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCredentialsHead indicates an expected call of GetCredentialsHead.
func (mr *MockStorageControllerMockRecorder) GetCredentialsHead(ctx, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCredentialsHead", reflect.TypeOf((*MockStorageController)(nil).GetCredentialsHead), ctx, uid)
}

// GetFile mocks base method.
func (m *MockStorageController) GetFile(ctx context.Context, uid, dataID string) (string, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFile", ctx, uid, dataID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFile indicates an expected call of GetFile.
func (mr *MockStorageControllerMockRecorder) GetFile(ctx, uid, dataID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFile", reflect.TypeOf((*MockStorageController)(nil).GetFile), ctx, uid, dataID)
}

// GetFileHead mocks base method.
func (m *MockStorageController) GetFileHead(ctx context.Context, uid string) (database.CategoryHead, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileHead", ctx, uid)
	ret0, _ := ret[0].(database.CategoryHead)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFileHead indicates an expected call of GetFileHead.
func (mr *MockStorageControllerMockRecorder) GetFileHead(ctx, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileHead", reflect.TypeOf((*MockStorageController)(nil).GetFileHead), ctx, uid)
}

// GetText mocks base method.
func (m *MockStorageController) GetText(ctx context.Context, uid, dataID string) (string, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetText", ctx, uid, dataID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetText indicates an expected call of GetText.
func (mr *MockStorageControllerMockRecorder) GetText(ctx, uid, dataID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetText", reflect.TypeOf((*MockStorageController)(nil).GetText), ctx, uid, dataID)
}

// GetTextHead mocks base method.
func (m *MockStorageController) GetTextHead(ctx context.Context, uid string) (database.CategoryHead, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTextHead", ctx, uid)
	ret0, _ := ret[0].(database.CategoryHead)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTextHead indicates an expected call of GetTextHead.
func (mr *MockStorageControllerMockRecorder) GetTextHead(ctx, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTextHead", reflect.TypeOf((*MockStorageController)(nil).GetTextHead), ctx, uid)
}

// GetUserData mocks base method.
func (m *MockStorageController) GetUserData(ctx context.Context, login string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserData", ctx, login)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUserData indicates an expected call of GetUserData.
func (mr *MockStorageControllerMockRecorder) GetUserData(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserData", reflect.TypeOf((*MockStorageController)(nil).GetUserData), ctx, login)
}

// Initialise mocks base method.
func (m *MockStorageController) Initialise(ctx context.Context, config *cfg.ConfigT) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialise", ctx, config)
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialise indicates an expected call of Initialise.
func (mr *MockStorageControllerMockRecorder) Initialise(ctx, config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialise", reflect.TypeOf((*MockStorageController)(nil).Initialise), ctx, config)
}

// RefreshSession mocks base method.
func (m *MockStorageController) RefreshSession(ctx context.Context, sid string) (string, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshSession", ctx, sid)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RefreshSession indicates an expected call of RefreshSession.
func (mr *MockStorageControllerMockRecorder) RefreshSession(ctx, sid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshSession", reflect.TypeOf((*MockStorageController)(nil).RefreshSession), ctx, sid)
}

// UpdateSession mocks base method.
func (m *MockStorageController) UpdateSession(ctx context.Context, uid, sid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSession", ctx, uid, sid)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSession indicates an expected call of UpdateSession.
func (mr *MockStorageControllerMockRecorder) UpdateSession(ctx, uid, sid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSession", reflect.TypeOf((*MockStorageController)(nil).UpdateSession), ctx, uid, sid)
}
