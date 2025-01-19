package main

import (
	"database/sql"
	"log"
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/chriskoorzen/go-subscription-webapp/cmd/web/db"
)

type Config struct {
	Session       *scs.SessionManager
	DB            *sql.DB
	InfoLog       *log.Logger
	SuccessLog    *log.Logger
	ErrorLog      *log.Logger
	Wait          *sync.WaitGroup
	Models        db.Models
	Mailer        Mail
	ErrorChan     chan error
	ErrorChanDone chan bool
}
