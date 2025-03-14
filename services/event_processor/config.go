package main

import "flag"

func GetDBCredentials() string {
	var dbCredentials string
	flag.StringVar(&dbCredentials, "db", "", "Database credentials")
	flag.Parse()
	return dbCredentials
}
