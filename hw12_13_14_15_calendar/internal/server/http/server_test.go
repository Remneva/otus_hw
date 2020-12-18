package internalhttp

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	http "net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	var result int64
	start := time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC)
	oneDayLater := start.AddDate(0, 0, 1)
	id := JSONID{
		ID: result,
	}
	request := Event{
		ID:          result,
		Owner:       1,
		Title:       "Title",
		Description: "result",
		StartDate:   "2020-03-01",
		StartTime:   start,
		EndDate:     "2020-03-01",
		EndTime:     oneDayLater,
	}
	t.Run("Create, update, get, delete event", func(t *testing.T) {
		jsonBody, _ := json.Marshal(&request)
		req, _ := http.NewRequest("POST", "http://localhost:8082/set",
			bytes.NewBuffer(jsonBody))
		resp, _ := http.DefaultClient.Do(req)
		respbody, _ := ioutil.ReadAll(resp.Body)
		r := JSONID{}
		result = r.ID
		json.Unmarshal(respbody, &result)
		assert.Equal(t, resp.StatusCode, 200)

		request.ID = result
		jsonBody, _ = json.Marshal(&request)
		req, _ = http.NewRequest("POST", "http://localhost:8082/update",
			bytes.NewBuffer(jsonBody))
		resp, _ = http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(t, resp.StatusCode, 200)
		require.NotNil(t, body)

		id.ID = result
		jsonBody, _ = json.Marshal(&id)

		req, _ = http.NewRequest("POST", "http://localhost:8082/get",
			bytes.NewBuffer(jsonBody))
		resp, _ = http.DefaultClient.Do(req)
		body, _ = ioutil.ReadAll(resp.Body)
		rb := Event{}
		json.Unmarshal(body, &rb)
		assert.Equal(t, resp.StatusCode, 200)
		assert.EqualValues(t, result, rb.ID)
		assert.EqualValues(t, 1, rb.Owner)
		assert.EqualValues(t, "Title", rb.Title)
		assert.EqualValues(t, "Description", rb.Description)
		assert.EqualValues(t, "2020-03-01", rb.StartDate)
		assert.EqualValues(t, "2020-03-01", rb.EndDate)

		req, _ = http.NewRequest("POST", "http://localhost:8082/delete",
			bytes.NewBuffer(jsonBody))
		resp, _ = http.DefaultClient.Do(req)
		body, _ = ioutil.ReadAll(resp.Body)
		assert.Equal(t, resp.StatusCode, 200)
	})

}
