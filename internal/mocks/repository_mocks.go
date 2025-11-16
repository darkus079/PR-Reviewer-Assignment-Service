package mocks

import (
	context "context"
	models "pr-reviewer-assignment-service/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

func (m *MockUserRepository) CreateUser(arg0 context.Context, arg1 *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockUserRepositoryMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepository)(nil).CreateUser), arg0, arg1)
}

func (m *MockUserRepository) GetUserByID(arg0 context.Context, arg1 string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", arg0, arg1)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockUserRepositoryMockRecorder) GetUserByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUserRepository)(nil).GetUserByID), arg0, arg1)
}

func (m *MockUserRepository) GetUsersByTeam(arg0 context.Context, arg1 string) ([]*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersByTeam", arg0, arg1)
	ret0, _ := ret[0].([]*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockUserRepositoryMockRecorder) GetUsersByTeam(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersByTeam", reflect.TypeOf((*MockUserRepository)(nil).GetUsersByTeam), arg0, arg1)
}

func (m *MockUserRepository) GetActiveUsersByTeam(arg0 context.Context, arg1 string) ([]*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActiveUsersByTeam", arg0, arg1)
	ret0, _ := ret[0].([]*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockUserRepositoryMockRecorder) GetActiveUsersByTeam(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveUsersByTeam", reflect.TypeOf((*MockUserRepository)(nil).GetActiveUsersByTeam), arg0, arg1)
}

func (m *MockUserRepository) SetUserActiveStatus(arg0 context.Context, arg1 string, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUserActiveStatus", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockUserRepositoryMockRecorder) SetUserActiveStatus(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUserActiveStatus", reflect.TypeOf((*MockUserRepository)(nil).SetUserActiveStatus), arg0, arg1, arg2)
}

func (m *MockUserRepository) UpdateUser(arg0 context.Context, arg1 *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockUserRepositoryMockRecorder) UpdateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserRepository)(nil).UpdateUser), arg0, arg1)
}

func (m *MockUserRepository) DeleteUser(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockUserRepositoryMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserRepository)(nil).DeleteUser), arg0, arg1)
}

func (m *MockUserRepository) UserExists(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserExists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockUserRepositoryMockRecorder) UserExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserExists", reflect.TypeOf((*MockUserRepository)(nil).UserExists), arg0, arg1)
}

type MockTeamRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTeamRepositoryMockRecorder
}

type MockTeamRepositoryMockRecorder struct {
	mock *MockTeamRepository
}

func NewMockTeamRepository(ctrl *gomock.Controller) *MockTeamRepository {
	mock := &MockTeamRepository{ctrl: ctrl}
	mock.recorder = &MockTeamRepositoryMockRecorder{mock}
	return mock
}

func (m *MockTeamRepository) EXPECT() *MockTeamRepositoryMockRecorder {
	return m.recorder
}

func (m *MockTeamRepository) CreateTeam(arg0 context.Context, arg1 *models.Team) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTeam", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockTeamRepositoryMockRecorder) CreateTeam(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTeam", reflect.TypeOf((*MockTeamRepository)(nil).CreateTeam), arg0, arg1)
}

func (m *MockTeamRepository) GetTeamByName(arg0 context.Context, arg1 string) (*models.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTeamByName", arg0, arg1)
	ret0, _ := ret[0].(*models.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockTeamRepositoryMockRecorder) GetTeamByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTeamByName", reflect.TypeOf((*MockTeamRepository)(nil).GetTeamByName), arg0, arg1)
}

func (m *MockTeamRepository) GetTeamWithMembers(arg0 context.Context, arg1 string) (*models.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTeamWithMembers", arg0, arg1)
	ret0, _ := ret[0].(*models.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockTeamRepositoryMockRecorder) GetTeamWithMembers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTeamWithMembers", reflect.TypeOf((*MockTeamRepository)(nil).GetTeamWithMembers), arg0, arg1)
}

func (m *MockTeamRepository) UpdateTeam(arg0 context.Context, arg1 *models.Team) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTeam", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockTeamRepositoryMockRecorder) UpdateTeam(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTeam", reflect.TypeOf((*MockTeamRepository)(nil).UpdateTeam), arg0, arg1)
}

func (m *MockTeamRepository) DeleteTeam(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTeam", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockTeamRepositoryMockRecorder) DeleteTeam(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTeam", reflect.TypeOf((*MockTeamRepository)(nil).DeleteTeam), arg0, arg1)
}

func (m *MockTeamRepository) TeamExists(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TeamExists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockTeamRepositoryMockRecorder) TeamExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TeamExists", reflect.TypeOf((*MockTeamRepository)(nil).TeamExists), arg0, arg1)
}

func (m *MockTeamRepository) GetAllTeams(arg0 context.Context) ([]*models.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTeams", arg0)
	ret0, _ := ret[0].([]*models.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockTeamRepositoryMockRecorder) GetAllTeams(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTeams", reflect.TypeOf((*MockTeamRepository)(nil).GetAllTeams), arg0)
}

type MockPullRequestRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPullRequestRepositoryMockRecorder
}

type MockPullRequestRepositoryMockRecorder struct {
	mock *MockPullRequestRepository
}

func NewMockPullRequestRepository(ctrl *gomock.Controller) *MockPullRequestRepository {
	mock := &MockPullRequestRepository{ctrl: ctrl}
	mock.recorder = &MockPullRequestRepositoryMockRecorder{mock}
	return mock
}

func (m *MockPullRequestRepository) EXPECT() *MockPullRequestRepositoryMockRecorder {
	return m.recorder
}

func (m *MockPullRequestRepository) CreatePullRequest(arg0 context.Context, arg1 *models.PullRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePullRequest", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockPullRequestRepositoryMockRecorder) CreatePullRequest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePullRequest", reflect.TypeOf((*MockPullRequestRepository)(nil).CreatePullRequest), arg0, arg1)
}

func (m *MockPullRequestRepository) GetPullRequestByID(arg0 context.Context, arg1 string) (*models.PullRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPullRequestByID", arg0, arg1)
	ret0, _ := ret[0].(*models.PullRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockPullRequestRepositoryMockRecorder) GetPullRequestByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPullRequestByID", reflect.TypeOf((*MockPullRequestRepository)(nil).GetPullRequestByID), arg0, arg1)
}

func (m *MockPullRequestRepository) GetPullRequestsByReviewer(arg0 context.Context, arg1 string) ([]*models.PullRequestShort, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPullRequestsByReviewer", arg0, arg1)
	ret0, _ := ret[0].([]*models.PullRequestShort)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockPullRequestRepositoryMockRecorder) GetPullRequestsByReviewer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPullRequestsByReviewer", reflect.TypeOf((*MockPullRequestRepository)(nil).GetPullRequestsByReviewer), arg0, arg1)
}

func (m *MockPullRequestRepository) UpdatePullRequest(arg0 context.Context, arg1 *models.PullRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePullRequest", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockPullRequestRepositoryMockRecorder) UpdatePullRequest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePullRequest", reflect.TypeOf((*MockPullRequestRepository)(nil).UpdatePullRequest), arg0, arg1)
}

func (m *MockPullRequestRepository) DeletePullRequest(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePullRequest", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockPullRequestRepositoryMockRecorder) DeletePullRequest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePullRequest", reflect.TypeOf((*MockPullRequestRepository)(nil).DeletePullRequest), arg0, arg1)
}

func (m *MockPullRequestRepository) MergePullRequest(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MergePullRequest", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockPullRequestRepositoryMockRecorder) MergePullRequest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MergePullRequest", reflect.TypeOf((*MockPullRequestRepository)(nil).MergePullRequest), arg0, arg1)
}

func (m *MockPullRequestRepository) PullRequestExists(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PullRequestExists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockPullRequestRepositoryMockRecorder) PullRequestExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PullRequestExists", reflect.TypeOf((*MockPullRequestRepository)(nil).PullRequestExists), arg0, arg1)
}

func (m *MockPullRequestRepository) GetAssignedReviewers(arg0 context.Context, arg1 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAssignedReviewers", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockPullRequestRepositoryMockRecorder) GetAssignedReviewers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAssignedReviewers", reflect.TypeOf((*MockPullRequestRepository)(nil).GetAssignedReviewers), arg0, arg1)
}

func (m *MockPullRequestRepository) SetAssignedReviewers(arg0 context.Context, arg1 string, arg2 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAssignedReviewers", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockPullRequestRepositoryMockRecorder) SetAssignedReviewers(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAssignedReviewers", reflect.TypeOf((*MockPullRequestRepository)(nil).SetAssignedReviewers), arg0, arg1, arg2)
}

func (m *MockPullRequestRepository) GetPRCountByStatus(arg0 context.Context) (map[string]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPRCountByStatus", arg0)
	ret0, _ := ret[0].(map[string]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockPullRequestRepositoryMockRecorder) GetPRCountByStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPRCountByStatus", reflect.TypeOf((*MockPullRequestRepository)(nil).GetPRCountByStatus), arg0)
}

func (m *MockPullRequestRepository) GetAssignmentsByUsers(arg0 context.Context) (map[string]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAssignmentsByUsers", arg0)
	ret0, _ := ret[0].(map[string]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockPullRequestRepositoryMockRecorder) GetAssignmentsByUsers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAssignmentsByUsers", reflect.TypeOf((*MockPullRequestRepository)(nil).GetAssignmentsByUsers), arg0)
}

func (m *MockPullRequestRepository) GetTeamPRCount(arg0 context.Context, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTeamPRCount", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockPullRequestRepositoryMockRecorder) GetTeamPRCount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTeamPRCount", reflect.TypeOf((*MockPullRequestRepository)(nil).GetTeamPRCount), arg0, arg1)
}
