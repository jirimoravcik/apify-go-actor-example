package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	apify "github.com/apify/apify-client-go"
)

var (
	KVStoreDefaultID string
	Token            string
	xlog             = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

type (
	KVStoreValue map[string]any
	KVStore      struct {
		ID     string
		Client *apify.KeyValueStoreClient
	}
)

func init() {
	Token = os.Getenv("APIFY_TOKEN")
	KVStoreDefaultID = os.Getenv("APIFY_DEFAULT_KEY_VALUE_STORE_ID")
	if strings.EqualFold(Token, "") || strings.EqualFold(KVStoreDefaultID, "") {
		xlog.Error("token or default kv store missing")
	}
}

func KVStoreNew(id string) KVStore {
	client := apify.NewClient(Token)
	return KVStore{ID: id, Client: client.KeyValueStore(id)}
}

func KVStoreDefault() KVStore {
	return KVStoreNew(KVStoreDefaultID)
}

func (kv KVStore) Get(key string) (KVStoreValue, error) {
	record, ok, err := kv.Client.GetRecord(context.Background(), key)
	if err != nil {
		xlog.Error("failed to get value", "key", key, "error", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("record %q not found in key-value store", key)
	}
	var value KVStoreValue
	if err := json.Unmarshal(record.Value, &value); err != nil {
		xlog.Error("failed to unmarshal record", "key", key, "error", err)
		return nil, err
	}
	return value, nil
}

func (kv KVStore) Put(key string, payload any) error {
	if err := kv.Client.SetRecordJSON(context.Background(), key, payload); err != nil {
		xlog.Error("failed to set value", "key", key, "error", err)
		return err
	}
	return nil
}

func (kv KVStore) Delete(key string) error {
	if err := kv.Client.DeleteRecord(context.Background(), key); err != nil {
		xlog.Error("failed to delete key", "key", key, "error", err)
		return err
	}
	return nil
}
