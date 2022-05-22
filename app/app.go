package app

import (
	"BeachAPI/api"
	//"github.com/mailgun/mailgun-go/v4"
	//"github.com/mailgun/mailgun-go/v4"
	"BeachAPI/app/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

//var MG *mailgun.MailgunImpl
var CFG *config.Config
var DB *sqlx.DB

func init() {
	CFG = config.InitCfg()
	//mg := mailgun.NewMailgun(CFG.Mailgun.Domain, CFG.Mailgun.PrivateAPIKey)
	//	mg.SetAPIBase(mailgun.APIBaseEU)
	//	MG = mg
	conn := `
           host=` + CFG.Db.Host + `
         dbname=` + CFG.Db.DbName + `
		   user=` + CFG.Db.UserName + `
        sslmode=disable
		   port=` + CFG.Db.Port + `
		password=` + CFG.Db.Password + `
`
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	DB = db
	initDb()
}

func initDb() {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS users
(
    id       serial not null primary key,
    phone_number    varchar,
 	role_id  INTEGER REFERENCES roles (id),
    balance numeric(18,2) default 0,
    balance_bonus numeric(18,2) default 0,
    confirmed boolean default false
);`)

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS roles
(
    id       serial not null primary key,
    name    varchar
);`)
	//role
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS balance_replenishment
(
    id       serial not null primary key,
    user_id  INTEGER REFERENCES users (id),
    sum    numeric(18,2),
	created timestamp
);`)

	api.CheckErrInfo(err, "init db 1")

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS codes
(
    id       serial not null primary key,
    user_id  INTEGER REFERENCES users (id),
    code    varchar,
	created timestamp
);`)

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS tokens
(
    id       serial not null primary key,
    user_id  INTEGER REFERENCES users (id),
    token    varchar,
	created timestamp
);`)

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS shops
(
    id       serial not null primary key,
    "name" varchar,
        description varchar default '',
    image varchar default 'https://alnezis.riznex.ru:1337/images/upload-191055778.png',
    balance numeric(18,2),
    user_id  INTEGER REFERENCES users (id),
    status integer default 0,
	created timestamp
);`)

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS withdrawal_balance
(
    id       serial not null primary key,
    shop_id  INTEGER REFERENCES shops (id),
    sum    numeric(18,2),
    status boolean default false,
	created timestamp
);`)

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS items
(
    id       serial not null primary key,
    shop_id  INTEGER REFERENCES shops (id),
    "name"    varchar,
    description varchar default '',
    count integer,
    price numeric(18,2),
    image varchar default 'https://alnezis.riznex.ru:1337/images/upload-191055778.png',
    status integer default 0,
    time_updated timestamp default now(),
	created timestamp default now()

);`) //status integer,  0, на модерации // 1 подтвержд // 2 отклонен

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS buys
(
    id        serial
        primary key,
    shop_id  INTEGER REFERENCES shops (id),
    item_id  INTEGER REFERENCES items (id),
    user_id  INTEGER REFERENCES users (id),
    price     double precision,
    created   timestamp default now(),
    hash_card varchar,
    count     integer
);`)

	api.CheckErrInfo(err, "init db 2")

}
