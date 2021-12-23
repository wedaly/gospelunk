package db

//go:generate protoc --go_out=. db.proto

import (
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"

	pb "github.com/wedaly/gospelunk/db/protobuf"
)

var pkgBucketName []byte

func init() {
	pkgBucketName = []byte("packages")
}

// DB stores a search index on disk.
type DB struct {
	boltDB *bolt.DB
}

// OpenReadWrite opens a database for both reading and writing.
// This creates the database if it does not yet exist.
// The caller is responsible for calling `db.Close()` when finished.
func OpenReadWrite(path string) (*DB, error) {
	boltDB, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "bolt.Open")
	}

	err = boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(pkgBucketName)
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "boltDB.Update")
	}

	return &DB{boltDB}, nil
}

// OpenReadOnly opens a database for reading, but not writing.
// The caller is responsible for calling `db.Close()` when finished.
func OpenReadOnly(path string) (*DB, error) {
	boltDB, err := bolt.Open(path, 0600, &bolt.Options{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrapf(err, "bolt.Open")
	}
	return &DB{boltDB}, nil
}

// Close closes the database.
func (db *DB) Close() {
	db.boltDB.Close()
}

// WritePackage writes serialized package data to the database.
// It overwrites any existing data for the same package directory.
func (db *DB) WritePackage(pkg *pb.Package) error {
	key := []byte(pkg.Dir)
	data, err := proto.Marshal(pkg)
	if err != nil {
		return errors.Wrapf(err, "proto.Marshal")
	}

	return db.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(pkgBucketName)
		if err := bucket.Put(key, data); err != nil {
			return errors.Wrapf(err, "bucket.Put")
		}
		return nil
	})
}

// ReadPackage reads package data from the database.
// If no package exists for the specified directory, it returns nil.
func (db *DB) ReadPackage(pkgDir string) (*pb.Package, error) {
	var pkg pb.Package
	key := []byte(pkgDir)
	err := db.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(pkgBucketName))
		data := bucket.Get(key)
		if data != nil {
			if err := proto.Unmarshal(data, &pkg); err != nil {
				return errors.Wrapf(err, "proto.Unmarshal")
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "boltDB.View")
	}
	return &pkg, nil
}
