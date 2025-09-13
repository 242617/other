package query_sql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/ollama/ollama/api"
)

func New(dsn string) *QuerySQL { return &QuerySQL{dsn: dsn} }

type QuerySQL struct{ dsn string }

func (QuerySQL) Name() string { return "query_sql" }
func (q QuerySQL) Tool() api.Tool {
	toolFunction := api.ToolFunction{
		Name:        q.Name(),
		Description: "Query database using SQL.",
	}
	toolFunction.Parameters.Properties = map[string]api.ToolProperty{
		"query": {
			Type:        []string{"string"},
			Description: "SQL query. SELECT only.",
		},
	}
	toolFunction.Parameters.Required = []string{"query"}

	return api.Tool{
		Type:     "function",
		Function: toolFunction,
	}
}

func (q *QuerySQL) Call(ctx context.Context, args map[string]any) (string, error) {
	query, ok := args["query"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected type for query: %T (%v)", args["query"], args["query"])
	}

	cfg, err := pgxpool.ParseConfig(q.dsn)
	if err != nil {
		return "", errors.Wrap(err, "pgxpool parse config")
	}

	cfg.MaxConnLifetime = time.Minute
	cfg.MaxConnIdleTime = 30 * time.Second
	cfg.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return "", errors.Wrap(err, "pgxpool new with config")
	}

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return "", errors.Wrap(err, "query")
	}
	defer rows.Close()

	cols := rows.FieldDescriptions()
	head := make([]string, len(cols))
	for i, column := range cols {
		head[i] = string(column.Name)
	}

	var sb strings.Builder
	sb.WriteString(strings.Join(head, "\t"))
	sb.WriteByte('\n')

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return "", errors.Wrap(err, "rows values")
		}

		strs := make([]string, len(values))
		for i, v := range values {
			if v == nil {
				strs[i] = "NULL"
				continue
			}
			switch vv := v.(type) {
			case [16]byte:
				strs[i] = uuid.UUID(vv).String()
			case []byte:
				strs[i] = string(vv)
			case string:
				strs[i] = vv
			case int64, int32, int16, int8, uint64, uint32, uint16, uint8:
				strs[i] = fmt.Sprintf("%v", vv)
			case float64, float32:
				strs[i] = fmt.Sprintf("%g", vv)
			case bool:
				strs[i] = fmt.Sprintf("%t", vv)
			case time.Time:
				strs[i] = vv.Format(time.RFC3339Nano)
			default:
				strs[i] = fmt.Sprintf("%v", vv)
			}
		}
		sb.WriteString(strings.Join(strs, ","))
		sb.WriteByte('\n')

	}

	if err := rows.Err(); err != nil {
		return "", errors.Wrap(err, "rows err")
	}

	return sb.String(), nil
}
