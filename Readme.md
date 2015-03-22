# context

[![Build Status](https://travis-ci.org/gomon/gomon.svg?branch=master)](https://travis-ci.org/gomon/gomon)
[![GoDoc](https://godoc.org/gopkg.in/mong-go/context.v1?status.svg)](http://godoc.org/gopkg.in/mong-go/context.v1)

Mgo session clone middleware using Gorilla/Context

## Install

  go get gopkg.in/mong-go/context.v1

## Usage

*Usage below uses Alice as the middlware library*

    import mgocontext "gopkg.in/mong-go/context.v1"

    session, err := mgo.Dial("127.0.0.1:27017")
    if err != nil {
      log.Fatal(err)
    }

    chain := alice.New(mgocontext.Handler(session, "db_name"), ...)
    ...

---

Context key will default the database name argument.

    mgocontext.Handler(session, "db_name")

    // db := context.Get(req, "db_name").(*mgo.Database)

Or you can provide a 3rd argument for a custom key to be used.

    mgocontext.Handler(session, "db_name", "myDatabase")

    // db := context.Get(req, "myDatabase").(*mgo.Database)

## License 

MIT
