package common

type VideoStatus string

const (
	VideoStatusActive  VideoStatus = "ACTIVE"
	VideoStatusDeleted VideoStatus = "DELETED"
)

type Provider string

const (
	ProviderLocal Provider = "LOCAL"
	ProviderS3    Provider = "S3"
)
