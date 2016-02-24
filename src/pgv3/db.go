package pgv3

import (
    "gopkg.in/pg.v3"
)

type User struct {
    Id   int64   `json:"id"`
    Name string  `json:"name"`
    Type string  `json:"type"`
}

type Relationship struct {
    Id      int64   `json:"-"`
    OtherId int64   `json:"id"`
    State   string  `json:"state"`
    Type    string  `json:"type"`
}

func CreateUser(db *pg.DB, user *User) error {
    _, err := db.QueryOne(user,
        `INSERT INTO users (name) VALUES (?name) RETURNING id`,
        user)
    return err
}

func GetUsers(db *pg.DB) ([]User, error) {
    users := []User{}
    _, err := db.Query(&users, `SELECT * FROM users`)
    return users, err
}

func UpdateUserRelationship(db *pg.DB, rel *Relationship) error {
    tx, err := db.Begin()
	if err != nil {
        return err
    }
    _, err = db.QueryOne(rel,
        `SELECT * FROM insert_or_update_relationship(?id, ?other_id, ?state)`,
        rel)
    if err != nil {
        tx.Rollback()
        return err
    }
    return tx.Commit()
}

func GetUserRelationships(db *pg.DB, id int64) ([]Relationship, error) {
    rels := []Relationship{}
    _, err := db.Query(&rels,
        `SELECT * FROM relationships WHERE id=?`,
        id)
    return rels, err
}

func CreateSchema(db *pg.DB) error {
    queries := []string{
        //`DROP TABLE IF EXISTS relationships`,
        //`DROP TABLE IF EXISTS users`,
        `CREATE TABLE IF NOT EXISTS users (
            id bigserial, name text,
            PRIMARY KEY(id)
        )`,
        `CREATE TABLE IF NOT EXISTS relationships (
            id bigint, other_id bigint, state text,
            PRIMARY KEY(id, other_id),
            CONSTRAINT fk_1
                FOREIGN KEY (id)
                REFERENCES users (id)
                ON DELETE NO ACTION
                ON UPDATE NO ACTION,
            CONSTRAINT fk_2
                FOREIGN KEY (other_id)
                REFERENCES users (id)
                ON DELETE NO ACTION
                ON UPDATE NO ACTION
        )`,
    }
    for _, q := range queries {
        _, err := db.Exec(q)
        if err != nil {
            return err
        }
    }
    return nil
}

func ConnectDatabase() *pg.DB {
    db := pg.Connect(&pg.Options{
        User: "postgres",
        Password: "123456",
    })
    return db
}

