package mydocstore

type Item struct {
	MyHashKey  string `docstore:"MyHashKey"`
	MyRangeKey int    `docstore:"MyRangeKey"`
	MyText     string `docstore:"MyText"`
}
