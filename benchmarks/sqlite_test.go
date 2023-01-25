package benchmarks

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/indexes/hash"
	"github.com/shamcode/simd/namespace"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/where"
	"log"
	"strconv"
	"testing"
)

func Benchmark_SIMDVsSQLite(b *testing.B) {
	usersCountForBenchmarking := []int{
		10,
		100,
		1_000,
		5_000,
		10_000,
		50_000,
	}

	for _, usersCount := range usersCountForBenchmarking {
		db, err := sql.Open("sqlite3", "file:cachedb?mode=memory&cache=shared")
		if err != nil {
			log.Fatal(err)
		}
		sqlStmt := `
	CREATE TABLE IF NOT EXISTS user (
		id INTEGER NOT NULL PRIMARY KEY, 
		name TEXT NOT NULL, 
		status INTEGER NOT NULL, 
		score INTEGER NOT NULL,
		is_online INTEGER NOT NULL
	);
	DELETE FROM user;
	CREATE INDEX IF NOT EXISTS id_idx ON user(id);
	`
		_, err = db.Exec(sqlStmt)
		if nil != err {
			b.Fatal(err)
		}

		simd := namespace.CreateNamespace()
		simd.AddIndex(hash.NewInt64HashIndex(userID, true))

		stmt, err := db.Prepare("INSERT INTO user (id, name, status, score, is_online) VALUES(?, ?, ?, ?, ?)")
		if nil != err {
			b.Fatal(err)
		}

		for i := 1; i < usersCount; i++ {
			user := &User{
				ID:       int64(i),
				Name:     "user_" + strconv.Itoa(i),
				Status:   StatusEnum(1 + i%2),
				Score:    i % 150,
				IsOnline: 0 == i%2,
			}
			err := simd.Upsert(user)
			if nil != err {
				b.Fatal(err)
			}
			_, err = stmt.Exec(user.ID, user.Name, user.Status, user.Score, user.IsOnline)
			if nil != err {
				b.Fatal(err)
			}
		}

		stmt.Close()

		qe := executor.CreateQueryExecutor(simd)
		b.Run(strconv.Itoa(usersCount)+"_simd", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for i := 1; i < usersCount/4; i++ {
					cur, err := qe.FetchAll(
						context.Background(),
						query.NewBuilder(
							query.WhereInt64(userID, where.EQ, int64(i)),
						).Query(),
					)
					if nil != err {
						b.Fatalf("query: %s", err)
					}
					u := cur.Item().(*User)
					if u.IsOnline != (0 == i%2) {
						b.Fatalf("wrong is_online: %d", i)
					}
				}
			}
		})

		stmt, err = db.Prepare("SELECT is_online FROM user WHERE id = ?")
		if nil != err {
			b.Fatal(err)
		}
		b.Run(strconv.Itoa(usersCount)+"_sqlite", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for i := 1; i < usersCount/4; i++ {
					rows, err := stmt.QueryContext(context.Background(), i)
					if nil != err {
						b.Fatal(err)
					}

					if !rows.Next() {
						b.Fatalf("not found user: %d", i)
					}
					var isOnline bool
					err = rows.Scan(&isOnline)
					if nil != err {
						b.Fatal(err)
					}
					if isOnline != (0 == i%2) {
						b.Fatalf("wrong is_online: %d", i)
					}
					rows.Close()
				}
			}
		})
		stmt.Close()
		db.Close()
	}
}
