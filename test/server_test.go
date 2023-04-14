package test

import "testing"

type ServerTest struct {
	Stage *ServerStage
}

func NewServerTest(t *testing.T) *ServerTest {
	return &ServerTest{
		Stage: NewServerStage(t),
	}
}

func Test_Server_UploadFile_Succeeds(t *testing.T) {
	s := NewServerTest(t).Stage

	s.Given().
		ServerIsUpAndRunning().And().
		FileDoesNotExist()

	s.When().
		ClientUploadsFile()

	s.Then().
		FileExists()
}

func Test_Server_UploadFile_DiscardsDuplicate(t *testing.T) {
	s := NewServerTest(t).Stage

	s.Given().
		ServerIsUpAndRunning().And().
		FileExists()

	s.When().
		ClientUploadsFile()

	s.Then().
		FileExists()
}
