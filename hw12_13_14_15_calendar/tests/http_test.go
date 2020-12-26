package test

//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/app"
//	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/logger"
//	s "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage"
//	"github.com/golang/mock/gomock"
//	"github.com/pkg/errors"
//	"github.com/stretchr/testify/suite"
//	"go.uber.org/zap/zapcore"
//	"io/ioutil"
//
//	_ "github.com/stretchr/testify/require"
//
//	"net/http/httptest"
//	"testing"
//	"time"
//)

//request.ID = result
//jsonBody, _ = json.Marshal(&request)
//req, _ = http.NewRequest("POST", "http://localhost:8082/update",
//	bytes.NewBuffer(jsonBody))
//resp, _ = http.DefaultClient.Do(req)
//body, _ := ioutil.ReadAll(resp.Body)
//assert.Equal(t, resp.StatusCode, 200)
//require.NotNil(t, body)
//
//id.ID = result
//jsonBody, _ = json.Marshal(&id)
//
//req, _ = http.NewRequest("POST", "http://localhost:8082/get",
//	bytes.NewBuffer(jsonBody))
//resp, _ = http.DefaultClient.Do(req)
//body, _ = ioutil.ReadAll(resp.Body)
//rb := Event{}
//json.Unmarshal(body, &rb)
//assert.Equal(t, resp.StatusCode, 200)
//assert.EqualValues(t, result, rb.ID)
//assert.EqualValues(t, 1, rb.Owner)
//assert.EqualValues(t, "Title", rb.Title)
//assert.EqualValues(t, "Description", rb.Description)
//assert.EqualValues(t, "2020-03-01", rb.StartDate)
//assert.EqualValues(t, "2020-03-01", rb.EndDate)
//
//req, _ = http.NewRequest("POST", "http://localhost:8082/delete",
//	bytes.NewBuffer(jsonBody))
//resp, _ = http.DefaultClient.Do(req)
//body, _ = ioutil.ReadAll(resp.Body)
//assert.Equal(t, resp.StatusCode, 200)
