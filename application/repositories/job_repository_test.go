package repositories_test

import (
	"encoder-project/application/repositories"
	"encoder-project/domain"
	"encoder-project/framework/database"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJobRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	job, err := domain.NewJob("output_path", "Pending", video)
	require.Nil(t, err)

	repo := repositories.JobRepositoryDb{Db: db}
	repo.Insert(job)

	j, err := repo.Find(job.ID)

	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.ID, job.ID)
	require.Equal(t, j.Video.ID, video.ID)
}

func TestJobRepositoryDbUpdate(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	job, err := domain.NewJob("output_path", "Pending", video)
	require.Nil(t, err)

	repo := repositories.JobRepositoryDb{Db: db}
	repo.Insert(job)

	job.Status = "Complete"

	repo.Update(job)

	j, err := repo.Find(job.ID)

	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.Status, job.Status)
}
