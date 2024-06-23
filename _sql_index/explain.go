package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"regexp"
	"sort"
	"strconv"
)

const (
	user_id      = "1"
	storageUntil = "2077-07-07 01:45:11.743128+03"
	issuedAt     = "2028-08-08  12:32:19.743128+03"
	hash         = "qwertyuiopasdfghjklyuasdfghjkzxcvbnm"
)

func main() {
	ctx := context.Background()
	connString := "postgres://avrigne:8679@localhost/_sql_index?sslmode=disable"
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := NewRepository(pool, ctx)

	//Run one by one, using comment to exclude
	repo.InsertExplain()
	//repo.SelectExistsExplain()
	//repo.SelectOrdersExplain()
	//repo.UpdateExplain()

	//Modify insert, to returned=true
	//repo.SelectReturnsExplain()
}

type repository struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func NewRepository(pool *pgxpool.Pool, ctx context.Context) *repository {
	return &repository{pool: pool, ctx: ctx}
}

func (repo *repository) AnalyzeQueryPlan(query string, args ...interface{}) (float64, float64, error) {
	conn, err := repo.pool.Acquire(repo.ctx)
	if err != nil {
		return 0, 0, err
	}
	defer conn.Release()

	explainQuery := "EXPLAIN (ANALYZE, VERBOSE) " + query
	rows, err := conn.Query(repo.ctx, explainQuery, args...)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	var prepTime, execTime float64
	prepTimeRegex := regexp.MustCompile(`Planning Time: (\d+\.\d+) ms`)
	execTimeRegex := regexp.MustCompile(`Execution Time: (\d+\.\d+) ms`)

	for rows.Next() {
		var plan string
		if err := rows.Scan(&plan); err != nil {
			return 0, 0, err
		}
		fmt.Println(plan)
		if matches := prepTimeRegex.FindStringSubmatch(plan); matches != nil {
			prepTime, err = strconv.ParseFloat(matches[1], 64)
			if err != nil {
				return 0, 0, err
			}
		}

		if matches := execTimeRegex.FindStringSubmatch(plan); matches != nil {
			execTime, err = strconv.ParseFloat(matches[1], 64)
			if err != nil {
				return 0, 0, err
			}
		}
	}

	if err := rows.Err(); err != nil {
		return 0, 0, err
	}

	return prepTime, execTime, nil
}

func median(values []float64) float64 {
	sort.Float64s(values)
	n := len(values)
	if n%2 == 0 {
		return (values[n/2-1] + values[n/2]) / 2
	}
	return values[n/2]
}

func (repo *repository) InsertExplain() {
	var prepTimes, execTimes []float64
	query := `INSERT INTO orders (id, user_id, storage_until, issued, issued_at, returned, hash) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	for i := 1; i <= 1000; i++ {
		prepTime, execTime, err := repo.AnalyzeQueryPlan(query,
			strconv.Itoa(i), user_id, storageUntil, false, issuedAt, true, hash,
		)
		if err != nil {
			log.Fatal(err)
		}

		prepTimes = append(prepTimes, prepTime)
		execTimes = append(execTimes, execTime)
	}
	fmt.Printf("Median Preparation Time: %.2f ms\n", median(prepTimes))
	fmt.Printf("Median Execution Time: %.2f ms\n", median(execTimes))
}

func (repo *repository) UpdateExplain() {
	var prepTimes, execTimes []float64
	query := `UPDATE orders SET issued=$1, issued_at=$2, returned=$3
              WHERE id=$4`
	for i := 1; i <= 1000; i++ {
		prepTime, execTime, err := repo.AnalyzeQueryPlan(query,
			true, issuedAt, false, strconv.Itoa(i),
		)
		if err != nil {
			log.Fatal(err)
		}

		prepTimes = append(prepTimes, prepTime)
		execTimes = append(execTimes, execTime)
	}
	fmt.Printf("Median Preparation Time: %.2f ms\n", median(prepTimes))
	fmt.Printf("Median Execution Time: %.2f ms\n", median(execTimes))
}

func (repo *repository) SelectExistsExplain() {
	var prepTimes, execTimes []float64
	query := `SELECT EXISTS(SELECT 1 FROM orders WHERE id=$1)`
	for i := 1; i <= 1000; i++ {
		prepTime, execTime, err := repo.AnalyzeQueryPlan(query,
			strconv.Itoa(i),
		)
		if err != nil {
			log.Fatal(err)
		}

		prepTimes = append(prepTimes, prepTime)
		execTimes = append(execTimes, execTime)
	}
	fmt.Printf("Median Preparation Time: %.2f ms\n", median(prepTimes))
	fmt.Printf("Median Execution Time: %.2f ms\n", median(execTimes))
}

func (repo *repository) SelectOrdersExplain() {
	var prepTimes, execTimes []float64
	query := `
        SELECT id, user_id, issued, storage_until, returned
        FROM orders
        WHERE user_id = $1 AND issued = FALSE
        ORDER BY storage_until DESC
        LIMIT $2
    `
	for i := 1; i <= 1000; i++ {
		prepTime, execTime, err := repo.AnalyzeQueryPlan(query,
			user_id, 1000,
		)
		if err != nil {
			log.Fatal(err)
		}

		prepTimes = append(prepTimes, prepTime)
		execTimes = append(execTimes, execTime)
	}
	fmt.Printf("Median Preparation Time: %.2f ms\n", median(prepTimes))
	fmt.Printf("Median Execution Time: %.2f ms\n", median(execTimes))
}

func (repo *repository) SelectReturnsExplain() {
	var prepTimes, execTimes []float64
	query := `
        SELECT id, user_id, storage_until, issued, issued_at, returned
        FROM orders
        WHERE returned = TRUE
        ORDER BY id
        LIMIT $1 OFFSET $2
    `
	for i := 1; i <= 1000; i++ {
		prepTime, execTime, err := repo.AnalyzeQueryPlan(query,
			1000, 0,
		)
		if err != nil {
			log.Fatal(err)
		}

		prepTimes = append(prepTimes, prepTime)
		execTimes = append(execTimes, execTime)
	}
	fmt.Printf("Median Preparation Time: %.2f ms\n", median(prepTimes))
	fmt.Printf("Median Execution Time: %.2f ms\n", median(execTimes))
}
