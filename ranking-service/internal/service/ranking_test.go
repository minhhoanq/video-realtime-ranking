package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/minhhoanq/video-realtime-ranking/ranking-service/internal/dataaccess/database"
	mockrd "github.com/minhhoanq/video-realtime-ranking/ranking-service/internal/dataaccess/redis/mock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestGetTopKVideoRanking(t *testing.T) {
	offset, limit := 0, 10
	key := "video:global_score"
	mockRedisData := []redis.Z{
		{Member: "video1", Score: 12},
		{Member: "video2", Score: 8},
	}

	expected := []database.Score{
		{VideoID: "video1", Score: 12},
		{VideoID: "video2", Score: 8},
	}

	testCases := []struct {
		name string
		body struct {
			key    string
			offset int
			limit  int
		}
		buildStubs    func(store *mockrd.MockRankingDataAccessor)
		checkResponse func(t *testing.T, res []database.Score, err error)
	}{
		{
			name: "OK",
			body: struct {
				key    string
				offset int
				limit  int
			}{
				key:    key,
				offset: offset,
				limit:  limit,
			},
			buildStubs: func(rankingRedisDataAccessor *mockrd.MockRankingDataAccessor) {
				start := int64(offset)
				stop := int64(offset + limit - 1)

				rankingRedisDataAccessor.EXPECT().
					GetTopRanked(gomock.Any(), gomock.Eq(key), gomock.Eq(start), gomock.Eq(stop)).
					Times(1).
					Return(mockRedisData, nil)
			},
			checkResponse: func(t *testing.T, res []database.Score, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, expected, res)
				require.ElementsMatch(t, expected, res)
			},
		},
		{
			name: "Error",
			body: struct {
				key    string
				offset int
				limit  int
			}{
				key:    key,
				offset: offset,
				limit:  limit,
			},
			buildStubs: func(rankingRedisDataAccessor *mockrd.MockRankingDataAccessor) {
				start := int64(offset)
				stop := int64(offset + limit - 1)

				rankingRedisDataAccessor.EXPECT().
					GetTopRanked(gomock.Any(), gomock.Eq(key), gomock.Eq(start), gomock.Eq(stop)).
					Times(1).
					Return(nil, fmt.Errorf("redis: connection refused")) // giả lập lỗi
			},
			checkResponse: func(t *testing.T, res []database.Score, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "redis: connection refused")
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		rankingRedisDataAccessor := mockrd.NewMockRankingDataAccessor(mockCtrl)

		tc.buildStubs(rankingRedisDataAccessor)

		service := NewRankingService(rankingRedisDataAccessor)
		res, err := service.GetTopKVideoRanking(context.Background(), tc.body.offset, tc.body.limit)
		tc.checkResponse(t, res, err)
	}
}

func TestGetTopKUserVideoRanking(t *testing.T) {
	offset, limit := 0, 10
	user_id := "user_id"
	userVideokey := fmt.Sprintf("%s:%s", "video", user_id)
	mockRedisData := []redis.Z{
		{Member: "video1", Score: 12},
		{Member: "video2", Score: 8},
	}

	expected := []database.Score{
		{VideoID: "video1", Score: 12},
		{VideoID: "video2", Score: 8},
	}

	testCases := []struct {
		name string
		body struct {
			key    string
			offset int
			limit  int
		}
		buildStubs    func(store *mockrd.MockRankingDataAccessor)
		checkResponse func(t *testing.T, res []database.Score, err error)
	}{
		{
			name: "OK",
			body: struct {
				key    string
				offset int
				limit  int
			}{
				key:    userVideokey,
				offset: offset,
				limit:  limit,
			},
			buildStubs: func(rankingRedisDataAccessor *mockrd.MockRankingDataAccessor) {
				start := int64(offset)
				stop := int64(offset + limit - 1)

				rankingRedisDataAccessor.EXPECT().
					GetTopRanked(gomock.Any(), gomock.Eq(userVideokey), gomock.Eq(start), gomock.Eq(stop)).
					Times(1).
					Return(mockRedisData, nil)
			},
			checkResponse: func(t *testing.T, res []database.Score, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, expected, res)
				require.ElementsMatch(t, expected, res)
			},
		},
		{
			name: "Error",
			body: struct {
				key    string
				offset int
				limit  int
			}{
				key:    userVideokey,
				offset: offset,
				limit:  limit,
			},
			buildStubs: func(rankingRedisDataAccessor *mockrd.MockRankingDataAccessor) {
				start := int64(offset)
				stop := int64(offset + limit - 1)

				rankingRedisDataAccessor.EXPECT().
					GetTopRanked(gomock.Any(), gomock.Eq(userVideokey), gomock.Eq(start), gomock.Eq(stop)).
					Times(1).
					Return(nil, fmt.Errorf("redis: connection refused")) // giả lập lỗi
			},
			checkResponse: func(t *testing.T, res []database.Score, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "redis: connection refused")
				require.Nil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		rankingRedisDataAccessor := mockrd.NewMockRankingDataAccessor(mockCtrl)

		tc.buildStubs(rankingRedisDataAccessor)

		service := NewRankingService(rankingRedisDataAccessor)
		res, err := service.GetTopKUserVideoRanking(context.Background(), user_id, tc.body.offset, tc.body.limit)
		tc.checkResponse(t, res, err)
	}
}
