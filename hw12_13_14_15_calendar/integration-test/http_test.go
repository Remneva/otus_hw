package test

import (
	"bytes"
	"encoding/json"
	h "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/server/http"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestServerHTTPset(t *testing.T) {

	t.Run("Create, update, get, delete event", func(t *testing.T) {
		start := time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC)
		oneDayLater := start.AddDate(0, 0, 1)
		request := h.Event{
			ID:          1,
			Owner:       218,
			Title:       "Xipe-Totec",
			Description: "qwerty",
			StartDate:   "",
			StartTime:   start,
			EndDate:     "",
			EndTime:     oneDayLater,
		}

		jsonBody, _ := json.Marshal(&request)
		req, _ := http.NewRequest("POST", "http://calendar:8887/set",
			bytes.NewBuffer(jsonBody))
		resp, _ := http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(t, resp.StatusCode, 200)
		require.NotNil(t, body)

		id := int64(1)
		request = h.Event{}
		request.ID = id
		jsonBody, _ = json.Marshal(&request)
		req, _ = http.NewRequest("POST", "http://calendar:8887/update",
			bytes.NewBuffer(jsonBody))
		resp, _ = http.DefaultClient.Do(req)
		body, _ = ioutil.ReadAll(resp.Body)
		assert.Equal(t, resp.StatusCode, 200)
		require.NotNil(t, body)

		jsonid := h.JSONID{}
		jsonid.ID = id
		jsonBody, _ = json.Marshal(&id)

		req, _ = http.NewRequest("POST", "http://localhost:8887/get",
			bytes.NewBuffer(jsonBody))
		resp, _ = http.DefaultClient.Do(req)
		body, _ = ioutil.ReadAll(resp.Body)
		rb := h.Event{}
		json.Unmarshal(body, &rb)
		assert.Equal(t, resp.StatusCode, 200)
		assert.EqualValues(t, jsonid.ID, rb.ID)
		assert.EqualValues(t, 1, rb.Owner)
		assert.EqualValues(t, "Title", rb.Title)
		assert.EqualValues(t, "Description", rb.Description)
		assert.EqualValues(t, "2020-03-01", rb.StartDate)
		assert.EqualValues(t, "2020-03-01", rb.EndDate)

		req, _ = http.NewRequest("POST", "http://localhost:8887/delete",
			bytes.NewBuffer(jsonBody))
		resp, _ = http.DefaultClient.Do(req)
		body, _ = ioutil.ReadAll(resp.Body)
		assert.Equal(t, resp.StatusCode, 200)

		request = h.Event{
			ID:          2,
			Owner:       2188,
			Title:       "Xipe-Totec",
			Description: "qwerty",
			StartDate:   "",
			StartTime:   start,
			EndDate:     "",
			EndTime:     oneDayLater,
		}

		jsonBody, _ = json.Marshal(&request)
		req, _ = http.NewRequest("POST", "http://calendar:8887/set",
			bytes.NewBuffer(jsonBody))
		resp, _ = http.DefaultClient.Do(req)
		body, _ = ioutil.ReadAll(resp.Body)
		assert.Equal(t, resp.StatusCode, 200)
		require.NotNil(t, body)
	})

}
