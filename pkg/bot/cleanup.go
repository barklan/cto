package bot

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/barklan/cto/pkg/storage"
	"github.com/dgraph-io/badger/v3"
	tb "gopkg.in/tucnak/telebot.v2"
)

// FIXME THIS IS DEPRECATED - do not use!
func CleanUp(data *storage.Data) {
	keysToDelete := make([][]byte, 0)
	tgMessagesToDelete := make([]tb.Message, 0)
	err := data.DB.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			domain := strings.Split(string(k), "-")[0]
			if domain == "botMsg" {
				keyCopy := append([]byte{}, k...)
				keysToDelete = append(keysToDelete, keyCopy)
				err := item.Value(func(v []byte) error {
					msgBytes := append([]byte{}, v...)
					var msg tb.Message
					err := json.Unmarshal(msgBytes, &msg)
					if err != nil {
						log.Printf("fucker")
					}
					tgMessagesToDelete = append(tgMessagesToDelete, msg)
					return nil
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Println("Cleanup Failed.")
		log.Println(err)
	}

	ticker := time.NewTicker(1 * time.Second)
	for _, msg := range tgMessagesToDelete {
		data.B.Delete(&msg)
		<-ticker.C
	}

	wb := data.DB.NewWriteBatch()
	defer wb.Cancel()
	for _, key := range keysToDelete {
		log.Println("Deleting key:", string(key))
		err := wb.Delete(key) // Will create txns as needed.
		if err != nil {
			log.Println("Failed to delete key:", key)
		}
	}
	err = wb.Flush() // Wait for all txns to finish.
	if err != nil {
		log.Println("DB flush failed.")
	}

	log.Println("CleanUP complete.")
}
