package models

type PageItems struct {
	req   *Request          // The req is Request object that contains the parsed result, which saved in PageItems.
	items map[string]string // The items is the container of parsed result.
	skip  bool              // The skip represents whether send ResultItems to scheduler or not.
}

// NewPageItems returns initialized PageItems object.
func NewPageItems(req *Request) *PageItems {
	items := make(map[string]string)
	return &PageItems{req: req, items: items, skip: false}
}

// SetSkip set skip true to make this page not to be processed by Pipeline.
func (this *PageItems) SetSkip(skip bool) *PageItems {
	this.skip = skip
	return this
}

// Request returns request of PageItems
func (this *PageItems) Request() *Request {
	return this.req
}

// Item returns value of the key.
func (this *PageItems) Item(key string) (string, bool) {
	t, ok := this.items[key]
	return t, ok
}

// All returns all the KVs result.
func (this *PageItems) All() map[string]string {
	return this.items
}

// Skip returns skip label.
func (this *PageItems) Skip() bool {
	return this.skip
}

// AddItem saves a KV result into PageItems.
func (this *PageItems) AddItem(key string, item string) {
	this.items[key] = item
}
