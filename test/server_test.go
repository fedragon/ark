package test

import (
	"testing"

	_ "github.com/fedragon/ark/testing"
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
		ClientUploadsFile("./test/testdata/a/image.heic")

	s.Then().
		UploadSucceeds()
}

func Test_Server_UploadHeicFile_DiscardsDuplicate(t *testing.T) {
	s := NewServerTest(t).Stage

	s.Given().
		ClientUploadsFile("./test/testdata/a/image.heic").And().
		FileExists()

	s.When().
		ClientUploadsFileAgain("./test/testdata/a/image.heic")

	s.Then().
		UploadIsSkipped()
}

func Test_Server_UploadJpgFile_Succeeds(t *testing.T) {
	s := NewServerTest(t).Stage

	s.Given().
		FileDoesNotExist()

	s.When().
		ClientUploadsFile("./test/testdata/a/image.jpg")

	s.Then().
		UploadSucceeds()
}

func Test_Server_UploadJpgFile_DiscardsDuplicate(t *testing.T) {
	s := NewServerTest(t).Stage

	s.Given().
		ClientUploadsFile("./test/testdata/a/image.jpg").And().
		FileExists()

	s.When().
		ClientUploadsFileAgain("./test/testdata/a/image.jpg")

	s.Then().
		UploadIsSkipped()
}

func Test_Server_UploadOrfFile_Succeeds(t *testing.T) {
	s := NewServerTest(t).Stage

	s.Given().
		FileDoesNotExist()

	s.When().
		ClientUploadsFile("./test/testdata/a/image.orf")

	s.Then().
		UploadSucceeds()
}

func Test_Server_UploadOrfFile_DiscardsDuplicate(t *testing.T) {
	s := NewServerTest(t).Stage

	s.Given().
		ClientUploadsFile("./test/testdata/a/image.orf").And().
		FileExists()

	s.When().
		ClientUploadsFileAgain("./test/testdata/a/image.orf")

	s.Then().
		UploadIsSkipped()
}
