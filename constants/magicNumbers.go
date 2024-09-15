package constants

type FileType struct {
	Ext  string
	Mime string
}

type MagicNumber string

var MagicNumbers = map[MagicNumber]FileType{
	"89504E47": {Ext: "png", Mime: "image/png"},
	"47494638": {Ext: "gif", Mime: "image/gif"},
	"25504446": {Ext: "pdf", Mime: "application/pdf"},
	"FFD8FFE0": {Ext: "jpg", Mime: "image/jpeg"},
	"424D":     {Ext: "bmp", Mime: "image/bmp"},
	"4949":     {Ext: "tif", Mime: "image/tiff"},
	"38425053": {Ext: "psd", Mime: "image/vnd.adobe.photoshop"},
	"252150532D41646F62652D332E3020455053462D332030": {Ext: "eps", Mime: "application/postscript"},
	"504B0304":                 {Ext: "docx", Mime: "application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
	"D0CF11E0A1B11AE1":         {Ext: "ppt", Mime: "application/vnd.ms-powerpoint"},
	"CAFEBABE":                 {Ext: "class", Mime: "application/java-vm"},
	"0000000C6A5020200D0A":     {Ext: "jp2", Mime: "image/jp2"},
	"52494646":                 {Ext: "wav", Mime: "audio/wav"},
	"D7CDC69A":                 {Ext: "wmf", Mime: "image/wmf"},
	"4D546864":                 {Ext: "mid", Mime: "audio/midi"},
	"00000100":                 {Ext: "ico", Mime: "image/vnd.microsoft.icon"},
	"494433":                   {Ext: "mp3", Mime: "audio/mpeg"},
	"465753":                   {Ext: "swf", Mime: "application/x-shockwave-flash"},
	"464C56":                   {Ext: "flv", Mime: "video/x-flv"},
	"00000018667479706D703432": {Ext: "mp4", Mime: "video/mp4"},
	"6D6F6F76":                 {Ext: "mov", Mime: "video/quicktime"},
	"3026B2758E66CF":           {Ext: "wmv", Mime: "video/x-ms-wmv"},
	"1F8B08":                   {Ext: "gz", Mime: "application/gzip"},
	"7573746172":               {Ext: "tar", Mime: "application/x-tar"},
	"4C01":                     {Ext: "obj", Mime: "application/octet-stream"},
	"4D534346":                 {Ext: "cab", Mime: "application/vnd.ms-cab-compressed"},
	"4D5A":                     {Ext: "exe", Mime: "application/x-msdownload"},
	"526172211A0700":           {Ext: "rar", Mime: "application/x-rar-compressed"},
	"3F5F0300":                 {Ext: "hlp", Mime: "application/winhlp"},
	"4B444D56":                 {Ext: "vmdk", Mime: "application/octet-stream"},
	"2142444E42":               {Ext: "pst", Mime: "application/octet-stream"},
	"7B5C72746631":             {Ext: "rtf", Mime: "application/rtf"},
	"5374616E64617264204A6574": {Ext: "mdb", Mime: "application/x-msaccess"},
	"2521":                     {Ext: "ps", Mime: "application/postscript"},
	"504B0304140008000800":     {Ext: "jar", Mime: "application/java-archive"},
	"4D6963726F736F66742056697375616C2053747564696F20536F6C7574696F6E2046696C65": {Ext: "sln", Mime: "text/plain"},
	"789C": {Ext: "zlib", Mime: "application/octet-stream"},
	// Add more magic numbers and corresponding file types
}
