package metrics

import (
	"context"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type StatsCollector struct {
	books    BookCounter
	users    UserCounter
	progress ProgressCounter

	booksTotal    *prometheus.Desc
	usersTotal    *prometheus.Desc
	booksByStatus *prometheus.Desc
}

type BookCounter interface {
	Count(ctx context.Context) (int64, error)
}
type UserCounter interface {
	Count(ctx context.Context) (int64, error)
}
type ProgressCounter interface {
	CountByStatus(ctx context.Context) (map[string]int64, error)
}

func NewStatsCollector(b BookCounter, u UserCounter, p ProgressCounter) *StatsCollector {
	return &StatsCollector{
		books:    b,
		users:    u,
		progress: p,
		booksTotal: prometheus.NewDesc(
			"pageturner_books_total",
			"Total number of books in the catalogue.",
			nil, nil,
		),
		usersTotal: prometheus.NewDesc(
			"pageturner_users_total",
			"Total number of registered users.",
			nil, nil,
		),
		booksByStatus: prometheus.NewDesc(
			"pageturner_books_by_status",
			"Number of book reading entries by status.",
			[]string{"status"}, nil,
		),
	}
}

func (c *StatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.booksTotal
	ch <- c.usersTotal
	ch <- c.booksByStatus
}

func (c *StatsCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if n, err := c.books.Count(ctx); err == nil {
		ch <- prometheus.MustNewConstMetric(c.booksTotal, prometheus.GaugeValue, float64(n))
	} else {
		log.Printf("metrics: books count failed: %v", err)
	}

	if n, err := c.users.Count(ctx); err == nil {
		ch <- prometheus.MustNewConstMetric(c.usersTotal, prometheus.GaugeValue, float64(n))
	} else {
		log.Printf("metrics: users count failed: %v", err)
	}

	if byStatus, err := c.progress.CountByStatus(ctx); err == nil {
		for status, n := range byStatus {
			ch <- prometheus.MustNewConstMetric(
				c.booksByStatus, prometheus.GaugeValue, float64(n), status)
		}
	} else {
		log.Printf("metrics: progress count failed: %v", err)
	}
}
