// +build integration

package test

import (
	"context"
	"fmt"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/server/pb"
	"github.com/golang/protobuf/ptypes"
	"github.com/tj/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
	"time"
)

func TestServerGRPC(t *testing.T) {
	start := time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC)
	starttime, err := ptypes.TimestampProto(start)
	oneDayLater := start.AddDate(0, 0, 1)
	endtime, err := ptypes.TimestampProto(oneDayLater)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("localhost:8887", opts...)

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
		Enddate:     "2020-03-01",
		Endtime:     endtime,
	}

	t.Run("Create, update, get, delete event", func(t *testing.T) {
		respId, err := client.SetEvent(context.Background(), request)
		if err != nil {
			fmt.Printf("fail to dial: %v\n", err)
		}
		request.Id = respId.Id
		respId, err = client.UpdateEvent(context.Background(), request)
		if err != nil {
			fmt.Printf("fail to dial: %v\n", err)
		}
		id := &pb.Id{
			Id: request.Id,
		}
		respEvent, err := client.GetEvent(context.Background(), id)
		if err != nil {
			fmt.Printf("fail to dial: %v\n", err)
		}
		assert.Equal(t, int64(1), respEvent.Owner)
		assert.Equal(t, "Title", respEvent.Title)
		assert.Equal(t, "result", respEvent.Description)
		assert.Equal(t, "2020-03-01T00:00:00Z", respEvent.Startdate)
		assert.Equal(t, "2020-03-01T00:00:00Z", respEvent.Enddate)

		respEmpty, err := client.DeleteEvent(context.Background(), id)
		if err != nil {
			fmt.Printf("fail")
		}
		assert.Equal(t, emptypb.Empty{}, respEmpty)

		request.Id = 0
		respId, err = client.UpdateEvent(context.Background(), request)
		e, _ := status.FromError(err)
		fmt.Println(e.Code().String())
		assert.Equal(t, codes.InvalidArgument, e.Code())

	})

}
