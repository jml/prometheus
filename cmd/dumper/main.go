// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"time"

	"github.com/prometheus/common/log"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/storage/local"
	"github.com/prometheus/prometheus/storage/remote"
	"github.com/prometheus/prometheus/storage/remote/generic"
)

func main() {
	// Storage.
	path := flag.String("storage.local.path", "data", "Base path for metrics storage.")
	genericURL := flag.String("storage.remote.generic-url", "http://localhost:9094/push", "The URL of the generic remote server to send samples to.")
	flag.Parse()

	c := generic.NewClient(*genericURL, 30*time.Second)
	qm := remote.NewStorageQueueManager(c, 100*1024, remote.Backlog)
	go qm.Run()

	storageOpts := &local.MemorySeriesStorageOptions{
		SyncStrategy:               local.Adaptive,
		PersistenceStoragePath:     *path,
		MemoryChunks:               1e6,
		PersistenceRetentionPeriod: 1e6 * time.Hour,
		MaxChunksToPersist:         1e6,
		CheckpointInterval:         1e6 * time.Hour,
		CheckpointDirtySeriesLimit: 1e6,
		MinShrinkRatio:             1,
	}
	memStorage := local.NewMemorySeriesStorage(storageOpts)
	sentSamples := 0
	err := memStorage.Dump(func(s *model.Sample) {
		qm.Append(s)
		sentSamples++
		if sentSamples%1000 == 0 {
			log.Infof("Queued %d samples for sending...", sentSamples)
		}
	})
	if err != nil {
		log.Fatal("error dumping storage:", err)
	}
	log.Infoln("Flushing remote storage queue...")
	qm.Stop()
	log.Infof("Sent %d samples...", sentSamples)
}
