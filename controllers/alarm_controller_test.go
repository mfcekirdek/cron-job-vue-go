package controllers

import (
	"github.com/magiconair/properties/assert"
	"net/http"
	"testing"
)

func TestAlarmController(t *testing.T) {
	bc := NewBaseController()
	userService := &userSvc{}
	jobService := &jobSvc{}
	awsClient := &awsClient{}
	telegramClient := &telegramClient{}

	alarmCtrl := NewAlarmController(bc, userService, jobService, awsClient, telegramClient, nil)

	t.Run("Is Get Not Allowed", func(t *testing.T) {
		w, req := createHttpReq(http.MethodGet, "/api/create-alarm", nil)
		alarmCtrl.CreateAlarm(w, req)

		res := parseBody(w)

		assert.Equal(t, res.Code, http.StatusNotFound)
		assert.Equal(t, res.Message, ErrMethodNotAllowed.Error())
	})
	t.Run("Getting Token error", func(t *testing.T) {
		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", nil)
		alarmCtrl.CreateAlarm(w, req)

		res := parseBody(w)

		assert.Equal(t, res.Code, http.StatusBadRequest)
		assert.Equal(t, res.Message, ErrTokenNotFound.Error())
	})
	t.Run("Non empty token but getting empty name error", func(t *testing.T) {
		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", nil)
		req.Form.Set("token", "token")

		alarmCtrl.CreateAlarm(w, req)

		res := parseBody(w)

		assert.Equal(t, res.Code, http.StatusBadRequest)
		assert.Equal(t, res.Message, ErrNameNotFound.Error())
	})
	t.Run("Non empty {token, name} but getting empty time error", func(t *testing.T) {
		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", nil)
		req.Form.Set("token", "token")
		req.Form.Set("name", "name")

		alarmCtrl.CreateAlarm(w, req)

		res := parseBody(w)

		assert.Equal(t, res.Code, http.StatusBadRequest)
		assert.Equal(t, res.Message, ErrTimeNotFound.Error())
	})

	t.Run("Non empty {token, name, time} but getting empty repeatType error", func(t *testing.T) {
		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", nil)
		req.Form.Set("token", "token")
		req.Form.Set("name", "name")
		req.Form.Set("time", "23:14")

		alarmCtrl.CreateAlarm(w, req)

		res := parseBody(w)

		assert.Equal(t, res.Code, http.StatusBadRequest)
		assert.Equal(t, res.Message, ErrRepeatTypeNotFound.Error())
	})
	t.Run("Non empty {token,name,time,repeatType} but getting reading image file err", func(t *testing.T) {
		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", nil)
		req.Form.Set("token", "token")
		req.Form.Set("name", "name")
		req.Form.Set("time", "23:14")
		req.Form.Set("repeatType", "5")

		alarmCtrl.CreateAlarm(w, req)

		res := parseBody(w)

		assert.Equal(t, res.Code, http.StatusBadRequest)
		assert.Equal(t, res.Message, ErrReadingFile.Error())
	})
	t.Run("Getting token err occured in db", func(t *testing.T) {
		body, contentType := fileUploadRequest()

		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", body)
		req.Form.Set("token", "db-err")
		req.Form.Set("name", "name")
		req.Form.Set("time", "23:14")
		req.Form.Set("repeatType", "5")
		req.Form.Set("fileName", "test")
		req.Form.Set("fileType", "image/png")

		req.Header.Add("Content-Type", contentType)

		alarmCtrl.CreateAlarm(w, req)

		res := parseBody(w)

		assert.Equal(t, res.Code, http.StatusBadRequest)
		assert.Equal(t, res.Message, ErrDb.Error())
	})

	t.Run("when job exist db error", func(t *testing.T) {
		body, contentType := fileUploadRequest()

		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", body)
		req.Form.Set("token", "job-already-exist-db-error")
		req.Form.Set("name", "name")
		req.Form.Set("time", "23:14")
		req.Form.Set("repeatType", "5")
		req.Form.Set("fileName", "arbitrary-name")
		req.Form.Set("fileType", "image/png")

		req.Header.Add("Content-Type", contentType)

		alarmCtrl.CreateAlarm(w, req)

		res := parseBody(w)

		assert.Equal(t, res.Code, http.StatusBadRequest)
		assert.Equal(t, res.Message, ErrGettingJob.Error())
	})
	t.Run("when job already exist error", func(t *testing.T) {
		body, contentType := fileUploadRequest()

		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", body)
		req.Form.Set("token", "job-already-exist")
		req.Form.Set("name", "name")
		req.Form.Set("time", "23:14")
		req.Form.Set("repeatType", "5")
		req.Form.Set("fileName", "arbitrary-name")
		req.Form.Set("fileType", "image/png")

		req.Header.Add("Content-Type", contentType)

		alarmCtrl.CreateAlarm(w, req)

		res := parseBody(w)

		assert.Equal(t, res.Code, http.StatusBadRequest)
		assert.Equal(t, res.Message, ErrJobAlreadyExist.Error())
	})

	t.Run("Getting non exist token error", func(t *testing.T) {
		body, contentType := fileUploadRequest()

		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", body)
		req.Form.Set("token", "not-exist-token")
		req.Form.Set("name", "name")
		req.Form.Set("time", "23:14")
		req.Form.Set("repeatType", "5")
		req.Form.Set("fileName", "test")
		req.Form.Set("fileType", "image/png")

		req.Header.Add("Content-Type", contentType)

		alarmCtrl.CreateAlarm(w, req)

		res := parseBody(w)

		assert.Equal(t, res.Code, http.StatusBadRequest)
		assert.Equal(t, res.Message, ErrTokenDoesNotExist.Error())
	})
	t.Run("Getting s3 upload error", func(t *testing.T) {
		body, contentType := fileUploadRequest()

		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", body)
		req.Form.Set("token", "sametintokeni")
		req.Form.Set("name", "name")
		req.Form.Set("time", "23:14")
		req.Form.Set("repeatType", "5")
		req.Form.Set("fileName", "badFileName")
		req.Form.Set("fileType", "image/png")

		req.Header.Add("Content-Type", contentType)

		alarmCtrl.CreateAlarm(w, req)

		res := parseBody(w)

		assert.Equal(t, res.Code, http.StatusBadRequest)
		assert.Equal(t, res.Message, ErrS3Upload.Error())
	})

	t.Run("When job created, job err occured, delete uploaded file in s3 also occured", func(t *testing.T) {
		body, contentType := fileUploadRequest()

		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", body)
		req.Form.Set("token", "sametintokeni")
		req.Form.Set("name", "name")
		req.Form.Set("time", "23:14")
		req.Form.Set("repeatType", "5")
		req.Form.Set("fileName", "error-scenario-with-s3")
		req.Form.Set("fileType", "image/png")

		req.Header.Add("Content-Type", contentType)

		alarmCtrl.CreateAlarm(w, req)

		res := parseBody(w)

		assert.Equal(t, res.Code, http.StatusBadRequest)
		assert.Equal(t, res.Message, ErrDeleteFileS3.Error())
	})
	t.Run("When job created, job err occured, delete uploaded file in s3 is success return add job error", func(t *testing.T) {
		body, contentType := fileUploadRequest()

		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", body)
		req.Form.Set("token", "sametintokeni")
		req.Form.Set("name", "name")
		req.Form.Set("time", "23:14")
		req.Form.Set("repeatType", "5")
		req.Form.Set("fileName", "error-scenario-job")
		req.Form.Set("fileType", "image/png")

		req.Header.Add("Content-Type", contentType)

		alarmCtrl.CreateAlarm(w, req)

		res := parseBody(w)

		assert.Equal(t, res.Code, http.StatusBadRequest)
		assert.Equal(t, res.Message, ErrAddingJob.Error())
	})
}
