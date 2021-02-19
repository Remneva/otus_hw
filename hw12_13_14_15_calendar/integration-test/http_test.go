package main

import (
	h "/hw12_13_14_15_calendar/pkg/server/http"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestServerHTTPset(t *testing.T) {

	t.Run("Create, update, get, delete event", func(t *testing.T) {
		request := Event{
			ID:          1,
			Owner:       218,
			Title:       "Xipe-Totec",
			Description: "qwerty",
			StartDate:   "2020-03-01",
			StartTime:   "2018-08-28T12:30:00+05:30",
			EndDate:     "2020-03-02",
			EndTime:     "2018-08-28T12:30:00+05:30",
		}

		jsonBody, _ := json.Marshal(&request)
		req, err := http.NewRequest("POST", "http://calendar:8082/set",
			bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		assert.NotNil(t, req)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.NotNil(t, resp)
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		require.NotNil(t, body)

		id := int64(1)
		request = Event{
			ID:          1,
			Owner:       218,
			Title:       "Title",
			Description: "qwerty",
			StartDate:   "2020-03-01",
			StartTime:   "2020-08-28T12:30:00+08:30",
			EndDate:     "2020-03-02",
			EndTime:     "2021-08-28T12:30:00+08:30",
		}
		request.ID = id
		jsonBody, _ = json.Marshal(&request)
		req, err = http.NewRequest("POST", "http://calendar:8082/update",
			bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		assert.NotNil(t, req)
		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.NotNil(t, resp)
		body, err = ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		jsonid := h.JSONID{}
		jsonid.ID = id
		jsonBody, _ = json.Marshal(&id)

		req, _ = http.NewRequest("POST", "http://calendar:8082/get",
			bytes.NewBuffer(jsonBody))
		resp, _ = http.DefaultClient.Do(req)
		body, _ = ioutil.ReadAll(resp.Body)
		rb := Event{}
		json.Unmarshal(body, &rb)
		assert.Equal(t, 200, resp.StatusCode)
		assert.EqualValues(t, jsonid.ID, rb.ID)
		assert.EqualValues(t, 218, rb.Owner)
		assert.EqualValues(t, "Title", rb.Title)
		assert.EqualValues(t, "qwerty", rb.Description)
		assert.EqualValues(t, "2020-03-01T00:00:00Z", rb.StartDate)
		assert.EqualValues(t, "2020-03-02T00:00:00Z", rb.EndDate)

		req, err = http.NewRequest("POST", "http://calendar:8082/delete",
			bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		body, err = ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, resp.StatusCode, 200)
	})

}

type Event struct {
	ID          int64  `json:"ID"`
	Owner       int64  `json:"Owner"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	StartDate   string `json:"StartDate"`
	StartTime   string `json:"StartTime"`
	EndDate     string `json:"EndDate"`
	EndTime     string `json:"EndTime"`
}
