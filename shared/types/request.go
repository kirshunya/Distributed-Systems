package types

type Request[T any] struct {
	CommandType int `json:"command_type"`
	Data        T   `json:"data"`
}

const (
	ECHO = iota
	TIME
	CLOSE
	UPLOAD
	DOWNLOAD
)

type EchoCommandData struct {
	Message string `json:"message"`
}

type TimeCommandData struct {
}

type CloseCommandData struct {
}

type UploadCommandData struct {
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	Status   string `json:"status"`
}

type DownloadCommandData struct {
	FileName string `json:"file_name"`
	Status   string `json:"status"`
}
