package headers

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

//Suffix = parse full
// Prefix = done 

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

}
