package main

type MockStore struct{}

func (m *MockStore) CreateUser(u *User) (*User, error) {
	return &User{}, nil
}

func (m *MockStore) GetUserByID(id string) (*User, error) {
	return &User{}, nil
}

func (m *MockStore) CreateTask(task *Task) (*Task, error) {
	return &Task{}, nil
}

func (m *MockStore) GetTask(id string) (*Task, error) {
	return &Task{}, nil
}
