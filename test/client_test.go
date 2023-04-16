package test

import "testing"

type ClientTest struct {
	Stage *ClientStage
}

func NewClientTest(t *testing.T) *ClientTest {
	return &ClientTest{
		Stage: NewClientStage(t),
	}
}

func Test_Client_UploadFile_Succeeds(t *testing.T) {
	s := NewClientTest(t).Stage

	s.When().
		ClientUploadsFile()

	s.Then().
		ImportSucceeds()
}

func Test_Client_UploadFile_IsSkipped(t *testing.T) {
	s := NewClientTest(t).Stage

	s.Given().
		UploadFileWillBeSkipped()

	s.When().
		ClientUploadsFile()

	s.Then().
		ImportIsSkipped()
}
