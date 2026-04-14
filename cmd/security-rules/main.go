package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"securityrules/security-rules/configs"
	"securityrules/security-rules/internal/app/handlers"
	"securityrules/security-rules/internal/middleware"
	"securityrules/security-rules/internal/routes"
	"securityrules/security-rules/internal/utils/azure"
	"securityrules/security-rules/internal/utils/log"
	"securityrules/security-rules/internal/utils/net"
	sf "securityrules/security-rules/internal/utils/snowflake"
	"securityrules/security-rules/internal/utils/types"
)

func main() {
	configs.Load()

	tp := newTracerProvider()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Logger.Error(fmt.Sprintf("unable to shutdown the trace provider: %v", err))
		}
	}()

	prepare()
	initializeFacade()

	app := fiber.New()

	middleware.FiberMiddleware(app)

	routes.PublicRoutes(app)
	routes.PrivateRoutes(app, privateRouteHandlers())
	routes.UtilityRoutes(app)

	net.StartServer(app)
}

func initializeFacade() {
	initializeVault()
	//gem := initializeGem()
	//aladdin := initializeAladdin()
	initializeSnowflake()
	//queue := types.NewQueue()

}

func newTracerProvider() *sdktrace.TracerProvider {
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		log.Logger.Fatal(fmt.Sprintf("unable to create a new OTEL tracer provider: %v", err))
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("security-rules"),
			)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp
}

func prepare() {
	//db.Connect()
}

// privateRouteHandlers configures specific handlers for the non-prod/prod environments.  For some external API calls, there is only a single environment.  For non-production
// it is recommended to simulate/mock the return responses to avoid resources being consumed/created on external platforms (i.e. GitLab, Permit.IO, etc.)
func privateRouteHandlers() routes.Handlers {
	if configs.EnvConfigs.GolangEnvironment.IsLocal() {
		return routes.Handlers{
			GetIdentity: handlers.GetIdentity,
		}
	}

	if configs.EnvConfigs.GolangEnvironment.IsProduction() {
		return routes.Handlers{
			GetIdentity: handlers.GetIdentity,
		}
	}

	return routes.Handlers{
		GetIdentity: handlers.GetIdentity,
	}
}

func initializeSnowflake() sf.Snowflake {

	authenticator, err := sf.ParseAuthType(configs.EnvConfigs.SnowflakeAuthenticator)

	if err != nil {
		msg := fmt.Sprintf("main.go: get Snowflake Authenticator - unable to parse with error: %v", err)
		log.Logger.Error(msg)
	}

	return sf.Snowflake{
		Account:       configs.EnvConfigs.SnowflakeAccount,
		User:          configs.EnvConfigs.SnowflakeUser,
		Role:          configs.EnvConfigs.SnowflakeRole,
		Warehouse:     configs.EnvConfigs.SnowflakeWarehouse,
		Database:      configs.EnvConfigs.SnowflakeDatabase,
		Schema:        configs.EnvConfigs.SnowflakeSchema,
		Authenticator: authenticator,
	}
}

func getSnowflakeSecrets() ([]byte, string, error) {
	keys := []string{configs.EnvConfigs.KeyVaultDerKey, configs.EnvConfigs.KeyVaultPwdKey}
	v := azure.NewVault(configs.EnvConfigs.KeyVaultUrl, !types.Environment.IsLocal(configs.EnvConfigs.GolangEnvironment))
	log.Logger.Info("main.go: getSnowflakeSecrets - retrieving values from Azure KeyVault for snowflake authentication...")
	m, err := v.GetMany(keys)
	if err != nil {
		msg := fmt.Sprintf("main.go: getSnowflakeSecrets - unable to get secrets from Azure KeyVault with error: %v", err)
		log.Logger.Error(msg)
		return nil, "", errors.Join(err, errors.New(msg))
	}
	decoded, err := base64.StdEncoding.DecodeString(m[configs.EnvConfigs.KeyVaultDerKey])
	if err != nil {
		msg := fmt.Sprintf("main.go: getSnowflakeSecrets - unable to decode DER value with error: %v", err)
		log.Logger.Error(msg)
		return nil, "", errors.Join(err, errors.New(msg))
	}

	log.Logger.Info("main.go: getSnowflakeSecrets - returning values from Azure KeyVault for snowflake authentication")
	return decoded, m[configs.EnvConfigs.KeyVaultPwdKey], nil
}

func initializeVault() azure.Vault {
	return azure.NewVault(
		configs.EnvConfigs.KeyVaultUrl,
		!types.Environment.IsLocal(configs.EnvConfigs.GolangEnvironment),
	)
}

func openSnowflakeConnection(facade *services.RefreshFacade) error {

	now := time.Now()

	log.Logger.Debug("main.go: openSnowflakeConnection - checking to see if the DB connection needs to be opened/re-opened...")
	if facade.DBConnection != nil {
		if facade.DBConnectionExpiresOn != nil &&
			!facade.DBConnectionExpiresOn.IsZero() &&
			now.After(*facade.DBConnectionExpiresOn) {
			log.Logger.Info(fmt.Sprintf("Connection expired: Closing it now: %v, expired on: %v", now, facade.DBConnectionExpiresOn))
			facade.DBConnection.Close()
			facade.DBConnection = nil
		} else {

			log.Logger.Debug("Connection reference is available")
			return nil
		}
	}

	der, pwd, err := getSnowflakeSecrets()
	if err != nil {
		msg := fmt.Sprintf("main.go openSnowflakeConnection - unable to retrieve secrets from Azure KeyVault with error: %v", err)
		log.Logger.Error(msg)
		return errors.Join(err, errors.New(msg))
	}
	log.Logger.Info("main.go: openSnowflakeConnection - connection is not open.  opening a new connection...")

	maxRetries := 2
	retries := 0

	var db *sql.DB
	for retries < maxRetries {
		db, err = facade.Snowflake.Open(der, pwd)
		if err != nil {
			msg := fmt.Sprintf("openSnowflakeConnection - unable to open snowflake connection with error: %v", err)
			log.Logger.Error(msg)
			return errors.Join(err, errors.New(msg))
		}

		err = db.Ping()
		if err == nil {
			facade.DBConnection = db
			connectionExpiresOn := now.Add(
				time.Duration(configs.EnvConfigs.SnowflakeConnectionTtlInMin * int(time.Minute)))

			facade.DBConnectionExpiresOn = &connectionExpiresOn

			log.Logger.Info(fmt.Sprintf("try [%d]: openSnowflakeConnection - DB connection successfully opened. Expires on: %v", retries, connectionExpiresOn))
			return nil

		} else {
			log.Logger.Warn(fmt.Sprintf("try [%d]: %s", retries, err.Error()))
		}

		retries++
	}

	return err
}

func closeSnowflakeConnection(facade *services.RefreshFacade) error {
	log.Logger.Debug("main.go: closeSnowflakeConnection - checking to see if a connection exists and should be closed...")
	if facade.DBConnection == nil {
		log.Logger.Debug("main.go: closeSnowflakeConnection - connection does not exist.  do not need to close the connection.")
		return nil
	}

	log.Logger.Debug("main.go: closeSnowflakeConnection - connection exists, attempting to close the DB connection...")
	if err := facade.DBConnection.Close(); err != nil {
		msg := fmt.Sprintf("main.go: closeSnowflakeConnection - unable to close snowflake connection with error: %v", err)
		log.Logger.Error(msg)
		return errors.Join(err, errors.New(msg))
	}

	log.Logger.Info("main.go: closeSnowflakeConnection - DB connection has been successfully closed.")
	log.Logger.Debug("main.go: closeSnowflakeConnection - clearing the reference to the DB connection from the facade")
	facade.DBConnection = nil
	return nil
}
