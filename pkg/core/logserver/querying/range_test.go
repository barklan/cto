package querying

import (
	"testing"
	"time"
)

func gotEqWant(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestValidPrefix(t *testing.T) {
	rq := RequestQuery{
		ProjectName: "nftg",
		Env:         "futurearts.ru",
		Service:     "celerybeat-logger",
		Date:        "2021-11-03",
		TimeQuery:   "13:42:",
	}
	got := rq.ValidPrefix(prefixSyntax)
	want := "nftg futurearts.ru celerybeat-logger 2021-11-03 13:42:"
	gotEqWant(t, got, want)
}

func TestBeaconToSeek(t *testing.T) {
	t.Run("minimal reqeust", func(t *testing.T) {
		rq := RequestQuery{
			ProjectName: "nftg",
			Env:         "futurearts.ru",
			Service:     "celerybeat-logger",
			Date:        "2021-11-03",
			TimeQuery:   "",
		}

		got, _ := rq.BeaconToSeek(prefixSyntax, time.Now())
		want := "nftg futurearts.ru celerybeat-logger 2021-11-03 24:00:00"

		gotEqWant(t, got, want)
	})
}

func TestTimeQueryBeaconToSeek(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
		want  string
	}{
		{
			"one hour range",
			"03",
			"04",
		},
		{
			"one hour range with colons",
			"03:",
			"04",
		},
		{
			"ten minute range",
			"03:1",
			"03:2",
		},
		{
			"one minute range",
			"03:05",
			"03:06",
		},
		{
			"one minute range edge",
			"03:59",
			"03:60",
		},
		{
			"one minute range with colon",
			"03:05:",
			"03:06",
		},
		{
			"one minute range with colon edge",
			"03:59:",
			"03:60",
		},
		{
			"ten second range",
			"03:11:5",
			"03:11:6",
		},
		{
			"one second range",
			"03:11:59",
			"03:11:60",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := TimeQueryBeaconToSeek(tC.input)
			gotEqWant(t, got, tC.want)
		})
	}
}

func TestRequestQuery_NorthStar(t *testing.T) {
	type fields struct {
		ProjectName string
		Env         string
		Service     string
		Date        string
		TimeQuery   string
	}
	tests := []struct {
		name   string
		fields fields
		syntax timeSyntax
		nowStr string
		want   string
	}{
		{
			"Last minute syntax simple",
			fields{
				ProjectName: "nftg",
				Env:         "futurearts.ru",
				Service:     "celerybeat-logger",
				Date:        "2021-11-03",
				TimeQuery:   "10m",
			},
			lastMinutesSyntax,
			"2021-11-03 15:04:05",
			"nftg futurearts.ru celerybeat-logger 2021-11-03 14:54:05",
		},
		{
			"Jump through day",
			fields{
				ProjectName: "nftg",
				Env:         "futurearts.ru",
				Service:     "celerybeat-logger",
				Date:        "2021-11-03",
				TimeQuery:   "15m",
			},
			lastMinutesSyntax,
			"2021-11-03 00:03:13",
			"nftg futurearts.ru celerybeat-logger 2021-11-02 23:48:13",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := RequestQuery{
				ProjectName: tt.fields.ProjectName,
				Env:         tt.fields.Env,
				Service:     tt.fields.Service,
				Date:        tt.fields.Date,
				TimeQuery:   tt.fields.TimeQuery,
			}
			now, err := time.Parse("2006-01-02 15:04:05", tt.nowStr)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if got, _ := rq.NorthStar(tt.syntax, now); got != tt.want {
				t.Errorf("RequestQuery.NorthStar() = %v, want %v", got, tt.want)
			}
		})
	}
}
