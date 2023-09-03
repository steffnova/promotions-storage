package promotion

type Stream <-chan Promotion
type Streamer func() (Stream, <-chan error)
