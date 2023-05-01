package test

import (
	"testing"
)

type ServerTest struct {
	Stage *ServerStage
}

func NewServerTest(t *testing.T) *ServerTest {
	return &ServerTest{
		Stage: NewServerStage(t),
	}
}

func Test_Server_UploadHeicFile_Succeeds(t *testing.T) {
	s := NewServerTest(t).Stage

	s.Given().
		FileDoesNotExist()

	s.When().
		ClientUploadsFile("./test/data/a/image.heic")

	s.Then().
		UploadSucceeds()
}

func Test_Server_UploadHeicFile_DiscardsDuplicate(t *testing.T) {
	s := NewServerTest(t).Stage

	s.Given().
		ClientUploadsFile("./test/data/a/image.heic").And().
		FileExists()

	s.When().
		ClientUploadsFileAgain("./test/data/a/image.heic")

	s.Then().
		UploadIsSkipped()
}

func Test_Server_UploadJpgFile_Succeeds(t *testing.T) {
	s := NewServerTest(t).Stage

	s.Given().
		FileDoesNotExist()

	s.When().
		ClientUploadsFile("./test/data/a/image.jpg")

	s.Then().
		UploadSucceeds()
}

func Test_Server_UploadJpgFile_DiscardsDuplicate(t *testing.T) {
	s := NewServerTest(t).Stage

	s.Given().
		ClientUploadsFile("./test/data/a/image.jpg").And().
		FileExists()

	s.When().
		ClientUploadsFileAgain("./test/data/a/image.jpg")

	s.Then().
		UploadIsSkipped()
}

func Test_Server_UploadOrfFile_Succeeds(t *testing.T) {
	s := NewServerTest(t).Stage

	s.Given().
		FileDoesNotExist()

	s.When().
		ClientUploadsFile("./test/data/a/image.orf")

	s.Then().
		UploadSucceeds()
}

func Test_Server_UploadOrfFile_DiscardsDuplicate(t *testing.T) {
	s := NewServerTest(t).Stage

	s.Given().
		ClientUploadsFile("./test/data/a/image.orf").And().
		FileExists()

	s.When().
		ClientUploadsFileAgain("./test/data/a/image.orf")

	s.Then().
		UploadIsSkipped()
}
