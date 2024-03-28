package constraints

const (
	AdvertsOnPage = 10

	AdvertHeaderMaxLen = 100
	AdvertBodyMaxLen   = 600
	MaxPrice           = 100000000000
	MinPrice           = 0

	ImageMaxWidth  = 1080
	ImageMinWidth  = 140
	ImageMaxHeight = 720
	ImageMinHeight = 60

	LoginMinLen    = 5
	LoginMaxLen    = 20
	PasswordMinLen = 8
)

var ImageExtentions = map[string]bool{
	".jpeg": true,
	".jpg":  true,
	".png":  true,
	".gif":  true,
}
