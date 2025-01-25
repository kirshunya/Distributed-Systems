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
	Message string `json:"message"`
}

type CloseCommandData struct {
	Message string `json:"message"`
}

type UploadCommandData struct {
	FileName string `json:"file_name"`
	Content  string `json:"content"`
}

type DownloadCommandData struct {
	FileName string `json:"file_name"`
}
