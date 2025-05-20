package main

import (
	"github.com/Markard/wordka/config"
	"github.com/Markard/wordka/internal/repo"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/Markard/wordka/pkg/postgres"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Job struct {
	Id  int
	Url string
}

type Result struct {
	JobId int
	Words []string
}

type FiveLetterNounParser struct {
	logger  logger.Interface
	pages   int
	workers int
	baseUrl string
	results chan<- *Result
}

func ParserWithPagination(lgr logger.Interface, pages int, workers int, baseUrl string, results chan<- *Result) *FiveLetterNounParser {
	return &FiveLetterNounParser{logger: lgr, pages: pages, workers: workers, baseUrl: baseUrl, results: results}
}

func (p *FiveLetterNounParser) Parse() {
	jobs := make(chan *Job, p.pages)

	for w := 1; w <= p.workers; w++ {
		worker := NewWorker(w, p.logger, jobs, p.results)
		go worker.Start()
	}

	for j := 1; j <= p.pages; j++ {
		job := &Job{
			Id:  j,
			Url: p.baseUrl + strconv.Itoa(j),
		}
		jobs <- job
	}
	close(jobs)
}

type Worker struct {
	Id      int
	logger  logger.Interface
	Jobs    <-chan *Job
	Results chan<- *Result
}

func NewWorker(id int, logger logger.Interface, jobs <-chan *Job, results chan<- *Result) *Worker {
	return &Worker{
		Id:      id,
		logger:  logger,
		Jobs:    jobs,
		Results: results,
	}
}

func (w *Worker) Start() {
	for job := range w.Jobs {
		w.Results <- w.DoWork(job)
	}
}

func (w *Worker) DoWork(job *Job) *Result {
	w.logger.Info("Fetching words from %s", job.Url)
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Get(job.Url)
	if err != nil {
		w.logger.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != 200 {
		w.logger.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		w.logger.Fatal(err)
	}

	result := &Result{
		JobId: job.Id,
		Words: make([]string, 0),
	}
	doc.Find("div.view").Each(func(i int, s *goquery.Selection) {
		re := regexp.MustCompile(`[а-яА-ЯёЁ-]{5}`)
		lines := strings.Split(s.Text(), "\n")
		for _, line := range lines {
			matches := re.FindStringSubmatch(line)
			if len(matches) == 1 {
				word := strings.ToLower(matches[0])
				result.Words = append(result.Words, word)
			}
		}
	})

	return result
}

func main() {
	setup := config.MustLoad()
	lgr := logger.New(setup.Config.Log.Level, setup.Config.Log.CallerSkipFrameCount, nil)
	pages := 36
	results := make(chan *Result, pages)
	parser := ParserWithPagination(
		lgr,
		pages,
		5,
		"https://bezbukv.ru/mask/%2A%2A%2A%2A%2A/noun?page=",
		results,
	)
	parser.Parse()

	db := postgres.New(setup.Env.PgDSN, lgr.ZerologLogger())
	defer func() {
		err := db.Close()
		if err != nil {
			lgr.Error(err)
		}
	}()
	gameRepo := repo.NewGameRepository(db)

	for r := 1; r <= pages; r++ {
		result := <-results
		err := gameRepo.SaveWords(result.Words)
		if err != nil {
			lgr.Error(err)
		}
		lgr.Info("Successfully fetched %d words", len(result.Words))
	}
	close(results)
}
