package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	httpRequestDurationSummary = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_duration_summary_seconds",
			Help: "HTTP request duration summary in seconds",
		},
		[]string{"method", "endpoint"},
	)

	// Database metrics
	dbConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Number of active database connections",
		},
	)

	dbConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)

	dbQueryErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_query_errors_total",
			Help: "Total number of database query errors",
		},
		[]string{"operation", "table"},
	)

	// Redis metrics
	redisConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "redis_connections_active",
			Help: "Number of active Redis connections",
		},
	)

	redisCommandsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_commands_total",
			Help: "Total number of Redis commands",
		},
		[]string{"command"},
	)

	redisCommandDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redis_command_duration_seconds",
			Help:    "Redis command duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"command"},
	)

	redisCacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "redis_cache_hits_total",
			Help: "Total number of Redis cache hits",
		},
	)

	redisCacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "redis_cache_misses_total",
			Help: "Total number of Redis cache misses",
		},
	)

	// Worker/Job metrics
	workerJobsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_jobs_total",
			Help: "Total number of worker jobs processed",
		},
		[]string{"type", "status"},
	)

	workerJobDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "worker_job_duration_seconds",
			Help:    "Worker job duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"type"},
	)

	workerQueueSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "worker_queue_size",
			Help: "Number of jobs in worker queue",
		},
		[]string{"queue"},
	)

	// Points system metrics
	pointsEarnedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "points_earned_total",
			Help: "Total points earned",
		},
		[]string{"source"},
	)

	pointsRedeemedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "points_redeemed_total",
			Help: "Total points redeemed",
		},
		[]string{"partner"},
	)

	pointsBalance = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "points_balance",
			Help: "Current points balance",
		},
		[]string{"user_id"},
	)

	// Business metrics
	activeUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users",
			Help: "Number of active users",
		},
	)

	newSignups = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "new_signups_total",
			Help: "Total number of new signups",
		},
	)

	postsCreated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "posts_created_total",
			Help: "Total number of posts created",
		},
	)

	connectionsMade = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "connections_made_total",
			Help: "Total number of connections made",
		},
	)
)

// HTTP metrics functions

func RecordHTTPRequest(method, endpoint string, status int, duration time.Duration) {
	httpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(status)).Inc()
	httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
	httpRequestDurationSummary.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// Database metrics functions

func SetDBConnections(active, idle int) {
	dbConnectionsActive.Set(float64(active))
	dbConnectionsIdle.Set(float64(idle))
}

func RecordDBQuery(operation, table string, duration time.Duration) {
	dbQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

func RecordDBQueryError(operation, table string) {
	dbQueryErrorsTotal.WithLabelValues(operation, table).Inc()
}

// Redis metrics functions

func SetRedisConnections(count int) {
	redisConnectionsActive.Set(float64(count))
}

func RecordRedisCommand(command string, duration time.Duration) {
	redisCommandsTotal.WithLabelValues(command).Inc()
	redisCommandDuration.WithLabelValues(command).Observe(duration.Seconds())
}

func RecordRedisCacheHit() {
	redisCacheHits.Inc()
}

func RecordRedisCacheMiss() {
	redisCacheMisses.Inc()
}

// Worker metrics functions

func RecordWorkerJob(jobType, status string, duration time.Duration) {
	workerJobsTotal.WithLabelValues(jobType, status).Inc()
	workerJobDuration.WithLabelValues(jobType).Observe(duration.Seconds())
}

func SetWorkerQueueSize(queue string, size int) {
	workerQueueSize.WithLabelValues(queue).Set(float64(size))
}

// Points metrics functions

func RecordPointsEarned(source string, amount float64) {
	pointsEarnedTotal.WithLabelValues(source).Add(amount)
}

func RecordPointsRedeemed(partner string, amount float64) {
	pointsRedeemedTotal.WithLabelValues(partner).Add(amount)
}

func SetPointsBalance(userID string, balance float64) {
	pointsBalance.WithLabelValues(userID).Set(balance)
}

// Business metrics functions

func SetActiveUsers(count int) {
	activeUsers.Set(float64(count))
}

func RecordNewSignup() {
	newSignups.Inc()
}

func RecordPostCreated() {
	postsCreated.Inc()
}

func RecordConnectionMade() {
	connectionsMade.Inc()
}
