package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	sqlstorage "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/require"
	"io/ioutil"
	http "net/http"
	"testing"
	"time"
)

type eventMatcher struct {
	sqlstorage.Event
}

func TestServer(t *testing.T) {
	start := time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC)
	oneDayLater := start.AddDate(0, 0, 1)
	var id int64
	request := Event{
		ID:          111,
		Title:       "test title",
		Description: "test test test",
		StartDate:   "2020-03-01",
		StartTime:   start,
		EndDate:     "2020-03-01",
		EndTime:     oneDayLater,
	}
	t.Run("Create event", func(t *testing.T) {

		mockCtl := gomock.NewController(t)
		defer mockCtl.Finish()

		mockDB := NewMockEventsStorage(mockCtl)
		//		store := memorystorage.New(mockDB)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		event := sqlstorage.Event{
			ID:          111,
			Title:       "test title",
			Description: "test test test",
			StartDate:   "2020-03-01",
			StartTime:   start,
			EndDate:     "2020-03-01",
			EndTime:     oneDayLater,
		}

		mockDB.EXPECT().AddEvent(ctx, eventMatcher{event}).Return(event.ID, nil)

		jsonBody, _ := json.Marshal(&request)
		req, _ := http.NewRequest("POST", "http://localhost:8082/set",
			bytes.NewBuffer(jsonBody))
		resp, _ := http.DefaultClient.Do(req)
		respbody, _ := ioutil.ReadAll(resp.Body)

		err := json.Unmarshal(respbody, &id)
		fmt.Println("id: ", id)
		if err != nil {
			require.NoError(t, err)
		}
		require.Equal(t, resp.StatusCode, 200)
		require.Equal(t, event.ID, id)
	})
}

//id := JSONID{
//	ID: result,
//}
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
