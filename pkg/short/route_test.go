package short

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func SetUpTestingEnv(t *testing.T) (*gin.Engine, *gorm.DB) {
	r := gin.Default()
	// Load testing env from local files
	godotenv.Load("../../.env.testing")

	dsn := fmt.Sprintf("host=%v port=%v dbname=%v password=%v user=%v sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"), os.Getenv("DB_PASS"), os.Getenv("DB_USER"))
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	Route(db, r.Group("/"))

	return r, db
}

func Test_Short_CanCreateShort(t *testing.T) {
	r, _ := SetUpTestingEnv(t)

	demoShort := ShortRequest{
		URL: "https://google.com",
	}
	jsonValue, _ := json.Marshal(demoShort)

	req, _ := http.NewRequest("POST", "/s/", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func Test_Short_InvalidRequest(t *testing.T) {
	r, _ := SetUpTestingEnv(t)

	badRequest := `{"url": 123}`
	jsonValue := io.NopCloser(strings.NewReader(badRequest))

	req, _ := http.NewRequest("POST", "/s/", jsonValue)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_Short_Redirect(t *testing.T) {
	r, db := SetUpTestingEnv(t)

	demoShort := Short{URL: "https://google.com"}
	db.Create(&demoShort)

	var preReqCount int64
	db.Model(&Audit{}).Count(&preReqCount)

	blankReq := ``
	blankRequest := io.NopCloser(strings.NewReader(blankReq))

	req, _ := http.NewRequest("GET", fmt.Sprintf("/s/%s", demoShort.ID), blankRequest)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusFound, w.Code)

	// Test that access was logged
	var accessCount int64
	db.Model(&AccessLog{}).Where("short_id = ?", demoShort.ID).Count(&accessCount)
	assert.Equal(t, int64(1), accessCount)

	// Test that audit was logged
	var postReqCount int64
	db.Model(&Audit{}).Count(&postReqCount)
	assert.Equal(t, int64(1), postReqCount-preReqCount)
	fmt.Println(preReqCount, postReqCount)
}

func Test_Short_NotFound(t *testing.T) {
	r, _ := SetUpTestingEnv(t)

	blankReq := ``
	blankRequest := io.NopCloser(strings.NewReader(blankReq))

	req, _ := http.NewRequest("GET", "/s/1231231", blankRequest)
	assert.NotNil(t, req)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func Test_Short_GetDetails(t *testing.T) {
	r, db := SetUpTestingEnv(t)

	demoShort := Short{URL: "https://google.com"}
	db.Create(&demoShort)

	withinLast24 := AccessLog{
		ShortID:    demoShort.ID,
		AccessTime: time.Now().UTC().Add(-5 * time.Hour),
	}
	withinWeek := AccessLog{
		ShortID:    demoShort.ID,
		AccessTime: time.Now().UTC().AddDate(0, 0, -3),
	}
	withinYear := AccessLog{
		ShortID:    demoShort.ID,
		AccessTime: time.Now().UTC().AddDate(0, -2, 0),
	}
	db.Create(&withinLast24)
	db.Create(&withinWeek)
	db.Create(&withinYear)

	expectedResponse := ShortResponse{
		Short:       demoShort,
		Last24Hours: 1,
		PastWeek:    2,
		AllTime:     3,
	}

	blankReq := ``
	blankRequest := io.NopCloser(strings.NewReader(blankReq))
	req, _ := http.NewRequest("GET", fmt.Sprintf("/s/%s/details", demoShort.ID), blankRequest)
	assert.NotNil(t, req)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var actualResponse ShortResponse
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, actualResponse)
}

func Test_Short_GetDetailsNotFound(t *testing.T) {
	r, _ := SetUpTestingEnv(t)

	blankReq := ``
	blankRequest := io.NopCloser(strings.NewReader(blankReq))

	req, _ := http.NewRequest("GET", "/s/1231231/details", blankRequest)
	assert.NotNil(t, req)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func Test_Short_Delete(t *testing.T) {
	r, db := SetUpTestingEnv(t)

	demoShort := Short{
		URL: "https://google.com",
	}
	db.Create(&demoShort)

	blankReq := ``
	blankRequest := io.NopCloser(strings.NewReader(blankReq))

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/s/%s", demoShort.ID), blankRequest)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
