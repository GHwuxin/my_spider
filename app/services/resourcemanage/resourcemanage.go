package resourcemanage

type ResourceManage interface {
	GetOne()
	FreeOne()
	Has() uint
	Left() uint
}
