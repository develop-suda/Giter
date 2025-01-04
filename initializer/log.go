package initializer

import (
	"context"
	"errors"
	"giter/public"
	"io"
	"os"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var RequestID string

// TODO: 日付を跨いで実行している場合、前日のファイルに書き込みをしてしまう
func Log() {
	// ginのログファイルを取得
	ginf, err := getLogFile("gin", "log")
	if err != nil {
		log.Error().Msg("error")
	}

	// zerologのログファイルを取得
	zerologf, err := getLogFile("zerolog", "log")
	if err != nil {
		log.Error().Msg("error")
	}

	// zerologのJSON形式のログファイルを取得
	zerologjsonf, err := getLogFile("zerolog", "json")
	if err != nil {
		log.Error().Msg("error")
	}
	// ginでのログをファイルと標準出力に出力
	gin.DefaultWriter = io.MultiWriter(ginf, os.Stdout)
	// zerologでのログをファイルと標準出力に出力
	log.Logger = log.Output(io.MultiWriter(zerologf, zerologjsonf))

	// logstudy()
}

// ログファイルを取得する関数
func getLogFile(fileName string, extension string) (*os.File, error) {
	yyyymmdd := public.FormatDate()

	filePath := "logs/" + fileName + "/" + yyyymmdd + "_" + fileName + "." + extension

	// ファイルが存在しない場合　err != nil
	// ファイルが存在する場合　　　　err = nil
	_, err := os.Stat(filePath)
	if err != nil {
		return createLogFile(filePath)
	} else {
		return openLogFile(filePath)
	}
}

// ログファイルを新規作成する関数
func createLogFile(filePath string) (*os.File, error) {
	return os.Create(filePath)
}

// 既存のログファイルを開く関数
func openLogFile(filePath string) (*os.File, error) {
	return os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

// リクエストIDを設定する関数
func SetRequestID(c *gin.Context) {
	RequestID = requestid.Get(c)
}

func logstudy() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.Trace().Send()
	log.Debug().Send()
	log.Info().Send()
	log.Warn().Send()
	log.Error().Send()
	// log.Fatal().Send()
	// log.Panic().Send()
	log.Info().Send()

	err := outer()
	// なんかスタックトレースでねえ
	log.Error().Stack().Err(err).Msg("")

	f()

}

func inner() error {
	return errors.New("seems we have an error here")
}

func middle() error {
	err := inner()
	if err != nil {
		return err
	}
	return nil
}

func outer() error {
	err := middle()
	if err != nil {
		return err
	}
	return nil
}

// こいつは使えそうや
func f() {
	logger := zerolog.New(os.Stdout)
	ctx := context.Background()

	// Attach the Logger to the context.Context
	ctx = logger.WithContext(ctx)
	someFunc(ctx)
}

func someFunc(ctx context.Context) {
	// Get Logger from the go Context. if it's nil, then
	// `zerolog.DefaultContextLogger` is returned, if
	// `DefaultContextLogger` is nil, then a disabled logger is returned.
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("Hello")
}
