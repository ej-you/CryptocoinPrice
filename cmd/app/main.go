// App binary starts HTTP-server.
package main

import (
	"github.com/sirupsen/logrus"

	"CryptocoinPrice/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		logrus.Fatal(err)
	}
	if err := application.Run(); err != nil {
		logrus.Fatal(err)
	}
}
