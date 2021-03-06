package main

import (
	"context"
	"fmt"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/server/pb"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServerGRPC(t *testing.T) {
	start := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	starttime, err := ptypes.TimestampProto(start)
	oneDayLater := start.AddDate(0, 0, 1)
	endtime, err := ptypes.TimestampProto(oneDayLater)
	conn, err := grpc.Dial("calendar:50051", grpc.WithInsecure())
	ctx := context.Background()
	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()

	client := pb.NewCalendarClient(conn)
	request := &pb.Event{
		Owner:       1,
		Title:       "Title",
		Description: "result",
		Startdate:   "2020-03-01",
		Starttime:   starttime,
		Enddate:     "2020-03-02",
		Endtime:     endtime,
	}

	t.Run("Create, update, get, delete event", func(t *testing.T) {
		respId, err := client.SetEvent(context.Background(), request)
		if err != nil {
			fmt.Printf("fail to dial: %v\n", err)
		}
		require.NoError(t, err)
		assert.NotNil(t, respId)

		request.Id = respId.Id
		id := &pb.Id{
			Id: respId.Id,
		}

		respEvent, err := client.GetEvent(ctx, id)
		if err != nil {
			fmt.Printf("fail to dial: %v\n", err)
		}
		require.NoError(t, err)
		assert.Equal(t, "Title", respEvent.Title)
		assert.Equal(t, "result", respEvent.Description)
		assert.Equal(t, "2020-03-01T00:00:00Z", respEvent.Startdate)
		assert.Equal(t, "2020-03-02T00:00:00Z", respEvent.Enddate)

		respEmpty, err := client.DeleteEvent(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "", respEmpty.String())

		request.Id = 0
		respId, err = client.UpdateEvent(context.Background(), request)
		e, ok := status.FromError(err)
		require.True(t, ok)
		fmt.Println(e.Code().String())
		assert.Equal(t, codes.InvalidArgument, e.Code())

	})

}
