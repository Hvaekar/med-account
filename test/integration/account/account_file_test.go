package account

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/test/integration/account/fixtures"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type FileTestSuite struct {
	TestSuite
}

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(FileTestSuite))
}

func (s *FileTestSuite) TestAddFile() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	file, err := os.Open("./fixtures/testfiles/test-file.png")
	s.Require().NoError(err)
	defer file.Close()

	s.uploaderAPIMock.On("UploadWithContext", mock.Anything, mock.Anything).
		Return(nil, nil).
		Once()

	f, err := s.client.AddFile(s.ctx, s.token.Access, file, filepath.Base(file.Name()))
	s.Require().NoError(err)

	s.NotEmpty(f.ID)
	s.NotEmpty(f.Name)

	//err = s.s3.DeleteObject(s.ctx, f.Name)
	//s.Require().NoError(err)
}

func (s *FileTestSuite) TestGetFiles() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	list, err := s.client.GetFiles(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Files, 2)
}

func (s *FileTestSuite) TestUpdateFile() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	description := "Some file description"
	req := model.UpdateFile{
		Description: &description,
	}

	file, err := s.client.UpdateFile(s.ctx, s.token.Access, 1, &req)
	s.Require().NoError(err)

	s.Equal(int64(1), file.ID)
	s.Equal(req.Description, file.Description)
}

func (s *FileTestSuite) TestUpdateFileNotFound() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	description := "Some file description"
	req := model.UpdateFile{
		Description: &description,
	}

	_, err := s.client.UpdateFile(s.ctx, s.token.Access, 100, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *FileTestSuite) TestDeleteFile() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	err := s.client.DeleteFile(s.ctx, s.token.Access, 1)
	s.Require().NoError(err)

	list, err := s.client.GetFiles(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Files, 1)
}
