-- name: CreateEntries :one
insert into entries (
    account_id,
    amount
) values ($1,$2) returning *;

-- name: GetEntries :one
select * from entries where id = $1 limit 1;

-- name: ListEntries :many
select * from entries where account_id = $1 order by id desc
limit $2 offset $3;