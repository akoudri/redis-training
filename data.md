# Authentication

AUTH Jn4haJxjuEDIO

# String

## Exercice 1

SET greeting "Hello"
GET greeting
SET greeting "Hello, World!"
STRLEN greeting

## Exercice 2

SET counter 0
INCRBY counter 10
DECRBY counter 3
GET counter

## Exercice 3

SET note "Ceci est une "
APPEND note " note"
APPEND note " de Redis"

## Exercice 4

GETRANGE greeting 0 4
GETRANGE greeting 7 -1

# Exercice 5

SET userData "name:John,age:30,city:Paris"
SET userData "name:John,age:30,city:Lyon"

# Transactions

## Exercice 1

MULTI
SET transactionTest "test"
APPEND transactionTest "ing"
EXEC

# Scripts

EVAL "return redis.call('SET', 'combined', redis.call('GET', KEYS[1]) .. redis.call('GET', KEYS[2]))" 2 greeting note
GET combined