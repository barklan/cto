package alembic

import (
	_ "github.com/lib/pq"
)

// func TestMain(m *testing.M) {
// 	db, pool, resource := PrepareTestDB()
// 	dbx = sqlx.NewDb(db, "postgres")

// 	// Run tests
// 	code := m.Run()

// 	// You can't defer this because os.Exit doesn't care for defer
// 	if err := pool.Purge(resource); err != nil {
// 		log.Fatalf("Could not purge resource: %s", err)
// 	}

// 	os.Exit(code)
// }

// func TestGetAlembicHeadFromImage(t *testing.T) {
// 	// TODO should be tested with dockertest artificial image
// 	os.Setenv("GITLAB_TOKEN_USERNAME", "maintoken")
// 	os.Setenv("GITLAB_TOKEN_PASSWORD", "6ptmCm43sWS_yb3Uy7s3")

// 	got := GetAlembicHeadFromImage("registry.gitlab.com/nftgalleryx/nftgallery_backend/backend", "0188b4f8")
// 	want := "0188b4f8"

// 	if got != want {
// 		t.Errorf("got %q want %q", got, want)
// 	}
// }

// func TestGetAlembicHistoryFromImage(t *testing.T) {
// 	got := GetAlembicHistoryFromImage("registry.gitlab.com/nftgalleryx/nftgallery_backend/backend", "0188b4f8")
// 	gotOne := got[0]
// 	want := "33474d0f7488"

// 	if gotOne != want {
// 		t.Errorf("got %q want %q", got, want)
// 	}
// }

// func TestGetAlembicVersionFromDB(t *testing.T) {
// 	want := "boom"

// 	tx := dbx.MustBegin()
// 	tx.MustExec("CREATE TABLE IF NOT EXISTS alembic_version(version_num VARCHAR (50) NOT NULL);")
// 	tx.MustExec("INSERT INTO alembic_version (version_num)"+
// 		"VALUES ($1)", want)

// 	err := tx.Commit()
// 	if err != nil {
// 		panic("failed to commit transaction")
// 	}

// 	got := GetAlembicVersionFromDB("stag")

// 	if got != want {
// 		t.Errorf("got %q want %q", got, want)
// 	}
// }

// var dbx *sqlx.DB

// func PrepareTestDB() (*sql.DB, *dockertest.Pool, *dockertest.Resource) {
// 	var db *sql.DB
// 	pool, err := dockertest.NewPool("")
// 	if err != nil {
// 		log.Fatalf("Could not connect to docker: %s", err)
// 	}

// 	// pulls an image, creates a container based on it and runs it
// 	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
// 		Repository: "postgres",
// 		Tag:        "13",
// 		Env: []string{
// 			"POSTGRES_PASSWORD=postgres",
// 			"POSTGRES_USER=postgres",
// 			"POSTGRES_DB=app",
// 			"listen_addresses = '*'",
// 		},
// 		Name: "stag_db",
// 	}, func(config *docker.HostConfig) {
// 		// set AutoRemove to true so that stopped container goes away by itself
// 		config.AutoRemove = true
// 		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
// 	})
// 	if err != nil {
// 		log.Fatalf("Could not start resource: %s", err)
// 	}

// 	hostAndPort := resource.GetHostPort("5432/tcp")
// 	databaseUrl := fmt.Sprintf("postgres://postgres:postgres@%s/app?sslmode=disable",
// 		hostAndPort)

// 	log.Println("Connecting to database on url: ", databaseUrl)

// 	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

// 	// exponential backoff-retry
// 	pool.MaxWait = 120 * time.Second
// 	if err = pool.Retry(func() error {
// 		db, err = sql.Open("postgres", databaseUrl)
// 		if err != nil {
// 			return err
// 		}
// 		return db.Ping()
// 	}); err != nil {
// 		log.Fatalf("Could not connect to docker: %s", err)
// 	}

// 	return db, pool, resource
// }
