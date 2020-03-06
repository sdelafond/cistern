package providers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/nbedos/cistern/utils"
)

func TestTravisClientfetchPipeline(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/build/609256446" {
			bs, err := ioutil.ReadFile("test_data/travis/travis_build_609256446.json")
			if err != nil {
				t.Fatal(err)
			}
			if _, err := fmt.Fprint(w, string(bs)); err != nil {
				t.Fatal(err)
			}
		}
	}))
	defer ts.Close()

	URL, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	client := TravisClient{
		baseURL:     *URL,
		httpClient:  ts.Client(),
		rateLimiter: time.Tick(time.Millisecond),
		token:       "token",
		provider: Provider{
			ID:   "id",
			Name: "name",
		},
		buildsPageSize: 10,
	}

	expectedPipeline := Pipeline{
		Number: "72",
		GitReference: GitReference{
			SHA:   "c824642cc7c3abf8abc2d522b58a345a98b95b9b",
			Ref:   "feature/travis_improvements",
			IsTag: false,
		},
		Step: Step{
			ID:        "609256446",
			Type:      StepPipeline,
			State:     Failed,
			CreatedAt: time.Date(2019, 11, 8, 14, 26, 21, 506000000, time.UTC),
			StartedAt: utils.NullTime{
				Valid: true,
				Time:  time.Date(2019, 11, 8, 20, 53, 52, 0, time.UTC),
			},
			FinishedAt: utils.NullTime{
				Valid: true,
				Time:  time.Date(2019, 11, 8, 20, 54, 18, 0, time.UTC),
			},
			UpdatedAt: time.Date(2019, 11, 8, 20, 54, 19, 108000000, time.UTC),
			Duration: utils.NullDuration{
				Valid:    true,
				Duration: 114 * time.Second,
			},
			WebURL: utils.NullString{
				String: fmt.Sprintf("%s/nbedos/cistern/builds/609256446", ts.URL),
				Valid:  true,
			},
		},
	}

	expectedPipeline.Children = []Step{
		{
			ID:        "11290169",
			Type:      StepStage,
			Name:      "Tests",
			State:     Failed,
			CreatedAt: time.Date(2019, 11, 8, 14, 26, 21, 506000000, time.UTC),
			StartedAt: utils.NullTime{
				Valid: true,
				Time:  time.Date(2019, 11, 8, 20, 53, 52, 0, time.UTC),
			},
			FinishedAt: utils.NullTime{
				Valid: true,
				Time:  time.Date(2019, 11, 8, 20, 54, 18, 0, time.UTC),
			},
			Duration: utils.NullDuration{
				Valid:    true,
				Duration: 26 * time.Second,
			},
			WebURL: utils.NullString{
				Valid:  true,
				String: fmt.Sprintf("%s/nbedos/cistern/builds/609256446", ts.URL),
			},
		},
	}
	expectedPipeline.Children[0].Children = []Step{
		{
			ID:        "609256447",
			Type:      StepJob,
			State:     Failed,
			Name:      "GoLang 1.13 on Ubuntu Bionic",
			CreatedAt: time.Date(2019, 11, 8, 14, 26, 21, 506000000, time.UTC),
			StartedAt: utils.NullTime{
				Valid: true,
				Time:  time.Date(2019, 11, 8, 20, 53, 52, 0, time.UTC),
			},
			FinishedAt: utils.NullTime{
				Valid: true,
				Time:  time.Date(2019, 11, 8, 20, 54, 18, 0, time.UTC),
			},
			Duration: utils.NullDuration{
				Valid:    true,
				Duration: 26 * time.Second,
			},
			Log: Log{},
			WebURL: utils.NullString{
				String: fmt.Sprintf("%s/nbedos/cistern/jobs/609256447", ts.URL),
				Valid:  true,
			},
			AllowFailure: false,
		},
		{
			ID:        "609256448",
			Type:      StepJob,
			State:     Failed,
			Name:      "GoLang 1.12 on Ubuntu Trusty",
			CreatedAt: time.Date(2019, 11, 8, 14, 26, 21, 509000000, time.UTC),
			StartedAt: utils.NullTime{
				Valid: true,
				Time:  time.Date(2019, 11, 8, 20, 32, 48, 0, time.UTC),
			},
			FinishedAt: utils.NullTime{
				Valid: true,
				Time:  time.Date(2019, 11, 8, 20, 33, 18, 0, time.UTC),
			},
			Duration: utils.NullDuration{
				Valid:    true,
				Duration: 30 * time.Second,
			},
			Log: Log{},
			WebURL: utils.NullString{
				String: fmt.Sprintf("%s/nbedos/cistern/jobs/609256448", ts.URL),
				Valid:  true,
			},
			AllowFailure: false,
		},
		{
			ID:        "609256449",
			Type:      StepJob,
			State:     Failed,
			Name:      "GoLang 1.13 on macOS 10.14",
			CreatedAt: time.Date(2019, 11, 8, 14, 26, 21, 512000000, time.UTC),
			StartedAt: utils.NullTime{
				Valid: true,
				Time:  time.Date(2019, 11, 8, 20, 33, 44, 0, time.UTC),
			},
			FinishedAt: utils.NullTime{
				Valid: true,
				Time:  time.Date(2019, 11, 8, 20, 34, 15, 0, time.UTC),
			},
			Duration: utils.NullDuration{
				Valid:    true,
				Duration: 31 * time.Second,
			},
			Log: Log{},
			WebURL: utils.NullString{
				String: fmt.Sprintf("%s/nbedos/cistern/jobs/609256449", ts.URL),
				Valid:  true,
			},
			AllowFailure: false,
		},
		{
			ID:        "609256450",
			Type:      StepJob,
			State:     Failed,
			Name:      "GoLang 1.12 on macOS 10.13",
			CreatedAt: time.Date(2019, 11, 8, 14, 26, 21, 514000000, time.UTC),
			StartedAt: utils.NullTime{
				Valid: true,
				Time:  time.Date(2019, 11, 8, 20, 33, 39, 0, time.UTC),
			},
			FinishedAt: utils.NullTime{
				Valid: true,
				Time:  time.Date(2019, 11, 8, 20, 34, 06, 0, time.UTC),
			},
			Duration: utils.NullDuration{
				Valid:    true,
				Duration: 27 * time.Second,
			},
			Log: Log{},
			WebURL: utils.NullString{
				String: fmt.Sprintf("%s/nbedos/cistern/jobs/609256450", ts.URL),
				Valid:  true,
			},
			AllowFailure: false,
		},
	}

	pipeline, err := client.fetchPipeline(context.Background(), "nbedos/cistern", "609256446")
	if err != nil {
		t.Fatal(err)
	}

	if diff := expectedPipeline.Diff(pipeline); diff != "" {
		t.Log(diff)
		t.Fatal("invalid pipeline")
	}
}

func TestParseTravisWebURL(t *testing.T) {
	u := "https://travis-ci.org/nbedos/termtosvg/builds/612815758"

	owner, repo, id, err := parseTravisWebURL(&TravisOrgURL, u)
	if err != nil {
		t.Fatal(err)
	}

	if owner != "nbedos" || repo != "termtosvg" || id != "612815758" {
		t.Fail()
	}
}
