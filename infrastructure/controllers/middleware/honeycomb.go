package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/config"
	"github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/otel-config-go/otelconfig"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Custom middleware for OpenTelemetry instrumentation
func HoneycombHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceName := config.Configuration.GetString("honeycomb.dev.name")
		apiKey := config.Configuration.GetString("honeycomb.dev.api_key")
	
		// Enable multi-span attributes
		bsp := honeycomb.NewBaggageSpanProcessor()
	
		// Use the Honeycomb distro to set up the OpenTelemetry SDK
		otelShutdown, err := otelconfig.ConfigureOpenTelemetry(
			otelconfig.WithSpanProcessor(bsp),
			otelconfig.WithServiceName(serviceName),
			honeycomb.WithApiKey(apiKey), 
		)
		if err != nil {
			log.Fatalf("error setting up OTel SDK - %s", err)
		} else {
			log.Println("Connected to honeycomb")
		}
		defer otelShutdown()
	
		// This is where you can put any custom logic you want to apply to all requests.
		// In this case, we're wrapping the request with OpenTelemetry instrumentation.
		otelgin.Middleware("gin-server")(c)
		c.Next()
	}
}