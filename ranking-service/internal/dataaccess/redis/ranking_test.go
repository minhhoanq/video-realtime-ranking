package redis

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestGetScore(t *testing.T) {
	key := GlobalScore
	member := "video1"
	expectedScore := 12.0

	testCases := []struct {
		name string
		body struct {
			key           string
			member        string
			expectedScore float64
		}
		buildStubs   func(mockStore redismock.ClientMock)
		checkReponse func(t *testing.T, score float64, err error)
	}{
		{
			name: "OK",
			body: struct {
				key           string
				member        string
				expectedScore float64
			}{
				key:           key,
				member:        member,
				expectedScore: expectedScore,
			},
			buildStubs: func(mockStore redismock.ClientMock) {
				mockStore.ExpectZScore(key, member).SetVal(expectedScore)
			},
			checkReponse: func(t *testing.T, score float64, err error) {
				require.NoError(t, err)
				require.Equal(t, expectedScore, score)
			},
		},
		{
			name: "NotFound",
			body: struct {
				key           string
				member        string
				expectedScore float64
			}{
				key:           key,
				member:        "non_existing",
				expectedScore: 0.0,
			},
			buildStubs: func(mockStore redismock.ClientMock) {
				mockStore.ExpectZScore(key, member).RedisNil()
			},
			checkReponse: func(t *testing.T, score float64, err error) {
				require.Error(t, err)
				require.ErrorIs(t, err, redis.Nil)
				require.Equal(t, 0.0, score)
			},
		},
		{
			name: "Error",
			body: struct {
				key           string
				member        string
				expectedScore float64
			}{
				key:           key,
				member:        "non_existing",
				expectedScore: 0.0,
			},
			buildStubs: func(mockStore redismock.ClientMock) {
				err := errors.New("connection refused")
				mockStore.ExpectZScore(key, member).SetErr(err)
			},
			checkReponse: func(t *testing.T, score float64, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "connection refused")
				require.Equal(t, 0.0, score)
			},
		},
	}

	for _, tc := range testCases {
		db, mock := redismock.NewClientMock()
		defer db.Close()

		accessor := NewRankingDataAccessor(db)
		ctx := context.Background()
		tc.buildStubs(mock)

		score, err := accessor.GetScore(ctx, key, member)
		tc.checkReponse(t, score, err)
	}
}

func TestGetTopRanked(t *testing.T) {
	start := int64(0)
	stop := int64(2)
	key := GlobalScore
	expected := []redis.Z{
		{Member: "video1", Score: 10.0},
		{Member: "video2", Score: 8.5},
		{Member: "video3", Score: 7.0},
	}

	testCases := []struct {
		name string
		body struct {
			start    int64
			stop     int64
			key      string
			member   string
			expected []redis.Z
		}
		buildStubs   func(mockStore redismock.ClientMock)
		checkReponse func(t *testing.T, result []redis.Z, err error)
	}{
		{
			name: "OK",
			body: struct {
				start    int64
				stop     int64
				key      string
				member   string
				expected []redis.Z
			}{
				start:    start,
				stop:     stop,
				key:      key,
				expected: expected,
			},
			buildStubs: func(mockStore redismock.ClientMock) {
				mockStore.ExpectZRevRangeWithScores(key, start, stop).SetVal(expected)
			},
			checkReponse: func(t *testing.T, result []redis.Z, err error) {
				require.NoError(t, err)
				require.Equal(t, int64(len(result)), stop-start+1)
				require.Equal(t, expected, result)
			},
		},
		{
			name: "Error",
			body: struct {
				start    int64
				stop     int64
				key      string
				member   string
				expected []redis.Z
			}{
				start:    start,
				stop:     stop,
				key:      key,
				expected: expected,
			},
			buildStubs: func(mockStore redismock.ClientMock) {
				err := errors.New("connection refused")
				mockStore.ExpectZRevRangeWithScores(key, start, stop).SetErr(err)
			},
			checkReponse: func(t *testing.T, result []redis.Z, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "connection refused")
				require.Nil(t, result)
			},
		},
	}

	for _, tc := range testCases {
		db, mock := redismock.NewClientMock()
		defer db.Close()

		accessor := NewRankingDataAccessor(db)
		ctx := context.Background()
		tc.buildStubs(mock)

		result, err := accessor.GetTopRanked(ctx, key, start, stop)
		tc.checkReponse(t, result, err)
	}
}
