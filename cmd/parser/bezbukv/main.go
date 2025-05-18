package main

import (
	"github.com/Markard/wordka/config"
	"github.com/Markard/wordka/internal/repo"
	"github.com/Markard/wordka/pkg/logger"
	"github.com/Markard/wordka/pkg/postgres"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type FiveLetterNounParser struct {
	logger  logger.Interface
	baseUrl string
}

func NewFiveLetterNounParser(baseUrl string, lgr logger.Interface) *FiveLetterNounParser {
	return &FiveLetterNounParser{
		logger:  lgr,
		baseUrl: baseUrl,
	}
}

func (p *FiveLetterNounParser) Parse() []string {
	var result []string
	page := 1
	for {
		url := p.getUrl(page)
		p.logger.Info("Fetching words from %s", url)
		resp, err := http.Get(url)
		if err != nil {
			p.logger.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			p.logger.Fatal(err)
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			p.logger.Fatal(err)
		}

		foundWords := 0
		doc.Find("div.view").Each(func(i int, s *goquery.Selection) {
			re := regexp.MustCompile(`[а-яА-ЯёЁ-]{5}`)
			lines := strings.Split(s.Text(), "\n")
			for _, line := range lines {
				matches := re.FindStringSubmatch(line)
				if len(matches) == 1 {
					word := strings.ToLower(matches[0])
					if !slices.Contains(result, word) {
						result = append(result, word)
						foundWords++
					}
				}
			}
		})

		if foundWords == 0 {
			break
		}
		page++
	}

	return result
}

func (p *FiveLetterNounParser) getUrl(page int) string {
	return p.baseUrl + strconv.Itoa(page)
}

const baseURL = "https://bezbukv.ru/mask/%2A%2A%2A%2A%2A/noun?page="

func main() {
	setup := config.MustLoad()
	lgr := logger.New(setup.Config.Log.Level, setup.Config.Log.CallerSkipFrameCount, nil)
	parser := NewFiveLetterNounParser(baseURL, lgr)
	words := parser.Parse()

	db := postgres.New(setup.Env.PgDSN, lgr.ZerologLogger())
	defer func() {
		err := db.Close()
		if err != nil {
			lgr.Error(err)
		}
	}()
	gameRepo := repo.NewGameRepository(db)
	err := gameRepo.SaveWords(words)
	if err != nil {
		lgr.Error(err)
	}
	lgr.Info("Successfully fetched %d words", len(words))
}
