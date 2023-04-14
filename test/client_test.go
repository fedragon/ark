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

	s.Given().
		FileDoesNotExist().And().
		UploadFileWillSucceed()

	s.When().
		ClientUploadsFile()

	s.Then().
		ImportSucceeds()
}

func Test_Client_UploadFile_Fails(t *testing.T) {
	s := NewClientTest(t).Stage

	s.Given().
		FileDoesNotExist().And().
		UploadFileWillFail()

	s.When().
		ClientUploadsFile()

	s.Then().
		ImportFails()
}
