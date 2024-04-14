# Installation du client

sudo apt install redis-tools

# Utilisatio du client

## Alternative 1

redis-cli

## Alternative 2

docker exec -ti redis-master redis-cli

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

## Exercice 5

SET userData "name:John,age:30,city:Paris"
SET userData "name:John,age:30,city:Lyon"

# Ensembles

SET greeting "Hello, World!"

## Exercice 1

SADD fruits pomme banane cerise
SISMEMBER fruits banane
SMEMBERS fruits
SREM fruits cerise

## Exercice 2

SADD ensembleA 1 2 3
SADD ensembleB 3 4 5
SINTER ensembleA ensembleB
SINTERSTORE intersection ensembleA ensembleB //Variante pour stocker
SUNION ensembleA ensembleB
SUNIONSTORE union ensembleA ensembleB //Variante pour stocker
SDIFF ensembleA ensembleB
SDIFFSTORE difference ensembleA ensembleB //Variante pour stocker

## Exercice 3

ZADD leaderboard 100 joueur1 150 joueur2 120 joueur3
ZINCRBY leaderboard 50 joueur1
ZRANK leaderboard joueur2
ZREVRANGE leaderboard 0 -1 WITHSCORES
ZREMRANGEBYSCORE leaderboard -inf (130

# Hashages

## Exercice 1

HSET utilisateur:100 nom "Alice" age "30" email "alice@example.com"
HGET utilisateur:100 email
HSET utilisateur:100 age "31"
HDEL utilisateur:100 email

# Exercice 2

HSET utilisateur:101 nom "Bob" age "25" email "bob@example.com"
HSET utilisateur:102 nom "Carol" age "28" email "carol@example.com"
HSET utilisateur:103 nom "Dave" age "32" email "dave@example.com"
HINCRBY utilisateur:101 age 1

# Listes

## Exercice 1

RPUSH maListeToDo "Faire les courses" "Apprendre Redis" "Faire du sport"
LINDEX maListeToDo 0
RPOP maListeToDo

## Exercice 2

LPUSH maListeToDo "Lire un livre"
LRANGE maListeToDo 0 -1
LREM maListeToDo 0 "Faire du sport"

# Données binaires

## Exercice 1

SETBIT presence:2024-04-15 5 1 // à répéter pour chaque employé et pour chaque jour
BITCOUNT presence:2024-04-15
BITOP AND presence:2024-04-all presence:2024-04-01 presence:2024-04-02 ... presence:2024-04-30
BITCOUNT presence:2024-04-all


## Exercice 2

BITFIELD game:scores SET u10 #3 500
BITFIELD game:scores GET u10 #3
// Pour la dernière question, passer par un script ou l'utilisation d'API en Go.


# Transactions

## Exercice 1

MULTI
SET transactionTest "test"
APPEND transactionTest "ing"
EXEC

## Exercice 2

MULTI
HSET utilisateur:101 ville "Paris"
HSET utilisateur:102 ville "Paris"
HSET utilisateur:103 ville "Paris"
EXEC

# Scripts

## Exercice 1

EVAL "return redis.call('SET', 'combined', redis.call('GET', KEYS[1]) .. redis.call('GET', KEYS[2]))" 2 greeting note
GET combined

## Exercice 2

local threshold = tonumber(ARGV[1])
local product_keys = redis.call('KEYS', 'produit:*')
local matching_products = {}

for i, key in ipairs(product_keys) do
    local price = tonumber(redis.call('HGET', key, 'prix'))
    if price and price < threshold then
        table.insert(matching_products, key)
    end
end

return matching_products

## Exercice 3

local key = KEYS[1] -- La clé de la liste Redis représentant la file d'attente
local message = redis.call('LPOP', key) -- Récupère et supprime le premier élément de la liste

-- Vérifie si un message a été récupéré
if message then
    -- Simule le traitement du message (dans cet exemple, nous allons simplement renvoyer le message traité)
    return message
else
    -- Si aucun message n'est disponible, renvoie une indication que la file d'attente est vide
    return 'La file d\'attente est vide'
end

EVAL "script" 1 fileDattenteMessages

