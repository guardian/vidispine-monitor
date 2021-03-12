package vsstoragecheck

type VSStorageMethod struct {
	Id          string `json:"id"`
	Uri         string `json:"uri"`
	Read        bool   `json:"read"`
	Write       bool   `json:"write"`
	Browse      bool   `json:"browse"`
	LastSuccess string `json:"lastSuccess"`
	Type        string `json:"type"`
}

// see https://apidoc.vidispine.com/latest/storage/storage.html#id1
type VSStorageType string

const (
	LocalStorage    VSStorageType = "LOCAL"
	SharedStorage                 = "SHARED"
	RemoteStorage                 = "REMOTE"
	ExternalStorage               = "EXTERNAL"
	ArchiveStorage                = "ARCHIVE"
	ExportStorage                 = "EXPORT"
)

type VSStorageState string

const (
	StorageStateNone       VSStorageState = "NONE"
	StorageStateReady                     = "READY"
	StorageStateOffline                   = "OFFLINE"
	StorageStateFailed                    = "FAILED"
	StorageStateDisabled                  = "DISABLED"
	StorageStateEvacuating                = "EVACUATING"
	StorageStateEvacuated                 = "EVACUATED"
)

type VSStorage struct {
	Id            string            `json:"id"`
	State         VSStorageState    `json:"state"`
	Type          VSStorageType     `json:"type"`
	Capacity      int64             `json:"capacity"`
	FreeCapacity  int64             `json:"freeCapacity"`
	Timestamp     string            `json:"timestamp"`
	Methods       []VSStorageMethod `json:"method"`
	LowWatermark  int64             `json:"lowWatermark"`
	HighWatermark int64             `json:"highWatermark"`
}

type VSStoragesResponse struct {
	Storage []VSStorage `json:"storage"`
}
