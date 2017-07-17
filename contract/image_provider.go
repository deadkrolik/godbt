package contract

//ImageProvider - transformer of raw input to Image
type ImageProvider interface {
	Parse(interface{}) (Image, error)
	CanParse(source interface{}) bool
}
