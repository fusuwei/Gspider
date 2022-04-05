package constant

const (
	ProducerModel = "Producer"
	ConsumerModel = "Consumer"
	//RUNNING       status        = "running"
	//ERROR         status        = "error"
	//WAITING       status        = "waiting"
)

type Node int

const (
	EmptyNode    Node = iota // 空节点
	ProducerNode             // 生产节点，用于生产单rabbitMq和其他队列中（选用这个节点必须实现queue）
	RequestNode              // 请求节点
	CallbackNode             // 回调节点，主要用于解析请求返回的数据
	SaveNode                 // 存储节点
)
