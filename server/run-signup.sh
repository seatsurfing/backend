#!/bin/sh
curl -X POST \
    -v \
    -H "Content-Type: application/json" \
    -d '{"firstname": "", "lastname": "", "email": "foo@bar.com",  "organization": "Test Org",  "domain": "testorg",  "contactFirstname": "Foo",  "contactLastname": "Bar",  "password": "12345678",  "country": "DE", "language": "de", "acceptTerms": true}' \
    http://localhost:8080/signup/