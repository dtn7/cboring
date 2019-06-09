package cboring

type MajorType = byte

const (
	UInt       MajorType = 0
	NegInt     MajorType = 1
	ByteString MajorType = 2
	TextString MajorType = 3
	Array      MajorType = 4
	Map        MajorType = 5
	Tagging    MajorType = 6
	Etc        MajorType = 7
)

func ParseMajor(b byte) (major MajorType, adds byte) {
	major = b >> 5
	adds = b & 0x1F
	return
}
