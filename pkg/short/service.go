package short

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

func BuildMetrics(conn *gorm.DB, entity Short) (*ShortResponse, error) {
	timeStamp24Hours := time.Now().UTC().Add(-24 * time.Hour)
	last24Hours, err := GetAccessed(conn, entity.ID, &timeStamp24Hours)
	if err != nil {
		return nil, err
	}
	timeStampLastWeek := time.Now().UTC().AddDate(0, 0, -7)
	lastWeek, err := GetAccessed(conn, entity.ID, &timeStampLastWeek)
	if err != nil {
		return nil, err
	}
	allTime, err := GetAccessed(conn, entity.ID, nil)
	if err != nil {
		return nil, err
	}

	res := &ShortResponse{
		Short:       entity,
		Last24Hours: last24Hours,
		PastWeek:    lastWeek,
		AllTime:     allTime,
	}

	return res, nil
}

func GetAccessed(conn *gorm.DB, id string, timestamp *time.Time) (int64, error) {
	var count int64

	query := conn.Model(&AccessLog{}).Where("short_id = ?", id)
	if timestamp != nil {
		query.Where("access_time > ?", *timestamp)
	}

	res := query.Count(&count)
	if res.Error != nil {
		return 0, res.Error
	}

	return count, nil
}

func JsonError(e error) gin.H {
	return gin.H{
		"error": e.Error(),
	}
}
