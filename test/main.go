package main

import (
	"context"
	"fmt"

	"github.com/nazarnovak/hobee-be/config"
	"github.com/nazarnovak/hobee-be/pkg/email"
	"github.com/nazarnovak/hobee-be/pkg/herrors2"
	"github.com/nazarnovak/hobee-be/pkg/log"
)

func main() {
	c, err := config.Load()
	if err != nil {
		fmt.Printf("Config init fail: %s", err.Error())
		return
	}

	if err := log.Init(c.Log.Out); err != nil {
		fmt.Printf("Log init fail: %s", err.Error())
		return
	}

	if err := email.Init(c.Email.ApiKey, c.Email.Domain); err != nil {
		log.Critical(context.Background(), herrors.New("Email init fail", "error", err))
		return
	}

	if err := email.Send("Subject", "Message"); err != nil {
		log.Critical(context.Background(), herrors.New("Email send fail", "error", err))
		return
	}

	fmt.Println("Done")
}
