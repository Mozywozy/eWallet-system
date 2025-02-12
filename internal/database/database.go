package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
)

type Service interface {
	Health() map[string]string
	Close() error
	GetDB() *gorm.DB
	GetRedis() *redis.Client
}

type service struct {
	db *gorm.DB
	redis *redis.Client
}

var (
	dbname     = os.Getenv("BLUEPRINT_DB_DATABASE")
	password   = os.Getenv("BLUEPRINT_DB_PASSWORD")
	username   = os.Getenv("BLUEPRINT_DB_USERNAME")
	port       = os.Getenv("BLUEPRINT_DB_PORT")
	host       = os.Getenv("BLUEPRINT_DB_HOST")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error saat membuka koneksi ke database: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "", // Kosong jika tanpa password
		DB:       0,  // Default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Fatalf("Gagal menghubungkan ke Redis: %v", err)
	}

	log.Println("✅ Sukses terhubung ke database & Redis!")

	dbInstance = &service{
		db: db,
		redis: redisClient,
	}
	return dbInstance
}

func (s *service) GetDB() *gorm.DB {
	return s.db
}

func (s *service) GetRedis() *redis.Client {
	return s.redis
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	sqlDB, err := s.db.DB() // 🔥 Ambil instance database dari GORM
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("Database tidak dapat diakses: %v", err)
		log.Printf("Database health check failed: %v", err)
		return stats
	}

	// Coba ping database
	err = sqlDB.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("Database tidak dapat diakses: %v", err)
		log.Printf("Database health check failed: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "Database dalam kondisi sehat"

	// Ambil statistik database
	dbStats := sqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	return stats
}

// Close menutup koneksi database
func (s *service) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		log.Printf("❌ Gagal mendapatkan koneksi database untuk ditutup: %v", err)
		return err
	}

	log.Println("🔻 Menutup koneksi database...")
	return sqlDB.Close()
}